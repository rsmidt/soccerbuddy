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
	ProjectionAccountName    eventing.ProjectionName = "accounts"
	ProjectionAccountIDXName                         = "projectionAccountV1Idx"
	ProjectionAccountPrefix                          = "projection:accounts:v1:"
)

type AccountProjection struct {
	ID        domain.AccountID `json:"id"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Email     string           `json:"email"`
	CreatedAt time.Time        `json:"created_at"`
	IsRoot    bool             `json:"is_root"`
	// TODO: Decide if we want to make it also a fat projection and include person details directly.
	LinkedPersons AccountLinkedPersonsSet `json:"linked_persons"`
}

type AccountLinkedPersonsSet map[domain.PersonID]*AccountLinkedPersonProjection

type AccountLinkedPersonProjection struct {
	PersonID domain.PersonID     `json:"person_id"`
	LinkedAs domain.AccountLink  `json:"linked_as"`
	LinkedAt time.Time           `json:"linked_at"`
	LinkedBy *OperatorProjection `json:"linked_by"`
}

type rdAccountProjector struct {
	rd rueidis.Client
}

func NewAccountProjector(rd rueidis.Client) eventing.Projector {
	return &rdAccountProjector{rd: rd}
}

func (r *rdAccountProjector) Init(ctx context.Context) error {
	ctx, span := tracing.Tracer.Start(ctx, "projector.Account.Init")
	defer span.End()

	cmd := r.rd.B().
		FtCreate().
		Index(ProjectionAccountIDXName).
		OnJson().
		Prefix(1).
		Prefix(ProjectionAccountPrefix).
		Schema().
		FieldName("$.first_name").As("first_name").Text().
		FieldName("$.last_name").As("last_name").Text().
		FieldName("$.email").As("email").Text().
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

func (r *rdAccountProjector) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.AccountAggregateType).
		Events(domain.AccountCreatedEventType, domain.RootAccountCreatedEventType, domain.AccountLinkedToPersonEventType).Finish().
		MustBuild()
}

func (r *rdAccountProjector) Projection() eventing.ProjectionName {
	return ProjectionAccountName
}

func (r *rdAccountProjector) Project(ctx context.Context, events ...*eventing.JournalEvent) error {
	var err error
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.AccountCreatedEvent:
			err = r.insertAccountCreatedEvent(ctx, event, e)
		case *domain.RootAccountCreatedEvent:
			err = r.insertRootAccountCreatedEvent(ctx, event, e)
		case *domain.AccountLinkedToPersonEvent:
			err = r.insertAccountLinkedToPersonEvent(ctx, event, e)
		}
		if err != nil {
			tracing.RecordError(ctx, err)
			return err
		}
	}
	return nil
}

func (r *rdAccountProjector) getProjection(ctx context.Context, id domain.AccountID) (*AccountProjection, error) {
	var p AccountProjection
	cmd := r.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", ProjectionAccountPrefix, id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdAccountProjector) insertAccountCreatedEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountCreatedEvent) error {
	p := AccountProjection{
		ID:            domain.AccountID(event.AggregateID()),
		FirstName:     e.FirstName.Value,
		LastName:      e.LastName.Value,
		Email:         e.Email.Value,
		IsRoot:        false,
		CreatedAt:     event.InsertedAt(),
		LinkedPersons: AccountLinkedPersonsSet{},
	}
	key := fmt.Sprintf("%s%s", ProjectionAccountPrefix, p.ID)
	return insertJSON(ctx, r.rd, key, &p)
}

func (r *rdAccountProjector) insertRootAccountCreatedEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.RootAccountCreatedEvent) error {
	p := AccountProjection{
		ID:            domain.AccountID(event.AggregateID()),
		FirstName:     e.FirstName,
		LastName:      e.LastName,
		Email:         e.Email,
		IsRoot:        true,
		CreatedAt:     event.InsertedAt(),
		LinkedPersons: AccountLinkedPersonsSet{},
	}
	key := fmt.Sprintf("%s%s", ProjectionAccountPrefix, p.ID)
	return insertJSON(ctx, r.rd, key, &p)
}

func (r *rdAccountProjector) insertAccountLinkedToPersonEvent(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountLinkedToPersonEvent) error {
	p, err := r.getProjection(ctx, domain.AccountID(e.AggregateID()))
	if err != nil {
		return err
	}
	var linkedBy *OperatorProjection
	if e.LinkedBy != nil {
		operator, err := r.getProjection(ctx, e.LinkedBy.ActorID)
		if err != nil {
			return err
		}
		linkedBy = &OperatorProjection{
			ActorID:       e.LinkedBy.ActorID,
			ActorFullName: fmt.Sprintf("%s %s", operator.FirstName, operator.LastName),
			OnBehalfOf:    e.LinkedBy.OnBehalfOf,
		}
	}
	p.LinkedPersons[e.PersonID] = &AccountLinkedPersonProjection{
		PersonID: e.PersonID,
		LinkedAs: e.LinkedAs,
		LinkedAt: event.InsertedAt(),
		LinkedBy: linkedBy,
	}
	key := fmt.Sprintf("%s%s", ProjectionAccountPrefix, p.ID)
	return insertJSON(ctx, r.rd, key, &p)
}
