package queries

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"time"
)

type ClubView struct {
	ID        eventing.AggregateID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time

	state domain.ClubState
}

func (c *ClubView) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate("club").
		AggregateID(c.ID).
		Events(domain.ClubCreatedEventType).
		Finish().MustBuild()
}

func (c *ClubView) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.ClubCreatedEvent:
			c.state = domain.ClubStateActive
			c.Name = e.Name
			c.Slug = e.Slug
			c.CreatedAt = event.InsertedAt()
		}
	}
}

type ClubByIDQuery struct {
	ID domain.ClubID
}

func (q *Queries) ClubByID(ctx context.Context, query ClubByIDQuery) (*ClubView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.ClubByID")
	defer span.End()

	if err := q.authorizer.Authorize(ctx, authz.ActionView, authz.NewClubResource(string(query.ID))); err != nil {
		return nil, err
	}
	view := &ClubView{ID: eventing.AggregateID(query.ID)}
	err := q.es.View(ctx, view)
	if err != nil {
		return nil, err
	}
	if view.state == domain.ClubStateUnspecified {
		return nil, nil
	}

	return view, nil
}

type ClubBySlugQuery struct {
	Slug string
}

func (q *Queries) ClubBySlug(ctx context.Context, query ClubBySlugQuery) (*ClubView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.ClubBySlug")
	defer span.End()

	clubID, err := q.es.OwnerLookup(ctx, eventing.LookupOpts{
		AggregateType: domain.ClubAggregateType,
		FieldName:     domain.ClubLookupSlug,
		FieldValue:    eventing.LookupFieldValue(query.Slug),
	})
	if err != nil {
		return nil, err
	}
	if err := q.authorizer.Authorize(ctx, authz.ActionView, authz.NewClubResource(clubID.Deref())); err != nil {
		return nil, err
	}
	view := &ClubView{ID: clubID}
	err = q.es.View(ctx, view)
	if err != nil {
		return nil, err
	}
	if view.state == domain.ClubStateUnspecified {
		return nil, nil
	}
	return view, nil
}
