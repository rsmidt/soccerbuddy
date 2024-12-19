package projector

import (
	"context"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

const (
	projectionTeamHomeLookupPrefix = "projection_lookup:team_homes:v1:"

	projectionTeamHomeAccountLookupPrefix = projectionTeamHomeLookupPrefix + "accounts:"
)

type teamHomeAccountLookup struct {
	ID       domain.AccountID `json:"id"`
	FullName string           `json:"full_name"`
}

func (r *rdTeamHomeProjector) handleAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountCreatedEvent) error {
	return r.insertAccountLookup(ctx, &teamHomeAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName.Value, e.LastName.Value),
	})
}

func (r *rdTeamHomeProjector) handleRootAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.RootAccountCreatedEvent) error {
	return r.insertAccountLookup(ctx, &teamHomeAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName, e.LastName),
	})
}
func (r *rdTeamHomeProjector) insertAccountLookup(ctx context.Context, lookup *teamHomeAccountLookup) error {
	return insertJSON(ctx, r.rd, r.accountLookupKey(lookup.ID), lookup)
}

func (r *rdTeamHomeProjector) lookupAccount(ctx context.Context, id domain.AccountID) (*teamHomeAccountLookup, error) {
	var p teamHomeAccountLookup
	cmd := r.rd.B().JsonGet().Key(r.accountLookupKey(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdTeamHomeProjector) accountLookupKey(id domain.AccountID) string {
	return fmt.Sprintf("%s%s", projectionTeamHomeAccountLookupPrefix, id)
}
