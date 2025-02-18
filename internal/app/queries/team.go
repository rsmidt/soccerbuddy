package queries

import (
	"context"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/projector"
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
	if err := q.authorizer.Authorize(ctx, authz.ActionView, authz.NewTeamResource(domain.TeamID(teamID))); err != nil {
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
			if e.OwningClubID != v.clubID || e.IsShredded() {
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

	if err := q.authorizer.Authorize(ctx, authz.ActionListPersons, authz.NewTeamResource(query.TeamID)); err != nil {
		return nil, err
	}
	clubIDRaw, err := q.es.Lookup(ctx, eventing.LookupOpts{
		AggregateID:   eventing.AggregateID(query.TeamID),
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

type ListTeamMembersTeamMemberView struct {
	ID        domain.TeamMemberID
	PersonID  domain.PersonID
	InviterID *domain.PersonID
	Name      string
	Role      domain.TeamMemberRole
	JoinedAt  time.Time
}

type ListTeamMembersView struct {
	MembersByPersonID map[domain.PersonID]ListTeamMembersTeamMemberView
}

type ListTeamMembersQuery struct {
	TeamID domain.TeamID
}

func (q *Queries) ListTeamMembers(ctx context.Context, query *ListTeamMembersQuery) (*ListTeamMembersView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.ListTeamMembers")
	defer span.End()

	if err := q.authorizer.Authorize(ctx, authz.ActionListPersons, authz.NewTeamResource(query.TeamID)); err != nil {
		return nil, err
	}

	var a []*projector.TeamMemberProjection
	cmd := q.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projector.ProjectionTeamPrefix, query.TeamID)).Path("$.members.*").Build()
	if err := q.rd.Do(ctx, cmd).DecodeJSON(&a); err != nil {
		return nil, err
	}
	membersByPersonID := make(map[domain.PersonID]ListTeamMembersTeamMemberView, len(a))
	for _, projection := range a {
		membersByPersonID[projection.PersonID] = ListTeamMembersTeamMemberView{
			ID:       projection.ID,
			PersonID: projection.PersonID,
			Name:     projection.Name,
			Role:     projection.Role,
			JoinedAt: projection.JoinedAt,
		}
	}

	return &ListTeamMembersView{
		MembersByPersonID: membersByPersonID,
	}, nil
}

type GetMyTeamHomeQuery struct {
	TeamID domain.TeamID
}

type MyTeamHomeView struct {
	ID           domain.TeamID
	Name         string
	Trainings    []*MyTeamHomeTrainingView
	OwningClubID domain.ClubID
}

type MyTeamHomeTrainingView struct {
	ID domain.TrainingID

	ScheduledAt     time.Time
	ScheduledAtIANA string
	EndsAt          time.Time
	EndsAtIANA      string

	GatheringPoint         *GatheringPointView
	AcknowledgmentSettings *AcknowledgmentSettingsView
	RatingSettings         RatingSettingsView

	// Nominations will only be set if enough rights are available.
	Nominations *NominationsView

	Description *string
	Location    *string
	FieldType   *string

	ScheduledBy operatorView
}

type GatheringPointView struct {
	Location        string
	GatherUntil     time.Time
	GatherUntilIANA string
}

type AcknowledgmentSettingsView struct {
	AcknowledgedUntil     time.Time
	AcknowledgedUntilIANA string
}

type RatingSettingsView struct {
	Policy domain.TrainingRatingPolicy
}

type NominationsView struct {
	Players []*TrainingNominationResponse
	Staff   []*TrainingNominationResponse
}

type TrainingNominationResponse struct {
	PersonID       domain.PersonID
	PersonName     string
	Type           domain.TrainingNominationAcknowledgmentType
	AcknowledgedAt *time.Time
	AcceptedAt     *time.Time
	TentativeAt    *time.Time
	DeclinedAt     *time.Time
	AcknowledgedBy *operatorView
	Reason         *string
	NominatedAt    time.Time
}

func (q *Queries) GetMyTeamHome(ctx context.Context, query *GetMyTeamHomeQuery) (*MyTeamHomeView, error) {
	ctx, span := tracing.Tracer.Start(ctx, "queries.GetMyTeamHome")
	defer span.End()

	permissions, err := q.authorizer.Permissions(ctx, authz.NewTeamResource(query.TeamID))
	if err != nil {
		return nil, err
	}
	if !permissions.Allows(authz.ActionView) {
		return nil, authz.ErrUnauthorized
	}

	p, err := q.getTeamProjection(ctx, query.TeamID)
	if err != nil {
		return nil, err
	}

	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrUnauthenticated
	}
	me, err := q.getMe(ctx, principal.AccountID)
	if err != nil {
		return nil, err
	}

	var (
		isCoach      bool
		personInTeam *GetMeLinkedPersonView
	)
outer:
	for _, person := range me.LinkedPersons {
		for _, membership := range person.TeamMemberships {
			if membership.ID != query.TeamID {
				continue
			}
			personInTeam = person
			if membership.Role == domain.TeamMemberRoleCoach {
				// If it's a coach, we can show all trainings regardless of other persons with other roles.
				isCoach = true
				break outer
			}
		}
	}
	if personInTeam == nil {
		// TODO(SOC-29): Handle the case a guest is requesting access to a team.
		return nil, domain.ErrTeamMemberNotFound
	}

	var trainings []*projector.TrainingProjection
	if isCoach {
		trainings, err = q.getTrainingProjectionsByTeamID(ctx, query.TeamID, time.Now())
	} else {
		trainings, err = q.getTrainingProjectionsByTeamIDAndPersonID(ctx, query.TeamID, personInTeam.ID, time.Now())
	}
	if err != nil {
		return nil, err
	}

	ts := make([]*MyTeamHomeTrainingView, len(trainings))
	i := 0
	for _, tp := range trainings {
		var gatheringPoint *GatheringPointView
		if tp.GatheringPoint != nil {
			gatheringPoint = &GatheringPointView{
				Location:        tp.GatheringPoint.Location,
				GatherUntil:     tp.GatheringPoint.GatherUntil,
				GatherUntilIANA: tp.GatheringPoint.GatherUntilIANA,
			}
		}
		var acknowledgmentSettings *AcknowledgmentSettingsView
		if tp.AcknowledgmentSettings != nil {
			acknowledgmentSettings = &AcknowledgmentSettingsView{
				AcknowledgedUntil:     tp.AcknowledgmentSettings.AcknowledgeUntil,
				AcknowledgedUntilIANA: tp.AcknowledgmentSettings.AcknowledgeUntilIANA,
			}
		}
		var nominations *NominationsView
		if permissions.Allows(authz.ActionEdit) {
			var playerResponses []*TrainingNominationResponse
			var staffResponses []*TrainingNominationResponse
			maybeMapOperator := func(operator *projector.OperatorProjection) *operatorView {
				if operator == nil {
					return nil
				}
				return &operatorView{
					FullName: operator.ActorFullName,
				}
			}
			for _, np := range tp.NominatedPlayers {
				playerResponses = append(playerResponses, &TrainingNominationResponse{
					PersonID:       np.ID,
					PersonName:     np.Name,
					Type:           np.Acknowledgment.Type,
					AcknowledgedAt: np.Acknowledgment.AcknowledgedAt,
					AcceptedAt:     np.Acknowledgment.AcceptedAt,
					TentativeAt:    np.Acknowledgment.TentativeAt,
					DeclinedAt:     np.Acknowledgment.DeclinedAt,
					Reason:         np.Acknowledgment.Reason,
					AcknowledgedBy: maybeMapOperator(np.Acknowledgment.AcknowledgedBy),
					NominatedAt:    np.NominatedAt,
				})
			}
			for _, ns := range tp.NominatedStaff {
				staffResponses = append(staffResponses, &TrainingNominationResponse{
					PersonID:       ns.ID,
					PersonName:     ns.Name,
					Type:           ns.Acknowledgment.Type,
					AcknowledgedAt: ns.Acknowledgment.AcknowledgedAt,
					AcceptedAt:     ns.Acknowledgment.AcceptedAt,
					TentativeAt:    ns.Acknowledgment.TentativeAt,
					DeclinedAt:     ns.Acknowledgment.DeclinedAt,
					Reason:         ns.Acknowledgment.Reason,
					AcknowledgedBy: maybeMapOperator(ns.Acknowledgment.AcknowledgedBy),
					NominatedAt:    ns.NominatedAt,
				})
			}
			nominations = &NominationsView{
				Players: playerResponses,
				Staff:   staffResponses,
			}
		}
		ts[i] = &MyTeamHomeTrainingView{
			ID:                     tp.ID,
			ScheduledAt:            tp.ScheduledAt,
			ScheduledAtIANA:        tp.ScheduledAtIANA,
			EndsAt:                 tp.EndsAt,
			EndsAtIANA:             tp.EndsAtIANA,
			Description:            tp.Description,
			Location:               tp.Location,
			FieldType:              tp.FieldType,
			GatheringPoint:         gatheringPoint,
			AcknowledgmentSettings: acknowledgmentSettings,
			RatingSettings: RatingSettingsView{
				Policy: tp.RatingSettings.Policy,
			},
			ScheduledBy: operatorView{
				FullName: tp.ScheduledBy.ActorFullName,
			},
			Nominations: nominations,
		}
		i++
	}
	return &MyTeamHomeView{
		ID:           p.ID,
		Name:         p.Name,
		Trainings:    ts,
		OwningClubID: p.OwningClubID,
	}, nil
}
