package queries

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"strings"
	"time"
)

type ListTeamsView struct {
	owningClubID domain.ClubID
	filterIds    map[string]struct{}
	TeamsById    map[domain.TeamID]struct {
		ID        domain.TeamID
		Name      string
		Slug      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
}

func (v *ListTeamsView) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.TeamAggregateType).
		Events(domain.TeamCreatedEventType, domain.TeamDeletedEventType).
		Finish().MustBuild()
}

func (v *ListTeamsView) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.TeamCreatedEvent:
			if e.OwningClubID != v.owningClubID {
				continue
			}
			if _, ok := v.filterIds[e.AggregateID().Deref()]; !ok {
				continue
			}
			teamID := domain.TeamID(e.AggregateID())
			v.TeamsById[teamID] = struct {
				ID        domain.TeamID
				Name      string
				Slug      string
				CreatedAt time.Time
				UpdatedAt time.Time
			}{
				ID:        teamID,
				Name:      e.Name,
				Slug:      e.Slug,
				CreatedAt: event.InsertedAt(),
				UpdatedAt: event.InsertedAt(),
			}
		case *domain.TeamDeletedEvent:
			delete(v.TeamsById, domain.TeamID(event.AggregateID()))
		}
	}
}

type ListTeamsQuery struct {
	OwningClubID domain.ClubID
}

func (q *Queries) ListTeams(ctx context.Context, query ListTeamsQuery) (*ListTeamsView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.ListTeams")
	defer span.End()

	idSet, err := q.authorizer.AuthorizedEntities(ctx, authz.ActionView, authz.ResourceTeamName)
	if err != nil {
		return nil, err
	}

	view := ListTeamsView{
		owningClubID: query.OwningClubID,
		filterIds:    idSet,
		TeamsById: make(map[domain.TeamID]struct {
			ID        domain.TeamID
			Name      string
			Slug      string
			CreatedAt time.Time
			UpdatedAt time.Time
		}),
	}
	err = q.es.View(ctx, &view)
	if err != nil {
		return nil, err
	}
	return &view, nil
}

type TeamOverviewView struct {
	Team struct {
		ID           domain.TeamID
		Name         string
		Slug         string
		OwningClubID domain.ClubID
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
}

type GetTeamOverviewQuery struct {
	TeamSlug string
}

func (q *Queries) GetTeamOverview(ctx context.Context, query GetTeamOverviewQuery) (*TeamOverviewView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.GetTeamOverview")
	defer span.End()

	teamID, err := q.es.OwnerLookup(ctx, eventing.LookupOpts{
		AggregateType: domain.TeamAggregateType,
		FieldName:     domain.TeamLookupSlug,
		FieldValue:    eventing.LookupFieldValue(query.TeamSlug),
	})
	if err != nil {
		return nil, err
	}
	if err := q.authorizer.Authorize(ctx, authz.ActionView, authz.NewTeamResource(teamID.Deref())); err != nil {
		return nil, err
	}
	t, err := q.repos.Team().FindByID(ctx, domain.TeamID(teamID))
	if err != nil {
		return nil, err
	}
	if t.State != domain.TeamStateActive {
		return nil, domain.ErrTeamNotFound
	}
	overview := TeamOverviewView{
		Team: struct {
			ID           domain.TeamID
			Name         string
			Slug         string
			OwningClubID domain.ClubID
			CreatedAt    time.Time
			UpdatedAt    time.Time
		}{
			ID:           t.ID,
			Name:         t.Name,
			Slug:         t.Slug,
			OwningClubID: t.OwningClubID,
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
		},
	}
	return &overview, nil
}

type PersonsNotInTeamViewPerson struct {
	ID        domain.PersonID
	FirstName string
	LastName  string
}

type PersonsNotInTeamView struct {
	teamID domain.TeamID
	clubID domain.ClubID
	query  string

	personsToRemove map[domain.PersonID]struct{}
	Persons         map[domain.PersonID]PersonsNotInTeamViewPerson
}

func (v *PersonsNotInTeamView) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.PersonAggregateType).
		Events(domain.PersonCreatedEventType).
		Finish().
		WithAggregate(domain.TeamMemberAggregateType).
		Events(domain.PersonInvitedToTeamEventType).
		Finish().MustBuild()
}

