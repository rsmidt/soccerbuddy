package projector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"time"
)

const (
	ProjectionTeamName    eventing.ProjectionName = "teams"
	ProjectionTeamIDXName                         = "projectionTeamV1Idx"
	ProjectionTeamPrefix                          = "projection:teams:v1:"
)

type TeamProjection struct {
	ID           domain.TeamID             `json:"id"`
	Name         string                    `json:"name"`
	Trainings    TeamTrainingProjectionSet `json:"events"`
	Members      TeamMemberSet             `json:"members"`
	OwningClubID domain.ClubID             `json:"owning_club_id"`
}

type TeamTrainingProjection struct {
	ID domain.TrainingID `json:"id"`

	ScheduledAt     time.Time `json:"scheduled_at"`
	ScheduledAtIANA string    `json:"scheduled_at_iana"`
	EndsAt          time.Time `json:"ends_at"`
	EndsAtIANA      string    `json:"ends_at_iana"`

	Description *string `json:"description"`
	Location    *string `json:"location"`
	FieldType   *string `json:"field_type"`

	// TODO: Add gathering point etc.

	NominatedPersons TeamTrainingNominatedPlayerSet `json:"nominated_persons"`
	NominatedStaff   TeamTrainingNominatedPlayerSet `json:"nominated_staff"`

	ScheduledBy OperatorProjection `json:"scheduled_by"`
}

type TeamTrainingNominatedPlayerProjection struct {
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

type TeamMemberProjection struct {
	ID       domain.TeamMemberID   `json:"id"`
	PersonID domain.PersonID       `json:"person_id"`
	Name     string                `json:"name"`
	Role     domain.TeamMemberRole `json:"role"`
	JoinedAt time.Time             `json:"joined_at"`
}

type (
	TeamTrainingProjectionSet      map[domain.TrainingID]TeamTrainingProjection
	TeamTrainingNominatedPlayerSet map[domain.PersonID]TeamTrainingNominatedPlayerProjection
	TeamMemberSet                  map[domain.PersonID]TeamMemberProjection
)

type rdTeamProjector struct {
	rd rueidis.Client
}

func NewTeamProjector(rd rueidis.Client) eventing.Projector {
	return &rdTeamProjector{rd: rd}
}

func (r *rdTeamProjector) Init(ctx context.Context) error {
	return nil
}

func (r *rdTeamProjector) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.TeamAggregateType).
		Events(domain.TeamCreatedEventType, domain.TeamDeletedEventType).Finish().
		WithAggregate(domain.TrainingAggregateType).
		Events(domain.TrainingScheduledEventType, domain.PersonsNominatedForTrainingEventType).Finish().
		WithAggregate(domain.AccountAggregateType).
		Events(domain.AccountCreatedEventType, domain.RootAccountCreatedEventType).Finish().
		WithAggregate(domain.PersonAggregateType).
		Events(domain.PersonCreatedEventType).Finish().
		WithAggregate(domain.TeamMemberAggregateType).
		Events(domain.PersonInvitedToTeamEventType).Finish().
		MustBuild()
}

func (r *rdTeamProjector) Projection() eventing.ProjectionName {
	return ProjectionTeamName
}

func (r *rdTeamProjector) Project(ctx context.Context, events ...*eventing.JournalEvent) error {
	var err error
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.TeamCreatedEvent:
			err = r.insertTeamCreatedEvent(ctx, event, e)
		case *domain.TeamDeletedEvent:
			err = r.insertTeamDeletedEvent(ctx, event, e)
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
			err = r.insertPersonInvitedToTeamEvent(ctx, event, e)
		}
		if err != nil {
			tracing.RecordError(ctx, err)
			return err
		}
	}
	return nil
}

func (r *rdTeamProjector) getProjection(ctx context.Context, id domain.TeamID) (*TeamProjection, error) {
	var p TeamProjection
	cmd := r.rd.B().JsonGet().Key(r.key(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdTeamProjector) insertTeamCreatedEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.TeamCreatedEvent) error {
	p := TeamProjection{
		ID:        domain.TeamID(event.AggregateID()),
		Name:      e.Name,
		Trainings: make(TeamTrainingProjectionSet),
		Members:   make(TeamMemberSet),
	}
	return insertJSON(ctx, r.rd, r.key(p.ID), &p)
}

func (r *rdTeamProjector) insertTrainingScheduledEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.TrainingScheduledEvent) error {
	p, err := r.getProjection(ctx, e.TeamID)
	if err != nil {
		return err
	}
	actor, err := r.lookupAccount(ctx, e.ScheduledBy.ActorID)
	if err != nil {
		return err
	}
	trainingID := domain.TrainingID(event.AggregateID())
	training := TeamTrainingProjection{
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
		NominatedStaff:   make(TeamTrainingNominatedPlayerSet),
		NominatedPersons: make(TeamTrainingNominatedPlayerSet),
	}
	p.Trainings[trainingID] = training
	return insertJSON(ctx, r.rd, r.key(p.ID), p)
}

