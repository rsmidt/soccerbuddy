package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
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
}

func NewClub(id ClubID) *Club {
	return &Club{
		BaseWriter: *eventing.NewBaseWriter(eventing.AggregateID(id), ClubAggregateType, eventing.VersionMatcherExact),
		ID:         id,
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
		}
		a.BaseWriter.Reduce(events)
	}
}

func (a *Club) Init(name, slug string) error {
	if a.State != ClubStateUnspecified {
		return NewInvalidAggregateStateError(a.Aggregate(), int(ClubStateUnspecified), int(a.State))
	}
	clubID := ClubID(a.Aggregate().AggregateID)
	a.Append(NewClubCreatedEvent(clubID, name, slug))
	return nil
}