func (v *PersonsNotInTeamView) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.PersonCreatedEvent:
			if e.OwningClubID != v.clubID || e.FirstName.IsShredded {
				continue
			}
			if !strings.Contains(strings.ToLower(e.FirstName.Value), strings.ToLower(v.query)) &&
				!strings.Contains(strings.ToLower(e.LastName.Value), strings.ToLower(v.query)) {
				continue
			}
			v.Persons[domain.PersonID(e.AggregateID())] = PersonsNotInTeamViewPerson{
				ID:        domain.PersonID(e.AggregateID()),
				FirstName: e.FirstName.Value,
				LastName:  e.LastName.Value,
			}
		case *domain.PersonInvitedToTeamEvent:
			if e.TeamID == v.teamID {
				v.personsToRemove[e.PersonID] = struct{}{}
			}
		}
	}

	// Because we can't now if after an invite event there will be another person removed event,
	// we cannot directly delete, but only afterwards.
	for id := range v.personsToRemove {
		delete(v.Persons, id)
	}
}

type SearchPersonsNotInTeamQuery struct {
	TeamID domain.TeamID
	Query  string
}

func (q *Queries) SearchPersonsNotInTeam(ctx context.Context, query SearchPersonsNotInTeamQuery) (*PersonsNotInTeamView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.SearchPersonsNotInTeam")
	defer span.End()

	if err := q.authorizer.Authorize(ctx, authz.ActionListPersons, authz.NewTeamResource(string(query.TeamID))); err != nil {
		return nil, err
	}
	clubIDRaw, err := q.es.Lookup(ctx, eventing.LookupOpts{
		AggregateType: domain.TeamAggregateType,
		FieldName:     domain.TeamLookupOwningClub,
	})
	if err != nil {
		return nil, err
	}
	clubID := domain.ClubID(*clubIDRaw)
	view := PersonsNotInTeamView{
		teamID:          query.TeamID,
		clubID:          clubID,
		query:           query.Query,
		personsToRemove: make(map[domain.PersonID]struct{}),
		Persons:         make(map[domain.PersonID]PersonsNotInTeamViewPerson),
	}
	err = q.es.View(ctx, &view)
	if err != nil {
		return nil, err
	}
	return &view, nil
}

type teamMembership struct {
	ID        domain.TeamMemberID
	PersonID  domain.PersonID
	InvitedBy domain.Operator
	JoinedAt  time.Time
	Role      domain.TeamMemberRoleRole
}

type listTeamMembershipsView struct {
	teamID domain.TeamID

	MembersByPersonID map[domain.PersonID]teamMembership
}

func (v *listTeamMembershipsView) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.TeamMemberAggregateType).
		Events(domain.PersonInvitedToTeamEventType).
		Finish().MustBuild()
}

func (v *listTeamMembershipsView) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.PersonInvitedToTeamEvent:
			if e.TeamID != v.teamID {
				continue
			}
			v.MembersByPersonID[e.PersonID] = teamMembership{
				ID:        domain.TeamMemberID(e.AggregateID()),
				PersonID:  e.PersonID,
				InvitedBy: e.InvitedBy,
				JoinedAt:  event.InsertedAt(),
				Role:      e.AssignedRole,
			}
		}
	}
}

type teamMember struct {
	ID        domain.TeamMemberID
	PersonID  domain.PersonID
	InviterID *domain.PersonID
	FirstName string
	LastName  string
	Role      domain.TeamMemberRoleRole
	JoinedAt  time.Time
}

type ListTeamMembersView struct {
	memberships map[domain.PersonID]teamMembership

	MembersByPersonID map[domain.PersonID]teamMember
}

func (v *ListTeamMembersView) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.PersonAggregateType).
		Events(domain.PersonCreatedEventType).
		Finish().MustBuild()
}

func (v *ListTeamMembersView) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.PersonCreatedEvent:
			personID := domain.PersonID(e.AggregateID())
			m, ok := v.memberships[personID]
			if !ok {
				continue
			}
			v.MembersByPersonID[personID] = teamMember{
				ID:        m.ID,
				PersonID:  personID,
				FirstName: e.FirstName.Value,
				LastName:  e.LastName.Value,
				JoinedAt:  m.JoinedAt,
				Role:      m.Role,
			}
		}
	}
}

type ListTeamMembersQuery struct {
	TeamID domain.TeamID
}

func (q *Queries) ListTeamMembers(ctx context.Context, query ListTeamMembersQuery) (*ListTeamMembersView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.ListTeamMembers")
	defer span.End()

	if err := q.authorizer.Authorize(ctx, authz.ActionListPersons, authz.NewTeamResource(string(query.TeamID))); err != nil {
		return nil, err
	}
	memberships := listTeamMembershipsView{
		teamID:            query.TeamID,
		MembersByPersonID: make(map[domain.PersonID]teamMembership),
	}
	err := q.es.View(ctx, &memberships)
	if err != nil {
		return nil, err
	}
	view := ListTeamMembersView{
		memberships:       memberships.MembersByPersonID,
		MembersByPersonID: make(map[domain.PersonID]teamMember),
	}
	err = q.es.View(ctx, &view)
	if err != nil {
		return nil, err
	}
	return &view, nil
}
