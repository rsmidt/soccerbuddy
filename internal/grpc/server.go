package grpc

import (
	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"errors"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/account/v1/accountv1connect"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/club/v1/clubv1connect"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/person/v1/personv1connect"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/team/v1/teamv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/grpc/middleware"
	"log/slog"
	"net/http"
)

type Server struct {
	cmds *commands.Commands
	qs   *queries.Queries
	log  *slog.Logger
}

func NewServer(cmds *commands.Commands, qs *queries.Queries, log *slog.Logger) *Server {
	return &Server{cmds: cmds, qs: qs, log: log}
}

type baseHandler struct {
	cmds *commands.Commands
	qs   *queries.Queries
	log  *slog.Logger
}

func (s *Server) Register(mux *http.ServeMux) error {
	authInterceptor := middleware.NewAuthenticationMiddleware(s.qs)

	base := &baseHandler{cmds: s.cmds, qs: s.qs, log: s.log}
	teamService := newTeamServiceHandler(base)
	accountService := newAccountServiceHandler(base)
	clubService := newClubServiceHandler(base)
	personService := newPersonServiceHandler(base)

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return err
	}
	commonOpts := connect.WithInterceptors(
		otelInterceptor,
		authInterceptor,
	)

	tspath, tshandler := teamv1connect.NewTeamServiceHandler(
		teamService,
		commonOpts,
	)
	aspath, ashandler := accountv1connect.NewAccountServiceHandler(
		accountService,
		commonOpts,
	)
	cpath, chandler := clubv1connect.NewClubServiceHandler(
		clubService,
		commonOpts,
	)
	ppath, phandler := personv1connect.NewPersonServiceHandler(
		personService,
		commonOpts,
	)

	reflector := grpcreflect.NewStaticReflector(
		teamv1connect.TeamServiceName,
		accountv1connect.AccountServiceName,
		clubv1connect.ClubServiceName,
		personv1connect.PersonServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(tspath, tshandler)
	mux.Handle(aspath, ashandler)
	mux.Handle(cpath, chandler)
	mux.Handle(ppath, phandler)
	return nil
}

func pbToAccountLink(s soccerbuddy.AccountLink) (domain.AccountLink, error) {
	switch s {
	case soccerbuddy.AccountLink_LINKED_AS_PARENT:
		return domain.AccountLinkParent, nil
	case soccerbuddy.AccountLink_LINKED_AS_SELF:
		return domain.AccountLinkSelf, nil
	default:
		return "", errors.New("unknown linked as")
	}
}

func accountLinkToPb(s domain.AccountLink) soccerbuddy.AccountLink {
	switch s {
	case domain.AccountLinkParent:
		return soccerbuddy.AccountLink_LINKED_AS_PARENT
	case domain.AccountLinkSelf:
		return soccerbuddy.AccountLink_LINKED_AS_SELF
	default:
		return soccerbuddy.AccountLink_LINKED_AS_UNSPECIFIED
	}
}
