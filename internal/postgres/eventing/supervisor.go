package eventing

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/postgres"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"github.com/shopspring/decimal"
	"log/slog"
	"sync"
	"sync/atomic"
)

const (
	pgCodeLockNotAvailable = "55P03"
)

var (
	errAlreadyLocked = errors.New("the projection state is already locked")

	_ eventing.EventListener = (*projectorSupervisor)(nil)
)

type PostgresProjector interface {
	eventing.Projector

	ProjectWithTx(ctx context.Context, tx pgx.Tx, events ...*eventing.JournalEvent) error
}

type projectorSupervisor struct {
	mu                  sync.RWMutex
	pool                *pgxpool.Pool
	es                  eventing.EventStore
	projectors          map[eventing.ProjectionName]eventing.Projector
	projectorByInterest map[eventing.EventInterest][]eventing.Projector
	log                 *slog.Logger
	started             atomic.Bool
}

func NewProjectorSupervisor(log *slog.Logger, pool *pgxpool.Pool, es eventing.EventStore) eventing.ProjectorSupervisor {
	return &projectorSupervisor{
		log:                 log,
		pool:                pool,
		es:                  es,
		projectors:          make(map[eventing.ProjectionName]eventing.Projector),
		projectorByInterest: make(map[eventing.EventInterest][]eventing.Projector),
	}
}

func (ps *projectorSupervisor) Register(projector eventing.Projector) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	interests := eventing.ProjectorToInterests(projector)
	for _, interest := range interests {
		ps.projectorByInterest[interest] = append(ps.projectorByInterest[interest], projector)
	}
	ps.projectors[projector.Projection()] = projector
}

func (ps *projectorSupervisor) Interests() eventing.EventInterestSet {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	interestSet := make(eventing.EventInterestSet)
	for interest := range ps.projectorByInterest {
		interestSet.Add(interest)
	}
	return interestSet
}

func (ps *projectorSupervisor) Notify(ctx context.Context, interests ...eventing.EventInterest) bool {
	if !ps.started.Load() {
		return false
	}
	ps.mu.RLock()
	ps.mu.RUnlock()

	var triggeredProjectors []eventing.Projector
	for _, interest := range interests {
		if projectors, ok := ps.projectorByInterest[interest]; ok {
			triggeredProjectors = append(triggeredProjectors, projectors...)
		}
	}
	if err := ps.trigger(ctx, false, triggeredProjectors...); err != nil {
		return false
	}
	return true
}

func (ps *projectorSupervisor) Enable() {
	ps.started.Store(true)
}

func (ps *projectorSupervisor) Trigger(ctx context.Context, projection ...eventing.ProjectionName) {
	ctx, span := tracing.Tracer.Start(ctx, "pg.ProjectorSupervisor.Trigger")
	defer span.End()

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	// Trigger all projectors if no projection is specified.
	var projectors []eventing.Projector
	if len(projection) == 0 {
		for _, projector := range ps.projectors {
			projectors = append(projectors, projector)
		}
	} else {
		for _, proj := range projection {
			projector, ok := ps.projectors[proj]
			if !ok {
				continue
			}
			projectors = append(projectors, projector)
		}
	}

	if err := ps.trigger(ctx, true, projectors...); err != nil {
		ps.log.Error("Failed to manually trigger projectors", slog.String("err", err.Error()))
		tracing.RecordError(ctx, err)
	}
}
func (ps *projectorSupervisor) trigger(ctx context.Context, wait bool, projectors ...eventing.Projector) (rerr error) {
	for _, projector := range projectors {
		projection := string(projector.Projection())

		ps.log.Debug("Advancing projector", slog.String("projection", projection))

		err := pgx.BeginFunc(ctx, ps.pool, func(tx pgx.Tx) error {
			ctx := postgres.WithTx(ctx, tx)

			// Enable by fetching and locking the projection state.
			state, err := ps.fetchProjectionState(ctx, wait, tx, projection)
			if errors.Is(err, pgx.ErrNoRows) {
				// Initialize the projection state if not found.
				state, err = ps.initProjectionState(ctx, tx, projection)
				if err != nil {
					return err
				}
			} else if errors.Is(err, errAlreadyLocked) {
				// The projection is currently being updated from somewhere else, we can skip.
				return nil
			} else if err != nil {
				return err
			}

			// Fetch the events from the event store.
			queryWithPos := eventing.NewJournalQueryBuilderFrom(projector.Query()).
				WithJournalPositionAfter(eventing.JournalPosition(state.GlobalPosition)).MustBuild()
			events, err := ps.es.Query(ctx, queryWithPos, eventing.WithLimitToOldestRunningTransaction())
			if err != nil {
				return err
			}

			txProjector, ok := projector.(PostgresProjector)
			if !ok {
				err = projector.Project(ctx, events...)
			} else {
				err = txProjector.ProjectWithTx(ctx, tx, events...)
			}
			if err != nil {
				return err
			}

			// Advance the projection state.
			return ps.updateProjectionState(ctx, tx, projection, events)
		})
		if err != nil {
			ps.log.
				With(slog.String("err", err.Error())).
				With(slog.String("projection", projection)).
				Error("Failed to advance projector")

			rerr = errors.Join(rerr, err)
		}
	}
	return
}

