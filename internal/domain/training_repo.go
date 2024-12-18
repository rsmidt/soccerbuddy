package domain

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type TrainingRepository interface {
	FindByID(ctx context.Context, id TrainingID) (*Training, error)

	Save(ctx context.Context, team *Training) error
}

type EventSourcedTrainingRepository struct {
	es eventing.EventStore
}

func NewEventSourcedTrainingRepository(es eventing.EventStore) *EventSourcedTrainingRepository {
	return &EventSourcedTrainingRepository{es: es}
}

func (e *EventSourcedTrainingRepository) FindByID(ctx context.Context, id TrainingID) (*Training, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.TrainingRepository.FindByID")
	defer span.End()

	team := NewTrainingByID(id)
	if err := e.es.View(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (e *EventSourcedTrainingRepository) Save(ctx context.Context, team *Training) error {
	ctx, span := tracing.Tracer.Start(ctx, "es.TrainingRepository.Save")
	defer span.End()

	return e.es.ProduceAppend(ctx, team)
}
