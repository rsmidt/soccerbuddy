package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	decimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/postgres"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

var (
	appendDurationHistogram metric.Float64Histogram
)

func init() {
	var err error
	appendDurationHistogram, err = tracing.Meter.Float64Histogram(
		"eventstore.append.duration",
		metric.WithDescription("The duration of the append operation."),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.05, 0.1, 1, 2.5, 5, 10),
	)
	if err != nil {
		panic(err)
	}
}

var uniqueConstraintRegex = regexp.MustCompile(`Key \(([^)]+)\)=\(([^)]+)\)`)

type pgEventStore struct {
	pool   *pgxpool.Pool
	mapper eventing.JournalEventMapper
	crypto eventing.EventCrypto
	rs     authz.RelationStore
	log    *slog.Logger
	hooks  []eventing.Hook
}

func NewEventStore(log *slog.Logger, pool *pgxpool.Pool, mapper eventing.JournalEventMapper, crypto eventing.EventCrypto, rs authz.RelationStore) eventing.EventStore {
	return &pgEventStore{
		pool:   pool,
		mapper: mapper,
		log:    log,
		crypto: crypto,
		rs:     rs,
	}
}

func (p *pgEventStore) Append(ctx context.Context, intents ...eventing.AggregateChangeIntent) error {
	ctx, span := tracing.Tracer.Start(ctx, "pg.EventStore.Append")
	defer span.End()

	start := time.Now()
	defer func() {
		appendDurationHistogram.Record(ctx, time.Since(start).Seconds())
	}()

	if len(intents) != 1 {
		return errors.New("only one intent is supported")
	}

	if len(intents[0].Events()) == 0 {
		return nil
	}

	err := pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		intent := intents[0]

		aggregateVersion, err := lockLatestAggregateVersion(ctx, tx, intent.AggregateID(), intent.AggregateType())
		if err != nil {
			return err
		}

		if !intent.VersionMatches(aggregateVersion) {
			p.log.Debug("Aggregate versions do not match", slog.Uint64("remote", uint64(aggregateVersion)), slog.Uint64("local", uint64(intent.LastKnownAggregateVersion())))
			return eventing.ErrVersionMismatch
		}

		err = p.handleUniqueConstraints(ctx, tx, intent)
		if err != nil {
			return err
		}

		err = p.handleLookups(ctx, tx, intent)
		if err != nil {
			return err
		}

		err = p.persistEvents(ctx, tx, intent)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	// Run all post persist hooks
	for _, hook := range p.hooks {
		if post, ok := hook.(eventing.PostPersist); ok {
			if err := post.PostPersist(ctx); err != nil {
				p.log.Error("Failed to run post persist hook", slog.String("err", err.Error()))
			}
		}
	}
	return nil
}

func (p *pgEventStore) ProduceAppend(ctx context.Context, producer eventing.ChangeProducer) error {
	ctx, span := tracing.Tracer.Start(ctx, "pg.EventStore.ProduceAppend")
	defer span.End()

	return p.Append(ctx, *producer.Changes())
}

func (p *pgEventStore) AddHook(hook eventing.Hook) {
	p.hooks = append(p.hooks, hook)
}

