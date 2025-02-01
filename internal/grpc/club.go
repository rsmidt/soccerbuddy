package grpc

import (
	connect "connectrpc.com/connect"
	"context"
	v1 "github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/club/v1"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/club/v1/clubv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type clubServer struct {
	*baseHandler
}

func newClubServiceHandler(base *baseHandler) clubv1connect.ClubServiceHandler {
	return &clubServer{baseHandler: base}
}

func (cs *clubServer) CreateClub(ctx context.Context, c *connect.Request[v1.CreateClubRequest]) (*connect.Response[v1.CreateClubResponse], error) {
	cmd := commands.CreateClubCommand{
		Name: c.Msg.Name,
	}
	id, err := cs.cmds.CreateClub(ctx, cmd)
	if err != nil {
		return nil, cs.handleCommonErrors(err)
	}
	club, err := cs.qs.ClubByID(ctx, queries.ClubByIDQuery{ID: *id})
	if err != nil {
		return nil, cs.handleCommonErrors(err)
	}
	return connect.NewResponse(&v1.CreateClubResponse{
		Id:   club.ID.Deref(),
		Name: club.Name,
		Slug: club.Slug,
	}), nil
}

func (cs *clubServer) GetClubBySlug(ctx context.Context, c *connect.Request[v1.GetClubBySlugRequest]) (*connect.Response[v1.GetClubBySlugResponse], error) {
	query := queries.ClubBySlugQuery{
		Slug: c.Msg.Slug,
	}
	view, err := cs.qs.ClubBySlug(ctx, query)
	if err != nil {
		return nil, cs.handleCommonErrors(err)
	}
	if view == nil {
		return nil, connect.NewError(connect.CodeNotFound, nil)
	}
	return connect.NewResponse(&v1.GetClubBySlugResponse{
		Id:        view.ID.Deref(),
		Name:      view.Name,
		Slug:      view.Slug,
		CreatedAt: timestamppb.New(view.CreatedAt),
		UpdatedAt: timestamppb.New(view.UpdatedAt),
	}), nil
}

func (cs *clubServer) ListClubs(ctx context.Context, c *connect.Request[v1.ListClubsRequest]) (*connect.Response[v1.ListClubsResponse], error) {
	query := queries.ListClubsQuery{}
	view, err := cs.qs.ListClubs(ctx, query)
	if err != nil {
		return nil, cs.handleCommonErrors(err)
	}
	clubs := make([]*v1.ListClubsResponse_Club, len(view))
	for i, club := range view {
		clubs[i] = &v1.ListClubsResponse_Club{
			Id:        string(club.ID),
			Name:      club.Name,
			Slug:      club.Slug,
			CreatedAt: timestamppb.New(club.CreatedAt),
			UpdatedAt: timestamppb.New(club.UpdatedAt),
		}
	}
	return connect.NewResponse(&v1.ListClubsResponse{
		Clubs: clubs,
	}), nil
}

func (cs *clubServer) PromoteUserToAdmin(ctx context.Context, c *connect.Request[v1.PromoteUserToAdminRequest]) (*connect.Response[v1.PromoteUserToAdminResponse], error) {
	cmd := commands.PromoteUserToClubAdminCommand{
		ClubID: domain.ClubID(c.Msg.ClubId),
		UserID: domain.AccountID(c.Msg.AccountId),
	}
	err := cs.cmds.PromoteUserToClubAdmin(ctx, &cmd)
	if err != nil {
		return nil, cs.handleCommonErrors(err)
	}
	return connect.NewResponse(&v1.PromoteUserToAdminResponse{}), nil
}
