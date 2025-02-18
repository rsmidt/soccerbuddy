package queries

import (
	"context"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/rsmidt/soccerbuddy/internal/redis"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"golang.org/x/exp/maps"
	"strings"
	"time"
)

type ListTeamsView struct {
	Teams []ListTeamsTeamView
}

type ListTeamsTeamView struct {
	ID           domain.TeamID
	Name         string
	Slug         string
	OwningClubID domain.ClubID
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	idFilter := strings.Join(maps.Keys(idSet), "|")

	cmd := q.rd.B().FtSearch().Index(projector.ProjectionTeamIDXName).
		Query(fmt.Sprintf("@id:{%s} @owning_club_id:{%s}", idFilter, query.OwningClubID)).
		Return("4").Identifier("$.id").Identifier("name").Identifier("owning_club_id").Identifier("slug").
		Dialect(4).
		Build()
	_, docs, err := q.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, fmt.Errorf("failed to query team projection: %w", err)
	}
	view := ListTeamsView{
		Teams: make([]ListTeamsTeamView, len(docs)),
	}
	for i, doc := range docs {
		view.Teams[i] = ListTeamsTeamView{
			ID:           domain.TeamID(doc.Doc["$.id"]),
			Name:         redis.FlattenToString(doc.Doc["name"]),
			OwningClubID: domain.ClubID(redis.FlattenToString(doc.Doc["owning_club_id"])),
			CreatedAt:    redis.FlattenToTime(doc.Doc["created_at"]),
			UpdatedAt:    redis.FlattenToTime(doc.Doc["updated_at"]),
		}
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
	Persons []PersonsNotInTeamViewPerson
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
	// Fetch the club ID just to double-check that we do not access persons of a different club.
	clubIDRaw, err := q.es.Lookup(ctx, eventing.LookupOpts{
		AggregateID:   eventing.AggregateID(query.TeamID),
		AggregateType: domain.TeamAggregateType,
		FieldName:     domain.TeamLookupOwningClub,
	})
	if err != nil {
		return nil, err
	}
	cmd := q.rd.B().FtSearch().Index(projector.ProjectionPersonIDXName).
		Query(fmt.Sprintf("-@team_id:{%s} @owning_club_id:{%s} (@first_name:%%%%%s%%%% | @last_name: %%%%%s%%%%)", query.TeamID, *clubIDRaw, query.Query, query.Query)).
		Return("3").Identifier("$.id").Identifier("first_name").Identifier("last_name").
		Dialect(4).
		Build()
	_, docs, err := q.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	v := make([]PersonsNotInTeamViewPerson, len(docs))
	for i, doc := range docs {
		v[i] = PersonsNotInTeamViewPerson{
			ID:        domain.PersonID(doc.Doc["$.id"]),
			FirstName: redis.FlattenToString(doc.Doc["first_name"]),
			LastName:  redis.FlattenToString(doc.Doc["last_name"]),
		}
	}
	return &PersonsNotInTeamView{
		Persons: v,
	}, nil
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
