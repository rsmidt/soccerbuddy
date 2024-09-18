package grpc

import (
	"connectrpc.com/connect"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"log/slog"
)

func (b *baseHandler) handleCommonErrors(err error) error {
	internalErr := connect.NewError(connect.CodeInternal, nil)
	var vErrs validation.Errors
	if errors.As(err, &vErrs) {
		fvs := make([]*errdetails.BadRequest_FieldViolation, len(vErrs))
		for i, vErr := range vErrs {
			fvs[i] = &errdetails.BadRequest_FieldViolation{
				Field:       vErr.Field(),
				Description: vErr.Type(),
			}
		}
		br := errdetails.BadRequest{FieldViolations: fvs}
		detail, err := connect.NewErrorDetail(&br)
		if err != nil {
			return internalErr
		}
		cErr := connect.NewError(connect.CodeInvalidArgument, errors.New("validation failed"))
		cErr.AddDetail(detail)
		return cErr
	}
	var vErr validation.FieldError
	if errors.As(err, &vErr) {
		fv := &errdetails.BadRequest_FieldViolation{
			Field:       vErr.Field(),
			Description: vErr.Type(),
		}
		br := errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{fv}}
		detail, err := connect.NewErrorDetail(&br)
		if err != nil {
			return internalErr
		}
		cErr := connect.NewError(connect.CodeInvalidArgument, errors.New("validation failed"))
		cErr.AddDetail(detail)
		return cErr
	}
	if errors.Is(err, domain.ErrUnauthenticated) {
		return connect.NewError(connect.CodeUnauthenticated, nil)
	}
	if errors.Is(err, authz.ErrUnauthorized) {
		return connect.NewError(connect.CodePermissionDenied, nil)
	}
	if errors.Is(err, eventing.ErrOwnerNotFound) {
		return connect.NewError(connect.CodeNotFound, nil)
	}
	if errors.Is(err, eventing.ErrValueNotFound) {
		return connect.NewError(connect.CodeNotFound, nil)
	}
	var eErr domain.InvalidAggregateStateError
	if errors.As(err, &eErr) {
		return connect.NewError(connect.CodeFailedPrecondition, errors.New("invalid aggregate state"))
	}
	b.log.Warn("Received unhandled error in GRPC server", slog.String("err", err.Error()))

	return internalErr
}
