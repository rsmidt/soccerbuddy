package permify

import (
	permify_payload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"context"
	"errors"
	permify_grpc "github.com/Permify/permify-go/grpc"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"log/slog"
)

type relationStore struct {
	client *permify_grpc.Client
	log    *slog.Logger
}

func NewRelationStore(log *slog.Logger, client *permify_grpc.Client) authz.RelationStore {
	return &relationStore{client: client, log: log}
}

func (r *relationStore) AddRelations(ctx context.Context, relations []authz.Relation) error {
	ctx, span := tracing.Tracer.Start(ctx, "permify.RelationStore.AddRelations")
	defer span.End()

	grpcRelations := make([]*permify_payload.Tuple, len(relations))
	for i, relation := range relations {
		r.log.Debug("Adding permify relation", slog.String("relation", relation.String()))

		grpcRelations[i] = &permify_payload.Tuple{
			Entity: &permify_payload.Entity{
				Type: relation.EntityType,
				Id:   relation.EntityID,
			},
			Relation: relation.Relation,
			Subject: &permify_payload.Subject{
				Type: relation.SubjectType,
				Id:   relation.SubjectID,
			},
		}
	}

	_, err := r.client.Data.Write(ctx, &permify_payload.DataWriteRequest{
		TenantId: "t1",
		Metadata: &permify_payload.DataWriteRequestMetadata{
			SchemaVersion: "",
		},
		Tuples:     grpcRelations,
		Attributes: nil,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *relationStore) RemoveRelations(ctx context.Context, relations []authz.Relation) error {
	ctx, span := tracing.Tracer.Start(ctx, "permify.RelationStore.RemoveRelations")
	defer span.End()

	var allErr error
	for _, relation := range relations {
		r.log.Debug("Removing permify relation", slog.String("relation", relation.String()))

		_, err := r.client.Data.Delete(ctx, &permify_payload.DataDeleteRequest{
			TenantId: "t1",
			TupleFilter: &permify_payload.TupleFilter{
				Entity: &permify_payload.EntityFilter{
					Type: relation.EntityType,
					Ids:  []string{relation.EntityID},
				},
				Relation: relation.Relation,
				Subject: &permify_payload.SubjectFilter{
					Type: relation.SubjectType,
					Ids:  []string{relation.SubjectID},
				},
			},
			AttributeFilter: &permify_payload.AttributeFilter{},
		})
		if err != nil {
			allErr = errors.Join(allErr, err)
		}
	}
	return allErr
}
