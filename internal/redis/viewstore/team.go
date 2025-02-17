package viewstore

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
	"github.com/rsmidt/soccerbuddy/internal/app/view"
	"github.com/rsmidt/soccerbuddy/internal/app/viewstore"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	"github.com/rsmidt/soccerbuddy/internal/redis"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"log/slog"
	"time"
)

var _ viewstore.TeamViewStore = (*RedisTeamViewStore)(nil)

type RedisTeamViewStore struct {
	rd     rueidis.Client
	logger *slog.Logger
}

func NewRedisTeamViewStore(rd rueidis.Client, logger *slog.Logger) *RedisTeamViewStore {
	return &RedisTeamViewStore{rd: rd, logger: logger}
}

func (r *RedisTeamViewStore) GetHome(ctx context.Context, permissions authz.PermissionsSet, teamID domain.TeamID) (*view.TeamHome, error) {
	ctx, span := tracing.Tracer.Start(ctx, "redisviewstore.GetHome")
	defer span.End()

	p, err := r.getTeamProjection(ctx, teamID)
	if err != nil {
		return nil, err
	}

	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return nil, domain.ErrUnauthenticated
	}
	me, err := getMe(ctx, r.rd, r.logger, principal.AccountID)
	if err != nil {
		return nil, err
	}

	var (
		isCoach      bool
		personInTeam *view.GetMeLinkedPerson
	)
outer:
	for _, person := range me.LinkedPersons {
		for _, membership := range person.TeamMemberships {
			if membership.ID != teamID {
				continue
			}
			personInTeam = person
			if membership.Role == domain.TeamMemberRoleCoach {
				// If it's a coach, we can show all trainings regardless of other persons with other roles.
				isCoach = true
				break outer
			}
		}
	}
	if personInTeam == nil {
		// TODO(SOC-29): Handle the case a guest is requesting access to a team.
		return nil, domain.ErrTeamMemberNotFound
	}

	var trainings []*projector.TrainingProjection
	if isCoach {
		trainings, err = r.getTrainingProjectionsByTeamID(ctx, teamID, time.Now())
	} else {
		trainings, err = r.getTrainingProjectionsByTeamIDAndPersonID(ctx, teamID, personInTeam.ID, time.Now())
	}
	if err != nil {
		return nil, err
	}

	ts := make([]*view.TeamHomeTraining, len(trainings))
	i := 0
	for _, tp := range trainings {
		var gatheringPoint *view.GatheringPoint
		if tp.GatheringPoint != nil {
			gatheringPoint = &view.GatheringPoint{
				Location:        tp.GatheringPoint.Location,
				GatherUntil:     tp.GatheringPoint.GatherUntil,
				GatherUntilIANA: tp.GatheringPoint.GatherUntilIANA,
			}
		}
		var acknowledgmentSettings *view.AcknowledgmentSettings
		if tp.AcknowledgmentSettings != nil {
			acknowledgmentSettings = &view.AcknowledgmentSettings{
				AcknowledgedUntil:     tp.AcknowledgmentSettings.AcknowledgeUntil,
				AcknowledgedUntilIANA: tp.AcknowledgmentSettings.AcknowledgeUntilIANA,
			}
		}
		var nominations *view.Nominations
		if permissions.Allows(authz.ActionEdit) {
			var playerResponses []*view.TrainingNominationResponse
			var staffResponses []*view.TrainingNominationResponse
			maybeMapOperator := func(operator *projector.OperatorProjection) *view.Operator {
				if operator == nil {
					return nil
				}
				return &view.Operator{
					FullName: operator.ActorFullName,
				}
			}
			for _, np := range tp.NominatedPlayers {
				playerResponses = append(playerResponses, &view.TrainingNominationResponse{
					PersonID:       np.ID,
					PersonName:     np.Name,
					Type:           np.Acknowledgment.Type,
					AcknowledgedAt: np.Acknowledgment.AcknowledgedAt,
					AcceptedAt:     np.Acknowledgment.AcceptedAt,
					TentativeAt:    np.Acknowledgment.TentativeAt,
					DeclinedAt:     np.Acknowledgment.DeclinedAt,
					Reason:         np.Acknowledgment.Reason,
					AcknowledgedBy: maybeMapOperator(np.Acknowledgment.AcknowledgedBy),
					NominatedAt:    np.NominatedAt,
				})
			}
			for _, ns := range tp.NominatedStaff {
				staffResponses = append(staffResponses, &view.TrainingNominationResponse{
					PersonID:       ns.ID,
					PersonName:     ns.Name,
					Type:           ns.Acknowledgment.Type,
					AcknowledgedAt: ns.Acknowledgment.AcknowledgedAt,
					AcceptedAt:     ns.Acknowledgment.AcceptedAt,
					TentativeAt:    ns.Acknowledgment.TentativeAt,
					DeclinedAt:     ns.Acknowledgment.DeclinedAt,
					Reason:         ns.Acknowledgment.Reason,
					AcknowledgedBy: maybeMapOperator(ns.Acknowledgment.AcknowledgedBy),
					NominatedAt:    ns.NominatedAt,
				})
			}
			nominations = &view.Nominations{
				Players: playerResponses,
				Staff:   staffResponses,
			}
		}
		ts[i] = &view.TeamHomeTraining{
			ID:                     tp.ID,
			ScheduledAt:            tp.ScheduledAt,
			ScheduledAtIANA:        tp.ScheduledAtIANA,
			EndsAt:                 tp.EndsAt,
			EndsAtIANA:             tp.EndsAtIANA,
			Description:            tp.Description,
			Location:               tp.Location,
			FieldType:              tp.FieldType,
			GatheringPoint:         gatheringPoint,
			AcknowledgmentSettings: acknowledgmentSettings,
			RatingSettings: view.RatingSettings{
				Policy: tp.RatingSettings.Policy,
			},
			ScheduledBy: view.Operator{
				FullName: tp.ScheduledBy.ActorFullName,
			},
			Nominations: nominations,
		}
		i++
	}
	return &view.TeamHome{
		ID:           p.ID,
		Name:         p.Name,
		Trainings:    ts,
		OwningClubID: p.OwningClubID,
	}, nil
}

