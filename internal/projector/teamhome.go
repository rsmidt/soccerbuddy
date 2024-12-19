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
	ProjectionTeamHomeName    eventing.ProjectionName = "team_home"
	ProjectionTeamHomeIDXName                         = "projectionTeamHomeV1Idx"
	ProjectionTeamHomePrefix                          = "projection:team_homes:v1:"
)

type TeamHomeProjection struct {
	ID           domain.TeamID                 `json:"id"`
	Name         string                        `json:"name"`
	Trainings    TeamHomeTrainingProjectionSet `json:"events"`
	OwningClubID domain.ClubID                 `json:"owning_club_id"`
}

type TeamHomeTrainingProjection struct {
	ID domain.TrainingID `json:"id"`

	ScheduledAt     time.Time `json:"scheduled_at"`
	ScheduledAtIANA string    `json:"scheduled_at_iana"`
	EndsAt          time.Time `json:"ends_at"`
	EndsAtIANA      string    `json:"ends_at_iana"`

	Description *string `json:"description"`
	Location    *string `json:"location"`
	FieldType   *string `json:"field_type"`

	// TODO: Add gathering point etc.

	ScheduledBy OperatorProjection `json:"scheduled_by"`
}

type TeamHomeTrainingProjectionSet map[domain.TrainingID]TeamHomeTrainingProjection

type rdTeamHomeProjector struct {
	rd rueidis.Client
}

func NewTeamHomeProjector(rd rueidis.Client) eventing.Projector {
	return &rdTeamHomeProjector{rd: rd}
}

func (r *rdTeamHomeProjector) Init(ctx context.Context) error {
	return nil
}

func (r *rdTeamHomeProjector) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.TeamAggregateType).
		Events(domain.TeamCreatedEventType, domain.TeamDeletedEventType).Finish().
		WithAggregate(domain.TrainingAggregateType).
		Events(domain.TrainingScheduledEventType).Finish().
		WithAggregate(domain.AccountAggregateType).
		Events(domain.AccountCreatedEventType, domain.RootAccountCreatedEventType).Finish().
		MustBuild()
}

func (r *rdTeamHomeProjector) Projection() eventing.ProjectionName {
	return ProjectionTeamHomeName
}

func (r *rdTeamHomeProjector) Project(ctx context.Context, events ...*eventing.JournalEvent) error {
	var err error
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.TeamCreatedEvent:
			err = r.insertTeamCreatedEvent(ctx, event, e)
		case *domain.TeamDeletedEvent:
			err = r.insertTeamDeletedEvent(ctx, event, e)
		case *domain.TrainingScheduledEvent:
			err = r.insertTrainingScheduledEvent(ctx, event, e)
		case *domain.AccountCreatedEvent:
			err = r.handleAccountLookup(ctx, event, e)
		case *domain.RootAccountCreatedEvent:
			err = r.handleRootAccountLookup(ctx, event, e)
		}
		if err != nil {
			tracing.RecordError(ctx, err)
			return err
		}
	}
	return nil
}

func (r *rdTeamHomeProjector) getProjection(ctx context.Context, id domain.TeamID) (*TeamHomeProjection, error) {
	var p TeamHomeProjection
	cmd := r.rd.B().JsonGet().Key(r.key(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdTeamHomeProjector) insertTeamCreatedEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.TeamCreatedEvent) error {
	p := TeamHomeProjection{
		ID:        domain.TeamID(event.AggregateID()),
		Name:      e.Name,
		Trainings: make(TeamHomeTrainingProjectionSet),
	}

	return insertJSON(ctx, r.rd, r.key(p.ID), &p)
}

func (r *rdTeamHomeProjector) insertTrainingScheduledEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.TrainingScheduledEvent) error {
	p, err := r.getProjection(ctx, e.TeamID)
	if err != nil {
		return err
	}
	actor, err := r.lookupAccount(ctx, e.ScheduledBy.ActorID)
	if err != nil {
		return err
	}
	trainingID := domain.TrainingID(event.AggregateID())
	training := TeamHomeTrainingProjection{
		ID:              trainingID,
		ScheduledAt:     e.ScheduledAt,
		ScheduledAtIANA: e.ScheduledAtIANA,
		EndsAt:          e.EndsAt,
		EndsAtIANA:      e.EndsAtIANA,
		Description:     e.Description,
		Location:        e.Location,
		FieldType:       e.FieldType,
		ScheduledBy: OperatorProjection{
			ActorID:       e.ScheduledBy.ActorID,
			ActorFullName: actor.FullName,
			OnBehalfOf:    e.ScheduledBy.OnBehalfOf,
		},
	}
	p.Trainings[trainingID] = training
	return insertJSON(ctx, r.rd, r.key(p.ID), p)
}

func (r *rdTeamHomeProjector) insertTeamDeletedEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.TeamDeletedEvent) error {
	key := r.key(domain.TeamID(e.AggregateID()))
	cmd := r.rd.B().Del().Key(key).Build()
	return r.rd.Do(ctx, cmd).Error()
}

func (r *rdTeamHomeProjector) key(id domain.TeamID) string {
	return fmt.Sprintf("%s%s", ProjectionTeamHomePrefix, id)
}
