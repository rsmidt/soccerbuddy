package projector

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"time"
)

const (
	ProjectionClubName    eventing.ProjectionName = "clubs"
	ProjectionClubIDXName                         = "projectionClubV1Idx"
	ProjectionClubPrefix                          = "projection:clubs:v1:"
)

type ClubProjection struct {
	ID        domain.ClubID `json:"id"`
	Name      string        `json:"name"`
	Slug      string        `json:"slug"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type rdClubProjector struct {
	rd rueidis.Client
}

func NewClubProjector(rd rueidis.Client) eventing.Projector {
	return &rdClubProjector{rd: rd}
}

func (r *rdClubProjector) Init(ctx context.Context) error {
	ctx, span := tracing.Tracer.Start(ctx, "projector.Club.Init")
	defer span.End()

	cmd := r.rd.B().
		FtCreate().
		Index(ProjectionClubIDXName).
		OnJson().
		Prefix(1).
		Prefix(ProjectionClubPrefix).
		Schema().
		FieldName("$.name").As("name").Text().Nostem().
		FieldName("$.slug").As("slug").Tag().
		Build()
	if err := r.rd.Do(ctx, cmd).Error(); err != nil {
		rderr, ok := rueidis.IsRedisErr(err)
		if ok && rderr.Error() == "Index already exists" {
			return nil
		}
		return err
	}
	return nil
}

func (r *rdClubProjector) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.ClubAggregateType).
		Events(domain.ClubCreatedEventType).Finish().
		MustBuild()
}

func (r *rdClubProjector) Projection() eventing.ProjectionName {
	return ProjectionClubName
}

func (r *rdClubProjector) Project(ctx context.Context, events ...*eventing.JournalEvent) error {
	var err error
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.ClubCreatedEvent:
			err = r.insertClubCreatedEvent(ctx, event, e)
		}
		if err != nil {
			tracing.RecordError(ctx, err)
			return err
		}
	}
	return nil
}

func (r *rdClubProjector) getProjection(ctx context.Context, id domain.ClubID) (*ClubProjection, error) {
	var p ClubProjection
	cmd := r.rd.B().JsonGet().Key(r.key(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdClubProjector) key(id domain.ClubID) string {
	return fmt.Sprintf("%s%s", ProjectionClubPrefix, id)
}

func (r *rdClubProjector) insertClubCreatedEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.ClubCreatedEvent) error {
	p := &ClubProjection{
		ID:        domain.ClubID(e.AggregateID()),
		Name:      e.Name,
		Slug:      e.Slug,
		CreatedAt: event.InsertedAt(),
		UpdatedAt: event.InsertedAt(),
	}
	return insertJSON(ctx, r.rd, r.key(p.ID), p)
}
