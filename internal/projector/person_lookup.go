package projector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

const (
	projectionPersonLookupPrefix = "projection_lookup:persons:v1:"

	projectionPersonAccountLookupPrefix = projectionPersonLookupPrefix + "accounts:"
	projectionPersonTeamLookupPrefix    = projectionPersonLookupPrefix + "teams:"
	projectionPersonClubLookupPrefix    = projectionPersonLookupPrefix + "clubs:"
)

type personAccountLookup struct {
	ID       domain.AccountID `json:"id"`
	FullName string           `json:"full_name"`
}

type personTeamLookup struct {
	ID           domain.TeamID `json:"id"`
	Name         string        `json:"name"`
	OwningClubID domain.ClubID `json:"owning_club_id"`
}

type personClubLookup struct {
	ID   domain.ClubID `json:"id"`
	Name string        `json:"name"`
}

func (r *rdPersonProjector) handleAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountCreatedEvent) error {
	return r.insertAccountLookup(ctx, &personAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName.Value, e.LastName.Value),
	})
}

func (r *rdPersonProjector) handleRootAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.RootAccountCreatedEvent) error {
	return r.insertAccountLookup(ctx, &personAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName, e.LastName),
	})
}

func (r *rdPersonProjector) handleRegisteredAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountRegisteredEvent) error {
	return r.insertAccountLookup(ctx, &personAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName.Value, e.LastName.Value),
	})
}

func (r *rdPersonProjector) insertAccountLookup(ctx context.Context, lookup *personAccountLookup) error {
	key := fmt.Sprintf("%s%s", projectionPersonAccountLookupPrefix, lookup.ID)
	return insertJSON(ctx, r.rd, key, lookup)
}

func (r *rdPersonProjector) lookupAccount(ctx context.Context, id domain.AccountID) (*personAccountLookup, error) {
	var p personAccountLookup
	cmd := r.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projectionPersonAccountLookupPrefix, id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdPersonProjector) handleTeamLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.TeamCreatedEvent) error {
	lookup := personTeamLookup{
		ID:           domain.TeamID(event.AggregateID()),
		Name:         e.Name,
		OwningClubID: e.OwningClubID,
	}
	key := fmt.Sprintf("%s%s", projectionPersonTeamLookupPrefix, lookup.ID)
	return insertJSON(ctx, r.rd, key, &lookup)
}

func (r *rdPersonProjector) handleClubLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.ClubCreatedEvent) error {
	lookup := personClubLookup{
		ID:   domain.ClubID(e.AggregateID()),
		Name: e.Name,
	}
	key := fmt.Sprintf("%s%s", projectionPersonClubLookupPrefix, lookup.ID)
	return insertJSON(ctx, r.rd, key, &lookup)
}

func (r *rdPersonProjector) lookupTeam(ctx context.Context, id domain.TeamID) (*personTeamLookup, error) {
	var p personTeamLookup
	cmd := r.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projectionPersonTeamLookupPrefix, id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdPersonProjector) lookupClub(ctx context.Context, id domain.ClubID) (*personClubLookup, error) {
	var p personClubLookup
	cmd := r.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projectionPersonClubLookupPrefix, id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func insertJSON(ctx context.Context, client rueidis.Client, key string, value any) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	cmd := client.B().JsonSet().Key(key).Path(".").Value(string(val)).Build()
	res := client.Do(ctx, cmd)
	return res.Error()
}
