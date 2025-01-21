package projector

import (
	"context"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

const (
	projectionTrainingLookupPrefix = "projection_lookup:trainings:v1:"

	projectionTrainingAccountLookupPrefix        = projectionTrainingLookupPrefix + "accounts:"
	projectionTrainingPersonLookupPrefix         = projectionTrainingLookupPrefix + "persons:"
	projectionTrainingPersonTeamRoleLookupPrefix = projectionTrainingLookupPrefix + "person_team_role:"
)

type trainingAccountLookup struct {
	ID       domain.AccountID `json:"id"`
	FullName string           `json:"full_name"`
}

type trainingPersonLookup struct {
	ID       domain.PersonID `json:"id"`
	FullName string          `json:"full_name"`
}

func (r *rdTrainingProjector) handleAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountCreatedEvent) error {
	return r.insertAccountLookup(ctx, &trainingAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName.Value, e.LastName.Value),
	})
}

func (r *rdTrainingProjector) handleRootAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.RootAccountCreatedEvent) error {
	return r.insertAccountLookup(ctx, &trainingAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName, e.LastName),
	})
}

func (r *rdTrainingProjector) handleRegisteredAccountLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountRegisteredEvent) error {
	return r.insertAccountLookup(ctx, &trainingAccountLookup{
		ID:       domain.AccountID(event.AggregateID()),
		FullName: fmt.Sprintf("%s %s", e.FirstName.Value, e.LastName.Value),
	})
}

func (r *rdTrainingProjector) insertAccountLookup(ctx context.Context, lookup *trainingAccountLookup) error {
	return insertJSON(ctx, r.rd, r.accountLookupKey(lookup.ID), lookup)
}

func (r *rdTrainingProjector) lookupAccount(ctx context.Context, id domain.AccountID) (*trainingAccountLookup, error) {
	var p trainingAccountLookup
	cmd := r.rd.B().JsonGet().Key(r.accountLookupKey(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdTrainingProjector) accountLookupKey(id domain.AccountID) string {
	return fmt.Sprintf("%s%s", projectionTrainingAccountLookupPrefix, id)
}

func (r *rdTrainingProjector) handlePersonLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonCreatedEvent) error {
	id := domain.PersonID(event.AggregateID())
	p := &trainingPersonLookup{
		ID:       id,
		FullName: fmt.Sprintf("%s %s", e.FirstName.Value, e.LastName.Value),
	}
	return insertJSON(ctx, r.rd, r.personLookupKey(id), p)
}

func (r *rdTrainingProjector) lookupPerson(ctx context.Context, id domain.PersonID) (*trainingPersonLookup, error) {
	var p trainingPersonLookup
	cmd := r.rd.B().JsonGet().Key(r.personLookupKey(id)).Path(".").Build()
	return &p, r.rd.Do(ctx, cmd).DecodeJSON(&p)
}

func (r *rdTrainingProjector) personLookupKey(id domain.PersonID) string {
	return fmt.Sprintf("%s%s", projectionTrainingPersonLookupPrefix, id)
}

func (r *rdTrainingProjector) handleTeamMemberLookup(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonInvitedToTeamEvent) error {
	cmd := r.rd.B().Set().Key(r.personTeamRoleLookupKey(e.PersonID, e.TeamID)).Value(e.AssignedRole.Deref()).Build()
	return r.rd.Do(ctx, cmd).Error()
}

func (r *rdTrainingProjector) personTeamRoleLookupKey(personID domain.PersonID, teamID domain.TeamID) string {
	return fmt.Sprintf("%s%s:%s", projectionTrainingPersonTeamRoleLookupPrefix, personID, teamID)
}

func (r *rdTrainingProjector) lookupPersonTeamRole(ctx context.Context, personID domain.PersonID, teamID domain.TeamID) (domain.TeamMemberRole, error) {
	cmd := r.rd.B().Get().Key(r.personTeamRoleLookupKey(personID, teamID)).Build()
	raw, err := r.rd.Do(ctx, cmd).ToString()
	if err != nil {
		return "", err
	}
	return domain.TeamMemberRole(raw), nil
}
