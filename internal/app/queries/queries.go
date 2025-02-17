package queries

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/app/viewstore"
	"github.com/rsmidt/soccerbuddy/internal/core"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/rsmidt/soccerbuddy/internal/redis"
	"log/slog"
)

type Queries struct {
	log        *slog.Logger
	es         eventing.EventStore
	authorizer authz.Authorizer
	rd         rueidis.Client
	vs         viewstore.ViewStores

	// Deprecated: use a proper view model.
	repos domain.Repositories
}

func NewQueries(
	log *slog.Logger,
	es eventing.EventStore,
	authorizer authz.Authorizer,
	rd rueidis.Client,
	repos domain.Repositories,
	vs viewstore.ViewStores,
) *Queries {
	return &Queries{
		log:        log,
		es:         es,
		authorizer: authorizer,
		rd:         rd,
		vs:         vs,
		repos:      repos,
	}
}

func (q *Queries) getAccountProjection(ctx context.Context, id domain.AccountID) (*projector.AccountProjection, error) {
	var a projector.AccountProjection
	cmd := q.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projector.ProjectionAccountPrefix, id)).Path(".").Build()
	return &a, q.rd.Do(ctx, cmd).DecodeJSON(&a)
}

func (q *Queries) getPersonProjection(ctx context.Context, id domain.PersonID) (*projector.PersonProjection, error) {
	var p projector.PersonProjection
	cmd := q.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projector.ProjectionPersonPrefix, id)).Path(".").Build()
	if err := q.rd.Do(ctx, cmd).DecodeJSON(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (q *Queries) getPersonProjectionByPendingToken(ctx context.Context, token domain.PersonLinkToken) ([]*projector.PersonProjection, error) {
	rdq := fmt.Sprintf("@pending_link_token:{%s}", token)
	cmd := q.rd.B().FtSearch().Index(projector.ProjectionPersonIDXName).Query(rdq).Dialect(4).Build()
	_, docs, err := q.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	return redis.UnmarshalDocs[projector.PersonProjection](docs)
}

func (q *Queries) getClubProjections(ctx context.Context, ds []domain.ClubID) ([]*projector.ClubProjection, error) {
	if len(ds) == 0 {
		return nil, nil
	}

	var p []*projector.ClubProjection
	keys := make([]string, len(ds))
	for i, d := range ds {
		keys[i] = fmt.Sprintf("%s%s", projector.ProjectionClubPrefix, d)
	}
	cmd := q.rd.B().JsonMget().Key(keys...).Path(".").Build()
	if err := rueidis.DecodeSliceOfJSON(q.rd.Do(ctx, cmd), &p); err != nil {
		return nil, err
	}
	return core.RemoveNils(p), nil
}