func (r *RedisTeamViewStore) getTeamProjection(ctx context.Context, id domain.TeamID) (*projector.TeamProjection, error) {
	var a projector.TeamProjection
	cmd := r.rd.B().JsonGet().Key(fmt.Sprintf("%s%s", projector.ProjectionTeamPrefix, id)).Path(".").Build()
	return &a, r.rd.Do(ctx, cmd).DecodeJSON(&a)
}

func (r *RedisTeamViewStore) getTrainingProjectionsByTeamID(ctx context.Context, teamID domain.TeamID, minTime time.Time) ([]*projector.TrainingProjection, error) {
	rdq := fmt.Sprintf("@owning_team_id:{%s} @scheduled_at_ts:[%d +inf]", teamID, minTime.Unix())
	cmd := r.rd.B().FtSearch().Index(projector.ProjectionTrainingIDXName).Query(rdq).Sortby("scheduled_at_ts").Asc().Dialect(4).Build()
	_, docs, err := r.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	return redis.UnmarshalDocs[projector.TrainingProjection](docs)
}

func (r *RedisTeamViewStore) getTrainingProjectionsByTeamIDAndPersonID(ctx context.Context, teamID domain.TeamID, personId domain.PersonID, minTime time.Time) ([]*projector.TrainingProjection, error) {
	rdq := fmt.Sprintf("@owning_team_id:{%s} @scheduled_at_ts:[%d +inf] @nominated_person_ids:{%s}", teamID, minTime.Unix(), personId)
	cmd := r.rd.B().FtSearch().Index(projector.ProjectionTrainingIDXName).Query(rdq).Sortby("scheduled_at_ts").Asc().Dialect(4).Build()
	_, docs, err := r.rd.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	return redis.UnmarshalDocs[projector.TrainingProjection](docs)
}
