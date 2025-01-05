package projector

import (
	"context"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

type Supervisors struct {
	Postgres eventing.ProjectorSupervisor
	Redis    eventing.ProjectorSupervisor
}

func (m *Supervisors) Register(ctx context.Context, relationStore authz.RelationStore, rd rueidis.Client) error {
	permProjector := NewPermissionProjector(relationStore)
	if err := permProjector.Init(ctx); err != nil {
		return err
	}
	personProjector := NewPersonProjector(rd)
	if err := personProjector.Init(ctx); err != nil {
		return err
	}
	accountProjector := NewAccountProjector(rd)
	if err := accountProjector.Init(ctx); err != nil {
		return err
	}
	teamProjector := NewTeamProjector(rd)
	if err := teamProjector.Init(ctx); err != nil {
		return err
	}
	trainingProjector := NewTrainingProjector(rd)
	if err := trainingProjector.Init(ctx); err != nil {
		return err
	}

	m.Postgres.Register(permProjector)
	m.Redis.Register(personProjector)
	m.Redis.Register(accountProjector)
	m.Redis.Register(teamProjector)
	m.Redis.Register(trainingProjector)
	return nil
}

func (m *Supervisors) Enable() {
	m.Postgres.Enable()
	m.Redis.Enable()
}
