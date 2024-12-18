package domain

import (
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
)

type (
	TrainingID string
)

const (
	TrainingAggregateType = eventing.AggregateType("training")
)

var (
	ErrTrainingOwningClubNotFound = errors.New("owning club not found")
	ErrTrainingNotFound           = errors.New("team not found")
)

type TrainingState int

const (
	TrainingStateUnspecified TrainingState = iota
	TrainingStateActive
	TrainingStateDeleted
	TrainingStateCompleted
)

type Training struct {
	eventing.BaseWriter

	ID           TrainingID
	TeamID       TeamID
	OwningClubID ClubID

	State TrainingState
}

func NewTraining(id TrainingID, teamID TeamID, owningClubID ClubID) *Training {
	return &Training{
		BaseWriter:   *eventing.NewBaseWriter(eventing.AggregateID(id), TrainingAggregateType, eventing.VersionMatcherExact),
		ID:           id,
		TeamID:       teamID,
		OwningClubID: owningClubID,
	}
}

// NewTrainingByID returns a new Training with the given ID.
// It's expected that the remaining attributes are filled by the events.
func NewTrainingByID(id TrainingID) *Training {
	return &Training{
		BaseWriter: *eventing.NewBaseWriter(eventing.AggregateID(id), TrainingAggregateType, eventing.VersionMatcherExact),
		ID:         id,
	}
}

type TrainingGatheringPoint struct {
	Location        string    `json:"location"`
	GatherUntil     time.Time `json:"gather_until"`
	GatherUntilIANA string    `json:"gather_until_iana"`
}

func NewTrainingGatheringPoint(location string, gatherAt time.Time, gatherAtIANA string) *TrainingGatheringPoint {
	return &TrainingGatheringPoint{Location: location, GatherUntil: gatherAt, GatherUntilIANA: gatherAtIANA}
}

type TrainingAcknowledgmentSettings struct {
	AcknowledgeUntil     time.Time `json:"acknowledge_until"`
	AcknowledgeUntilIANA string    `json:"acknowledge_until_iana"`
}

func NewTrainingAcknowledgmentSettings(acknowledgeUntil time.Time, acknowledgeUntilIANA string) *TrainingAcknowledgmentSettings {
	return &TrainingAcknowledgmentSettings{AcknowledgeUntil: acknowledgeUntil, AcknowledgeUntilIANA: acknowledgeUntilIANA}
}

type TrainingRatingPolicy int

const (
	TrainingRatingPolicyUnspecified TrainingRatingPolicy = iota
	TrainingRatingPolicyForbidden
	TrainingRatingPolicyAllowed
	TrainingRatingPolicyRequired
)

type TrainingRatingSettings struct {
	Policy TrainingRatingPolicy `json:"policy"`
}

func NewTrainingRatingSettings(policy TrainingRatingPolicy) *TrainingRatingSettings {
	return &TrainingRatingSettings{Policy: policy}
}

func (t *Training) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.WithAggregate(TrainingAggregateType).
		AggregateID(eventing.AggregateID(t.ID)).
		Finish().MustBuild()
}

func (t *Training) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *TrainingScheduledEvent:
			t.State = TrainingStateActive
			t.ID = TrainingID(e.AggregateID())
			t.TeamID = e.TeamID
			t.OwningClubID = e.OwningClubID

		}
	}
	t.BaseWriter.Reduce(events)
}

func (t *Training) Schedule(
	scheduledAt time.Time,
	scheduledAtIANA string,
	endsAt time.Time,
	endsAtIANA string,
	description *string,
	location *string,
	fieldType *string,
	gatheringPoint *TrainingGatheringPoint,
	acknowledgmentSettings *TrainingAcknowledgmentSettings,
	ratingSettings TrainingRatingSettings,
	scheduledBy Operator,
) error {
	if t.State != TrainingStateUnspecified {
		return NewInvalidAggregateStateError(t.Aggregate(), int(TrainingStateUnspecified), int(t.State))
	}
	event := NewTrainingScheduledEvent(t.ID, scheduledAt, scheduledAtIANA, endsAt, endsAtIANA, description, location, fieldType, gatheringPoint, acknowledgmentSettings, ratingSettings, t.TeamID, t.OwningClubID, scheduledBy)
	t.Append(event)
	return nil
}
