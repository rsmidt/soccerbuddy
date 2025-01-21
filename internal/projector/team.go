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
	ID           domain.TeamID `json:"id"`
	Name         string        `json:"name"`
	Members      TeamMemberSet `json:"members"`
	OwningClubID domain.ClubID `json:"owning_club_id"`
}

func (p *TeamProjection) FindMember(personID domain.PersonID) (TeamMemberProjection, bool) {
	member, ok := p.Members[personID]
	return member, ok
}

type TeamMemberProjection struct {
	ID       domain.TeamMemberID   `json:"id"`
	PersonID domain.PersonID       `json:"person_id"`
	Name     string                `json:"name"`
	Role     domain.TeamMemberRole `json:"role"`
	JoinedAt time.Time             `json:"joined_at"`
}

type (
	TeamMemberSet map[domain.PersonID]TeamMemberProjection
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
		WithAggregate(domain.AccountAggregateType).
		Events(domain.AccountCreatedEventType, domain.RootAccountCreatedEventType, domain.AccountRegisteredEventType).Finish().
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
		case *domain.AccountCreatedEvent:
			err = r.handleAccountLookup(ctx, event, e)
		case *domain.RootAccountCreatedEvent:
			err = r.handleRootAccountLookup(ctx, event, e)
		case *domain.PersonCreatedEvent:
			err = r.handlePersonLookup(ctx, event, e)
		case *domain.PersonInvitedToTeamEvent:
			err = r.insertPersonInvitedToTeamEvent(ctx, event, e)
		case *domain.AccountRegisteredEvent:
			err = r.handleRegisteredAccountLookup(ctx, event, e)
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
		ID:      domain.TeamID(event.AggregateID()),
		Name:    e.Name,
		Members: make(TeamMemberSet),
	}
	return insertJSON(ctx, r.rd, r.key(p.ID), &p)
}

func (r *rdTeamProjector) insertTeamDeletedEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.TeamDeletedEvent) error {
	key := r.key(domain.TeamID(e.AggregateID()))
	cmd := r.rd.B().Del().Key(key).Build()
	return r.rd.Do(ctx, cmd).Error()
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
