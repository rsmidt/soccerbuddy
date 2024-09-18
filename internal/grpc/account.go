package grpc

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	v1 "github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/account/v1"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/account/v1/accountv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"net/http"
)

type accountServer struct {
	*baseHandler
}

func newAccountServiceHandler(base *baseHandler) accountv1connect.AccountServiceHandler {
	return &accountServer{baseHandler: base}
}

func (a *accountServer) GetMe(ctx context.Context, c *connect.Request[v1.GetMeRequest]) (*connect.Response[v1.GetMeResponse], error) {
	var teams []*v1.GetMeResponse_Team
	teams = append(teams, &v1.GetMeResponse_Team{
		Id:   "7a528a53-e701-466c-b64d-4bd240298b39",
		Name: "Team 1",
	})

	var linkedProfiles []*v1.GetMeResponse_LinkedProfile
	linkedProfiles = append(linkedProfiles, &v1.GetMeResponse_LinkedProfile{
		Id:          "c186ed9e-169a-4dd0-b6db-80981fba7547",
		FirstName:   "Max",
		LastName:    "Musterfrau",
		ProfileType: "PLAYER",
		LinkedAs:    v1.GetMeResponse_LINKED_AS_PARENT,
		Team:        teams,
	})

	resp := &v1.GetMeResponse{
		Id:             "ceac884c-2912-4406-a60d-f9052cd8487e",
		Email:          "hello@world.com",
		FirstName:      "Marie",
		LastName:       "Musterfrau",
		LinkedProfiles: linkedProfiles,
	}

	return connect.NewResponse(resp), nil
}

func (a *accountServer) CreateAccount(ctx context.Context, c *connect.Request[v1.CreateAccountRequest]) (*connect.Response[v1.CreateAccountResponse], error) {
	cmd := commands.CreateAccountCommand{
		FirstName: c.Msg.FirstName,
		LastName:  c.Msg.LastName,
		Email:     c.Msg.Email,
		Password:  c.Msg.Password,
	}
	id, err := a.cmds.CreateAccount(ctx, cmd)
	if err != nil {
		return nil, a.handleCommonErrors(err)
	}

	return connect.NewResponse(&v1.CreateAccountResponse{
		Id:        string(id),
		Email:     c.Msg.Email,
		FirstName: c.Msg.FirstName,
		LastName:  c.Msg.LastName,
	}), nil
}

func (a *accountServer) Login(ctx context.Context, c *connect.Request[v1.LoginRequest]) (*connect.Response[v1.LoginResponse], error) {
	cmd := commands.LoginAccountCommand{
		Email:     c.Msg.Email,
		Password:  c.Msg.Password,
		UserAgent: c.Msg.UserAgent,
		IPAddress: c.Peer().Addr,
	}
	result, err := a.cmds.Login(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrWrongCredentials) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, a.handleCommonErrors(err)
	}
	res := connect.NewResponse(&v1.LoginResponse{
		SessionId: string(result.Token),
	})
	cookie := http.Cookie{
		Name:  "ID",
		Value: string(result.Token),
		// TODO: set domain
		Domain:   "",
		Path:     "/",
		Expires:  result.ExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	res.Header().Set("Set-Cookie", cookie.String())
	return res, nil
}

func (a *accountServer) RegisterAccount(ctx context.Context, c *connect.Request[v1.RegisterAccountRequest]) (*connect.Response[v1.RegisterAccountResponse], error) {
	//TODO implement me
	panic("implement me")
}
