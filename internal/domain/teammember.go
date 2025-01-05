package domain

import (
	"errors"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
)

type (
	TeamMemberID string

	// TeamMemberRole is the function of the team member (e.g. coach, player, etc.).
	TeamMemberRole string
)

func (t TeamMemberRole) Deref() string {
	return string(t)
}

const (
	TeamMemberAggregateType = eventing.AggregateType("team_member")

	// TeamMembershipUniqueConstraint guarantees that a person with a given ID cannot be added twice to a team.
	TeamMembershipUniqueConstraint = "team_membership"

	// TeamMembershipLookup allows looking up the owner ID of a membership by team ID + person ID.
	TeamMembershipLookup = "team_membership"

	// A default set of roles for team members.
	TeamMemberRoleCoach  TeamMemberRole = "COACH"
	TeamMemberRolePlayer TeamMemberRole = "PLAYER"
	TeamMemberRoleGuest  TeamMemberRole = "GUEST"
)

var (
	ErrTeamMemberNotFound = errors.New("team member not found")
)

type TeamMemberState int

const (
	TeamMemberStateUnspecified TeamMemberState = iota
	TeamMemberStateActive
)

type TeamMember struct {
	eventing.BaseWriter

	State     TeamMemberState
	ID        TeamMemberID
	PersonID  PersonID
	TeamID    TeamID
	InvitedBy Operator
	JoinedAt  time.Time
}

func NewTeamMemberByID(id TeamMemberID) *TeamMember {
	return &TeamMember{
		BaseWriter: *eventing.NewBaseWriter(eventing.AggregateID(id), TeamMemberAggregateType, eventing.VersionMatcherExact),
		ID:         id,
	}
}

func NewTeamMember(id TeamMemberID, teamID TeamID, personID PersonID) *TeamMember {
	return &TeamMember{
		BaseWriter: *eventing.NewBaseWriter(eventing.AggregateID(id), TeamMemberAggregateType, eventing.VersionMatcherExact),
		ID:         id,
		TeamID:     teamID,
		PersonID:   personID,
	}
}

func (m *TeamMember) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.WithAggregate(TeamMemberAggregateType).
		AggregateID(m.Aggregate().AggregateID).
		Finish().MustBuild()
}

func (m *TeamMember) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *PersonInvitedToTeamEvent:
			m.State = TeamMemberStateActive
			m.InvitedBy = e.InvitedBy
			m.PersonID = e.PersonID
			m.TeamID = e.TeamID
			m.JoinedAt = event.InsertedAt()
		}
	}
	m.BaseWriter.Reduce(events)
}

func (m *TeamMember) Invite(operator Operator, assignedRole TeamMemberRole) error {
	if m.State != TeamMemberStateUnspecified {
		return NewInvalidAggregateStateError(m.Aggregate(), int(TeamMemberStateUnspecified), int(m.State))
	}
	m.Append(NewPersonInvitedToTeamEvent(m.ID, m.PersonID, m.TeamID, operator, assignedRole))
	return nil
}

func createTeamMembershipLookupValue(teamID TeamID, personID PersonID) string {
	return fmt.Sprintf("%s:%s", teamID, personID)
}
