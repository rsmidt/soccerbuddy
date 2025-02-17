package queries

import (
	"context"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/app/view"
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

	rdq := fmt.Sprintf("@owning_club_id:{%s}", query.OwningClubID)
	cmd := q.rd.B().FtSearch().Index(projector.ProjectionPersonIDXName).Query(rdq).Dialect(4).Build()
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

type GetPersonOverviewQuery struct {
	ID domain.PersonID
}

func (q *Queries) GetPersonOverview(ctx context.Context, query GetPersonOverviewQuery) (*view.PersonOverview, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.GetOverview")
	defer span.End()

	if err := q.authorizer.Authorize(ctx, authz.ActionView, authz.NewPersonResource(query.ID)); err != nil {
		return nil, err
	}
	return q.vs.Person().GetOverview(ctx, query.ID)
}

type PendingPersonLinkView struct {
	FullName  string
	LinkAs    domain.AccountLink
	InvitedBy view.Operator
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
	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrPrincipalNotFound
	}

	// Get the account details for the user that wants to link.
	account, err := q.vs.Account().GetMe(ctx, principal.AccountID)
	if err != nil {
		return nil, err
	}

	persons, err := q.getPersonProjectionByPendingToken(ctx, query.LinkToken)
	if err != nil {
		return nil, err
	}
	var p *projector.PersonProjection
	if len(persons) == 0 {
		// If we couldn't find any person by projection, it most likely mean that the token was already used or is invalid.
		// Next, we check if the token was used by this account.
		// If this link is already used, we can signal that to frontend.
		for _, linkedPerson := range account.LinkedPersons {
			if linkedPerson.UsedLinkToken != nil && *linkedPerson.UsedLinkToken == query.LinkToken {
				return nil, domain.ErrAccountAlreadyLinkedToPerson
			}
		}
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
	if pl == nil {
		// Should never happen.
		q.log.ErrorContext(ctx, "was able to fetch person projection by pending link token but data does not contain pending link")
		return nil, domain.ErrPersonInvalidLinkToken
	}

	// Check if the requesting account was already linked to this person before.
	if _, ok := account.LinkedPersons[p.ID]; ok {
		return nil, domain.ErrAccountAlreadyLinkedToPerson
	}

	return &PendingPersonLinkView{
		FullName: fmt.Sprintf("%s %s", p.FirstName, p.LastName),
		LinkAs:   pl.LinkAs,
		Club: &pendingPersonLinkClubView{
			ID:   p.Club.ID,
			Name: p.Club.Name,
		},
		InvitedBy: view.Operator{
			FullName: pl.InvitedBy.ActorFullName,
		},
	}, nil
}
