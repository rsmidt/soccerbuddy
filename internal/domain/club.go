package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
)

type (
	ClubID string
)

const (
	ClubAggregateType = eventing.AggregateType("club")

	ClubNameUniqueConstraint = "club_name"
	ClubSlugUniqueConstraint = "club_slug"

	ClubLookupSlug = "slug"
	ClubLookupName = "name"
)

type ClubState int

const (
	ClubStateUnspecified ClubState = iota
	ClubStateActive
)

type Club struct {
	eventing.BaseWriter

	State ClubState

	ID   ClubID
	Name string
	Slug string

	Admins AdminsSet
}

type AdminsSet map[AccountID]struct{}

func NewClub(id ClubID) *Club {
	return &Club{
		BaseWriter: *eventing.NewBaseWriter(eventing.AggregateID(id), ClubAggregateType, eventing.VersionMatcherExact),
		ID:         id,
		Admins:     make(AdminsSet),
	}
}

func (a *Club) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.WithAggregate(ClubAggregateType).
		AggregateID(eventing.AggregateID(a.ID)).
		Finish().MustBuild()
}

func (a *Club) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *ClubCreatedEvent:
			a.State = ClubStateActive
			a.Name = e.Name
			a.Slug = e.Slug
		case *ClubAdminAddedEvent:
			a.Admins[e.AddedUserID] = struct{}{}
		}
		a.BaseWriter.Reduce(events)
	}
}

func (a *Club) Init(name, slug string, createdAt time.Time) error {
	if a.State != ClubStateUnspecified {
		return NewInvalidAggregateStateError(a.Aggregate(), int(ClubStateUnspecified), int(a.State))
	}
	clubID := ClubID(a.Aggregate().AggregateID)
	a.Append(NewClubCreatedEvent(clubID, name, slug, createdAt))
	return nil
}

func (a *Club) AddAdmin(id AccountID, addedAt time.Time, addedBy Operator) error {
	if a.State != ClubStateActive {
		return NewInvalidAggregateStateError(a.Aggregate(), int(ClubStateActive), int(a.State))
	}
	// Prevent adding the same admin twice.
	if _, ok := a.Admins[id]; ok {
		return nil
	}
	a.Append(NewClubAdminAddedEvent(a.ID, id, addedAt, addedBy))
	return nil
}