func (p *pgEventStore) queryConn(ctx context.Context, conn *pgx.Conn, query eventing.JournalQuery, opts ...eventing.QueryOpts) ([]*eventing.JournalEvent, error) {
	span := trace.SpanFromContext(ctx)
	var config eventing.QueryConfig
	for _, opt := range opts {
		config = opt.Apply(config)
	}

	var stmtBuilder strings.Builder
	stmtBuilder.WriteString("SELECT id, aggregate_id, aggregate_type, aggregate_version, global_position, event_type, event_version, payload, created_at FROM event_journal WHERE (")
	var args []any
	var argI int
	for aggregateType, aggregateQuery := range query.AggQueriesByType() {
		if argI > 0 {
			stmtBuilder.WriteString(" OR ")
		}
		stmtBuilder.WriteString(fmt.Sprintf("((aggregate_type = $%d)", argI+1))
		args = append(args, aggregateType)
		argI++

		if aggregateQuery.ID() != "" {
			stmtBuilder.WriteString(fmt.Sprintf(" AND (aggregate_id = $%d)", argI+1))
			args = append(args, aggregateQuery.ID())
			argI++
		}
		if aggregateQuery.Version() > 0 {
			stmtBuilder.WriteString(fmt.Sprintf(" AND (aggregate_version >= $%d)", argI+1))
			args = append(args, aggregateQuery.Version())
			argI++
		}
		if len(aggregateQuery.Events()) == 1 {
			stmtBuilder.WriteString(fmt.Sprintf(" AND (event_type = $%d)", argI+1))
			args = append(args, aggregateQuery.Events()[0])
			argI++
		} else if len(aggregateQuery.Events()) > 1 {
			stmtBuilder.WriteString(fmt.Sprintf(" AND (event_type = ANY ($%d))", argI+1))
			args = append(args, aggregateQuery.Events())
			argI++
		}
		stmtBuilder.WriteString(")")
	}
	stmtBuilder.WriteString(")")
	if query.JournalPositionAfter() != nil {
		if argI > 0 {
			stmtBuilder.WriteString(" AND ")
		}
		stmtBuilder.WriteString(fmt.Sprintf("global_position > $%d", argI+1))
		args = append(args, query.JournalPositionAfter().Deref())
		argI++
	}

	// Only return rows that have been inserted before the oldest running transaction started.
	// Again, many thanks to Zitadel for pointing that out.
	if config.LimitToOldestRunningTransaction {
		stmtBuilder.WriteString(" AND global_position < (")
		stmtBuilder.WriteString(oldestRunningTransactionQuery)
		stmtBuilder.WriteString(")")
	}

	// Order by global position.
	stmtBuilder.WriteString(" ORDER BY global_position ASC")

	stmt := stmtBuilder.String()
	span.SetAttributes(semconv.DBQueryText(stmt))
	rows, err := conn.Query(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	journalEvents, err := postgres.CollectRowsNonNil(rows, func(row pgx.CollectableRow) (*eventing.JournalEvent, error) {
		var (
			id               eventing.EventID
			aggregateID      eventing.AggregateID
			aggregateType    eventing.AggregateType
			aggregateVersion eventing.AggregateVersion
			globalPosition   decimal.Decimal
			eventType        eventing.EventType
			eventVersion     eventing.EventVersion
			payload          []byte
			createdAt        time.Time
		)
		err := row.Scan(&id, &aggregateID, &aggregateType, &aggregateVersion, &globalPosition, &eventType, &eventVersion, &payload, &createdAt)
		if err != nil {
			return nil, err
		}
		event, err := p.mapper.MapFrom(aggregateID, aggregateType, eventVersion, eventType, id, aggregateVersion, eventing.JournalPosition(globalPosition), createdAt, payload)
		if err != nil {
			return nil, err
		}
		return event, err
	})
	if err != nil {
		return nil, err
	}
	events := make([]eventing.Event, len(journalEvents))
	for i, journalEvent := range journalEvents {
		events[i] = journalEvent.Event
	}
	err = p.crypto.DecryptEvents(ctx, events)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt events: %w", err)
	}
	return journalEvents, nil
}

func (p *pgEventStore) Query(ctx context.Context, query eventing.JournalQuery, opts ...eventing.QueryOpts) ([]*eventing.JournalEvent, error) {
	ctx, span := tracing.Tracer.Start(ctx, "pg.EventStore.Query")
	defer span.End()

	// Check if there's already a connection in the context we can use.
	conn, ok := ctx.Value("conn").(*pgx.Conn)
	if ok {
		return p.queryConn(ctx, conn, query, opts...)
	}

	// If not, acquire a connection from the pool.
	var result []*eventing.JournalEvent
	err := p.pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		events, err := p.queryConn(ctx, conn.Conn(), query, opts...)
		if err != nil {
			return err
		}
		result = events
		return nil
	})
	return result, err
}

func (p *pgEventStore) View(ctx context.Context, view eventing.JournalViewer) error {
	ctx, span := tracing.Tracer.Start(ctx, "pg.EventStore.View")
	defer span.End()

	query := view.Query()
	events, err := p.Query(ctx, query)
	if err != nil {
		return err
	}
	view.Reduce(events)
	return nil
}

