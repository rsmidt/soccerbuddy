package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTeamMember_Invite(t *testing.T) {
	teamMemberID := idgen.New[TeamMemberID]()
	teamID := idgen.New[TeamID]()
	personID := idgen.New[PersonID]()
	operator := NewOperator(idgen.New[AccountID](), nil)
	role := TeamMemberRoleRole("player")

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		operator      Operator
		role          TeamMemberRoleRole
		expectedError error
	}{
		{
			name:          "Succeeds if team member is invited correctly",
			initialEvents: createInitialEvents(),
			emittedEvents: []eventing.Event{
				NewPersonInvitedToTeamEvent(teamMemberID, personID, teamID, operator, role),
			},
			operator:      operator,
			role:          role,
			expectedError: nil,
		},
		{
			name: "Fails if team member is already invited",
			initialEvents: createInitialEvents(
				NewPersonInvitedToTeamEvent(teamMemberID, personID, teamID, operator, role),
			),
			operator:      operator,
			role:          role,
			expectedError: NewInvalidAggregateStateError(NewTeamMember(teamMemberID, teamID, personID).Aggregate(), int(TeamMemberStateUnspecified), int(TeamMemberStateActive)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			teamMember := NewTeamMember(teamMemberID, teamID, personID)
			teamMember.Reduce(tt.initialEvents)
			err := teamMember.Invite(tt.operator, tt.role)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, teamMember.Changes().Events())
		})
	}
}

func TestTeamMember_Reduce(t *testing.T) {
	teamMemberID := idgen.New[TeamMemberID]()
	teamID := idgen.New[TeamID]()
	personID := idgen.New[PersonID]()
	operator := NewOperator(idgen.New[AccountID](), nil)
	role := TeamMemberRoleRole("player")

	tests := []struct {
		name           string
		initialEvents  []*eventing.JournalEvent
		expectedState  TeamMemberState
		expectedRole   TeamMemberRoleRole
		expectedID     TeamMemberID
		expectedTeam   TeamID
		expectedPerson PersonID
	}{
		{
			name: "Succeeds if team member state is updated correctly",
			initialEvents: createInitialEvents(
				NewPersonInvitedToTeamEvent(teamMemberID, personID, teamID, operator, role),
			),
			expectedState:  TeamMemberStateActive,
			expectedRole:   role,
			expectedID:     teamMemberID,
			expectedTeam:   teamID,
			expectedPerson: personID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			teamMember := NewTeamMember(teamMemberID, teamID, personID)
			teamMember.Reduce(tt.initialEvents)
			assert.Equal(t, tt.expectedState, teamMember.State)
			assert.Equal(t, tt.expectedRole, role)
			assert.Equal(t, tt.expectedID, teamMember.ID)
			assert.Equal(t, tt.expectedTeam, teamMember.TeamID)
			assert.Equal(t, tt.expectedPerson, teamMember.PersonID)
		})
	}
}
