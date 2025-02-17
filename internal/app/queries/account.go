package queries

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/app/view"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

func (q *Queries) GetMe(ctx context.Context) (*view.GetMe, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.GetMe")
	defer span.End()

	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrUnauthenticated
	}

	return q.vs.Account().GetMe(ctx, principal.AccountID)
}
