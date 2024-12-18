package grpc

import (
	"connectrpc.com/connect"
	"context"
	teamv1 "github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/team/v1"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/team/v1/teamv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
	"github.com/rsmidt/soccerbuddy/internal/core"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Role:     domain.TeamMemberRoleRole(c.Msg.Role),
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
	members, err := t.qs.ListTeamMembers(ctx, query)
	if err != nil {
		return nil, t.handleCommonErrors(err)
	}
	var respMembers []*teamv1.ListTeamMembersResponse_Member
	for _, member := range members.MembersByPersonID {
		var inviterID *string
		if member.InviterID != nil {
			inviterID = core.PTR(string(*member.InviterID))
		}
		respMembers = append(respMembers, &teamv1.ListTeamMembersResponse_Member{
			Id:        string(member.ID),
			PersonId:  string(member.PersonID),
			InviterId: inviterID,
			FirstName: member.FirstName,
			LastName:  member.LastName,
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
	}
	if err := t.cmds.ScheduleTraining(ctx, &cmd); err != nil {
		return nil, err
	}
	// TODO: Return from projection?
	return connect.NewResponse(&teamv1.ScheduleTrainingResponse{}), nil
}
