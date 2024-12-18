package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
)

// ========================================================
// TrainingScheduledEvent
// ========================================================

const (
	TrainingScheduledEventType    = eventing.EventType("training_scheduled")
	TrainingScheduledEventVersion = eventing.EventVersion("v1")
)

var _ eventing.Event = (*TrainingScheduledEvent)(nil)

type TrainingScheduledEvent struct {
	*eventing.EventBase

	ScheduledAt     time.Time `json:"scheduled_at"`
	ScheduledAtIANA string    `json:"scheduled_at_iana"`
	EndsAt          time.Time `json:"ends_at"`
	EndsAtIANA      string    `json:"ends_at_iana"`

	Description            *string                         `json:"description"`
	Location               *string                         `json:"location"`
	FieldType              *string                         `json:"field_type"`
	GatheringPoint         *TrainingGatheringPoint         `json:"gathering_point"`
	AcknowledgmentSettings *TrainingAcknowledgmentSettings `json:"acknowledgment_settings"`
	RatingSettings         TrainingRatingSettings          `json:"rating_settings"`

	TeamID       TeamID   `json:"team_id"`
	OwningClubID ClubID   `json:"owning_club_id"`
	ScheduledBy  Operator `json:"scheduled_by"`
}

func NewTrainingScheduledEvent(
	id TrainingID,
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
	teamID TeamID,
	owningClubID ClubID,
	scheduledBy Operator,
) *TrainingScheduledEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), TrainingAggregateType, TrainingScheduledEventVersion, TrainingScheduledEventType)

	return &TrainingScheduledEvent{
		EventBase:              base,
		ScheduledAt:            scheduledAt,
		ScheduledAtIANA:        scheduledAtIANA,
		EndsAt:                 endsAt,
		EndsAtIANA:             endsAtIANA,
		Description:            description,
		Location:               location,
		FieldType:              fieldType,
		GatheringPoint:         gatheringPoint,
		AcknowledgmentSettings: acknowledgmentSettings,
		RatingSettings:         ratingSettings,
		TeamID:                 teamID,
		OwningClubID:           owningClubID,
		ScheduledBy:            scheduledBy,
	}
}

func (t *TrainingScheduledEvent) IsShredded() bool {
	return false
}
