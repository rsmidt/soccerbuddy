package commands

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

type CreateTeamCommand struct {
	Name string

	OwningClubID domain.ClubID
	Subject      *domain.PersonID
}

func (e *CreateTeamCommand) Validate() error {
	var errs validation.Errors
	if err := validation.ValidateStringRequiredWithLength(e.Name, "name", 5, 100); err != nil {
		errs = append(errs, err)
	}
	if e.OwningClubID == "" {
		errs = append(errs, validation.NewFieldError("owning_club_id", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) CreateTeam(ctx context.Context, cmd CreateTeamCommand) (*domain.Team, error) {
	ctx, span := tracing.Tracer.Start(ctx, "Commands.CreateTeam")
	defer span.End()

	// Make sure the user is allowed to create a team.
	if err := c.authorizer.Authorize(ctx, authz.ActionCreateTeam, authz.NewClubResource(cmd.OwningClubID)); err != nil {
		return nil, err
	}
	operator, err := c.authorizer.OptionalActingOperator(ctx, cmd.Subject)
	if err != nil {
		return nil, err
	}
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Make sure the team name is unique.
	if exists, err := c.repos.Team().ExistsByName(ctx, cmd.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, validation.NewFieldError("name", validation.ErrAlreadyExists)
	}

	// Make sure the owning club exists.
	if exists, err := c.repos.Club().ExistsByID(ctx, cmd.OwningClubID); err != nil {
		return nil, err
	} else if !exists {
		return nil, domain.ErrTeamOwningClubNotFound
	}
	id := idgen.New[domain.TeamID]()
	team, err := c.repos.Team().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	slug := domain.Slugify(cmd.Name)
	if err := team.Init(cmd.Name, slug, cmd.OwningClubID, operator); err != nil {
		return nil, err
	}
	if err := c.repos.Team().Save(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

type DeleteTeamCommand struct {
	ID      domain.TeamID
	Subject *domain.PersonID
}

func (c *Commands) DeleteTeam(ctx context.Context, cmd DeleteTeamCommand) error {
	ctx, span := tracing.Tracer.Start(ctx, "Commands.DeleteTeam")
	defer span.End()

	// Make sure the user is allowed to delete the team.
	if err := c.authorizer.Authorize(ctx, authz.ActionDelete, authz.NewTeamResource(cmd.ID)); err != nil {
		return err
	}
	operator, err := c.authorizer.RequiredActingOperator(ctx, cmd.Subject)
	if err != nil {
		return err
	}
	team, err := c.repos.Team().FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if err := team.Delete(operator); err != nil {
		return err
	}
	return c.repos.Team().Save(ctx, team)
}