func (ps *projectorSupervisor) fetchProjectionState(ctx context.Context, wait bool, tx pgx.Tx, projection string) (*eventing.ProjectionState, error) {
	ctx, span := tracing.Tracer.Start(ctx, "pg.ProjectorSupervisor.fetchProjectionState")
	defer span.End()

	query := getProjectionStateSQL
	if wait {
		query = fmt.Sprintf("%s NOWAIT", query)
	}
	rows, err := tx.Query(ctx, getProjectionStateSQL, projection)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgCodeLockNotAvailable {
			return nil, errAlreadyLocked
		}
		return nil, err
	}
	state, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (eventing.ProjectionState, error) {
		var state eventing.ProjectionState
		var position decimal.Decimal
		err := row.Scan(&state.Name, &state.LastProcessedEventID, &state.LastProcessedTimestamp, &state.AggregateVersion, &position, &state.UpdatedAt)
		state.GlobalPosition = position
		return state, err
	})
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (ps *projectorSupervisor) updateProjectionState(ctx context.Context, tx pgx.Tx, projection string, events []*eventing.JournalEvent) error {
	ctx, span := tracing.Tracer.Start(ctx, "pg.ProjectorSupervisor.updateProjectionState")
	defer span.End()

	// If there are no events, just update the state timestamps.
	if len(events) == 0 {
		_, err := tx.Exec(
			ctx,
			updateProjectionStateEmptySQL,
			projection,
		)
		return err
	}

	// Find the last event.
	var lastEvent *eventing.JournalEvent
	for _, event := range events {
		if lastEvent == nil || event.JournalPosition().Deref().GreaterThan(lastEvent.JournalPosition().Deref()) {
			lastEvent = event
		}
	}

	_, err := tx.Exec(
		ctx,
		updateProjectionStateSQL,
		lastEvent.EventID(),
		lastEvent.InsertedAt(),
		lastEvent.AggregateVersion(),
		lastEvent.JournalPosition().Deref(),
		projection,
	)
	return err
}

func (ps *projectorSupervisor) initProjectionState(ctx context.Context, tx pgx.Tx, projection string) (*eventing.ProjectionState, error) {
	_, err := tx.Exec(ctx, "INSERT INTO projection_state (projection_name) VALUES ($1)", projection)
	if err != nil {
		return nil, err
	}
	return ps.fetchProjectionState(ctx, false, tx, projection)
}

const getProjectionStateSQL = `
SELECT
	projection_name, last_processed_event_id, last_processed_timestamp, aggregate_version, global_position, updated_at
FROM projection_state
WHERE projection_name = $1
FOR UPDATE
`

const updateProjectionStateSQL = `
UPDATE projection_state 
SET last_processed_event_id = $1, last_processed_timestamp = $2, aggregate_version = $3, global_position = $4, updated_at = NOW() 
WHERE projection_name = $5
`

const updateProjectionStateEmptySQL = `
UPDATE projection_state 
SET updated_at = NOW() 
WHERE projection_name = $1
`
