package commands

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type CreateRootAccountCommand struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func (c CreateRootAccountCommand) Validate() error {
	var errs validation.Errors
	if c.Email == "" {
		errs = append(errs, validation.NewFieldError("email", validation.ErrRequired))
	}
	if c.Password == "" {
		errs = append(errs, validation.NewFieldError("password", validation.ErrRequired))
	} else if len(c.Password) < 8 {
		errs = append(errs, validation.NewMinLengthError("password", 8))
	}
	if c.FirstName == "" {
		errs = append(errs, validation.NewFieldError("firstName", validation.ErrRequired))
	}
	if c.LastName == "" {
		errs = append(errs, validation.NewFieldError("lastName", validation.ErrRequired))
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) CreateRootAccount(ctx context.Context, cmd CreateRootAccountCommand) error {
	ctx, span := tracing.Tracer.Start(ctx, "commands.CreateRootAccount")
	defer span.End()

	if err := cmd.Validate(); err != nil {
		return err
	}

	exists, err := c.repos.Account().ExistsByEmail(ctx, cmd.Email)
	if err != nil && !errors.Is(err, domain.ErrAccountNotFound) {
		return err
	} else if exists {
		return domain.ErrRootAccountAlreadyInitialized
	}

	hashedPW, err := domain.Argon2idHashPassword(cmd.Password)
	if err != nil {
		return err
	}
	id := idgen.New[domain.AccountID]()
	account, err := c.repos.Account().FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := account.InitAsRoot(cmd.Email, hashedPW, cmd.FirstName, cmd.LastName); err != nil {
		return err
	}
	return c.repos.Account().Save(ctx, account)
}
