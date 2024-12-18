package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTraining_Schedule(t *testing.T) {
	trainingID := idgen.New[TrainingID]()
	teamID := idgen.New[TeamID]()
	clubID := idgen.New[ClubID]()
	scheduledBy := NewOperator(idgen.New[AccountID](), nil)

	scheduledAt := time.Now()
	endsAt := scheduledAt.Add(2 * time.Hour)
	description := "Team practice"
	location := "Main field"
	fieldType := "Grass"

	gatheringPoint := NewTrainingGatheringPoint(
		"Meeting point",
		scheduledAt.Add(-30*time.Minute),
		"Europe/Berlin",
	)

	acknowledgmentSettings := NewTrainingAcknowledgmentSettings(
		scheduledAt.Add(-1*time.Hour),
		"Europe/Berlin",
	)

	ratingSettings := *NewTrainingRatingSettings(TrainingRatingPolicyAllowed)

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		expectedError error
	}{
		{
			name:          "Succeeds if training is scheduled correctly",
			initialEvents: createInitialEvents(),
			emittedEvents: []eventing.Event{
				NewTrainingScheduledEvent(
					trainingID,
					scheduledAt,
					"Europe/Berlin",
					endsAt,
					"Europe/Berlin",
					&description,
					&location,
					&fieldType,
					gatheringPoint,
					acknowledgmentSettings,
					ratingSettings,
					teamID,
					clubID,
					scheduledBy,
				),
			},
			expectedError: nil,
		},
		{
			name: "Fails if training is already scheduled",
			initialEvents: createInitialEvents(
				NewTrainingScheduledEvent(
					trainingID,
					scheduledAt,
					"Europe/Berlin",
					endsAt,
					"Europe/Berlin",
					&description,
					&location,
					&fieldType,
					gatheringPoint,
					acknowledgmentSettings,
					ratingSettings,
					teamID,
					clubID,
					scheduledBy,
				),
			),
			expectedError: NewInvalidAggregateStateError(
				NewTraining(trainingID, teamID, clubID).Aggregate(),
				int(TrainingStateUnspecified),
				int(TrainingStateActive),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			training := NewTraining(trainingID, teamID, clubID)
			training.Reduce(tt.initialEvents)

			err := training.Schedule(
				scheduledAt,
				"Europe/Berlin",
				endsAt,
				"Europe/Berlin",
				&description,
				&location,
				&fieldType,
				gatheringPoint,
				acknowledgmentSettings,
				ratingSettings,
				scheduledBy,
			)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, training.Changes().Events())
		})
	}
}

func TestTraining_Reduce(t *testing.T) {
	trainingID := idgen.New[TrainingID]()
	teamID := idgen.New[TeamID]()
	clubID := idgen.New[ClubID]()
	scheduledBy := NewOperator(idgen.New[AccountID](), nil)

	scheduledAt := time.Now()
	endsAt := scheduledAt.Add(2 * time.Hour)
	description := "Team practice"
	location := "Main field"
	fieldType := "Grass"

	gatheringPoint := NewTrainingGatheringPoint(
		"Meeting point",
		scheduledAt.Add(-30*time.Minute),
		"Europe/Berlin",
	)

	acknowledgmentSettings := NewTrainingAcknowledgmentSettings(
		scheduledAt.Add(-1*time.Hour),
		"Europe/Berlin",
	)

	ratingSettings := *NewTrainingRatingSettings(TrainingRatingPolicyAllowed)

	tests := []struct {
		name           string
		initialEvents  []*eventing.JournalEvent
		expectedState  TrainingState
		expectedTeamID TeamID
		expectedClubID ClubID
	}{
		{
			name: "Succeeds if training state is updated correctly",
			initialEvents: createInitialEvents(
				NewTrainingScheduledEvent(
					trainingID,
					scheduledAt,
					"Europe/Berlin",
					endsAt,
					"Europe/Berlin",
					&description,
					&location,
					&fieldType,
					gatheringPoint,
					acknowledgmentSettings,
					ratingSettings,
					teamID,
					clubID,
					scheduledBy,
				),
			),
			expectedState:  TrainingStateActive,
			expectedTeamID: teamID,
			expectedClubID: clubID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			training := NewTrainingByID(trainingID)
			training.Reduce(tt.initialEvents)
			assert.Equal(t, tt.expectedState, training.State)
			assert.Equal(t, tt.expectedTeamID, training.TeamID)
			assert.Equal(t, tt.expectedClubID, training.OwningClubID)
		})
	}
}
