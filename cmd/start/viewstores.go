package main

import (
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/app/viewstore"
	rdvs "github.com/rsmidt/soccerbuddy/internal/redis/viewstore"
	"log/slog"
)

func assembleViewStores(logger *slog.Logger, rd rueidis.Client) viewstore.ViewStores {
	return &viewStores{
		teamVS:    rdvs.NewRedisTeamViewStore(rd, logger),
		accountVS: rdvs.NewRedisAccountViewStore(rd, logger),
	}
}

type viewStores struct {
	teamVS    viewstore.TeamViewStore
	accountVS viewstore.AccountViewStore
}

func (v viewStores) Team() viewstore.TeamViewStore {
	return v.teamVS
}

func (v viewStores) Account() viewstore.AccountViewStore {
	return v.accountVS
}
