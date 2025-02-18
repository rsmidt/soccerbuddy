package domain

import (
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
)

type (
	TeamID string
)

const (
	TeamAggregateType = eventing.AggregateType("team")

	// TeamNameInClubUniqueConstraint guarantees that only one team with a given name can exist in a given club.
	TeamNameInClubUniqueConstraint = "team_name_club"
	TeamSlugUniqueConstraint       = "team_slug"

	TeamLookupSlug       = eventing.LookupFieldName("slug")
	TeamLookupName       = eventing.LookupFieldName("name")
	TeamLookupOwningClub = eventing.LookupFieldName("owning_club")
)

var (
	ErrTeamOwningClubNotFound = errors.New("owning club not found")
	ErrTeamNotFound           = errors.New("team not found")
)

type TeamState int

const (
	TeamStateUnspecified TeamState = iota
	TeamStateActive
	TeamStateDeleted
)

type Team struct {
	eventing.BaseWriter

	ID           TeamID
	Name         string
	Slug         string
	OwningClubID ClubID
	CreatedAt    time.Time
	CreatedBy    Operator
	UpdatedAt    time.Time
	State        TeamState
}

func NewTeam(id TeamID) *Team {
	return &Team{
		BaseWriter: *eventing.NewBaseWriter(eventing.AggregateID(id), TeamAggregateType, eventing.VersionMatcherExact),
		ID:         id,
	}
}

func (t *Team) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.WithAggregate(TeamAggregateType).
		AggregateID(eventing.AggregateID(t.ID)).
		Finish().MustBuild()
}

func (t *Team) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *TeamCreatedEvent:
			t.State = TeamStateActive
			t.Name = e.Name
			t.Slug = e.Slug
			t.OwningClubID = e.OwningClubID
			t.CreatedAt = event.InsertedAt()
			t.UpdatedAt = event.InsertedAt()
		case *TeamDeletedEvent:
			t.State = TeamStateDeleted
		}
	}
	t.BaseWriter.Reduce(events)
}

func (t *Team) Init(name, slug string, owningClubID ClubID, createdBy Operator, createdAt time.Time) error {
	if t.State != TeamStateUnspecified {
		return NewInvalidAggregateStateError(t.Aggregate(), int(TeamStateUnspecified), int(t.State))
	}
	event := NewTeamCreatedEvent(t.ID, name, slug, owningClubID, createdBy, createdAt)
	t.Append(event)
	return nil
}

func (t *Team) Delete(deletedBy Operator) error {
	if t.State != TeamStateActive {
		return NewInvalidAggregateStateError(t.Aggregate(), int(TeamStateActive), int(t.State))
	}
	t.Append(NewTeamDeletedEvent(t.ID, t.OwningClubID, deletedBy))
	return nil
}
