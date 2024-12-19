package grpc

import (
	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"errors"
	"fmt"
	_type "github.com/rsmidt/soccerbuddy/gen/go/google/type"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/account/v1/accountv1connect"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/club/v1/clubv1connect"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/person/v1/personv1connect"
	teamv1 "github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/team/v1"
	"github.com/rsmidt/soccerbuddy/gen/go/soccerbuddy/team/v1/teamv1connect"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/grpc/middleware"
	"log/slog"
	"net/http"
	"time"
)

// The default location to interpret times in.
var defaultLocation *time.Location

func init() {
	var err error
	defaultLocation, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(fmt.Errorf("failed to load location: %w", err))
	}
}

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

func pbToLocalTime(at *_type.DateTime, loc *time.Location) time.Time {
	return time.Date(int(at.GetYear()), time.Month(at.GetMonth()), int(at.GetDay()), int(at.GetHours()), int(at.GetMinutes()), int(at.GetSeconds()), int(at.GetNanos()), loc)
}

func localTimeToPb(time *time.Time) *_type.DateTime {
	if time == nil {
		return nil
	}
	return &_type.DateTime{
		Year:       int32(time.Year()),
		Month:      int32(time.Month()),
		Day:        int32(time.Day()),
		Hours:      int32(time.Hour()),
		Minutes:    int32(time.Minute()),
		Seconds:    int32(time.Second()),
		Nanos:      int32(time.Nanosecond()),
		TimeOffset: nil,
	}
}

func pbToGatheringPoint(point *teamv1.GatheringPoint) *domain.TrainingGatheringPoint {
	if point == nil {
		return nil
	}
	return domain.NewTrainingGatheringPoint(point.Location, pbToLocalTime(point.GatheringUntil, defaultLocation), defaultLocation.String())
}

func pbToAcknowledgementSettings(settings *teamv1.AcknowledgementSettings) *domain.TrainingAcknowledgmentSettings {
	if settings == nil {
		return nil
	}
	return domain.NewTrainingAcknowledgmentSettings(pbToLocalTime(settings.Deadline, defaultLocation), defaultLocation.String())
}

func pbToRatingSettings(settings *teamv1.RatingSettings) *domain.TrainingRatingSettings {
	if settings == nil {
		return nil
	}
	return domain.NewTrainingRatingSettings(pbToTrainingsPolicy(settings.Policy))
}

func pbToTrainingsPolicy(policy soccerbuddy.RatingPolicy) domain.TrainingRatingPolicy {
	switch policy {
	case soccerbuddy.RatingPolicy_RATING_POLICY_UNSPECIFIED:
		return domain.TrainingRatingPolicyUnspecified
	case soccerbuddy.RatingPolicy_RATING_POLICY_ALLOWED:
		return domain.TrainingRatingPolicyAllowed
	case soccerbuddy.RatingPolicy_RATING_POLICY_FORBIDDEN:
		return domain.TrainingRatingPolicyForbidden
	case soccerbuddy.RatingPolicy_RATING_POLICY_REQUIRED:
		return domain.TrainingRatingPolicyRequired
	default:
		return domain.TrainingRatingPolicyUnspecified
	}
}
