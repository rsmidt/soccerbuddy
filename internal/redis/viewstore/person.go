package viewstore

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/app/view"
	"github.com/rsmidt/soccerbuddy/internal/app/viewstore"
	"github.com/rsmidt/soccerbuddy/internal/core"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"log/slog"
)

var _ viewstore.PersonViewStore = (*RedisPersonViewStore)(nil)

func NewRedisPersonViewStore(rd rueidis.Client, logger *slog.Logger) *RedisPersonViewStore {
	return &RedisPersonViewStore{rd: rd, logger: logger}
}

type RedisPersonViewStore struct {
	rd     rueidis.Client
	logger *slog.Logger
}

func (r *RedisPersonViewStore) DescribePendingPersonLink(ctx context.Context, linkToken domain.PersonLinkToken) (*view.PendingPersonLink, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisPersonViewStore) GetOverview(ctx context.Context, id domain.PersonID) (*view.PersonOverview, error) {
	ctx, span := tracing.Tracer.Start(ctx, "redisviewstore.GetOverview")
	defer span.End()

	projection, err := getPersonProjection(ctx, r.rd, id)
	if err != nil {
		return nil, err
	}

	ts := make([]*view.PersonOverviewTeam, len(projection.Teams))
	for i, t := range projection.Teams {
		ts[i] = &view.PersonOverviewTeam{
			ID:       t.ID,
			Name:     t.Name,
			Role:     t.Role,
			JoinedAt: t.JoinedAt,
		}
	}
	pl := make([]*view.PersonOverviewPendingAccountLink, len(projection.PendingLinks))
	for i, p := range projection.PendingLinks {
		pl[i] = &view.PersonOverviewPendingAccountLink{
			LinkedAs:  p.LinkAs,
			InvitedBy: view.Operator{FullName: p.InvitedBy.ActorFullName},
			InvitedAt: p.InvitedAt,
			ExpiresAt: p.ExpiresAt,
		}
	}
	la := make([]*view.PersonOverviewLinkedAccount, len(projection.LinkedAccounts))
	for i, l := range projection.LinkedAccounts {
		var invitedBy *view.Operator
		if l.InvitedBy != nil {
			invitedBy = &view.Operator{FullName: l.InvitedBy.ActorFullName}
		}
		var linkedBy *view.Operator
		if l.LinkedBy != nil {
			linkedBy = &view.Operator{FullName: l.LinkedBy.ActorFullName}
		}
		la[i] = &view.PersonOverviewLinkedAccount{
			FullName:  l.FullName,
			LinkedAs:  l.LinkedAs,
			LinkedAt:  l.LinkedAt,
			InvitedBy: invitedBy,
			InvitedAt: l.InvitedAt,
			LinkedBy:  linkedBy,
		}
	}
	return &view.PersonOverview{
		ID:        projection.ID,
		FirstName: projection.FirstName,
		LastName:  projection.LastName,
		Birthdate: projection.BirthDate,
		CreatedAt: projection.CreatedAt,
		CreatedBy: view.Operator{
			FullName: projection.CreatedBy.ActorFullName,
		},
		Teams:               ts,
		LinkedAccounts:      la,
		PendingAccountLinks: pl,
	}, nil
}

func getPersonProjection(ctx context.Context, rd rueidis.Client, id domain.PersonID) (*projector.PersonProjection, error) {
	var p projector.PersonProjection
	cmd := rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projector.ProjectionPersonPrefix, id)).Path(".").Build()
	if err := rd.Do(ctx, cmd).DecodeJSON(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func getPersonProjections(ctx context.Context, rd rueidis.Client, ds []domain.PersonID) ([]*projector.PersonProjection, error) {
	if len(ds) == 0 {
		return nil, nil
	}

	var p []*projector.PersonProjection
	keys := make([]string, len(ds))
	for i, d := range ds {
		keys[i] = fmt.Sprintf("%s%s", projector.ProjectionPersonPrefix, d)
	}
	cmd := rd.B().JsonMget().Key(keys...).Path(".").Build()
	if err := rueidis.DecodeSliceOfJSON(rd.Do(ctx, cmd), &p); err != nil {
		return nil, err
	}
	return core.RemoveNils(p), nil
}
