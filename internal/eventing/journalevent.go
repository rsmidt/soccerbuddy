package eventing

import "time"

// JournalEventMapper is responsible for mapping journal events to domain events.
type JournalEventMapper interface {
	MapFrom(
		aggregateID AggregateID,
		aggregateType AggregateType,
		eventVersion EventVersion,
		eventType EventType,
		eventID EventID,
		aggregateVersion AggregateVersion,
		journalPosition JournalPosition,
		insertedAt time.Time,
		payload []byte,
	) (*JournalEvent, error)
}

type JournalEvent struct {
	Event

	eventID          EventID
	aggregateVersion AggregateVersion
	journalPosition  JournalPosition
	insertedAt       time.Time
}

func NewJournalEvent(event Event, eventID EventID, aggregateVersion AggregateVersion, journalPosition JournalPosition, insertedAt time.Time) *JournalEvent {
	return &JournalEvent{
		Event:            event,
		eventID:          eventID,
		aggregateVersion: aggregateVersion,
		journalPosition:  journalPosition,
		insertedAt:       insertedAt,
	}
}

func (e *JournalEvent) EventID() EventID {
	return e.eventID
}

func (e *JournalEvent) AggregateVersion() AggregateVersion {
	return e.aggregateVersion
}

func (e *JournalEvent) JournalPosition() JournalPosition {
	return e.journalPosition
}

func (e *JournalEvent) InsertedAt() time.Time {
	return e.insertedAt
}
