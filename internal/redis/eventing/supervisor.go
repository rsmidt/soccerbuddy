package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type redisSupervisor struct {
	mu                  sync.RWMutex
	log                 *slog.Logger
	projectors          map[eventing.ProjectionName]eventing.Projector
	projectorByInterest map[eventing.EventInterest][]eventing.Projector
	started             atomic.Bool
	es                  eventing.EventStore

	rd     rueidis.Client
	locker rueidislock.Locker
}

func NewProjectorSupervisor(log *slog.Logger, es eventing.EventStore, rd rueidis.Client, locker rueidislock.Locker) eventing.ProjectorSupervisor {
	return &redisSupervisor{
		log:                 log,
		projectors:          make(map[eventing.ProjectionName]eventing.Projector),
		projectorByInterest: make(map[eventing.EventInterest][]eventing.Projector),
		es:                  es,
		rd:                  rd,
		locker:              locker,
	}
}

func (r *redisSupervisor) Register(projector eventing.Projector) {
	r.mu.Lock()
	defer r.mu.Unlock()

	interests := eventing.ProjectorToInterests(projector)
	for _, interest := range interests {
		r.projectorByInterest[interest] = append(r.projectorByInterest[interest], projector)
	}
	r.projectors[projector.Projection()] = projector
}

func (r *redisSupervisor) Interests() eventing.EventInterestSet {
	r.mu.RLock()
	defer r.mu.RUnlock()

	interestSet := make(eventing.EventInterestSet)
	for interest := range r.projectorByInterest {
		interestSet.Add(interest)
	}
	return interestSet
}

func (r *redisSupervisor) Notify(ctx context.Context, interests ...eventing.EventInterest) bool {
	ctx, span := tracing.Tracer.Start(ctx, "rd.ProjectorSupervisor.Notify")
	defer span.End()

	if !r.started.Load() {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	var triggeredProjectors []eventing.Projector
	for _, interest := range interests {
		if projectors, ok := r.projectorByInterest[interest]; ok {
			triggeredProjectors = append(triggeredProjectors, projectors...)
		}
	}
	if err := r.trigger(ctx, false, triggeredProjectors...); err != nil {
		return false
	}
	return true
}

func (r *redisSupervisor) Enable() {
	r.started.Store(true)
}

func (r *redisSupervisor) Trigger(ctx context.Context, projection ...eventing.ProjectionName) {
	ctx, span := tracing.Tracer.Start(ctx, "rd.ProjectorSupervisor.Trigger")
	defer span.End()

	var projectors []eventing.Projector
	for _, proj := range projection {
		projector, ok := r.projectors[proj]
		if !ok {
			r.log.Error("Failed to manually trigger projectors", slog.String("projection", string(proj)), slog.String("err", "projector not found"))
			continue
		}
		projectors = append(projectors, projector)
	}

	if err := r.trigger(ctx, true, projectors...); err != nil {
		r.log.Error("Failed to manually trigger projectors", slog.String("err", err.Error()))
		tracing.RecordError(ctx, err)
	}
}

func (r *redisSupervisor) trigger(ctx context.Context, wait bool, projectors ...eventing.Projector) error {
	for _, projector := range projectors {
		if err := r.triggerProjector(ctx, wait, projector); err != nil {
			r.log.ErrorContext(ctx, "Failed to trigger projector", slog.String("projector", string(projector.Projection())), slog.String("projector_type", "redis"), slog.String("err", err.Error()))
		}
	}
	return nil
}

func (r *redisSupervisor) triggerProjector(ctx context.Context, wait bool, projector eventing.Projector) error {
	ctx, span := tracing.Tracer.Start(ctx, "rd.ProjectorSupervisor.triggerProjector")
	defer span.End()

	name := string(projector.Projection())
	r.log.DebugContext(ctx, "Advancing projector", slog.String("projection", name), slog.String("projector_type", "redis"))

	// Acquire a lock either by waiting or not.
	var (
		cancel context.CancelFunc
		err    error
	)
	if wait {
		ctx, cancel, err = r.locker.WithContext(ctx, name)
	} else {
		ctx, cancel, err = r.locker.TryWithContext(ctx, name)
	}
	if errors.Is(err, rueidislock.ErrNotLocked) {
		return nil
	} else if err != nil {
		return err
	}
	defer cancel()

	// Get the current state.
	var state eventing.ProjectionState
	stateKey := fmt.Sprintf("projection:state:%s:v1", projector.Projection())
	cmd := r.rd.B().JsonGet().Key(stateKey).Path(".").Build().Pin()
	if err := r.rd.Do(ctx, cmd).DecodeJSON(&state); rueidis.IsRedisNil(err) {
		state = eventing.ProjectionState{Name: projector.Projection()}
	} else if err != nil {
		return err
	}

	// Fetch the events from the event store.
	queryWithPos := eventing.NewJournalQueryBuilderFrom(projector.Query()).
		WithJournalPositionAfter(eventing.JournalPosition(state.GlobalPosition)).MustBuild()
	events, err := r.es.Query(ctx, queryWithPos, eventing.WithLimitToOldestRunningTransaction())
	if err != nil {
		return fmt.Errorf("failed to query events for projection: %v", err)
	}
	if err := projector.Project(ctx, events...); err != nil {
		return fmt.Errorf("failed to project: %v", err)
	}

	updatedState := updateState(state, events)
	val, err := json.Marshal(&updatedState)
	if err != nil {
		return err
	}
	cmd = r.rd.B().JsonSet().Key(stateKey).Path(".").Value(string(val)).Build()
	if err := r.rd.Do(ctx, cmd).Error(); err != nil {
		return err
	}
	return nil
}

func updateState(state eventing.ProjectionState, events []*eventing.JournalEvent) eventing.ProjectionState {
	if len(events) == 0 {
		return eventing.ProjectionState{
			Name:                   state.Name,
			LastProcessedEventID:   state.LastProcessedEventID,
			LastProcessedTimestamp: state.LastProcessedTimestamp,
			AggregateVersion:       state.AggregateVersion,
			GlobalPosition:         state.GlobalPosition,
			UpdatedAt:              time.Now(),
		}
	}

	// Find the last event.
	var lastEvent *eventing.JournalEvent
	for _, event := range events {
		if lastEvent == nil || event.JournalPosition().Deref().GreaterThan(lastEvent.JournalPosition().Deref()) {
			lastEvent = event
		}
	}
	eventID := lastEvent.EventID()
	processedTimestamp := lastEvent.InsertedAt()
	return eventing.ProjectionState{
		Name:                   state.Name,
		LastProcessedEventID:   &eventID,
		LastProcessedTimestamp: &processedTimestamp,
		AggregateVersion:       lastEvent.AggregateVersion(),
		GlobalPosition:         lastEvent.JournalPosition().Deref(),
		UpdatedAt:              time.Now(),
	}
}
