package eventing

// EventDescriptor describes the current state of an aggregate.
type EventDescriptor interface {
	AggregateID() AggregateID
	AggregateType() AggregateType
	EventVersion() EventVersion
	EventType() EventType
}

type Event interface {
	EventDescriptor

	IsShredded() bool
}

type EventBase struct {
	aggregateID   AggregateID
	aggregateType AggregateType
	eventVersion  EventVersion
	eventType     EventType
}

func NewEventBase(aggregateID AggregateID, aggregateType AggregateType, eventVersion EventVersion, eventType EventType) *EventBase {
	return &EventBase{
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		eventVersion:  eventVersion,
		eventType:     eventType,
	}
}

func (e *EventBase) AggregateID() AggregateID {
	return e.aggregateID
}

func (e *EventBase) AggregateType() AggregateType {
	return e.aggregateType
}

func (e *EventBase) EventVersion() EventVersion {
	return e.eventVersion
}

func (e *EventBase) EventType() EventType {
	return e.eventType
}
