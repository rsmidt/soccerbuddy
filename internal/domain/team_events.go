package domain

import "github.com/rsmidt/soccerbuddy/internal/eventing"

// ========================================================
// TeamCreatedEvent
// ========================================================

const (
	TeamCreatedEventType    = eventing.EventType("team_created")
	TeamCreatedEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event                 = (*TeamCreatedEvent)(nil)
	_ eventing.UniqueConstraintAdder = (*TeamCreatedEvent)(nil)
	_ eventing.LookupProvider        = (*TeamCreatedEvent)(nil)
)

type TeamCreatedEvent struct {
	*eventing.EventBase

	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	CreatedBy    Operator `json:"created_by"`
	OwningClubID ClubID   `json:"club_id"`
}

func NewTeamCreatedEvent(id TeamID, name, slug string, owningClubID ClubID, createdBy Operator) *TeamCreatedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), TeamAggregateType, TeamCreatedEventVersion, TeamCreatedEventType)

	return &TeamCreatedEvent{
		EventBase:    base,
		Name:         name,
		Slug:         slug,
		CreatedBy:    createdBy,
		OwningClubID: owningClubID,
	}
}

func (t *TeamCreatedEvent) IsShredded() bool {
	return false
}

func (t *TeamCreatedEvent) UniqueConstraintsToAdd() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewUniqueConstraint(t.AggregateID(), TeamNameUniqueConstraint, t.Name),
		eventing.NewUniqueConstraint(t.AggregateID(), TeamSlugUniqueConstraint, t.Slug),
	}
}

func (t *TeamCreatedEvent) LookupValues() eventing.LookupMap {
	return eventing.LookupMap{
		TeamLookupSlug:       eventing.LookupFieldValue(t.Slug),
		TeamLookupName:       eventing.LookupFieldValue(t.Name),
		TeamLookupOwningClub: eventing.LookupFieldValue(t.OwningClubID),
	}
}

// ========================================================
// TeamDeletedEvent
// ========================================================

const (
	TeamDeletedEventType    = eventing.EventType("team_deleted")
	TeamDeletedEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event         = (*TeamDeletedEvent)(nil)
	_ eventing.LookupRemover = (*TeamDeletedEvent)(nil)
)

type TeamDeletedEvent struct {
	*eventing.EventBase

	OwningClubID ClubID   `json:"owning_club_id"`
	DeletedBy    Operator `json:"deleted_by"`
}

func NewTeamDeletedEvent(id TeamID, owningClubID ClubID, deletedBy Operator) *TeamDeletedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), TeamAggregateType, TeamDeletedEventVersion, TeamDeletedEventType)

	return &TeamDeletedEvent{
		EventBase:    base,
		OwningClubID: owningClubID,
		DeletedBy:    deletedBy,
	}
}

func (t *TeamDeletedEvent) IsShredded() bool {
	return false
}

func (t *TeamDeletedEvent) UniqueConstraintsToRemove() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewDeleteAllConstraint(t.AggregateID()),
	}
}

func (t *TeamDeletedEvent) LookupRemoves() []eventing.LookupFieldName {
	return []eventing.LookupFieldName{TeamLookupSlug, TeamLookupOwningClub, TeamLookupName}
}
