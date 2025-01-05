package permify

import (
	permify_payload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"context"
	"fmt"
	permify_grpc "github.com/Permify/permify-go/grpc"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type authorizer struct {
	client   *permify_grpc.Client
	tenantID string
	log      *slog.Logger
}

func NewAuthorizer(log *slog.Logger, client *permify_grpc.Client, tenantID string) authz.Authorizer {
	return &authorizer{client: client, tenantID: tenantID, log: log}
}

func (a *authorizer) Authorize(ctx context.Context, action string, resource *authz.Resource) error {
	ctx, span := tracing.Tracer.Start(ctx, "permify.Authorizer.Authorize")
	defer span.End()

	// Extract the authentication principal from the context.
	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return domain.ErrUnauthenticated
	}

	a.log.
		With(slog.String("subject", fmt.Sprintf("%s:%s", authz.ResourceUserName, principal.AccountID))).
		With(slog.String("permission", action)).
		With(slog.String("entity", fmt.Sprintf("%s:%s", resource.Name, resource.ID))).
		Debug("Authorizing")

	// Check if the principal is authorized to perform the action on the resource.
	cr, err := a.client.Permission.Check(ctx, &permify_payload.PermissionCheckRequest{
		TenantId: a.tenantID,
		Metadata: &permify_payload.PermissionCheckRequestMetadata{
			SchemaVersion: "",
			SnapToken:     "",
			Depth:         30,
		},
		Entity: &permify_payload.Entity{
			Type: resource.Name,
			Id:   resource.ID,
		},
		Permission: action,
		Subject: &permify_payload.Subject{
			Type: authz.ResourceUserName,
			Id:   string(principal.AccountID),
		},
	})
	if err != nil {
		tracing.RecordError(ctx, err)
		return authz.ErrUnauthorized
	}
	isAllowed := cr.Can == permify_payload.CheckResult_CHECK_RESULT_ALLOWED
	span.AddEvent("Permissions evaluated", trace.WithAttributes(attribute.Bool("allowed", isAllowed)))
	if isAllowed {
		return nil
	}
	return authz.ErrUnauthorized
}

func (a *authorizer) AuthorizedEntities(ctx context.Context, action, resourceName string) (authz.EntityIDSet, error) {
	ctx, span := tracing.Tracer.Start(ctx, "permify.Authorizer.AuthorizedEntities")
	defer span.End()

	// Extract the authentication principal from the context.
	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrUnauthenticated
	}

	a.log.
		With(slog.String("subject", fmt.Sprintf("%s:%s", authz.ResourceUserName, principal.AccountID))).
		With(slog.String("permission", action)).
		With(slog.String("entity_type", resourceName)).
		Debug("Listing authorized entities")

	// Check if the principal is authorized to perform the action on the resource.
	// TODO: pagination?
	cr, err := a.client.Permission.LookupEntity(ctx, &permify_payload.PermissionLookupEntityRequest{
		TenantId: a.tenantID,
		Metadata: &permify_payload.PermissionLookupEntityRequestMetadata{
			SchemaVersion: "",
			SnapToken:     "",
			Depth:         20,
		},
		EntityType: resourceName,
		Permission: action,
		Subject: &permify_payload.Subject{
			Type: authz.ResourceUserName,
			Id:   string(principal.AccountID),
		},
		PageSize: 100,
	})
	if err != nil {
		tracing.RecordError(ctx, err)
		return nil, authz.ErrUnauthorized
	}
	idSet := make(authz.EntityIDSet, len(cr.EntityIds))
	for _, id := range cr.EntityIds {
		idSet[id] = struct{}{}
	}
	return idSet, nil
}

func (a *authorizer) RequiredActingOperator(ctx context.Context, personID *domain.PersonID) (domain.Operator, error) {
	ctx, span := tracing.Tracer.Start(ctx, "permify.Authorizer.RequiredActingOperator")
	defer span.End()

	// Extract the authentication principal from the context.
	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return domain.Operator{}, domain.ErrUnauthenticated
	}

	// Only the root principal can act without a person ID.
	if personID == nil && principal.Role != domain.PrincipalRoleRoot {
		return domain.Operator{}, domain.ErrMissingSubject
	}

	return a.OptionalActingOperator(ctx, personID)
}

func (a *authorizer) OptionalActingOperator(ctx context.Context, personID *domain.PersonID) (domain.Operator, error) {
	ctx, span := tracing.Tracer.Start(ctx, "permify.Authorizer.OptionalActingOperator")
	defer span.End()

	// Extract the authentication principal from the context.
	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return domain.Operator{}, domain.ErrUnauthenticated
	}

	// If we're not acting on behalf of someone, we do not need to authorize.
	if personID == nil {
		return domain.NewOperator(principal.AccountID, nil), nil
	}

	a.log.
		With(slog.String("subject", fmt.Sprintf("%s:%s", authz.ResourceUserName, principal.AccountID))).
		With(slog.String("permission", authz.RelationUser)).
		With(slog.String("entity", fmt.Sprintf("%s:%s", authz.ResourcePersonName, *personID))).
		Debug("Authorizing operator")

	cr, err := a.client.Permission.Check(ctx, &permify_payload.PermissionCheckRequest{
		TenantId: a.tenantID,
		Metadata: &permify_payload.PermissionCheckRequestMetadata{
			SchemaVersion: "",
			SnapToken:     "",
			Depth:         20,
		},
		Entity: &permify_payload.Entity{
			Type: authz.ResourcePersonName,
			Id:   string(*personID),
		},
		Permission: authz.RelationUser,
		Subject: &permify_payload.Subject{
			Type: authz.ResourceUserName,
			Id:   string(principal.AccountID),
		},
	})
	if err != nil {
		tracing.RecordError(ctx, err)
		return domain.Operator{}, authz.ErrUnauthorized
	}
	isAllowed := cr.Can == permify_payload.CheckResult_CHECK_RESULT_ALLOWED
	if !isAllowed {
		return domain.Operator{}, authz.ErrUnauthorized
	}

	return domain.NewOperator(principal.AccountID, personID), nil
}

func (a *authorizer) Permissions(ctx context.Context, resource *authz.Resource) (authz.PermissionsSet, error) {
	ctx, span := tracing.Tracer.Start(ctx, "permify.Authorizer.Permissions")
	defer span.End()

	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrUnauthenticated
	}

	a.log.
		With(slog.String("subject", fmt.Sprintf("%s:%s", authz.ResourceUserName, principal.AccountID))).
		With(slog.String("permission", authz.RelationUser)).
		Debug("Requesting permissions")

	cr, err := a.client.Permission.SubjectPermission(context.Background(), &permify_payload.PermissionSubjectPermissionRequest{
		TenantId: "t1",
		Metadata: &permify_payload.PermissionSubjectPermissionRequestMetadata{
			SchemaVersion:  "",
			OnlyPermission: false,
			Depth:          30,
		},
		Entity: &permify_payload.Entity{
			Type: resource.Name,
			Id:   resource.ID,
		},
		Subject: &permify_payload.Subject{
			Type: authz.ResourceUserName,
			Id:   string(principal.AccountID),
		},
	})
	if err != nil {
		return nil, err
	}

	permissions := make(authz.PermissionsSet)
	for perm, decision := range cr.Results {
		if decision == permify_payload.CheckResult_CHECK_RESULT_ALLOWED {
			permissions[perm] = struct{}{}
		}
	}
	return permissions, err
}
