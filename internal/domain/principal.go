package domain

import (
	"context"
	"errors"
)

var (
	ErrPrincipalNotFound = errors.New("principal not found")
	ErrUnauthenticated   = errors.New("unauthenticated")
	ErrMissingSubject    = errors.New("the subject of the action was not specified")
)

// PrincipalRole specifies the type of principal that is accessing the system.
type PrincipalRole int

const (
	PrincipalRoleUnspecified PrincipalRole = iota

	// PrincipalRoleRoot is allowed to skip subject verification.
	PrincipalRoleRoot

	// PrincipalRoleRegular will always require subject verification is necessary.
	PrincipalRoleRegular
)

// Principal is the authenticated principal of a request.
type Principal struct {
	AccountID    AccountID
	SessionToken SessionToken
	Role         PrincipalRole
}

func NewPrincipal(accountID AccountID, sessionToken SessionToken, role PrincipalRole) *Principal {
	return &Principal{
		AccountID:    accountID,
		SessionToken: sessionToken,
		Role:         role,
	}
}

type ctxKey int

const principalCtxKey ctxKey = iota

// PrincipalFromContext extracts the Principal from the context.
func PrincipalFromContext(ctx context.Context) (*Principal, bool) {
	p, ok := ctx.Value(principalCtxKey).(*Principal)
	return p, ok
}

// NewContextWithPrincipal adds the Principal to the context.
func NewContextWithPrincipal(ctx context.Context, p *Principal) context.Context {
	return context.WithValue(ctx, principalCtxKey, p)
}
