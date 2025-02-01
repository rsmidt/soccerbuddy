package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
)

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

	CreatedAt time.Time `json:"created_at"`
}

func NewClubCreatedEvent(clubID ClubID, name, slug string, createdAt time.Time) *ClubCreatedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(clubID), ClubAggregateType, ClubCreatedEventVersion, ClubCreatedEventType)

	return &ClubCreatedEvent{
		EventBase: base,
		Name:      name,
		Slug:      slug,
		CreatedAt: createdAt,
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

// ========================================================
// ClubAdminAddedEvent
// ========================================================

const (
	ClubAdminAddedEventType    = eventing.EventType("club_admin_added")
	ClubAdminAddedEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event = (*ClubAdminAddedEvent)(nil)
)

type ClubAdminAddedEvent struct {
	*eventing.EventBase

	AddedUserID AccountID `json:"added_user_id"`
	AddedAt     time.Time `json:"created_at"`
	AddedBy     Operator  `json:"added_by"`
}

func NewClubAdminAddedEvent(clubID ClubID, addedUserID AccountID, addedAt time.Time, addedBy Operator) *ClubAdminAddedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(clubID), ClubAggregateType, ClubAdminAddedEventVersion, ClubAdminAddedEventType)

	return &ClubAdminAddedEvent{
		EventBase:   base,
		AddedUserID: addedUserID,
		AddedAt:     addedAt,
		AddedBy:     addedBy,
	}
}

func (e *ClubAdminAddedEvent) IsShredded() bool {
	return false
}
