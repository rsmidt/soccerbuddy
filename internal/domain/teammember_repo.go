package domain

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type TeamMemberRepository interface {
	FindByID(ctx context.Context, id TeamMemberID) (*TeamMember, error)
	FindByTeamAndPerson(ctx context.Context, teamID TeamID, person PersonID) (*TeamMember, error)

	Save(ctx context.Context, member *TeamMember) error
}

type EventSourcedTeamMemberRepository struct {
	es eventing.EventStore
}

func NewEventSourcedTeamMemberRepository(es eventing.EventStore) *EventSourcedTeamMemberRepository {
	return &EventSourcedTeamMemberRepository{es: es}
}

func (e *EventSourcedTeamMemberRepository) FindByID(ctx context.Context, id TeamMemberID) (*TeamMember, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.TeamMemberRepository.FindByID")
	defer span.End()

	member := NewTeamMemberByID(id)
	if err := e.es.View(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (e *EventSourcedTeamMemberRepository) FindByTeamAndPerson(ctx context.Context, teamID TeamID, personID PersonID) (*TeamMember, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.TeamMemberRepository.FindByTeamAndPerson")
	defer span.End()

	ownerID, err := e.es.OwnerLookup(ctx, eventing.LookupOpts{
		AggregateType: TeamMemberAggregateType,
		FieldName:     TeamMembershipLookup,
		FieldValue:    eventing.LookupFieldValue(createTeamMembershipLookupValue(teamID, personID)),
	})
	if errors.Is(err, eventing.ErrOwnerNotFound) {
		return nil, ErrTeamMemberNotFound
	} else if err != nil {
		return nil, err
	}
	member := NewTeamMemberByID(TeamMemberID(ownerID))
	if err := e.es.View(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}

func (e *EventSourcedTeamMemberRepository) Save(ctx context.Context, member *TeamMember) error {
	ctx, span := tracing.Tracer.Start(ctx, "es.TeamMemberRepository.Save")
	defer span.End()

	return e.es.ProduceAppend(ctx, member)
}
