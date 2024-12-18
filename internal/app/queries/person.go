package queries

import (
	"context"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/rsmidt/soccerbuddy/internal/redis"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"log/slog"
	"time"
)

type personInClubView struct {
	ID        domain.PersonID `json:"id"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Birthdate time.Time       `json:"birthdate"`
}

type PersonsInClubView struct {
	Persons []*personInClubView
}

type ListPersonsInClubQuery struct {
	OwningClubID domain.ClubID
}

func (q *Queries) ListPersonsInClub(ctx context.Context, query ListPersonsInClubQuery) (*PersonsInClubView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.ListPersonsInClub")
	defer span.End()

	// TODO: use more appropriate relation than edit.
	if err := q.authorizer.Authorize(ctx, authz.ActionEdit, authz.NewClubResource(query.OwningClubID)); err != nil {
		return nil, err
	}

	// TODO: this should be more abstracted.
	rdq := fmt.Sprintf("@owning_club_id:(%s)", query.OwningClubID)
	cmd := q.rd.B().FtSearch().Index(projector.ProjectionPersonIDXName).Query(rdq).Build()
	_, docs, err := q.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	persons, err := redis.UnmarshalDocs[personInClubView](docs)
	if err != nil {
		return nil, err
	}
	return &PersonsInClubView{Persons: persons}, nil
}

type operatorView struct {
	FullName string
}

type teamView struct {
	ID       domain.TeamID
	Name     string
	Role     domain.TeamMemberRoleRole
	JoinedAt time.Time
}

type linkedAccountView struct {
	FullName  string
	LinkedAs  domain.AccountLink
	LinkedAt  time.Time
	InvitedBy *operatorView
	InvitedAt *time.Time
	LinkedBy  *operatorView
}

type pendingAccountLinkView struct {
	LinkedAs  domain.AccountLink
	InvitedBy operatorView
	InvitedAt time.Time
	ExpiresAt time.Time
}

type PersonOverview struct {
	ID                  domain.PersonID
	FirstName           string
	LastName            string
	Birthdate           time.Time
	CreatedAt           time.Time
	CreatedBy           operatorView
	LinkedAccounts      []*linkedAccountView
	PendingAccountLinks []*pendingAccountLinkView
	Teams               []*teamView
}

type GetPersonOverviewQuery struct {
	ID domain.PersonID
}

func (q *Queries) GetPersonOverview(ctx context.Context, query GetPersonOverviewQuery) (*PersonOverview, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.GetPersonOverview")
	defer span.End()

	if err := q.authorizer.Authorize(ctx, authz.ActionView, authz.NewPersonResource(query.ID)); err != nil {
		return nil, err
	}

	projection, err := q.getPersonProjection(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	ts := make([]*teamView, len(projection.Teams))
	for i, t := range projection.Teams {
		ts[i] = &teamView{
			ID:       t.ID,
			Name:     t.Name,
			Role:     t.Role,
			JoinedAt: t.JoinedAt,
		}
	}
	pl := make([]*pendingAccountLinkView, len(projection.PendingLinks))
	for i, p := range projection.PendingLinks {
		pl[i] = &pendingAccountLinkView{
			LinkedAs:  p.LinkAs,
			InvitedBy: operatorView{FullName: p.InvitedBy.ActorFullName},
			InvitedAt: p.InvitedAt,
			ExpiresAt: p.ExpiresAt,
		}
	}
	la := make([]*linkedAccountView, len(projection.LinkedAccounts))
	for i, l := range projection.LinkedAccounts {
		var invitedBy *operatorView
		if l.InvitedBy != nil {
			invitedBy = &operatorView{FullName: l.InvitedBy.ActorFullName}
		}
		var linkedBy *operatorView
		if l.LinkedBy != nil {
			linkedBy = &operatorView{FullName: l.LinkedBy.ActorFullName}
		}
		la[i] = &linkedAccountView{
			FullName:  l.FullName,
			LinkedAs:  l.LinkedAs,
			LinkedAt:  l.LinkedAt,
			InvitedBy: invitedBy,
			InvitedAt: l.InvitedAt,
			LinkedBy:  linkedBy,
		}
	}
	return &PersonOverview{
		ID:        projection.ID,
		FirstName: projection.FirstName,
		LastName:  projection.LastName,
		Birthdate: projection.BirthDate,
		CreatedAt: projection.CreatedAt,
		CreatedBy: operatorView{
			FullName: projection.CreatedBy.ActorFullName,
		},
		Teams:               ts,
		LinkedAccounts:      la,
		PendingAccountLinks: pl,
	}, nil
}

type PendingPersonLinkView struct {
	FullName  string
	LinkAs    domain.AccountLink
	InvitedBy operatorView
	Club      *pendingPersonLinkClubView
}

type pendingPersonLinkClubView struct {
	ID   domain.ClubID
	Name string
}

type DescribePendingPersonLinkQuery struct {
	LinkToken domain.PersonLinkToken
}

func (q *Queries) DescribePendingPersonLink(ctx context.Context, query DescribePendingPersonLinkQuery) (*PendingPersonLinkView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.DescribePendingPersonLink")
	defer span.End()

	// If not authenticated, the user can either register or login.
	_, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrPrincipalNotFound
	}

	persons, err := q.getPersonProjectionByPendingToken(ctx, query.LinkToken)
	if err != nil {
		return nil, err
	}
	var p *projector.PersonProjection
	if len(persons) == 0 {
		return nil, domain.ErrPersonInvalidLinkToken
	} else if len(persons) > 1 {
		q.log.WarnContext(ctx, "found multiple persons for the same link token; taking first person now", slog.String("link_token", string(query.LinkToken)))
	}
	p = persons[0]
	var pl *projector.PendingLinkProjection
	for _, link := range p.PendingLinks {
		if link.Token == query.LinkToken {
			pl = link
			break
		}
	}
	return &PendingPersonLinkView{
		FullName: fmt.Sprintf("%s %s", p.FirstName, p.LastName),
		LinkAs:   pl.LinkAs,
		Club: &pendingPersonLinkClubView{
			ID:   p.Club.ID,
			Name: p.Club.Name,
		},
		InvitedBy: operatorView{
			FullName: pl.InvitedBy.ActorFullName,
		},
	}, nil
}
