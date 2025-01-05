package commands

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"time"
)

type ScheduleTrainingCommand struct {
	ScheduledAt     time.Time
	ScheduledAtIANA string
	EndsAt          time.Time
	EndsAtIANA      string

	Description            *string
	Location               *string
	FieldType              *string
	GatheringPoint         *domain.TrainingGatheringPoint
	AcknowledgmentSettings *domain.TrainingAcknowledgmentSettings
	RatingSettings         *domain.TrainingRatingSettings
	Nominations            *domain.TrainingNominations

	TeamID domain.TeamID
}

func (s *ScheduleTrainingCommand) Validate() error {
	var errs validation.Errors
	if s.ScheduledAt.IsZero() {
		errs = append(errs, validation.NewFieldError("scheduled_at", validation.ErrRequired))
	}
	if s.ScheduledAtIANA == "" {
		errs = append(errs, validation.NewFieldError("scheduled_at_iana", validation.ErrRequired))
	}
	if s.EndsAt.IsZero() {
		errs = append(errs, validation.NewFieldError("ends_at", validation.ErrRequired))
	}
	if s.EndsAtIANA == "" {
		errs = append(errs, validation.NewFieldError("ends_at_iana", validation.ErrRequired))
	}
	if s.ScheduledAt.After(s.EndsAt) {
		errs = append(errs, validation.NewFieldError("ends_at", validation.ErrDateBefore))
	}
	if s.Description != nil && *s.Description == "" {
		errs = append(errs, validation.NewFieldError("description", validation.ErrNotEmpty))
	}
	if s.Location != nil && *s.Location == "" {
		errs = append(errs, validation.NewFieldError("location", validation.ErrNotEmpty))
	}
	if s.FieldType != nil && *s.FieldType == "" {
		errs = append(errs, validation.NewFieldError("field_type", validation.ErrNotEmpty))
	}
	if s.GatheringPoint != nil {
		if s.GatheringPoint.Location == "" {
			errs = append(errs, validation.NewFieldError("gathering_point.location", validation.ErrNotEmpty))
		}
		if s.GatheringPoint.GatherUntil.IsZero() {
			errs = append(errs, validation.NewFieldError("gathering_point.gather_until", validation.ErrRequired))
		}
		if s.GatheringPoint.GatherUntilIANA == "" {
			errs = append(errs, validation.NewFieldError("gathering_point.gather_until_iana", validation.ErrRequired))
		}
	}
	if s.RatingSettings != nil {
		if s.RatingSettings.Policy == domain.TrainingRatingPolicyUnspecified {
			errs = append(errs, validation.NewFieldError("rating_settings.policy", validation.ErrRequired))
		}
	}
	if s.AcknowledgmentSettings != nil {
		if s.AcknowledgmentSettings.AcknowledgeUntil.IsZero() {
			errs = append(errs, validation.NewFieldError("acknowledgment_settings.acknowledge_until", validation.ErrRequired))
		}
		if s.AcknowledgmentSettings.AcknowledgeUntilIANA == "" {
			errs = append(errs, validation.NewFieldError("acknowledgment_settings.acknowledge_until_iana", validation.ErrRequired))
		}
	}
	if s.Nominations != nil {
		if len(s.Nominations.PlayerIDs) == 0 && len(s.Nominations.StaffIDs) == 0 {
			errs = append(errs, validation.NewFieldError("nominations", validation.ErrRequired))
		}
	}
	if s.TeamID == "" {
		errs = append(errs, validation.NewFieldError("team_id", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) ScheduleTraining(ctx context.Context, cmd *ScheduleTrainingCommand) error {
	ctx, span := tracing.Tracer.Start(ctx, "commands.ScheduleTraining")
	defer span.End()

	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := c.authorizer.Authorize(ctx, authz.ActionScheduleTraining, authz.NewTeamResource(cmd.TeamID)); err != nil {
		return err
	}
	operator, err := c.authorizer.OptionalActingOperator(ctx, nil)
	if err != nil {
		return err
	}

	// TODO: Move to either command validation (need to make side-effect free).
	if cmd.ScheduledAt.Before(time.Now()) {
		return validation.NewFieldError("scheduled_at", validation.ErrMinDate)
	}

	team, err := c.repos.Team().FindByID(ctx, cmd.TeamID)
	if err != nil {
		return err
	}
	if team.State != domain.TeamStateActive {
		return errors.New("team is not active")
	}

	training := domain.NewTraining(idgen.New[domain.TrainingID](), cmd.TeamID, team.OwningClubID)

	ratingSettings := domain.TrainingRatingSettings{
		Policy: domain.DefaultTrainingRatingPolicy,
	}
	if cmd.RatingSettings != nil {
		ratingSettings = *cmd.RatingSettings
	}

	// TODO: Test for date collision?
	err = training.Schedule(
		cmd.ScheduledAt,
		cmd.ScheduledAtIANA,
		cmd.EndsAt,
		cmd.EndsAtIANA,
		cmd.Description,
		cmd.Location,
		cmd.FieldType,
		cmd.GatheringPoint,
		cmd.AcknowledgmentSettings,
		ratingSettings,
		operator,
	)
	if err != nil {
		return err
	}
	if err := c.repos.Training().Save(ctx, training); err != nil {
		return err
	}

	if cmd.Nominations != nil {
		// TODO: Replace with a single operation.
		training, err = c.repos.Training().FindByID(ctx, training.ID)
		if err != nil {
			return err
		}
		// TODO: make policy configurable.
		if err := training.NominatePersons(cmd.Nominations.PlayerIDs, cmd.Nominations.StaffIDs, operator, domain.TrainingNominationNotificationPolicySilent); err != nil {
			return err
		}
		return c.repos.Training().Save(ctx, training)
	}
	return nil
}

type NominatePersonsForTrainingCommand struct {
	TrainingID domain.TrainingID
	PlayerIDs  []domain.PersonID
	StaffIDs   []domain.PersonID
}

func (n *NominatePersonsForTrainingCommand) Validate() error {
	var errs validation.Errors
	if n.TrainingID == "" {
		errs = append(errs, validation.NewFieldError("training_id", validation.ErrRequired))
	}
	if len(n.PlayerIDs) == 0 && len(n.StaffIDs) == 0 {
		errs = append(errs, validation.NewFieldError("player_ids", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) NominatePersonsForTraining(ctx context.Context, cmd *NominatePersonsForTrainingCommand) error {
	ctx, span := tracing.Tracer.Start(ctx, "commands.NominatePersonsForTraining")
	defer span.End()

	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := c.authorizer.Authorize(ctx, authz.ActionEdit, authz.NewTrainingResource(cmd.TrainingID)); err != nil {
		return err
	}
	operator, err := c.authorizer.OptionalActingOperator(ctx, nil)
	if err != nil {
		return err
	}

	training, err := c.repos.Training().FindByID(ctx, cmd.TrainingID)
	if err != nil {
		return err
	}
	if err := training.NominatePersons(cmd.PlayerIDs, cmd.StaffIDs, operator, domain.TrainingNominationNotificationPolicySilent); err != nil {
		return err
	}
	return c.repos.Training().Save(ctx, training)
}
