package domain

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type SessionRepository interface {
	FindByID(ctx context.Context, id SessionID) (*Session, error)
	FindByToken(ctx context.Context, token SessionToken) (*Session, error)

	Save(ctx context.Context, session *Session) error
}

type EventSourcedSessionRepository struct {
	es eventing.EventStore
}

func NewEventSourcedSessionRepository(es eventing.EventStore) *EventSourcedSessionRepository {
	return &EventSourcedSessionRepository{es: es}
}

func (e *EventSourcedSessionRepository) FindByID(ctx context.Context, id SessionID) (*Session, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.SessionRepository.FindByID")
	defer span.End()

	session := NewSession(id)
	if err := e.es.View(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (e *EventSourcedSessionRepository) FindByToken(ctx context.Context, token SessionToken) (*Session, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.SessionRepository.FindByEmail")
	defer span.End()

	ownerID, err := e.es.OwnerLookup(ctx, eventing.LookupOpts{
		AggregateType: SessionAggregateType,
		FieldName:     SessionLookupToken,
		FieldValue:    eventing.LookupFieldValue(token),
	})
	if errors.Is(err, eventing.ErrOwnerNotFound) {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}
	session := NewSession(SessionID(ownerID.Deref()))
	if err := e.es.View(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (e *EventSourcedSessionRepository) Save(ctx context.Context, session *Session) error {
	ctx, span := tracing.Tracer.Start(ctx, "es.SessionRepository.Save")
	defer span.End()

	return e.es.ProduceAppend(ctx, session)
}
