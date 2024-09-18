package projector

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"slices"
	"time"
)

const (
	ProjectionPersonName    eventing.ProjectionName = "persons"
	ProjectionPersonIDXName                         = "projectionPersonV1Idx"
	ProjectionPersonPrefix                          = "projection:persons:v1:"
)

type OperatorProjection struct {
	ActorID       domain.AccountID `json:"actor_id"`
	ActorFullName string           `json:"actor_full_name"`
	OnBehalfOf    *domain.PersonID `json:"acting_on_behalf_of"`
}

type PersonProjection struct {
	ID           domain.PersonID    `json:"id"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	BirthDate    time.Time          `json:"birth_date"`
	OwningClubID domain.ClubID      `json:"owning_club_id"`
	CreatedAt    time.Time          `json:"created_at"`
	CreatedBy    OperatorProjection `json:"created_by"`
	// TODO: convert all to maps for better idempotence.
	PendingLinks   []*PendingLinkProjection   `json:"pending_links"`
	LinkedAccounts []*LinkedAccountProjection `json:"linked_accounts"`
	Teams          []*teamProjection          `json:"teams"`
	Club           clubProjection             `json:"club"`
}

type teamProjection struct {
	ID       domain.TeamID             `json:"id"`
	Name     string                    `json:"name"`
	Role     domain.TeamMemberRoleRole `json:"role"`
	JoinedAt time.Time                 `json:"joined_at"`
}

type clubProjection struct {
	ID   domain.ClubID `json:"id"`
	Name string        `json:"name"`
}

type PendingLinkProjection struct {
	Token     domain.PersonLinkToken `json:"token"`
	LinkAs    domain.AccountLink     `json:"link_as"`
	ExpiresAt time.Time              `json:"expires_at"`
	InvitedBy OperatorProjection     `json:"invited_by"`
	InvitedAt time.Time              `json:"invited_at"`
}

type LinkedAccountProjection struct {
	LinkedAs domain.AccountLink `json:"linked_as"`
	LinkedAt time.Time          `json:"linked_at"`
	FullName string             `json:"full_name"`
	// LinkedBy is only set if InvitedBy is not set and vice versa.
	LinkedBy  *OperatorProjection `json:"linked_by"`
	InvitedBy *OperatorProjection `json:"invited_by"`
	InvitedAt *time.Time          `json:"invited_at"`
}

type rdPersonProjector struct {
	rd rueidis.Client
}

func NewPersonProjector(rd rueidis.Client) eventing.Projector {
	return &rdPersonProjector{rd: rd}
}

func (r *rdPersonProjector) Init(ctx context.Context) error {
	ctx, span := tracing.Tracer.Start(ctx, "projector.Person.Init")
	defer span.End()

	cmd := r.rd.B().
		FtCreate().
		Index(ProjectionPersonIDXName).
		OnJson().
		Prefix(1).
		Prefix(ProjectionPersonPrefix).
		Schema().
		FieldName("$.owning_club_id").As("owning_club_id").Text().
		FieldName("$.first_name").As("first_name").Text().
		FieldName("$.last_name").As("last_name").Text().
		FieldName("$.pending_links[0:].token").As("pending_link_token").Text().
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

func (r *rdPersonProjector) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.PersonAggregateType).
		Events(domain.PersonCreatedEventType, domain.PersonLinkInitiatedEventType, domain.PersonLinkClaimedEventType).Finish().
		WithAggregate(domain.AccountAggregateType).
		Events(domain.AccountCreatedEventType, domain.RootAccountCreatedEventType).Finish().
		WithAggregate(domain.TeamAggregateType).
		Events(domain.TeamCreatedEventType).Finish().
		WithAggregate(domain.TeamMemberAggregateType).
		Events(domain.PersonInvitedToTeamEventType).Finish().
		WithAggregate(domain.ClubAggregateType).
		Events(domain.ClubCreatedEventType).Finish().
		MustBuild()
}

func (r *rdPersonProjector) Projection() eventing.ProjectionName {
	return ProjectionPersonName
}

func (r *rdPersonProjector) Project(ctx context.Context, events ...*eventing.JournalEvent) error {
	var err error
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.PersonCreatedEvent:
			err = r.insertPerson(ctx, event, e)
		case *domain.PersonLinkInitiatedEvent:
			err = r.insertPendingLink(ctx, event, e)
		case *domain.PersonLinkClaimedEvent:
			err = r.handleLinkClaimed(ctx, event, e)
		case *domain.AccountCreatedEvent:
			err = r.handleAccountLookup(ctx, event, e)
		case *domain.RootAccountCreatedEvent:
			err = r.handleRootAccountLookup(ctx, event, e)
		case *domain.TeamCreatedEvent:
			err = r.handleTeamLookup(ctx, event, e)
		case *domain.ClubCreatedEvent:
			err = r.handleClubLookup(ctx, event, e)
		case *domain.PersonInvitedToTeamEvent:
			err = r.insertTeamMember(ctx, event, e)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rdPersonProjector) insertPerson(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonCreatedEvent) error {
	acc, err := r.lookupAccount(ctx, e.Creator.ActorID)
	if err != nil {
		return err
	}
	c, err := r.lookupClub(ctx, e.OwningClubID)
	if err != nil {
		return err
	}

	projection := PersonProjection{
		ID:           domain.PersonID(e.AggregateID()),
		FirstName:    e.FirstName.Value,
		LastName:     e.LastName.Value,
		BirthDate:    maybeParseTime(e.Birthdate.Value),
		OwningClubID: e.OwningClubID,
		CreatedBy: OperatorProjection{
			ActorID:       e.Creator.ActorID,
			ActorFullName: acc.FullName,
			OnBehalfOf:    e.Creator.OnBehalfOf,
		},
		CreatedAt: event.InsertedAt(),
		Club: clubProjection{
			ID:   c.ID,
			Name: c.Name,
		},
	}
	key := fmt.Sprintf("%s%s", ProjectionPersonPrefix, event.AggregateID())
	return insertJSON(ctx, r.rd, key, &projection)
}

func (r *rdPersonProjector) insertTeamMember(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonInvitedToTeamEvent) error {
	t, err := r.lookupTeam(ctx, e.TeamID)
	if err != nil {
		return err
	}
	projection, err := r.getProjection(ctx, e.PersonID)
	if err != nil {
		return err
	}
	projection.Teams = append(projection.Teams, &teamProjection{
		ID:       e.TeamID,
		Name:     t.Name,
		Role:     e.AssignedRole,
		JoinedAt: event.InsertedAt(),
	})
	key := fmt.Sprintf("%s%s", ProjectionPersonPrefix, projection.ID)
	return insertJSON(ctx, r.rd, key, projection)
}

func (r *rdPersonProjector) insertPendingLink(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonLinkInitiatedEvent) error {
	projection, err := r.getProjection(ctx, domain.PersonID(e.AggregateID()))
	if err != nil {
		return err
	}
	inviter, err := r.lookupAccount(ctx, e.InvitedBy.ActorID)
	if err != nil {
		return err
	}
	projection.PendingLinks = append(projection.PendingLinks, &PendingLinkProjection{
		Token:     e.Token,
		LinkAs:    e.LinkAs,
		ExpiresAt: e.ExpiresAt,
		InvitedBy: OperatorProjection{
			ActorID:       e.InvitedBy.ActorID,
			ActorFullName: inviter.FullName,
			OnBehalfOf:    e.InvitedBy.OnBehalfOf,
		},
		InvitedAt: event.InsertedAt(),
	})
	key := fmt.Sprintf("%s%s", ProjectionPersonPrefix, projection.ID)
	return insertJSON(ctx, r.rd, key, projection)

}

func (r *rdPersonProjector) getProjection(ctx context.Context, id domain.PersonID) (*PersonProjection, error) {
	var p PersonProjection
	cmd := r.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", ProjectionPersonPrefix, id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdPersonProjector) handleLinkClaimed(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonLinkClaimedEvent) error {
	projection, err := r.getProjection(ctx, domain.PersonID(e.AggregateID()))
	if err != nil {
		return err
	}
	var pl *PendingLinkProjection
	for _, p := range projection.PendingLinks {
		if p.Token == e.UsedToken {
			pl = p
		}
	}
	if pl == nil {
		return fmt.Errorf("no pending link found for token %s", e.UsedToken)
	}
	projection.PendingLinks = slices.DeleteFunc(projection.PendingLinks, func(projection *PendingLinkProjection) bool {
		return projection.Token == e.UsedToken
	})
	linkedAccount, err := r.lookupAccount(ctx, e.AccountID)
	if err != nil {
		return err
	}
	projection.LinkedAccounts = append(projection.LinkedAccounts, &LinkedAccountProjection{
		LinkedAs:  e.LinkedAs,
		LinkedAt:  event.InsertedAt(),
		FullName:  linkedAccount.FullName,
		LinkedBy:  nil,
		InvitedBy: &pl.InvitedBy,
		InvitedAt: &pl.InvitedAt,
	})
	key := fmt.Sprintf("%s%s", ProjectionPersonPrefix, projection.ID)
	return insertJSON(ctx, r.rd, key, projection)
}

func maybeParseTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
