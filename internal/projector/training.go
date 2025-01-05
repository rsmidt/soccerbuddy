package projector

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"time"
)

const (
	ProjectionTrainingName    eventing.ProjectionName = "trainings"
	ProjectionTrainingIDXName                         = "projectionTrainingV1Idx"
	ProjectionTrainingPrefix                          = "projection:trainings:v1:"
)

type TrainingProjection struct {
	ID domain.TrainingID `json:"id"`

	ScheduledAt     time.Time `json:"scheduled_at"`
	ScheduledAtIANA string    `json:"scheduled_at_iana"`
	ScheduledAtTS   int64     `json:"scheduled_at_ts"`
	EndsAt          time.Time `json:"ends_at"`
	EndsAtIANA      string    `json:"ends_at_iana"`
	EndsAtTS        int64     `json:"ends_at_ts"`

	Description *string `json:"description"`
	Location    *string `json:"location"`
	FieldType   *string `json:"field_type"`

	// TODO: Add gathering point etc.
	GatheringPoint         *TrainingGatheringPointProjection         `json:"gathering_point"`
	AcknowledgmentSettings *TrainingAcknowledgmentSettingsProjection `json:"acknowledgment_settings"`
	RatingSettings         TrainingRatingSettingsProjection          `json:"rating_settings"`

	NominatedPlayers TrainingNominatedPlayerSet `json:"nominated_players"`
	NominatedStaff   TrainingNominatedPlayerSet `json:"nominated_staff"`

	ScheduledBy OperatorProjection `json:"scheduled_by"`

	OwningClubID domain.ClubID  `json:"owning_club_id"`
	OwningTeamID *domain.TeamID `json:"owning_team_id"`
}

type TrainingGatheringPointProjection struct {
	Location        string    `json:"location"`
	GatherUntil     time.Time `json:"gather_until"`
	GatherUntilIANA string    `json:"gather_until_iana"`
}

type TrainingAcknowledgmentSettingsProjection struct {
	AcknowledgeUntil     time.Time `json:"acknowledge_until"`
	AcknowledgeUntilIANA string    `json:"acknowledge_until_iana"`
}

type TrainingRatingSettingsProjection struct {
	Policy domain.TrainingRatingPolicy `json:"policy"`
}

type TrainingNominatedPersonProjection struct {
	ID             domain.PersonID                            `json:"id"`
	Name           string                                     `json:"name"`
	Role           domain.TeamMemberRole                      `json:"role"`
	Acknowledgment TrainingNominationAcknowledgmentProjection `json:"acknowledgment_status"`
	NominatedAt    time.Time                                  `json:"nominated_at"`
	NominatedBy    OperatorProjection                         `json:"nominated_by"`
}

type TrainingNominationAcknowledgmentProjection struct {
	Type           domain.TrainingNominationAcknowledgmentType `json:"type,omitempty"`
	AcknowledgedAt *time.Time                                  `json:"acknowledged_at,omitempty"`
	AcceptedAt     *time.Time                                  `json:"accepted_at,omitempty"`
	DeclinedAt     *time.Time                                  `json:"declined_at,omitempty"`
	AcknowledgedBy *OperatorProjection                         `json:"acknowledged_by,omitempty"`
	Reason         *string                                     `json:"reason,omitempty"`
}

type (
	TrainingNominatedPlayerSet map[domain.PersonID]TrainingNominatedPersonProjection
)

type rdTrainingProjector struct {
	rd rueidis.Client
}

func NewTrainingProjector(rd rueidis.Client) eventing.Projector {
	return &rdTrainingProjector{rd: rd}
}

func (r *rdTrainingProjector) Init(ctx context.Context) error {
	ctx, span := tracing.Tracer.Start(ctx, "projector.redis.Training.Init")
	defer span.End()

	cmd := r.rd.B().
		FtCreate().
		Index(ProjectionTrainingIDXName).
		OnJson().
		Prefix(1).
		Prefix(ProjectionTrainingPrefix).
		Schema().
		FieldName("$.owning_club_id").As("owning_club_id").Text().
		FieldName("$.owning_team_id").As("owning_team_id").Text().
		FieldName("$.scheduled_at_ts").As("scheduled_at_ts").Numeric().Sortable().
		FieldName("$.ends_at_ts").As("ends_at_ts").Numeric().Sortable().
		Build()
	if err := r.rd.Do(ctx, cmd).Error(); err != nil {
		rderr, ok := rueidis.IsRedisErr(err)
		if ok && rderr.Error() == "Index already exists" {
			return nil
		}
		return err
	}
	return nil
}

func (r *rdTrainingProjector) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.AccountAggregateType).
		Events(domain.AccountCreatedEventType, domain.RootAccountCreatedEventType).Finish().
		WithAggregate(domain.PersonAggregateType).
		Events(domain.PersonCreatedEventType).Finish().
		WithAggregate(domain.TeamMemberAggregateType).
		Events(domain.PersonInvitedToTeamEventType).Finish().
		WithAggregate(domain.TrainingAggregateType).
		Events(domain.TrainingScheduledEventType, domain.PersonsNominatedForTrainingEventType).Finish().
		MustBuild()
}

func (r *rdTrainingProjector) Projection() eventing.ProjectionName {
	return ProjectionTrainingName
}

