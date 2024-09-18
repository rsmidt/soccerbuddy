package domain

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type ClubRepository interface {
	FindByID(ctx context.Context, id ClubID) (*Club, error)

	Save(ctx context.Context, club *Club) error

	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByID(ctx context.Context, id ClubID) (bool, error)
}

type EventSourcedClubRepository struct {
	es eventing.EventStore
}

func NewEventSourcedClubRepository(es eventing.EventStore) *EventSourcedClubRepository {
	return &EventSourcedClubRepository{es: es}
}

func (e *EventSourcedClubRepository) FindByID(ctx context.Context, id ClubID) (*Club, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.ClubRepository.FindByID")
	defer span.End()

	club := NewClub(id)
	if err := e.es.View(ctx, club); err != nil {
		return nil, err
	}
	return club, nil
}

func (e *EventSourcedClubRepository) Save(ctx context.Context, club *Club) error {
	ctx, span := tracing.Tracer.Start(ctx, "es.ClubRepository.Save")
	defer span.End()

	return e.es.ProduceAppend(ctx, club)
}

func (e *EventSourcedClubRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.ClubRepository.ExistsByName")
	defer span.End()

	_, err := e.es.OwnerLookup(ctx, eventing.LookupOpts{
		AggregateType: ClubAggregateType,
		FieldName:     ClubLookupName,
		FieldValue:    eventing.LookupFieldValue(name),
	})
	if errors.Is(err, eventing.ErrOwnerNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (e *EventSourcedClubRepository) ExistsByID(ctx context.Context, id ClubID) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.ClubRepository.ExistsByID")
	defer span.End()

	club, err := e.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return club.State != ClubStateUnspecified, nil
}
