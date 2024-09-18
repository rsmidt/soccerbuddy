package middleware

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"net/http"
	"strings"
)

func NewAuthenticationMiddleware(qs *queries.Queries) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Client is not supported.
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			// Get session ID cookie or authorization header.
			var rawSessionToken string
			tempReq := http.Request{Header: req.Header()}
			cookie, err := tempReq.Cookie("ID")
			if errors.Is(err, http.ErrNoCookie) {
				// If there's no cookie, try the authorization header.
				rawSessionToken = extractFromBearer(req.Header().Get("Authorization"))
				if rawSessionToken == "" {
					return next(ctx, req)
				}
			} else if err != nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, nil)
			} else {
				rawSessionToken = cookie.Value
			}
			if rawSessionToken == "" {
				return next(ctx, req)
			}

			// Get principal by session ID.
			query := queries.PrincipalBySessionTokenQuery{Token: domain.SessionToken(rawSessionToken)}
			principal, err := qs.PrincipalBySessionToken(ctx, query)
			if errors.Is(err, domain.ErrPrincipalNotFound) {
				// Principal not found. Let's ignore this cookie.
				return next(ctx, req)
			} else if err != nil {
				tracing.RecordError(ctx, err)
				return nil, connect.NewError(connect.CodeInternal, nil)
			}
			ctx = domain.NewContextWithPrincipal(ctx, principal)
			return next(ctx, req)
		}
	}
}

func extractFromBearer(header string) string {
	if header == "" {
		return ""
	}
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}
