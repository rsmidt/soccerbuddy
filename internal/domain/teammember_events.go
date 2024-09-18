package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

// ========================================================
// PersonCreatedEvent
// ========================================================

const (
	PersonInvitedToTeamEventType    = eventing.EventType("person_invited_to_team")
	PersonInvitedToTeamEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event                 = (*PersonInvitedToTeamEvent)(nil)
	_ eventing.UniqueConstraintAdder = (*PersonInvitedToTeamEvent)(nil)
	_ eventing.LookupProvider        = (*PersonInvitedToTeamEvent)(nil)
)

type PersonInvitedToTeamEvent struct {
	*eventing.EventBase

	TeamID       TeamID             `json:"team_id"`
	PersonID     PersonID           `json:"person_id"`
	AssignedRole TeamMemberRoleRole `json:"assigned_role"`
	InvitedBy    Operator           `json:"invited_by"`
}

func NewPersonInvitedToTeamEvent(id TeamMemberID, personID PersonID, teamID TeamID, invitedBy Operator, assignedRole TeamMemberRoleRole) *PersonInvitedToTeamEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), TeamMemberAggregateType, PersonInvitedToTeamEventVersion, PersonInvitedToTeamEventType)

	return &PersonInvitedToTeamEvent{
		EventBase:    base,
		PersonID:     personID,
		TeamID:       teamID,
		InvitedBy:    invitedBy,
		AssignedRole: assignedRole,
	}
}

func (p *PersonInvitedToTeamEvent) IsShredded() bool {
	return false
}

func (p *PersonInvitedToTeamEvent) UniqueConstraintsToAdd() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewUniqueConstraint(p.AggregateID(), TeamMembershipUniqueConstraint, createTeamMembershipLookupValue(p.TeamID, p.PersonID)),
	}
}

func (p *PersonInvitedToTeamEvent) LookupValues() eventing.LookupMap {
	return eventing.LookupMap{
		TeamMembershipUniqueConstraint: eventing.LookupFieldValue(createTeamMembershipLookupValue(p.TeamID, p.PersonID)),
	}
}
