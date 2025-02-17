package view

import (
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"time"
)

type TeamHome struct {
	ID           domain.TeamID
	Name         string
	Trainings    []*TeamHomeTraining
	OwningClubID domain.ClubID
}

type TeamHomeTraining struct {
	ID domain.TrainingID

	ScheduledAt     time.Time
	ScheduledAtIANA string
	EndsAt          time.Time
	EndsAtIANA      string

	GatheringPoint         *GatheringPoint
	AcknowledgmentSettings *AcknowledgmentSettings
	RatingSettings         RatingSettings

	// Nominations will only be set if enough rights are available.
	Nominations *Nominations

	Description *string
	Location    *string
	FieldType   *string

	ScheduledBy Operator
}

type GatheringPoint struct {
	Location        string
	GatherUntil     time.Time
	GatherUntilIANA string
}

type AcknowledgmentSettings struct {
	AcknowledgedUntil     time.Time
	AcknowledgedUntilIANA string
}

type RatingSettings struct {
	Policy domain.TrainingRatingPolicy
}

type Nominations struct {
	Players []*TrainingNominationResponse
	Staff   []*TrainingNominationResponse
}

type TrainingNominationResponse struct {
	PersonID       domain.PersonID
	PersonName     string
	Type           domain.TrainingNominationAcknowledgmentType
	AcknowledgedAt *time.Time
	AcceptedAt     *time.Time
	TentativeAt    *time.Time
	DeclinedAt     *time.Time
	AcknowledgedBy *Operator
	Reason         *string
	NominatedAt    time.Time
}
