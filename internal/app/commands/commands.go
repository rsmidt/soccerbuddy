package commands

import (
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"log/slog"
)

type Commands struct {
	log        *slog.Logger
	es         eventing.EventStore
	authorizer authz.Authorizer
	rd         rueidis.Client
	repos      domain.Repositories
}

func NewCommands(
	log *slog.Logger,
	es eventing.EventStore,
	authorizer authz.Authorizer,
	rd rueidis.Client,
	repos domain.Repositories,
) *Commands {
	return &Commands{log: log, es: es, authorizer: authorizer, rd: rd, repos: repos}
}
