package eventing

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
)

type (
	ProjectionName string
)

// Projector is used to project events into a projection.
type Projector interface {
	JournalInquirer

	// Init runs any initialization logic. Can be called multiple times.
	Init(ctx context.Context) error

	// Projection specifies the name of the projection.
	Projection() ProjectionName

	// Project is used to project events into the projection.
	Project(ctx context.Context, events ...*JournalEvent) error
}

// ProjectorSupervisor is used to manage projectors.
type ProjectorSupervisor interface {
	EventListener

	// Enable is used to enable the event consumption of this supervisor.
	Enable()

	// Trigger is used to manually advance a specific projection.
	Trigger(ctx context.Context, projection ...ProjectionName)

	// Register is used to register a projector.
	Register(projector Projector)
}

// EventInterest is used to specify the events that a listener is interested in.
type EventInterest struct {
	AggType   AggregateType
	EventType EventType
}

type EventInterestSet map[EventInterest]struct{}

// IsInterestedIn returns true if interest is present in the set.
func (i EventInterestSet) IsInterestedIn(interest EventInterest) bool {
	_, ok := i[interest]
	return ok
}

// Add adds a new interest to the set.
func (i EventInterestSet) Add(interest EventInterest) {
	i[interest] = struct{}{}
}

// ProjectorToInterests returns the interests of a projector.
func ProjectorToInterests(projector Projector) []EventInterest {
	var interests []EventInterest
	query := projector.Query()
	for agg, aggQuery := range query.AggQueriesByType() {
		for _, event := range aggQuery.Events() {
			interest := EventInterest{
				AggType:   agg,
				EventType: event,
			}
			interests = append(interests, interest)
		}
	}
	return interests
}

type EventListener interface {
	// Interests returns the events that the listener is interested in.
	Interests() EventInterestSet

	// Notify is called when events were appended to the journal.
	// If false is returned, the listener will be removed.
	Notify(ctx context.Context, interest ...EventInterest) bool
}

type EventNotifier interface {
	// AddListener adds a listener to the notifier.
	AddListener(listener EventListener)

	// Start starts the notifier.
	Start(ctx context.Context) error
}

type ProjectionState struct {
	// Name is the name of the projection.
	Name ProjectionName `json:"name"`

	// LastProcessedEventID is the ID of the last event processed by the projection.
	LastProcessedEventID *EventID `json:"last_processed_event_id"`

	// LastProcessedTimestamp is the timestamp of the last event processed by the projection.
	LastProcessedTimestamp *time.Time `json:"last_processed_timestamp"`

	// AggregateVersion is the version of the aggregate at the time of the last event processed by the projection.
	AggregateVersion AggregateVersion `json:"aggregate_version"`

	// GlobalPosition is the global position of the last event processed by the projection.
	GlobalPosition decimal.Decimal `json:"global_position"`

	// UpdatedAt is the timestamp of the last update to the projection state.
	UpdatedAt time.Time `json:"updated_at"`
}
