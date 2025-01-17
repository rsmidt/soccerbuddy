package grpc

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	v1 "github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/account/v1"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/account/v1/accountv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

type accountServer struct {
	*baseHandler
}

func newAccountServiceHandler(base *baseHandler) accountv1connect.AccountServiceHandler {
	return &accountServer{baseHandler: base}
}

func (a *accountServer) GetMe(ctx context.Context, c *connect.Request[v1.GetMeRequest]) (*connect.Response[v1.GetMeResponse], error) {
	me, err := a.qs.GetMe(ctx)
	if err != nil {
		return nil, a.handleCommonErrors(err)
	}
	persons := make([]*v1.GetMeResponse_LinkedPerson, len(me.LinkedPersons))
	for i, p := range me.LinkedPersons {
		var linkedBy *v1.GetMeResponse_Operator
		if p.LinkedBy != nil {
			linkedBy = &v1.GetMeResponse_Operator{
				FullName: p.LinkedBy.FullName,
				IsMe:     p.LinkedBy.IsMe,
			}
		}

		memberships := make([]*v1.GetMeResponse_TeamMembership, len(p.TeamMemberships))
		for j, m := range p.TeamMemberships {
			memberships[j] = &v1.GetMeResponse_TeamMembership{
				Id:           string(m.ID),
				Name:         m.Name,
				Role:         string(m.Roles),
				JoinedAt:     timestamppb.New(m.JoinedAt),
				OwningClubId: string(m.OwningClubID),
			}
		}

		persons[i] = &v1.GetMeResponse_LinkedPerson{
			Id:              string(p.ID),
			LinkedAs:        accountLinkToPb(p.LinkedAs),
			FirstName:       p.FirstName,
			LastName:        p.LastName,
			LinkedAt:        timestamppb.New(p.LinkedAt),
			LinkedBy:        linkedBy,
			TeamMemberships: memberships,
			OwningClubId:    string(p.OwningClubID),
		}
	}
	return connect.NewResponse(&v1.GetMeResponse{
		Id:            string(me.ID),
		Email:         me.Email,
		FirstName:     me.FirstName,
		LastName:      me.LastName,
		LinkedPersons: persons,
	}), nil
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

func (a *accountServer) AttachMobileDevice(ctx context.Context, c *connect.Request[v1.AttachMobileDeviceRequest]) (*connect.Response[v1.AttachMobileDeviceResponse], error) {
	cmd := commands.AttachMobileDeviceCommand{
		InstallationID:          domain.InstallationID(c.Msg.InstallationId),
		NotificationDeviceToken: domain.NotificationDeviceToken(c.Msg.DeviceNotificationToken),
	}
	err := a.cmds.AttachMobileDevice(ctx, &cmd)
	if err != nil {
		return nil, a.handleCommonErrors(err)
	}
	return connect.NewResponse(&v1.AttachMobileDeviceResponse{}), nil
}