func (r *rdTeamProjector) insertTeamDeletedEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.TeamDeletedEvent) error {
	key := r.key(domain.TeamID(e.AggregateID()))
	cmd := r.rd.B().Del().Key(key).Build()
	return r.rd.Do(ctx, cmd).Error()
}

func (r *rdTeamProjector) insertPersonsNominatedForTrainingEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonsNominatedForTrainingEvent) error {
	// Skip this training as it is not associated with a team.
	if e.TeamID == nil {
		return nil
	}
	p, err := r.getProjection(ctx, *e.TeamID)
	if err != nil {
		return err
	}

	trainingID := domain.TrainingID(event.AggregateID())
	training, ok := p.Trainings[trainingID]
	if !ok {
		return fmt.Errorf("training %s not found in team %s", trainingID, *e.TeamID)
	}

	nominator, err := r.lookupAccount(ctx, e.NominatedBy.ActorID)
	if err != nil {
		return err
	}

	for _, playerID := range e.NominatedPlayers {
		// TODO: Batch the lookup.
		person, err := r.lookupPerson(ctx, playerID)
		if err != nil {
			return err
		}
		// TODO: Properly handle guest nominations.
		var role domain.TeamMemberRole = "GUEST"
		if member, ok := p.Members[playerID]; ok {
			role = member.Role
		}
		training.NominatedPersons[playerID] = TeamTrainingNominatedPlayerProjection{
			ID:   playerID,
			Name: person.FullName,
			Role: role,
			Acknowledgment: TrainingNominationAcknowledgmentProjection{
				Type: domain.TrainingNominationUnacknowledged,
			},
			NominatedAt: event.InsertedAt(),
			NominatedBy: OperatorProjection{
				ActorID:       e.NominatedBy.ActorID,
				ActorFullName: nominator.FullName,
				OnBehalfOf:    e.NominatedBy.OnBehalfOf,
			},
		}
	}

	for _, staffID := range e.NominatedStaff {
		// TODO: Batch the lookup.
		person, err := r.lookupPerson(ctx, staffID)
		if err != nil {
			return err
		}
		// TODO: Properly handle guest nominations.
		var role domain.TeamMemberRole = "GUEST"
		if member, ok := p.Members[staffID]; ok {
			role = member.Role
		}
		training.NominatedStaff[staffID] = TeamTrainingNominatedPlayerProjection{
			ID:   staffID,
			Name: person.FullName,
			Role: role,
			Acknowledgment: TrainingNominationAcknowledgmentProjection{
				Type: domain.TrainingNominationUnacknowledged,
			},
			NominatedAt: event.InsertedAt(),
			NominatedBy: OperatorProjection{
				ActorID:       e.NominatedBy.ActorID,
				ActorFullName: nominator.FullName,
				OnBehalfOf:    e.NominatedBy.OnBehalfOf,
			},
		}
	}

	p.Trainings[trainingID] = training
	return insertJSON(ctx, r.rd, r.key(p.ID), p)
}

func (r *rdTeamProjector) key(id domain.TeamID) string {
	return fmt.Sprintf("%s%s", ProjectionTeamPrefix, id)
}

func (r *rdTeamProjector) insertPersonInvitedToTeamEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonInvitedToTeamEvent) error {
	lookup, err := r.lookupPerson(ctx, e.PersonID)
	if err != nil {
		return err
	}
	person := TeamMemberProjection{
		ID:       domain.TeamMemberID(e.AggregateID()),
		PersonID: e.PersonID,
		Name:     lookup.FullName,
		Role:     e.AssignedRole,
		JoinedAt: event.InsertedAt(),
	}

	val, err := json.Marshal(&person)
	if err != nil {
		return err
	}

	// Only update the member property.
	cmd := r.rd.B().JsonSet().Key(r.key(e.TeamID)).Path(fmt.Sprintf(".members.%s", person.PersonID)).Value(string(val)).Build()
	res := r.rd.Do(ctx, cmd)
	return res.Error()
}
