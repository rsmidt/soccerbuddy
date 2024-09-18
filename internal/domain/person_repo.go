package domain

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type PersonRepository interface {
	FindByID(ctx context.Context, id PersonID) (*Person, error)

	Save(ctx context.Context, person *Person) error

	ExistsByID(ctx context.Context, id PersonID) (bool, error)
}

type EventSourcedPersonRepository struct {
	es eventing.EventStore
}

func NewEventSourcedPersonRepository(es eventing.EventStore) *EventSourcedPersonRepository {
	return &EventSourcedPersonRepository{es: es}
}

func (e *EventSourcedPersonRepository) FindByID(ctx context.Context, id PersonID) (*Person, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.PersonRepository.FindByID")
	defer span.End()

	person := NewPerson(id)
	if err := e.es.View(ctx, person); err != nil {
		return nil, err
	}
	return person, nil
}

func (e *EventSourcedPersonRepository) Save(ctx context.Context, person *Person) error {
	ctx, span := tracing.Tracer.Start(ctx, "es.PersonRepository.Save")
	defer span.End()

	return e.es.ProduceAppend(ctx, person)
}

func (e *EventSourcedPersonRepository) ExistsByID(ctx context.Context, id PersonID) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.PersonRepository.ExistsByID")
	defer span.End()

	person, err := e.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return person.State != PersonStateUnspecified, nil
}
