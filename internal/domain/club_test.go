package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClub_Init(t *testing.T) {
	clubID := idgen.New[ClubID]()
	createdAt := time.Now()

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
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome", createdAt),
			},
			clubName:      "FC Awesome",
			clubSlug:      "fc-awesome",
			expectedError: nil,
		},
		{
			name: "Fails if already initialized",
			initialEvents: createInitialEvents(
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome", createdAt),
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
			err := club.Init(tt.clubName, tt.clubSlug, createdAt)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, club.Changes().Events())
		})
	}
}

func TestClub_AddAdmin(t *testing.T) {
	clubID := idgen.New[ClubID]()
	accountID := idgen.New[AccountID]()
	now := time.Now()
	operator := NewOperator(idgen.New[AccountID](), nil)

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		accountID     AccountID
		addedAt       time.Time
		addedBy       Operator
		expectedError error
	}{
		{
			name: "Succeeds if club is active",
			initialEvents: createInitialEvents(
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome", now),
			),
			emittedEvents: []eventing.Event{
				NewClubAdminAddedEvent(clubID, accountID, now, operator),
			},
			accountID:     accountID,
			addedAt:       now,
			addedBy:       operator,
			expectedError: nil,
		},
		{
			name:          "Fails if club is not initialized",
			initialEvents: createInitialEvents(),
			accountID:     accountID,
			addedAt:       now,
			addedBy:       operator,
			expectedError: NewInvalidAggregateStateError(NewClub(clubID).Aggregate(), int(ClubStateActive), int(ClubStateUnspecified)),
		},
		{
			name: "No event emitted when adding same admin twice",
			initialEvents: createInitialEvents(
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome", now),
				NewClubAdminAddedEvent(clubID, accountID, now, operator),
			),
			accountID:     accountID,
			addedAt:       now,
			addedBy:       operator,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			club := NewClub(clubID)
			club.Reduce(tt.initialEvents)
			err := club.AddAdmin(tt.accountID, tt.addedAt, tt.addedBy)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, club.Changes().Events())
		})

	}
}

func TestClub_Reduce(t *testing.T) {
	clubID := idgen.New[ClubID]()
	accountID := idgen.New[AccountID]()
	now := time.Now()
	operator := NewOperator(idgen.New[AccountID](), nil)
	testBase := eventing.NewBaseWriter(eventing.AggregateID(clubID), ClubAggregateType, eventing.VersionMatcherExact)

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		expectedState Club
	}{
		{
			name: "Reduce ClubCreatedEvent",
			initialEvents: createInitialEvents(
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome", now),
			),
			expectedState: Club{
				BaseWriter: *testBase,
				State:      ClubStateActive,
				ID:         clubID,
				Name:       "FC Awesome",
				Slug:       "fc-awesome",
				Admins:     make(AdminsSet),
			},
		},
		{
			name: "Reduce ClubCreatedEvent and ClubAdminAddedEvent",
			initialEvents: createInitialEvents(
				NewClubCreatedEvent(clubID, "FC Awesome", "fc-awesome", now),
				NewClubAdminAddedEvent(clubID, accountID, now, operator),
			),
			expectedState: Club{
				BaseWriter: *testBase,
				State:      ClubStateActive,
				ID:         clubID,
				Name:       "FC Awesome",
				Slug:       "fc-awesome",
				Admins: AdminsSet{
					accountID: struct{}{},
				},
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
			assert.Equal(t, tt.expectedState.Admins, club.Admins)
			tt.expectedState.Reduce(tt.initialEvents)
			assert.True(t, tt.expectedState.BaseWriter.Equals(&club.BaseWriter))
		})
	}
}
