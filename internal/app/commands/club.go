package commands

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type CreateClubCommand struct {
	Name string `json:"name"`
}

func (c CreateClubCommand) Validate() error {
	var errs validation.Errors
	if err := validation.ValidateStringRequiredWithLength(c.Name, "name", 3, 50); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) CreateClub(ctx context.Context, cmd CreateClubCommand) (*domain.ClubID, error) {
	ctx, span := tracing.Tracer.Start(ctx, "commands.CreateClub")
	defer span.End()

	if err := c.authorizer.Authorize(ctx, authz.ActionCreateClub, authz.SystemResource); err != nil {
		return nil, err
	}
	err := cmd.Validate()
	if err != nil {
		return nil, err
	}

	if exists, err := c.repos.Club().ExistsByName(ctx, cmd.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, validation.NewExistsError("name")
	}

	id := idgen.New[domain.ClubID]()
	slug := domain.Slugify(cmd.Name)
	club, err := c.repos.Club().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := club.Init(cmd.Name, slug); err != nil {
		return nil, err
	}
	if err := c.repos.Club().Save(ctx, club); err != nil {
		return nil, err
	}
	return &id, nil
}
