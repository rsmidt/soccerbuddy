package queries

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"log/slog"
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
	IsSuper       bool
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
	Role         domain.TeamMemberRole
}

func (q *Queries) GetMe(ctx context.Context) (*GetMeView, error) {
	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrUnauthenticated
	}

	return q.getMe(ctx, principal.AccountID)
}

func (q *Queries) getMe(ctx context.Context, accountID domain.AccountID) (*GetMeView, error) {
	account, err := q.getAccountProjection(ctx, accountID)
	if err != nil {
		return nil, err
	}
	personIDs := slices.Collect(maps.Keys(account.LinkedPersons))
	persons, err := q.getPersonProjections(ctx, personIDs)
	if err != nil {
		return nil, err
	}
	linkedPersons := make([]*GetMeLinkedPersonView, 0, len(persons))
	for _, p := range persons {
		link, ok := account.LinkedPersons[p.ID]
		if !ok {
			// Received a result from projection that is not linked to this account.
			q.log.WarnContext(ctx, "Received person (%s) from projection that is not linked to account (%s).", slog.String(string(p.ID), string(account.ID)))
			continue
		}
		memberships := make([]*GetMeLinkedPersonTeamMembershipsView, len(p.Teams))
		for i, t := range p.Teams {
			memberships[i] = &GetMeLinkedPersonTeamMembershipsView{
				ID:           t.ID,
				Name:         t.Name,
				JoinedAt:     t.JoinedAt,
				Role:         t.Role,
				OwningClubID: t.OwningClubID,
			}
		}
		view := &GetMeLinkedPersonView{
			ID:              p.ID,
			FirstName:       p.FirstName,
			LastName:        p.LastName,
			LinkedAs:        link.LinkedAs,
			LinkedAt:        link.LinkedAt,
			LinkedBy:        mapOperatorToGetMeOperatorView(account.ID, link.LinkedBy),
			TeamMemberships: memberships,
			OwningClubID:    p.OwningClubID,
		}
		linkedPersons = append(linkedPersons, view)
	}
	return &GetMeView{
		ID:            account.ID,
		Email:         account.Email,
		FirstName:     account.FirstName,
		LastName:      account.LastName,
		LinkedPersons: linkedPersons,
		IsSuper:       account.IsRoot,
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
