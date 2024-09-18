package domain

import "github.com/rsmidt/soccerbuddy/internal/eventing"

// ========================================================
// ClubCreatedEvent
// ========================================================

const (
	ClubCreatedEventType    = eventing.EventType("club_created")
	ClubCreatedEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event                 = (*ClubCreatedEvent)(nil)
	_ eventing.UniqueConstraintAdder = (*ClubCreatedEvent)(nil)
	_ eventing.LookupProvider        = (*ClubCreatedEvent)(nil)
)

type ClubCreatedEvent struct {
	*eventing.EventBase

	Name string `json:"name"`
	Slug string `json:"slug"`
}

func NewClubCreatedEvent(clubID ClubID, name, slug string) *ClubCreatedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(clubID), ClubAggregateType, ClubCreatedEventVersion, ClubCreatedEventType)

	return &ClubCreatedEvent{
		EventBase: base,
		Name:      name,
		Slug:      slug,
	}
}

func (e *ClubCreatedEvent) IsShredded() bool {
	return false
}

func (e *ClubCreatedEvent) UniqueConstraintsToAdd() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewUniqueConstraint(e.AggregateID(), ClubNameUniqueConstraint, e.Name),
		eventing.NewUniqueConstraint(e.AggregateID(), ClubSlugUniqueConstraint, e.Slug),
	}
}

func (e *ClubCreatedEvent) LookupValues() eventing.LookupMap {
	return eventing.LookupMap{
		ClubLookupSlug: eventing.LookupFieldValue(e.Slug),
		ClubLookupName: eventing.LookupFieldValue(e.Name),
	}
}
