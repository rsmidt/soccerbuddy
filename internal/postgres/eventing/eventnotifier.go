package eventing

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"log/slog"
	"strings"
	"sync"
)

type eventNotifier struct {
	mu        sync.RWMutex
	subs      []eventing.EventListener
	allAggSet map[eventing.AggregateType]struct{}
	pool      *pgxpool.Pool
	log       *slog.Logger
}

func NewEventNotifier(log *slog.Logger, pool *pgxpool.Pool) eventing.EventNotifier {
	return &eventNotifier{pool: pool, log: log, allAggSet: make(map[eventing.AggregateType]struct{})}
}

func (e *eventNotifier) AddListener(listener eventing.EventListener) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.subs = append(e.subs, listener)
	for interest := range listener.Interests() {
		e.allAggSet[interest.AggType] = struct{}{}
	}
}

func (e *eventNotifier) Start(ctx context.Context) error {
	conn, err := e.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	err = e.prepareListeners(ctx, conn)
	if err != nil {
		return err
	}

	// Run the notification loop.
	var retries int
	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			// For whatever reason, pg wraps the context cancelled error.
			if errors.Is(err, context.Canceled) {
				return nil
			}
			e.log.Error("Failed to wait for notification", slog.String("err", err.Error()))
			if retries > 3 {
				e.log.Error("shutting down because of too many retries")
				return err
			}
			retries++
			continue
		}
		e.log.Debug("Received notification", slog.String("channel", notification.Channel), slog.String("payload", notification.Payload))

		e.handleNotification(ctx, notification)
	}
}

func (e *eventNotifier) prepareListeners(ctx context.Context, conn *pgxpool.Conn) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Create a new listener for each interested channel.
	for agg := range e.allAggSet {
		e.log.Debug("Setting up channel", slog.String("channel", string(agg)))
		stmt := fmt.Sprintf("LISTEN event_store_%s", string(agg))
		_, err := conn.Exec(ctx, stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: Batch notifications.
func (e *eventNotifier) handleNotification(ctx context.Context, notification *pgconn.Notification) {
	ctx, span := tracing.Tracer.Start(ctx, "pg.EventNotifier.handleNotification")
	defer span.End()

	aggregateType := eventing.AggregateType(strings.TrimPrefix(notification.Channel, "event_store_"))
	eventType := eventing.EventType(notification.Payload)
	interestFromNotification := eventing.EventInterest{
		AggType:   aggregateType,
		EventType: eventType,
	}

	for _, sub := range e.subs {
		if !sub.Interests().IsInterestedIn(interestFromNotification) {
			continue
		}
		if !sub.Notify(ctx, interestFromNotification) {
			// TODO: Remove the listener.
		}
	}
}
