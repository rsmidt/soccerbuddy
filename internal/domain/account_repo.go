package domain

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type AccountRepository interface {
	FindByID(ctx context.Context, id AccountID) (*Account, error)
	FindByEmail(ctx context.Context, email string) (*Account, error)

	Save(ctx context.Context, account *Account) error

	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type EventSourcedAccountRepository struct {
	es eventing.EventStore
}

func NewEventSourcedAccountRepository(es eventing.EventStore) *EventSourcedAccountRepository {
	return &EventSourcedAccountRepository{es: es}
}

func (e *EventSourcedAccountRepository) FindByID(ctx context.Context, id AccountID) (*Account, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.AccountRepository.FindByID")
	defer span.End()

	account := NewAccount(id)
	if err := e.es.View(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (e *EventSourcedAccountRepository) FindByEmail(ctx context.Context, email string) (*Account, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.AccountRepository.FindByEmail")
	defer span.End()

	ownerID, err := e.es.OwnerLookup(ctx, eventing.LookupOpts{
		AggregateType: AccountAggregateType,
		FieldName:     AccountLookupEmail,
		FieldValue:    eventing.LookupFieldValue(email),
	})
	if errors.Is(err, eventing.ErrOwnerNotFound) {
		return nil, ErrAccountNotFound
	} else if err != nil {
		return nil, err
	}
	account := NewAccount(AccountID(ownerID.Deref()))
	if err := e.es.View(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (e *EventSourcedAccountRepository) Save(ctx context.Context, account *Account) error {
	ctx, span := tracing.Tracer.Start(ctx, "es.AccountRepository.Save")
	defer span.End()

	return e.es.ProduceAppend(ctx, account)
}

func (e *EventSourcedAccountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, span := tracing.Tracer.Start(ctx, "es.AccountRepository.ExistsByEmail")
	defer span.End()

	_, err := e.es.OwnerLookup(ctx, eventing.LookupOpts{
		AggregateType: AccountAggregateType,
		FieldName:     AccountLookupEmail,
		FieldValue:    eventing.LookupFieldValue(email),
	})
	if errors.Is(err, eventing.ErrOwnerNotFound) {
		return false, ErrAccountNotFound
	} else if err != nil {
		return false, err
	}
	return true, nil
}
