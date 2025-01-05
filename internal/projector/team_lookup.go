package projector

import (
	"context"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

const (
	projectionTeamLookupPrefix = "projection_lookup:teams:v1:"

	projectionTeamAccountLookupPrefix = projectionTeamLookupPrefix + "accounts:"
	projectionTeamPersonLookupPrefix  = projectionTeamLookupPrefix + "persons:"
)

type teamAccountLookup struct {
	ID       domain.AccountID `json:"id"`
	FullName string           `json:"full_name"`
}

type teamPersonLookup struct {
	ID       domain.PersonID `json:"id"`
	FullName string          `json:"full_name"`
}

func (r *rdTeamProjector) handleAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountCreatedEvent) error {
	return r.insertAccountLookup(ctx, &teamAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName.Value, e.LastName.Value),
	})
}

func (r *rdTeamProjector) handleRootAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.RootAccountCreatedEvent) error {
	return r.insertAccountLookup(ctx, &teamAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName, e.LastName),
	})
}

func (r *rdTeamProjector) insertAccountLookup(ctx context.Context, lookup *teamAccountLookup) error {
	return insertJSON(ctx, r.rd, r.accountLookupKey(lookup.ID), lookup)
}

func (r *rdTeamProjector) lookupAccount(ctx context.Context, id domain.AccountID) (*teamAccountLookup, error) {
	var p teamAccountLookup
	cmd := r.rd.B().JsonGet().Key(r.accountLookupKey(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdTeamProjector) accountLookupKey(id domain.AccountID) string {
	return fmt.Sprintf("%s%s", projectionTeamAccountLookupPrefix, id)
}

func (r *rdTeamProjector) handlePersonLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonCreatedEvent) error {
	id := domain.PersonID(event.AggregateID())
	p := &teamPersonLookup{
		ID:       id,
		FullName: fmt.Sprintf("%s %s", e.FirstName.Value, e.LastName.Value),
	}
	return insertJSON(ctx, r.rd, r.personLookupKey(id), p)
}

func (r *rdTeamProjector) lookupPerson(ctx context.Context, id domain.PersonID) (*teamPersonLookup, error) {
	var p teamPersonLookup
	cmd := r.rd.B().JsonGet().Key(r.personLookupKey(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdTeamProjector) personLookupKey(id domain.PersonID) string {
	return fmt.Sprintf("%s%s", projectionTeamPersonLookupPrefix, id)
}
