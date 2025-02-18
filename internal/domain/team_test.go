package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTeam_Init(t *testing.T) {
	teamID := idgen.New[TeamID]()
	clubID := idgen.New[ClubID]()
	createdBy := NewOperator(idgen.New[AccountID](), nil)
	createdAt := time.Now()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		teamName      string
		teamSlug      string
		owningClubID  ClubID
		createdBy     Operator
		createdAt     time.Time
		expectedError error
	}{
		{
			name:          "Succeeds if team is initialized correctly",
			initialEvents: createInitialEvents(),
			emittedEvents: []eventing.Event{
				NewTeamCreatedEvent(teamID, "Team Awesome", "team-awesome", clubID, createdBy, createdAt),
			},
			teamName:      "Team Awesome",
			teamSlug:      "team-awesome",
			owningClubID:  clubID,
			createdBy:     createdBy,
			createdAt:     createdAt,
			expectedError: nil,
		},
		{
			name: "Fails if team is already initialized",
			initialEvents: createInitialEvents(
				NewTeamCreatedEvent(teamID, "Team Awesome", "team-awesome", clubID, createdBy, createdAt),
			),
			teamName:      "Team Awesome",
			teamSlug:      "team-awesome",
			owningClubID:  clubID,
			createdBy:     createdBy,
			createdAt:     createdAt,
			expectedError: NewInvalidAggregateStateError(NewTeam(teamID).Aggregate(), int(TeamStateUnspecified), int(TeamStateActive)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			team := NewTeam(teamID)
			team.Reduce(tt.initialEvents)
			err := team.Init(tt.teamName, tt.teamSlug, tt.owningClubID, tt.createdBy, tt.createdAt)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, team.Changes().Events())
		})
	}
}

func TestTeam_Reduce(t *testing.T) {
	teamID := idgen.New[TeamID]()
	clubID := idgen.New[ClubID]()
	createdBy := NewOperator(idgen.New[AccountID](), nil)
	createdAt := time.Now()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		expectedState TeamState
		expectedName  string
		expectedSlug  string
		expectedClub  ClubID
	}{
		{
			name: "Succeeds if team state is updated correctly",
			initialEvents: createInitialEvents(
				NewTeamCreatedEvent(teamID, "Team Awesome", "team-awesome", clubID, createdBy, createdAt),
			),
			expectedState: TeamStateActive,
			expectedName:  "Team Awesome",
			expectedSlug:  "team-awesome",
			expectedClub:  clubID,
		},
		{
			name: "Succeeds if team is deleted",
			initialEvents: createInitialEvents(
				NewTeamCreatedEvent(teamID, "Team Awesome", "team-awesome", clubID, createdBy, createdAt),
				NewTeamDeletedEvent(teamID, clubID, createdBy),
			),
			expectedState: TeamStateDeleted,
			expectedName:  "Team Awesome",
			expectedSlug:  "team-awesome",
			expectedClub:  clubID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			team := NewTeam(teamID)
			team.Reduce(tt.initialEvents)
			assert.Equal(t, tt.expectedState, team.State)
			assert.Equal(t, tt.expectedName, team.Name)
			assert.Equal(t, tt.expectedSlug, team.Slug)
			assert.Equal(t, tt.expectedClub, team.OwningClubID)
		})
	}
}

func TestTeam_Delete(t *testing.T) {
	teamID := idgen.New[TeamID]()
	clubID := idgen.New[ClubID]()
	createdBy := NewOperator(idgen.New[AccountID](), nil)
	createdAt := time.Now()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		expectedError error
	}{
		{
			name: "Succeeds if team is deleted correctly",
			initialEvents: createInitialEvents(
				NewTeamCreatedEvent(teamID, "Team Awesome", "team-awesome", clubID, createdBy, createdAt),
			),
			emittedEvents: []eventing.Event{
				NewTeamDeletedEvent(teamID, clubID, createdBy),
			},
			expectedError: nil,
		},
		{
			name: "Fails if team is not active",
			initialEvents: createInitialEvents(
				NewTeamCreatedEvent(teamID, "Team Awesome", "team-awesome", clubID, createdBy, createdAt),
				NewTeamDeletedEvent(teamID, clubID, createdBy),
			),
			// TODO: Find a better way to construct this.
			expectedError: &InvalidAggregateStateError{
				Aggregate: &eventing.Aggregate{
					AggregateID:   eventing.AggregateID(teamID),
					AggregateType: TeamAggregateType,
					Version:       1,
				},
				ExpectedState: int(TeamStateActive),
				ActualState:   int(TeamStateDeleted),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			team := NewTeam(teamID)
			team.Reduce(tt.initialEvents)
			err := team.Delete(createdBy)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, team.Changes().Events())
		})
	}
}
