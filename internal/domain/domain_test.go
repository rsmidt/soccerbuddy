package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/shopspring/decimal"
	"time"
)

func createJournalEvent(event eventing.Event, journalPosition eventing.JournalPosition, aggregateVersion eventing.AggregateVersion, insertedAt time.Time) *eventing.JournalEvent {
	return eventing.NewJournalEvent(
		event,
		eventing.EventID(journalPosition.Deref().String()),
		aggregateVersion,
		journalPosition,
		insertedAt,
	)
}

func createInitialEvents(events ...eventing.Event) []*eventing.JournalEvent {
	journalPositionCounter := decimal.NewFromInt(0)
	aggregateVersionCounter := eventing.AggregateVersion(0)
	insertedAtCounter := time.Now()

	journalEvents := make([]*eventing.JournalEvent, len(events))
	for i, event := range events {
		journalEvents[i] = createJournalEvent(event, eventing.JournalPosition(journalPositionCounter), aggregateVersionCounter, insertedAtCounter)
		aggregateVersionCounter++
		journalPositionCounter.Add(decimal.NewFromInt(1))
		insertedAtCounter = insertedAtCounter.Add(time.Second)
	}
	return journalEvents
}
