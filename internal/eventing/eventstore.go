//go:generate go run github.com/rsmidt/soccerbuddy/cmd/eventing-gen
package eventing

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
)

var (
	ErrVersionMismatch = errors.New("version mismatch")
	ErrIntentOutdated  = errors.New("intent is outdated")
)

type (
	AggregateID      string
	AggregateType    string
	AggregateVersion uint
	EventVersion     string
	EventType        string
	EventID          string
	JournalPosition  decimal.Decimal
)

func (a AggregateID) Deref() string {
	return string(a)
}

func (j JournalPosition) Deref() decimal.Decimal {
	return decimal.Decimal(j)
}

type EventStore interface {
	// Append is used to append events to the event store.
	Append(ctx context.Context, intents ...AggregateChangeIntent) ([]*JournalEvent, error)

	// ProduceAppend is used to produce events to the event store and apply them to the producer.
	ProduceAppend(ctx context.Context, producer Writer) error

	// Query is used to query events from the event store.
	Query(ctx context.Context, query JournalQuery, opts ...QueryOpts) ([]*JournalEvent, error)

	// View is used to query and reduce events into a view.
	View(ctx context.Context, view JournalViewer) error

	// Lookup allows for finding transactional consistent values of aggregates.
	// Look at [LookupProvider] for more information.
	Lookup(ctx context.Context, opts LookupOpts) (*LookupFieldValue, error)

	// OwnerLookup is a special case of Lookup that is used to find the owner of an aggregate.
	// Returns [ErrOwnerNotFound] if the entry does not exist.
	OwnerLookup(ctx context.Context, opts LookupOpts) (AggregateID, error)

	// AddHook adds a hook to the event store.
	AddHook(hook Hook)
}

// Hook is a marker interface for all possible hooks that can run during the event store lifecycle.
type Hook interface {
}

// PostPersist hooks run after events were persisted in the store.
type PostPersist interface {
	Hook

	PostPersist(ctx context.Context) error
}

type PostPersistFunc func(ctx context.Context) error

func (p PostPersistFunc) PostPersist(ctx context.Context) error {
	return p(ctx)
}

func NewPostPersistHook(fn PostPersistFunc) PostPersist {
	return fn
}
