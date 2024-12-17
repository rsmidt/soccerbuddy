package grpc

import (
	"connectrpc.com/connect"
	"context"
	v1 "github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/person/v1"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/person/v1/personv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type personServer struct {
	*baseHandler
}

func newPersonServiceHandler(base *baseHandler) personv1connect.PersonServiceHandler {
	return &personServer{baseHandler: base}
}

func (p *personServer) CreatePerson(ctx context.Context, c *connect.Request[v1.CreatePersonRequest]) (*connect.Response[v1.CreatePersonResponse], error) {
	cmd := commands.CreatePersonCommand{
		FirstName:    c.Msg.FirstName,
		LastName:     c.Msg.LastName,
		Birthdate:    c.Msg.Birthdate.AsTime(),
		OwningClubID: domain.ClubID(c.Msg.OwningClubId),
	}
	newPerson, err := p.cmds.CreatePerson(ctx, cmd)
	if err != nil {
		return nil, p.handleCommonErrors(err)
	}
	return &connect.Response[v1.CreatePersonResponse]{Msg: &v1.CreatePersonResponse{
		Id:        string(newPerson.ID),
		FirstName: newPerson.Firstname,
		LastName:  newPerson.Lastname,
		Birthdate: timestamppb.New(newPerson.Birthdate),
	}}, nil
}

func (p *personServer) ListPersonsInClub(ctx context.Context, c *connect.Request[v1.ListPersonsInClubRequest]) (*connect.Response[v1.ListPersonsInClubResponse], error) {
	query := queries.ListPersonsInClubQuery{
		OwningClubID: domain.ClubID(c.Msg.OwningClubId),
	}
	view, err := p.qs.ListPersonsInClub(ctx, query)
	if err != nil {
		return nil, p.handleCommonErrors(err)
	}
	resp := make([]*v1.ListPersonsInClubResponse_Person, len(view.Persons))
	for i, p := range view.Persons {
		resp[i] = &v1.ListPersonsInClubResponse_Person{
			Id:        string(p.ID),
			FirstName: p.FirstName,
			LastName:  p.LastName,
		}
	}
	return &connect.Response[v1.ListPersonsInClubResponse]{Msg: &v1.ListPersonsInClubResponse{
		Persons: resp,
	}}, nil

}

func (p *personServer) GetPersonOverview(ctx context.Context, c *connect.Request[v1.GetPersonOverviewRequest]) (*connect.Response[v1.GetPersonOverviewResponse], error) {
	query := queries.GetPersonOverviewQuery{
		ID: domain.PersonID(c.Msg.Id),
	}
	view, err := p.qs.GetPersonOverview(ctx, query)
	if err != nil {
		return nil, p.handleCommonErrors(err)
	}
	teams := make([]*v1.GetPersonOverviewResponse_Team, len(view.Teams))
	for i, t := range view.Teams {
		teams[i] = &v1.GetPersonOverviewResponse_Team{
			Id:       string(t.ID),
			Name:     t.Name,
			Role:     string(t.Role),
			JoinedAt: timestamppb.New(t.JoinedAt),
		}
	}
	pendingLinks := make([]*v1.GetPersonOverviewResponse_PendingAccountLink, len(view.PendingAccountLinks))
	for i, l := range view.PendingAccountLinks {
		pendingLinks[i] = &v1.GetPersonOverviewResponse_PendingAccountLink{
			LinkedAs:  accountLinkToPb(l.LinkedAs),
			InvitedBy: &v1.GetPersonOverviewResponse_Operator{FullName: l.InvitedBy.FullName},
			InvitedAt: timestamppb.New(l.InvitedAt),
			ExpiresAt: timestamppb.New(l.ExpiresAt),
		}
	}
	links := make([]*v1.GetPersonOverviewResponse_LinkedAccount, len(view.LinkedAccounts))
	for i, l := range view.LinkedAccounts {
		var invitedBy *v1.GetPersonOverviewResponse_LinkedAccount_Invite
		if l.InvitedBy != nil {
			invitedBy = &v1.GetPersonOverviewResponse_LinkedAccount_Invite{
				Invite: &v1.GetPersonOverviewResponse_LinkedAccount_OwnerLinked{
					InvitedBy: &v1.GetPersonOverviewResponse_Operator{FullName: l.InvitedBy.FullName},
					InvitedAt: timestamppb.New(*l.InvitedAt),
				},
			}
		}
		var linkedBy *v1.GetPersonOverviewResponse_LinkedAccount_External
		if l.LinkedBy != nil {
			linkedBy = &v1.GetPersonOverviewResponse_LinkedAccount_External{
				External: &v1.GetPersonOverviewResponse_LinkedAccount_ExternallyLinked{
					LinkedBy: &v1.GetPersonOverviewResponse_Operator{FullName: l.LinkedBy.FullName},
				},
			}
		}
		acc := &v1.GetPersonOverviewResponse_LinkedAccount{
			LinkedAs: accountLinkToPb(l.LinkedAs),
			FullName: l.FullName,
			LinkedAt: timestamppb.New(l.LinkedAt),
			Actor:    invitedBy,
		}
		if linkedBy != nil {
			acc.Actor = linkedBy
		} else if invitedBy != nil {
			acc.Actor = invitedBy
		}
		links[i] = acc
	}
	return &connect.Response[v1.GetPersonOverviewResponse]{Msg: &v1.GetPersonOverviewResponse{
		Id:                  string(view.ID),
		FirstName:           view.FirstName,
		LastName:            view.LastName,
		Birthdate:           timestamppb.New(view.Birthdate),
		CreatedAt:           timestamppb.New(view.CreatedAt),
		CreatedBy:           &v1.GetPersonOverviewResponse_Operator{FullName: view.CreatedBy.FullName},
		Teams:               teams,
		LinkedAccounts:      links,
		PendingAccountLinks: pendingLinks,
	}}, nil
}

func (p *personServer) InitiatePersonAccountLink(ctx context.Context, c *connect.Request[v1.InitiatePersonAccountLinkRequest]) (*connect.Response[v1.InitiatePersonAccountLinkResponse], error) {
	linkAs, err := pbToAccountLink(c.Msg.LinkAs)
	if err != nil {
		return nil, p.handleCommonErrors(err)
	}
	cmd := commands.InitiatePersonAccountLinkCommand{
		PersonID: domain.PersonID(c.Msg.PersonId),
		LinkAs:   linkAs,
	}
	resp, err := p.cmds.InitiatePersonAccountLink(ctx, cmd)
	if err != nil {
		return nil, p.handleCommonErrors(err)
	}
	return &connect.Response[v1.InitiatePersonAccountLinkResponse]{Msg: &v1.InitiatePersonAccountLinkResponse{
		LinkToken: string(resp.LinkToken),
		ExpiresAt: timestamppb.New(resp.ExpiresAt),
	}}, nil
}

func (p *personServer) DescribePendingPersonLink(ctx context.Context, c *connect.Request[v1.DescribePendingPersonLinkRequest]) (*connect.Response[v1.DescribePendingPersonLinkResponse], error) {
	query := queries.DescribePendingPersonLinkQuery{
		LinkToken: domain.PersonLinkToken(c.Msg.LinkToken),
	}
	view, err := p.qs.DescribePendingPersonLink(ctx, query)
	if err != nil {
		return nil, p.handleCommonErrors(err)
	}
	return connect.NewResponse(&v1.DescribePendingPersonLinkResponse{
		Person: &v1.DescribePendingPersonLinkResponse_Person{
			FullName:  view.FullName,
			LinkAs:    accountLinkToPb(view.LinkAs),
			InvitedBy: view.InvitedBy.FullName,
			ClubName:  view.Club.Name,
		},
	}), nil
}

func (p *personServer) ClaimPersonLink(ctx context.Context, c *connect.Request[v1.ClaimPersonLinkRequest]) (*connect.Response[v1.ClaimPersonLinkResponse], error) {
	cmd := commands.ClaimPersonLinkCommand{
		LinkToken: domain.PersonLinkToken(c.Msg.LinkToken),
	}
	err := p.cmds.ClaimPersonLink(ctx, cmd)
	if err != nil {
		return nil, p.handleCommonErrors(err)
	}
	return connect.NewResponse(&v1.ClaimPersonLinkResponse{}), nil
}
