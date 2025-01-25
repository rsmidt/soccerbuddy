package domain

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type TeamRepository interface {
	FindByID(ctx context.Context, id TeamID) (*Team, error)

	Save(ctx context.Context, team *Team) error

	ExistsByNameInClub(ctx context.Context, name string, id ClubID) (bool, error)
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

func (e *EventSourcedTeamRepository) ExistsByNameInClub(ctx context.Context, name string, clubID ClubID) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.TeamRepository.ExistsByNameInClub")
	defer span.End()

	team, err := e.FindByID(ctx, TeamID(name))
	if err != nil {
		return false, err
	}
	if team.State == TeamStateUnspecified {
		return false, nil
	}
	if team.OwningClubID != clubID {
		return false, nil
	}
	return team.Name == name, nil
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
