package queries

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/rsmidt/soccerbuddy/internal/redis"
	"github.com/sourcegraph/conc/iter"
	"log/slog"
	"slices"
	"time"
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

func (q *Queries) getPersonProjections(ctx context.Context, ds []domain.PersonID) ([]*projector.PersonProjection, error) {
	if len(ds) == 0 {
		return nil, nil
	}

	var p []*projector.PersonProjection
	keys := iter.Map(ds, func(id *domain.PersonID) string {
		return fmt.Sprintf("%s%s", projector.ProjectionPersonPrefix, *id)
	})
	cmd := q.rd.B().JsonMget().Key(keys...).Path(".").Build()
	if err := rueidis.DecodeSliceOfJSON(q.rd.Do(ctx, cmd), &p); err != nil {
		return nil, err
	}
	return removeNils(p), nil
}

func (q *Queries) getTeamProjection(ctx context.Context, id domain.TeamID) (*projector.TeamProjection, error) {
	var a projector.TeamProjection
	cmd := q.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projector.ProjectionTeamPrefix, id)).Path(".").Build()
	return &a, q.rd.Do(ctx, cmd).DecodeJSON(&a)
}

func (q *Queries) getTrainingProjectionsByTeamID(ctx context.Context, teamID domain.TeamID, minTime time.Time) ([]*projector.TrainingProjection, error) {
	rdq := fmt.Sprintf("@owning_team_id:{%s} @scheduled_at_ts:[%d +inf]", teamID, minTime.Unix())
	cmd := q.rd.B().FtSearch().Index(projector.ProjectionTrainingIDXName).Query(rdq).Sortby("scheduled_at_ts").Asc().Dialect(4).Build()
	_, docs, err := q.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	return redis.UnmarshalDocs[projector.TrainingProjection](docs)
}

func (q *Queries) getTrainingProjectionsByTeamIDAndPersonID(ctx context.Context, teamID domain.TeamID, personId domain.PersonID, minTime time.Time) ([]*projector.TrainingProjection, error) {
	rdq := fmt.Sprintf("@owning_team_id:{%s} @scheduled_at_ts:[%d +inf] @nominated_person_ids:{%s}", teamID, minTime.Unix(), personId)
	cmd := q.rd.B().FtSearch().Index(projector.ProjectionTrainingIDXName).Query(rdq).Sortby("scheduled_at_ts").Asc().Dialect(4).Build()
	_, docs, err := q.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	return redis.UnmarshalDocs[projector.TrainingProjection](docs)
}

func (q *Queries) getClubProjections(ctx context.Context, ds []domain.ClubID) ([]*projector.ClubProjection, error) {
	if len(ds) == 0 {
		return nil, nil
	}

	var p []*projector.ClubProjection
	keys := iter.Map(ds, func(id *domain.ClubID) string {
		return fmt.Sprintf("%s%s", projector.ProjectionClubPrefix, *id)
	})
	cmd := q.rd.B().JsonMget().Key(keys...).Path(".").Build()
	if err := rueidis.DecodeSliceOfJSON(q.rd.Do(ctx, cmd), &p); err != nil {
		return nil, err
	}
	return removeNils(p), nil
}

func removeNils[T any](s []*T) []*T {
	return slices.DeleteFunc(s, func(t *T) bool {
		return t == nil
	})
}
