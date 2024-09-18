package grpc

import (
	connect "connectrpc.com/connect"
	"context"
	v1 "github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/club/v1"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/club/v1/clubv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
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
