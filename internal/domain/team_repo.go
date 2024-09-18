package domain

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type TeamRepository interface {
	FindByID(ctx context.Context, id TeamID) (*Team, error)

	Save(ctx context.Context, team *Team) error

	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByID(ctx context.Context, id TeamID) (bool, error)
}

type EventSourcedTeamRepository struct {
	es eventing.EventStore
}

func NewEventSourcedTeamRepository(es eventing.EventStore) *EventSourcedTeamRepository {
	return &EventSourcedTeamRepository{es: es}
}

func (e *EventSourcedTeamRepository) FindByID(ctx context.Context, id TeamID) (*Team, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.TeamRepository.FindByID")
	defer span.End()

	team := NewTeam(id)
	if err := e.es.View(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (e *EventSourcedTeamRepository) Save(ctx context.Context, team *Team) error {
	ctx, span := tracing.Tracer.Start(ctx, "es.TeamRepository.Save")
	defer span.End()

	return e.es.ProduceAppend(ctx, team)
}

func (e *EventSourcedTeamRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.TeamRepository.ExistsByName")
	defer span.End()

	_, err := e.es.OwnerLookup(ctx, eventing.LookupOpts{
		AggregateType: TeamAggregateType,
		FieldName:     TeamLookupName,
		FieldValue:    eventing.LookupFieldValue(name),
	})
	if errors.Is(err, eventing.ErrOwnerNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (e *EventSourcedTeamRepository) ExistsByID(ctx context.Context, id TeamID) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.TeamRepository.ExistsByID")
	defer span.End()

	team, err := e.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return team.State != TeamStateUnspecified, nil
}
