package viewstore

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/app/view"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
)

type ViewStores interface {
	Team() TeamViewStore
	Account() AccountViewStore
	Person() PersonViewStore
}

type TeamViewStore interface {
	GetHome(ctx context.Context, permissions authz.PermissionsSet, teamID domain.TeamID) (*view.TeamHome, error)
}

type AccountViewStore interface {
	GetMe(ctx context.Context, id domain.AccountID) (*view.GetMe, error)
}

type PersonViewStore interface {
	GetOverview(ctx context.Context, id domain.PersonID) (*view.PersonOverview, error)
	DescribePendingPersonLink(ctx context.Context, linkToken domain.PersonLinkToken) (*view.PendingPersonLink, error)
}
