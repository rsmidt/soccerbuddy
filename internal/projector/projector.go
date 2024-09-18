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

	m.Postgres.Register(permProjector)
	m.Redis.Register(personProjector)
	return nil
}

func (m *Supervisors) Enable() {
	m.Postgres.Enable()
	m.Redis.Enable()
}
