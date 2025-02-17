package grpc

import (
	"connectrpc.com/connect"
	"context"
	teamv1 "github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/team/v1"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/team/v1/teamv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
	"github.com/rsmidt/soccerbuddy/internal/app/view"
	"github.com/rsmidt/soccerbuddy/internal/core"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

type teamServer struct {
	*baseHandler
}

func newTeamServiceHandler(base *baseHandler) teamv1connect.TeamServiceHandler {
	return &teamServer{baseHandler: base}
}

func (t *teamServer) CreateTeam(ctx context.Context, req *connect.Request[teamv1.CreateTeamRequest]) (*connect.Response[teamv1.CreateTeamResponse], error) {
	cmd := commands.CreateTeamCommand{
		Name:         req.Msg.Name,
		OwningClubID: domain.ClubID(req.Msg.OwningClubId),
	}
	newTeam, err := t.cmds.CreateTeam(ctx, cmd)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	return connect.NewResponse(&teamv1.CreateTeamResponse{
		Id:           string(newTeam.ID),
		Name:         newTeam.Name,
		Slug:         newTeam.Slug,
		OwningClubId: string(newTeam.OwningClubID),
		CreatedAt:    timestamppb.New(newTeam.CreatedAt),
		UpdatedAt:    timestamppb.New(newTeam.UpdatedAt),
	}), nil
}

func (t *teamServer) ListTeams(ctx context.Context, req *connect.Request[teamv1.ListTeamsRequest]) (*connect.Response[teamv1.ListTeamsResponse], error) {
	query := queries.ListTeamsQuery{
		OwningClubID: domain.ClubID(req.Msg.OwningClubId),
	}
	teams, err := t.qs.ListTeams(ctx, query)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	var respTeams []*teamv1.ListTeamsResponse_Team
	for _, team := range teams.TeamsById {
		respTeams = append(respTeams, &teamv1.ListTeamsResponse_Team{
			Id:        string(team.ID),
			Name:      team.Name,
			Slug:      team.Slug,
			CreatedAt: timestamppb.New(team.CreatedAt),
			UpdatedAt: timestamppb.New(team.UpdatedAt),
		})
	}
	return connect.NewResponse(&teamv1.ListTeamsResponse{
		Teams: respTeams,
	}), nil
}

func (t *teamServer) GetTeamOverview(ctx context.Context, c *connect.Request[teamv1.GetTeamOverviewRequest]) (*connect.Response[teamv1.GetTeamOverviewResponse], error) {
	query := queries.GetTeamOverviewQuery{
		TeamSlug: c.Msg.TeamSlug,
	}
	teamOverview, err := t.qs.GetTeamOverview(ctx, query)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	return connect.NewResponse(&teamv1.GetTeamOverviewResponse{
		Id:           string(teamOverview.Team.ID),
		Name:         teamOverview.Team.Name,
		Slug:         teamOverview.Team.Slug,
		OwningClubId: string(teamOverview.Team.OwningClubID),
		CreatedAt:    timestamppb.New(teamOverview.Team.CreatedAt),
		UpdatedAt:    timestamppb.New(teamOverview.Team.UpdatedAt),
	}), nil
}

func (t *teamServer) DeleteTeam(ctx context.Context, c *connect.Request[teamv1.DeleteTeamRequest]) (*connect.Response[teamv1.DeleteTeamResponse], error) {
	cmd := commands.DeleteTeamCommand{
		ID: domain.TeamID(c.Msg.TeamId),
	}
	err := t.cmds.DeleteTeam(ctx, cmd)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	return connect.NewResponse(&teamv1.DeleteTeamResponse{}), nil
}

func (t *teamServer) SearchPersonsNotInTeam(ctx context.Context, c *connect.Request[teamv1.SearchPersonsNotInTeamRequest]) (*connect.Response[teamv1.SearchPersonsNotInTeamResponse], error) {
	query := queries.SearchPersonsNotInTeamQuery{
		TeamID: domain.TeamID(c.Msg.TeamId),
		Query:  c.Msg.Query,
	}
	persons, err := t.qs.SearchPersonsNotInTeam(ctx, query)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	var respPersons []*teamv1.SearchPersonsNotInTeamResponse_Person
	for _, person := range persons.Persons {
		respPersons = append(respPersons, &teamv1.SearchPersonsNotInTeamResponse_Person{
			Id:        string(person.ID),
			FirstName: person.FirstName,
			LastName:  person.LastName,
		})
	}
	return connect.NewResponse(&teamv1.SearchPersonsNotInTeamResponse{
		Persons: respPersons,
	}), nil
}

func (t *teamServer) AddPersonToTeam(ctx context.Context, c *connect.Request[teamv1.AddPersonToTeamRequest]) (*connect.Response[teamv1.AddPersonToTeamResponse], error) {
	cmd := commands.AddPersonToTeamCommand{
		TeamID:   domain.TeamID(c.Msg.TeamId),
		PersonID: domain.PersonID(c.Msg.PersonId),
		Role:     domain.TeamMemberRole(c.Msg.Role),
	}
	err := t.cmds.AddPersonToTeam(ctx, cmd)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	return connect.NewResponse(&teamv1.AddPersonToTeamResponse{}), nil
}

func (t *teamServer) ListTeamMembers(ctx context.Context, c *connect.Request[teamv1.ListTeamMembersRequest]) (*connect.Response[teamv1.ListTeamMembersResponse], error) {
	query := queries.ListTeamMembersQuery{
		TeamID: domain.TeamID(c.Msg.TeamId),
	}
	members, err := t.qs.ListTeamMembers(ctx, &query)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	var respMembers []*teamv1.ListTeamMembersResponse_Member
	for _, member := range members.MembersByPersonID {
		var inviterID *string
		if member.InviterID != nil {
			inviterID = core.PTR(string(*member.InviterID))
		}
		nameParts := strings.Split(member.Name, " ")
		firstName := strings.Join(nameParts[:len(nameParts)-1], " ")
		lastName := nameParts[len(nameParts)-1]

		respMembers = append(respMembers, &teamv1.ListTeamMembersResponse_Member{
			Id:        string(member.ID),
			PersonId:  string(member.PersonID),
			InviterId: inviterID,
			FirstName: firstName,
			LastName:  lastName,
			JoinedAt:  timestamppb.New(member.JoinedAt),
			Role:      string(member.Role),
		})
	}
	return connect.NewResponse(&teamv1.ListTeamMembersResponse{
		Members: respMembers,
	}), nil
}

func (t *teamServer) ScheduleTraining(ctx context.Context, c *connect.Request[teamv1.ScheduleTrainingRequest]) (*connect.Response[teamv1.ScheduleTrainingResponse], error) {
	cmd := commands.ScheduleTrainingCommand{
		ScheduledAt:            pbToLocalTime(c.Msg.ScheduledAt, defaultLocation),
		ScheduledAtIANA:        defaultLocation.String(),
		EndsAt:                 pbToLocalTime(c.Msg.EndsAt, defaultLocation),
		EndsAtIANA:             defaultLocation.String(),
		Description:            c.Msg.Description,
		Location:               c.Msg.Location,
		FieldType:              c.Msg.FieldType,
		GatheringPoint:         pbToGatheringPoint(c.Msg.GatheringPoint),
		AcknowledgmentSettings: pbToAcknowledgementSettings(c.Msg.AcknowledgmentSettings),
		RatingSettings:         pbToRatingSettings(c.Msg.RatingSettings),
		TeamID:                 domain.TeamID(c.Msg.TeamId),
		Nominations:            pbToNominations(c.Msg.Nominations),
	}
	if err := t.cmds.ScheduleTraining(ctx, &cmd); err != nil {
		return nil, t.handleCommonErrors(err)
	}
	// TODO: Return from projection?
	return connect.NewResponse(&teamv1.ScheduleTrainingResponse{}), nil
}

func (t *teamServer) GetMyTeamHome(ctx context.Context, c *connect.Request[teamv1.GetMyTeamHomeRequest]) (*connect.Response[teamv1.GetMyTeamHomeResponse], error) {
	query := queries.GetMyTeamHomeQuery{
		TeamID: domain.TeamID(c.Msg.TeamId),
	}
	th, err := t.qs.GetMyTeamHome(ctx, &query)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	ts := make([]*teamv1.GetMyTeamHomeResponse_Training, len(th.Trainings))
	for i, training := range th.Trainings {
		ts[i] = &teamv1.GetMyTeamHomeResponse_Training{
			Id:                     string(training.ID),
			ScheduledAt:            localTimeToPb(&training.ScheduledAt),
			EndsAt:                 localTimeToPb(&training.EndsAt),
			Location:               training.Location,
			FieldType:              training.FieldType,
			Description:            training.Description,
			GatheringPoint:         gatheringPointToPb(training.GatheringPoint),
			AcknowledgmentSettings: acknowledgmentSettingsToPb(training.AcknowledgmentSettings),
			RatingSettings:         ratingSettingsToPb(training.RatingSettings),
			Nominations:            nominationResponsesToPb(training.Nominations),
		}
	}
	return connect.NewResponse(&teamv1.GetMyTeamHomeResponse{
		TeamId:    string(th.ID),
		TeamName:  th.Name,
		Trainings: ts,
	}), nil
}

func (t *teamServer) NominatePersonsForTraining(ctx context.Context, c *connect.Request[teamv1.NominatePersonsForTrainingRequest]) (*connect.Response[teamv1.NominatePersonsForTrainingResponse], error) {
	playerIDs := make([]domain.PersonID, len(c.Msg.PlayerIds))
	for i, id := range c.Msg.PlayerIds {
		playerIDs[i] = domain.PersonID(id)
	}
	staffIDs := make([]domain.PersonID, len(c.Msg.StaffIds))
	for i, id := range c.Msg.StaffIds {
		staffIDs[i] = domain.PersonID(id)
	}

	cmd := commands.NominatePersonsForTrainingCommand{
		TrainingID: domain.TrainingID(c.Msg.TrainingId),
		PlayerIDs:  playerIDs,
		StaffIDs:   staffIDs,
	}
	if err := t.cmds.NominatePersonsForTraining(ctx, &cmd); err != nil {
		return nil, t.handleCommonErrors(err)
	}
	return connect.NewResponse(&teamv1.NominatePersonsForTrainingResponse{}), nil
}

func gatheringPointToPb(point *view.GatheringPoint) *teamv1.GatheringPoint {
	if point == nil {
		return nil
	}
	return &teamv1.GatheringPoint{
		Location:       point.Location,
		GatheringUntil: localTimeToPb(&point.GatherUntil),
	}
}
func ratingSettingsToPb(settings view.RatingSettings) *teamv1.RatingSettings {
	return &teamv1.RatingSettings{
		Policy: trainingRatingPolicyToPb(settings.Policy),
	}
}

func acknowledgmentSettingsToPb(settings *view.AcknowledgmentSettings) *teamv1.AcknowledgementSettings {
	if settings == nil {
		return nil
	}
	return &teamv1.AcknowledgementSettings{
		Deadline: localTimeToPb(&settings.AcknowledgedUntil),
	}
}

func nominationResponsesToPb(nominations *view.Nominations) *teamv1.GetMyTeamHomeResponse_Nominations {
	if nominations == nil {
		return nil
	}
	playerResponses := core.Map(nominations.Players, mapTrainingNominationResponseToPb)
	staffResponses := core.Map(nominations.Staff, mapTrainingNominationResponseToPb)
	return &teamv1.GetMyTeamHomeResponse_Nominations{
		Players: playerResponses,
		Staff:   staffResponses,
	}
}

func mapTrainingNominationResponseToPb(ns *view.TrainingNominationResponse) *teamv1.GetMyTeamHomeResponse_Nomination {
	nomination := &teamv1.GetMyTeamHomeResponse_Nomination{
		PersonId:   string(ns.PersonID),
		PersonName: ns.PersonName,
		RsvpAt:     timestamppb.New(ns.NominatedAt),
	}
	switch ns.Type {
	case domain.TrainingNominationAccepted:
		nomination.Response = &teamv1.GetMyTeamHomeResponse_Nomination_Accepted_{
			Accepted: &teamv1.GetMyTeamHomeResponse_Nomination_Accepted{AcceptedAt: timestamppb.New(*ns.AcceptedAt)},
		}
	case domain.TrainingNominationDeclined:
		nomination.Response = &teamv1.GetMyTeamHomeResponse_Nomination_Declined_{
			Declined: &teamv1.GetMyTeamHomeResponse_Nomination_Declined{
				DeclinedAt: timestamppb.New(*ns.DeclinedAt),
				Reason:     ns.Reason,
			},
		}
	case domain.TrainingNominationTentative:
		nomination.Response = &teamv1.GetMyTeamHomeResponse_Nomination_Tentative_{
			Tentative: &teamv1.GetMyTeamHomeResponse_Nomination_Tentative{
				TentativeAt: timestamppb.New(*ns.TentativeAt),
				Reason:      ns.Reason,
			},
		}
	case domain.TrainingNominationUnacknowledged:
		nomination.Response = &teamv1.GetMyTeamHomeResponse_Nomination_NotAnswered_{
			NotAnswered: &teamv1.GetMyTeamHomeResponse_Nomination_NotAnswered{},
		}
	}

	return nomination
}
