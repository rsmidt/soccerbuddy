package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClub_Init(t *testing.T) {
	clubID := idgen.New[ClubID]()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		clubName      string
		clubSlug      string
		expectedError error
	}{
		{
			name:          "Succeeds if not yet initialized",
			initialEvents: createInitialEvents(),
			emittedEvents: []eventing.Event{
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome"),
			},
			clubName:      "FC Awesome",
			clubSlug:      "fc-awesome",
			expectedError: nil,
		},
		{
			name: "Fails if already initialized",
			initialEvents: createInitialEvents(
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome"),
			),
			clubName:      "FC Awesome",
			clubSlug:      "fc-awesome",
			expectedError: NewInvalidAggregateStateError(NewClub(clubID).Aggregate(), int(ClubStateUnspecified), int(ClubStateActive)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			club := NewClub(clubID)
			club.Reduce(tt.initialEvents)
			err := club.Init(tt.clubName, tt.clubSlug)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, club.Changes().Events())
		})
	}
}

func TestClub_Reduce(t *testing.T) {
	clubID := idgen.New[ClubID]()
	testBase := eventing.NewBaseWriter(eventing.AggregateID(clubID), ClubAggregateType, eventing.VersionMatcherExact)

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		expectedState Club
	}{
		{
			name: "Reduce ClubCreatedEvent",
			initialEvents: createInitialEvents(
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome"),
			),
			expectedState: Club{
				BaseWriter: *testBase,
				State:      ClubStateActive,
				ID:         clubID,
				Name:       "FC Awesome",
				Slug:       "fc-awesome",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			club := NewClub(clubID)
			club.Reduce(tt.initialEvents)
			assert.Equal(t, tt.expectedState.ID, club.ID)
			assert.Equal(t, tt.expectedState.Name, club.Name)
			assert.Equal(t, tt.expectedState.State, club.State)
			assert.Equal(t, tt.expectedState.Slug, club.Slug)
			assert.True(t, tt.expectedState.BaseWriter.Equals(&club.BaseWriter))
		})
	}
}
