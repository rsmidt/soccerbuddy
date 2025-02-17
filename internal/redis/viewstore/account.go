package viewstore

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/app/view"
	"github.com/rsmidt/soccerbuddy/internal/app/viewstore"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"log/slog"
	"maps"
	"slices"
)

var _ viewstore.AccountViewStore = (*RedisAccountViewStore)(nil)

type RedisAccountViewStore struct {
	rd     rueidis.Client
	logger *slog.Logger
}

func NewRedisAccountViewStore(rd rueidis.Client, logger *slog.Logger) *RedisAccountViewStore {
	return &RedisAccountViewStore{rd: rd, logger: logger}
}

func (r *RedisAccountViewStore) GetMe(ctx context.Context, id domain.AccountID) (*view.GetMe, error) {
	ctx, span := tracing.Tracer.Start(ctx, "redisviewstore.GetMe")
	defer span.End()

	return getMe(ctx, r.rd, r.logger, id)
}

func getMe(ctx context.Context, rd rueidis.Client, logger *slog.Logger, accountID domain.AccountID) (*view.GetMe, error) {
	account, err := getAccountProjection(ctx, rd, accountID)
	if err != nil {
		return nil, err
	}
	personIDs := slices.Collect(maps.Keys(account.LinkedPersons))
	persons, err := getPersonProjections(ctx, rd, personIDs)
	if err != nil {
		return nil, err
	}
	linkedPersons := make([]*view.GetMeLinkedPerson, 0, len(persons))
	for _, p := range persons {
		link, ok := account.LinkedPersons[p.ID]
		if !ok {
			// Received a result from projection that is not linked to this account.
			logger.WarnContext(ctx, "Received person (%s) from projection that is not linked to account (%s).", slog.String(string(p.ID), string(account.ID)))
			continue
		}
		memberships := make([]*view.GetMeLinkedPersonTeamMemberships, len(p.Teams))
		for i, t := range p.Teams {
			memberships[i] = &view.GetMeLinkedPersonTeamMemberships{
				ID:           t.ID,
				Name:         t.Name,
				JoinedAt:     t.JoinedAt,
				Role:         t.Role,
				OwningClubID: t.OwningClubID,
			}
		}
		v := &view.GetMeLinkedPerson{
			ID:              p.ID,
			FirstName:       p.FirstName,
			LastName:        p.LastName,
			LinkedAs:        link.LinkedAs,
			LinkedAt:        link.LinkedAt,
			LinkedBy:        mapOperatorToGetMeOperatorView(account.ID, link.LinkedBy),
			TeamMemberships: memberships,
			OwningClubID:    p.OwningClubID,
		}
		linkedPersons = append(linkedPersons, v)
	}
	return &view.GetMe{
		ID:            account.ID,
		Email:         account.Email,
		FirstName:     account.FirstName,
		LastName:      account.LastName,
		LinkedPersons: linkedPersons,
		IsSuper:       account.IsRoot,
	}, nil
}

func mapOperatorToGetMeOperatorView(accountID domain.AccountID, op *projector.OperatorProjection) *view.GetMeOperator {
	if op == nil {
		return nil
	}
	return &view.GetMeOperator{
		FullName: op.ActorFullName,
		IsMe:     op.ActorID == accountID,
	}
}

func getAccountProjection(ctx context.Context, rd rueidis.Client, id domain.AccountID) (*projector.AccountProjection, error) {
	var a projector.AccountProjection
	cmd := rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projector.ProjectionAccountPrefix, id)).Path(".").Build()
	return &a, rd.Do(ctx, cmd).DecodeJSON(&a)
}
