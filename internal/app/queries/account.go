package queries

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/sourcegraph/conc/iter"
	"maps"
	"slices"
	"time"
)

type GetMeView struct {
	ID            domain.AccountID
	Email         string
	FirstName     string
	LastName      string
	LinkedPersons []*GetMeLinkedPersonView
}

type GetMeOperatorView struct {
	FullName string
	// If the operator that performed this account is this account owner.
	IsMe bool
}

type GetMeLinkedPersonView struct {
	ID              domain.PersonID
	FirstName       string
	LastName        string
	LinkedAs        domain.AccountLink
	LinkedAt        time.Time
	LinkedBy        *GetMeOperatorView
	TeamMemberships []*GetMeLinkedPersonTeamMembershipsView
	OwningClubID    domain.ClubID
}

type GetMeLinkedPersonTeamMembershipsView struct {
	ID           domain.TeamID
	Name         string
	OwningClubID domain.ClubID
	JoinedAt     time.Time
	Roles        domain.TeamMemberRoleRole
}

func (q *Queries) GetMe(ctx context.Context) (*GetMeView, error) {
	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrUnauthenticated
	}

	account, err := q.getAccountProjection(ctx, principal.AccountID)
	if err != nil {
		return nil, err
	}
	personIDs := slices.Collect(maps.Keys(account.LinkedPersons))
	persons, err := q.getPersonProjections(ctx, personIDs)
	if err != nil {
		return nil, err
	}
	linkedPersons, err := iter.MapErr(persons, func(t **projector.PersonProjection) (*GetMeLinkedPersonView, error) {
		p := *t
		link, ok := account.LinkedPersons[p.ID]
		if !ok {
			// Received a result from projection that is not linked to this account.
			q.log.WarnContext(ctx, "Received person (%s) from projection that is not linked to account (%s).", p.ID, account.ID)
			return nil, nil
		}
		memberships := make([]*GetMeLinkedPersonTeamMembershipsView, len(p.Teams))
		for i, t := range p.Teams {
			memberships[i] = &GetMeLinkedPersonTeamMembershipsView{
				ID:           t.ID,
				Name:         t.Name,
				JoinedAt:     t.JoinedAt,
				Roles:        t.Role,
				OwningClubID: t.OwningClubID,
			}
		}
		return &GetMeLinkedPersonView{
			ID:              p.ID,
			FirstName:       p.FirstName,
			LastName:        p.LastName,
			LinkedAs:        link.LinkedAs,
			LinkedAt:        link.LinkedAt,
			LinkedBy:        mapOperatorToGetMeOperatorView(account.ID, link.LinkedBy),
			TeamMemberships: memberships,
			OwningClubID:    p.OwningClubID,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return &GetMeView{
		ID:            account.ID,
		Email:         account.Email,
		FirstName:     account.FirstName,
		LastName:      account.LastName,
		LinkedPersons: linkedPersons,
	}, nil
}

func mapOperatorToGetMeOperatorView(accountID domain.AccountID, op *projector.OperatorProjection) *GetMeOperatorView {
	if op == nil {
		return nil
	}
	return &GetMeOperatorView{
		FullName: op.ActorFullName,
		IsMe:     op.ActorID == accountID,
	}
}