func (r *rdTrainingProjector) Project(ctx context.Context, events ...*eventing.JournalEvent) error {
	var err error
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.TrainingScheduledEvent:
			err = r.insertTrainingScheduledEvent(ctx, event, e)
		case *domain.PersonsNominatedForTrainingEvent:
			err = r.insertPersonsNominatedForTrainingEvent(ctx, event, e)
		case *domain.AccountCreatedEvent:
			err = r.handleAccountLookup(ctx, event, e)
		case *domain.RootAccountCreatedEvent:
			err = r.handleRootAccountLookup(ctx, event, e)
		case *domain.PersonCreatedEvent:
			err = r.handlePersonLookup(ctx, event, e)
		case *domain.PersonInvitedToTeamEvent:
			err = r.handleTeamMemberLookup(ctx, event, e)
		}
		if err != nil {
			tracing.RecordError(ctx, err)
			return err
		}
	}
	return nil
}

func (r *rdTrainingProjector) getProjection(ctx context.Context, id domain.TrainingID) (*TrainingProjection, error) {
	var p TrainingProjection
	cmd := r.rd.B().JsonGet().Key(r.key(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdTrainingProjector) insertTrainingScheduledEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.TrainingScheduledEvent) error {
	actor, err := r.lookupAccount(ctx, e.ScheduledBy.ActorID)
	if err != nil {
		return err
	}
	trainingID := domain.TrainingID(event.AggregateID())
	var gatheringPoint *TrainingGatheringPointProjection
	if e.GatheringPoint != nil {
		gatheringPoint = &TrainingGatheringPointProjection{
			Location:        e.GatheringPoint.Location,
			GatherUntil:     e.GatheringPoint.GatherUntil,
			GatherUntilIANA: e.GatheringPoint.GatherUntilIANA,
		}
	}
	var acknowledgmentSettings *TrainingAcknowledgmentSettingsProjection
	if e.AcknowledgmentSettings != nil {
		acknowledgmentSettings = &TrainingAcknowledgmentSettingsProjection{
			AcknowledgeUntil:     e.AcknowledgmentSettings.AcknowledgeUntil,
			AcknowledgeUntilIANA: e.AcknowledgmentSettings.AcknowledgeUntilIANA,
		}
	}
	training := TrainingProjection{
		ID:                     trainingID,
		ScheduledAt:            e.ScheduledAt,
		ScheduledAtIANA:        e.ScheduledAtIANA,
		ScheduledAtTS:          e.ScheduledAt.Unix(),
		EndsAt:                 e.EndsAt,
		EndsAtIANA:             e.EndsAtIANA,
		EndsAtTS:               e.EndsAt.Unix(),
		Description:            e.Description,
		Location:               e.Location,
		FieldType:              e.FieldType,
		GatheringPoint:         gatheringPoint,
		AcknowledgmentSettings: acknowledgmentSettings,
		RatingSettings: TrainingRatingSettingsProjection{
			Policy: e.RatingSettings.Policy,
		},
		ScheduledBy: OperatorProjection{
			ActorID:       e.ScheduledBy.ActorID,
			ActorFullName: actor.FullName,
			OnBehalfOf:    e.ScheduledBy.OnBehalfOf,
		},
		NominatedStaff:   make(TrainingNominatedPlayerSet),
		NominatedPlayers: make(TrainingNominatedPlayerSet),
		OwningClubID:     e.OwningClubID,
		OwningTeamID:     &e.TeamID,
	}
	return insertJSON(ctx, r.rd, r.key(trainingID), training)
}

func (r *rdTrainingProjector) key(id domain.TrainingID) string {
	return fmt.Sprintf("%s%s", ProjectionTrainingPrefix, id)
}

func (r *rdTrainingProjector) insertPersonsNominatedForTrainingEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonsNominatedForTrainingEvent) error {
	trainingID := domain.TrainingID(event.AggregateID())
	projection, err := r.getProjection(ctx, trainingID)
	if err != nil {
		return err
	}

	actor, err := r.lookupAccount(ctx, e.NominatedBy.ActorID)
	if err != nil {
		return err
	}

	for _, playerID := range e.NominatedPlayers {
		person, err := r.lookupPerson(ctx, playerID)
		if err != nil {
			return err
		}
		projection.NominatedPlayers[playerID] = TrainingNominatedPersonProjection{
			ID:   playerID,
			Name: person.FullName,
			Role: domain.TeamMemberRolePlayer,
			Acknowledgment: TrainingNominationAcknowledgmentProjection{
				Type: domain.TrainingNominationUnacknowledged,
			},
			NominatedAt: event.InsertedAt(),
			NominatedBy: OperatorProjection{
				ActorID:       e.NominatedBy.ActorID,
				ActorFullName: actor.FullName,
				OnBehalfOf:    e.NominatedBy.OnBehalfOf,
			},
		}
	}

	for _, staffID := range e.NominatedStaff {
		person, err := r.lookupPerson(ctx, staffID)
		if err != nil {
			return err
		}
		role := domain.TeamMemberRoleGuest
		if e.TeamID != nil {
			role, err = r.lookupPersonTeamRole(ctx, person.ID, *e.TeamID)
			if err != nil {
				return err
			}
		}
		projection.NominatedStaff[staffID] = TrainingNominatedPersonProjection{
			ID:   staffID,
			Name: person.FullName,
			Role: role,
			Acknowledgment: TrainingNominationAcknowledgmentProjection{
				Type: domain.TrainingNominationUnacknowledged,
			},
			NominatedAt: event.InsertedAt(),
			NominatedBy: OperatorProjection{
				ActorID:       e.NominatedBy.ActorID,
				ActorFullName: actor.FullName,
				OnBehalfOf:    e.NominatedBy.OnBehalfOf,
			},
		}
	}

	return insertJSON(ctx, r.rd, r.key(trainingID), projection)
}
