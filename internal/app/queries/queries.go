package queries

import (
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"log/slog"
)

type Queries struct {
	log        *slog.Logger
	es         eventing.EventStore
	authorizer authz.Authorizer
	rd         rueidis.Client

	// Deprecated: use a proper view model.
	repos domain.Repositories
}

func NewQueries(
	log *slog.Logger,
	es eventing.EventStore,
	authorizer authz.Authorizer,
	rd rueidis.Client,
	repos domain.Repositories,
) *Queries {
	return &Queries{
		log:        log,
		es:         es,
		authorizer: authorizer,
		rd:         rd,
		repos:      repos,
	}
}
