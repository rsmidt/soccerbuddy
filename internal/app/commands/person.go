package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/rsmidt/soccerbuddy/internal/redis"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"time"
)

type CreatePersonCommand struct {
	FirstName    string
	LastName     string
	Birthdate    time.Time
	OwningClubID domain.ClubID
}

func (c *CreatePersonCommand) Validate() error {
	var errs validation.Errors
	if err := validation.ValidateStringRequiredWithLength(c.FirstName, "firstname", 1, 50); err != nil {
		errs = append(errs, err)
	}
	if err := validation.ValidateStringRequiredWithLength(c.LastName, "lastname", 1, 50); err != nil {
		errs = append(errs, err)
	}
	if err := validation.ValidateDateRequiredInRange(c.Birthdate, "birthdate", domain.PersonMinBirthdate, domain.PersonMaxBirthdate); err != nil {
		errs = append(errs, err)
	}
	if c.OwningClubID == "" {
		errs = append(errs, validation.NewFieldError("owning_club_id", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) CreatePerson(ctx context.Context, cmd CreatePersonCommand) (*domain.Person, error) {
	ctx, span := tracing.Tracer.Start(ctx, "commands.CreatePerson")
	defer span.End()

	if err := c.authorizer.Authorize(ctx, authz.ActionCreatePerson, authz.NewClubResource(cmd.OwningClubID)); err != nil {
		return nil, err
	}
	operator, err := c.authorizer.OptionalActingOperator(ctx, nil)
	if err != nil {
		return nil, err
	}
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	if exists, err := c.repos.Club().ExistsByID(ctx, cmd.OwningClubID); err != nil {
		return nil, err
	} else if !exists {
		return nil, domain.ErrOwningClubNotFound
	}

	person, err := c.repos.Person().FindByID(ctx, idgen.New[domain.PersonID]())
	if err != nil {
		return nil, err
	}
	person.Init(cmd.FirstName, cmd.LastName, cmd.Birthdate, operator, cmd.OwningClubID)
	if err := c.repos.Person().Save(ctx, person); err != nil {
		return nil, err
	}
	return person, nil
}

type AddPersonToTeamCommand struct {
	TeamID    domain.TeamID
	PersonID  domain.PersonID
	InviterID *domain.PersonID
	Role      domain.TeamMemberRole
}

func (c *AddPersonToTeamCommand) Validate() error {
	var errs validation.Errors
	if c.TeamID == "" {
		errs = append(errs, validation.NewFieldError("team_id", validation.ErrRequired))
	}
	if c.PersonID == "" {
		errs = append(errs, validation.NewFieldError("person_id", validation.ErrRequired))
	}
	if c.Role != "COACH" && c.Role != "PLAYER" {
		errs = append(errs, validation.NewFieldError("role", validation.ErrInvalidChoice))
	}
	if c.InviterID != nil && *c.InviterID == "" {
		errs = append(errs, validation.NewFieldError("inviter_id", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) AddPersonToTeam(ctx context.Context, cmd AddPersonToTeamCommand) error {
	ctx, span := tracing.Tracer.Start(ctx, "commands.AddPersonToTeam")
	defer span.End()

	if err := c.authorizer.Authorize(ctx, authz.ActionEdit, authz.NewTeamResource(cmd.TeamID)); err != nil {
		return err
	}
	operator, err := c.authorizer.OptionalActingOperator(ctx, cmd.InviterID)
	if err != nil {
		return err
	}
	if err := cmd.Validate(); err != nil {
		return err
	}
	if exists, err := c.repos.Team().ExistsByID(ctx, cmd.TeamID); err != nil {
		return err
	} else if !exists {
		return domain.ErrTeamNotFound
	}
	if exists, err := c.repos.Person().ExistsByID(ctx, cmd.PersonID); err != nil {
		return err
	} else if !exists {
		return domain.ErrPersonNotFound
	}
	member, err := c.repos.TeamMember().FindByTeamAndPerson(ctx, cmd.TeamID, cmd.PersonID)
	if errors.Is(err, domain.ErrTeamMemberNotFound) {
		// If the team member does not exist, create a new one.
		member = domain.NewTeamMember(idgen.New[domain.TeamMemberID](), cmd.TeamID, cmd.PersonID)
	} else if err != nil {
		return err
	}
	if err := member.Invite(operator, cmd.Role); err != nil {
		return err
	}
	if err := c.repos.TeamMember().Save(ctx, member); err != nil {
		return err
	}
	return nil
}

type InitiatePersonAccountLinkCommand struct {
	PersonID domain.PersonID
	LinkAs   domain.AccountLink
}

func (c *InitiatePersonAccountLinkCommand) Validate() error {
	var errs validation.Errors
	if c.PersonID == "" {
		errs = append(errs, validation.NewFieldError("person_id", validation.ErrRequired))
	}
	if c.LinkAs != domain.AccountLinkParent && c.LinkAs != domain.AccountLinkSelf {
		errs = append(errs, validation.NewFieldError("link_as", validation.ErrInvalidChoice))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type InitiatePersonAccountLinkResult struct {
	LinkToken domain.PersonLinkToken
	ExpiresAt time.Time
}

func (c *Commands) InitiatePersonAccountLink(ctx context.Context, command InitiatePersonAccountLinkCommand) (*InitiatePersonAccountLinkResult, error) {
	ctx, span := tracing.Tracer.Start(ctx, "commands.InitiatePersonAccountLink")
	defer span.End()

	if err := c.authorizer.Authorize(ctx, authz.ActionPersonInitiateLink, authz.NewPersonResource(command.PersonID)); err != nil {
		return nil, err
	}
	operator, err := c.authorizer.OptionalActingOperator(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := command.Validate(); err != nil {
		return nil, err
	}
	// If the desired role is self, we need to verify that there does not yet exist a link for it.
	person, err := c.repos.Person().FindByID(ctx, command.PersonID)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(time.Hour * 24 * 7)
	token, err := randomString(16)
	if err != nil {
		return nil, err
	}
	linkToken := domain.PersonLinkToken(token)
	if err := person.InitiateNewLink(operator, command.LinkAs, linkToken, expiresAt); err != nil {
		return nil, err
	}
	if err := c.repos.Person().Save(ctx, person); err != nil {
		return nil, err
	}
	return &InitiatePersonAccountLinkResult{
		LinkToken: linkToken,
		ExpiresAt: expiresAt,
	}, nil
}

type ClaimPersonLinkCommand struct {
	LinkToken domain.PersonLinkToken
}

func (c *ClaimPersonLinkCommand) Validate() error {
	var errs validation.Errors
	if c.LinkToken == "" {
		errs = append(errs, validation.NewFieldError("link_token", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) ClaimPersonLink(ctx context.Context, cmd ClaimPersonLinkCommand) error {
	ctx, span := tracing.Tracer.Start(ctx, "commands.ClaimPersonLink")
	defer span.End()

	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return domain.ErrUnauthenticated
	}
	// First query the projection to get the person ID and token information.
	pers, err := c.getPersonProjectionByPendingToken(ctx, cmd.LinkToken)
	if err != nil {
		return err
	}
	if len(pers) == 0 {
		return domain.ErrPersonInvalidLinkToken
	}
	persProjection := pers[0]

	// Claim the token on the person side.
	person, err := c.repos.Person().FindByID(ctx, persProjection.ID)
	if err != nil {
		return err
	}
	if err := person.Claim(cmd.LinkToken, principal.AccountID); err != nil {
		return err
	}
	pl, err := person.FindPendingLink(cmd.LinkToken)
	if err != nil {
		return err
	}

	// Make the link on the account side.
	account, err := c.repos.Account().FindByID(ctx, principal.AccountID)
	if err != nil {
		return err
	}
	if err := account.Link(person.ID, pl.LinkAs, nil, persProjection.OwningClubID); err != nil {
		return err
	}

	// Save both aggregates.
	if err := c.repos.Person().Save(ctx, person); err != nil {
		return err
	}
	// TODO: Emit a compensation event on person if the append fails or append both in the same transaction (basically we need a saga...).
	if err := c.repos.Account().Save(ctx, account); err != nil {
		return err
	}
	return nil
}

func (c *Commands) getPersonProjectionByPendingToken(ctx context.Context, token domain.PersonLinkToken) ([]*projector.PersonProjection, error) {
	// TODO: this should be more abstracted.
	rdq := fmt.Sprintf("@pending_link_token:(%s)", token)
	cmd := c.rd.B().FtSearch().Index(projector.ProjectionPersonIDXName).Query(rdq).Build()
	_, docs, err := c.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	return redis.UnmarshalDocs[projector.PersonProjection](docs)
}
