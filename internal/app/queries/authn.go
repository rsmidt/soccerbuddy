package queries

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type PrincipalBySessionTokenQuery struct {
	Token domain.SessionToken
}

// PrincipalBySessionToken constructs the current authentication principal given a session ID.
func (q *Queries) PrincipalBySessionToken(ctx context.Context, query PrincipalBySessionTokenQuery) (*domain.Principal, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.PrincipalBySessionToken")
	defer span.End()

	session, err := q.repos.Session().FindByToken(ctx, query.Token)
	if errors.Is(err, domain.ErrSessionNotFound) {
		return nil, domain.ErrPrincipalNotFound
	} else if err != nil {
		return nil, err
	} else if session.State != domain.SessionStateActive {
		return nil, domain.ErrPrincipalNotFound
	}
	return domain.NewPrincipal(session.AccountID, session.Token, session.Role), nil
}