func (p *pgEventStore) Lookup(ctx context.Context, opts eventing.LookupOpts) (*eventing.LookupFieldValue, error) {
	ctx, span := tracing.Tracer.Start(ctx, "pg.EventStore.Lookup")
	defer span.End()

	stmt := "SELECT field_value FROM event_journal_lookup WHERE owner_aggregate_type = $1 AND field_name = $2"
	row, err := p.pool.Query(ctx, stmt, opts.AggregateType, opts.FieldName)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup: %w", err)
	}
	val, err := pgx.CollectOneRow(row, func(row pgx.CollectableRow) (*eventing.LookupFieldValue, error) {
		var value eventing.LookupFieldValue
		err := row.Scan(&value)
		return &value, err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, eventing.ErrValueNotFound
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func (p *pgEventStore) OwnerLookup(ctx context.Context, opts eventing.LookupOpts) (eventing.AggregateID, error) {
	ctx, span := tracing.Tracer.Start(ctx, "pg.EventStore.OwnerLookup")
	defer span.End()

	stmt := "SELECT owner_aggregate_id FROM event_journal_lookup WHERE owner_aggregate_type = $1 AND field_name = $2 AND field_value = $3"
	row, err := p.pool.Query(ctx, stmt, opts.AggregateType, opts.FieldName, opts.FieldValue)
	if err != nil {
		return "", fmt.Errorf("failed to lookup owner: %w", err)
	}
	ownerID, err := pgx.CollectOneRow(row, func(row pgx.CollectableRow) (eventing.AggregateID, error) {
		var ownerID eventing.AggregateID
		err := row.Scan(&ownerID)
		return ownerID, err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return "", eventing.ErrOwnerNotFound
	} else if err != nil {
		return "", err
	}
	return ownerID, nil
}

func lockLatestAggregateVersion(ctx context.Context, tx pgx.Tx, aggregateID eventing.AggregateID, aggregateType eventing.AggregateType) (eventing.AggregateVersion, error) {
	ctx, span := tracing.Tracer.Start(ctx, "pg.lockLatestAggregateVersion")
	defer span.End()

	stmt := `
SELECT ej.aggregate_version
FROM event_journal ej
WHERE ej.aggregate_id = $1
AND ej.aggregate_type = $2
ORDER BY ej.aggregate_version DESC
LIMIT 1
FOR UPDATE
`
	rows, err := tx.Query(ctx, stmt, aggregateID, aggregateType)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "55P03" {
		return 0, eventing.ErrIntentOutdated
	} else if err != nil {
		return 0, fmt.Errorf("failed to lock latest aggregate version: %w", err)
	}
	version, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (eventing.AggregateVersion, error) {
		var aggregateVersion eventing.AggregateVersion
		err := row.Scan(&aggregateVersion)
		return aggregateVersion, err
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return version, nil
}

func (p *pgEventStore) handleUniqueConstraints(ctx context.Context, tx pgx.Tx, intent eventing.AggregateChangeIntent) error {
	ctx, span := tracing.Tracer.Start(ctx, "pg.handleUniqueConstraints")
	defer span.End()

	var addArgs []any
	var addStmtBuilder strings.Builder
	var addEventI int
	var hasAddStmt bool
	for _, event := range intent.Events() {
		adder, ok := event.(eventing.UniqueConstraintAdder)
		if !ok {
			continue
		}
		for _, toAdd := range adder.UniqueConstraintsToAdd() {
			hasAddStmt = true
			if addEventI > 0 {
				addStmtBuilder.WriteString(", ")
			} else {
				addStmtBuilder.WriteString("INSERT INTO unique_constraint (field, value, owner_aggregate_id) VALUES ")
			}

			addStmtBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d)", addEventI+1, addEventI+2, addEventI+3))
			addArgs = append(addArgs, toAdd.ConstrainedField(), toAdd.ConstrainedValue(), intent.AggregateID())
			addEventI += 3
		}
	}
	addStmt := addStmtBuilder.String()
	if hasAddStmt {
		_, err := tx.Exec(ctx, addStmt, addArgs...)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				field, value, err := extractViolatedConstraint(pgErr)
				if err != nil {
					return fmt.Errorf("failed to extract violated unique constraint: %w", err)
				}
				return eventing.NewUniqueConstraintError(eventing.NewUniqueConstraint(intent.AggregateID(), field, value))
			}

			return fmt.Errorf("failed to add unique constraints: %w", err)
		}
	}

	var rmArgs []any
	var rmStmtBuilder strings.Builder
	var rmEventI int
	var hasRmStmt bool
	for _, event := range intent.Events() {
		remover, ok := event.(eventing.UniqueConstraintRemover)
		if !ok {
			continue
		}
		for _, toRemove := range remover.UniqueConstraintsToRemove() {
			hasRmStmt = true
			if rmEventI > 0 {
				addStmtBuilder.WriteString(" OR ")
			} else {
				rmStmtBuilder.WriteString("DELETE FROM unique_constraint uc WHERE ")
			}

			if toRemove.ConstrainedField() == "" && toRemove.ConstrainedValue() == "" {
				rmStmtBuilder.WriteString(fmt.Sprintf("uc.owner_aggregate_id = $%d", rmEventI+1))
				rmArgs = append(addArgs, intent.AggregateID())
			} else {
				rmStmtBuilder.WriteString(fmt.Sprintf("(uc.field = $%d AND uc.value = $%d AND uc.owner_aggregate_id = $%d)", rmEventI+1, rmEventI+2, rmEventI+3))
				rmArgs = append(addArgs, toRemove.ConstrainedField(), toRemove.ConstrainedValue(), intent.AggregateID())
				rmEventI += 3
			}

		}
	}
	rmStmt := rmStmtBuilder.String()
	if !hasRmStmt {
		return nil
	}
	p.log.Debug("Removing unique constraints", "stmt", rmStmt, "args", rmArgs)
	_, err := tx.Exec(ctx, rmStmt, rmArgs...)
	if err != nil {
		return fmt.Errorf("failed to remove unique constraints: %w", err)
	}

	return nil
}

func extractViolatedConstraint(pgErr *pgconn.PgError) (field, value string, err error) {
	matches := uniqueConstraintRegex.FindStringSubmatch(pgErr.Detail)
	if len(matches) != 3 {
		return "", "", errors.New("regex failed")
	}
	values := strings.Split(matches[2], ",")
	valuesTrimmed := make([]string, 0, len(values))
	for _, value := range values {
		valuesTrimmed = append(valuesTrimmed, strings.TrimSpace(value))
	}
	return valuesTrimmed[0], valuesTrimmed[1], nil
}

func (p *pgEventStore) handleLookups(ctx context.Context, tx pgx.Tx, intent eventing.AggregateChangeIntent) error {
	ctx, span := tracing.Tracer.Start(ctx, "pg.handleLookups")
	defer span.End()

	for _, event := range intent.Events() {
		lookupProvider, ok := event.(eventing.LookupProvider)
		if ok {
			// TODO: Optimize this to use a single query.
			for fieldName, fieldValue := range lookupProvider.LookupValues() {
				id := idgen.NewString()
				stmt := "INSERT INTO event_journal_lookup (id, owner_aggregate_id, owner_aggregate_type, field_name, field_value) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (owner_aggregate_id, field_name) DO UPDATE SET field_value = $3"
				_, err := tx.Exec(ctx, stmt, id, event.AggregateID(), event.AggregateType(), fieldName, fieldValue)
				if err != nil {
					return fmt.Errorf("failed to insert lookups: %w", err)
				}
			}
		}
		lookupRemover, ok := event.(eventing.LookupRemover)
		if ok {
			for _, fieldName := range lookupRemover.LookupRemoves() {
				stmt := "DELETE FROM event_journal_lookup WHERE owner_aggregate_id = $1 AND field_name = $2"
				_, err := tx.Exec(ctx, stmt, intent.AggregateID(), fieldName)
				if err != nil {
					return fmt.Errorf("failed to remove lookups: %w", err)
				}
			}
		}
	}
	return nil
}

func (p *pgEventStore) persistEvents(ctx context.Context, tx pgx.Tx, intent eventing.AggregateChangeIntent) error {
	ctx, span := tracing.Tracer.Start(ctx, "pg.persistEvents")
	defer span.End()

	err := p.crypto.EncryptEvents(ctx, intent.Events())
	if err != nil {
		return fmt.Errorf("failed to encrypt events: %w", err)
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"event_journal"}, []string{"id", "aggregate_id", "aggregate_type", "aggregate_version", "event_type", "event_version", "payload"},
		pgx.CopyFromSlice(len(intent.Events()), func(i int) ([]interface{}, error) {
			event := intent.Events()[i]
			payload, err := json.Marshal(event)
			if err != nil {
				return nil, err
			}
			id := idgen.NewString()

			return []any{
				id,
				event.AggregateID(),
				event.AggregateType(),
				uint(intent.LastKnownAggregateVersion()) + uint(i+1),
				event.EventType(),
				event.EventVersion(),
				payload,
			}, nil
		}))

	return err
}

const oldestRunningTransactionQuery = `
SELECT COALESCE(EXTRACT(EPOCH FROM min(xact_start)), EXTRACT(EPOCH FROM now()))
FROM pg_stat_activity 
WHERE datname = current_database() 
	AND application_name = 'soccerbuddy' 
	AND state <> 'idle'
`
