package main

import (
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

var _ domain.Repositories = (*repositories)(nil)

// repositories implements domain.Repositories.
// It groups all repositories in this application for easier access.
type repositories struct {
	accountRepo    domain.AccountRepository
	clubRepo       domain.ClubRepository
	personRepo     domain.PersonRepository
	sessionRepo    domain.SessionRepository
	teamRepo       domain.TeamRepository
	teamMemberRepo domain.TeamMemberRepository
	trainingRepo   domain.TrainingRepository
}

func assembleRepositories(es eventing.EventStore) *repositories {
	return &repositories{
		accountRepo:    domain.NewEventSourcedAccountRepository(es),
		clubRepo:       domain.NewEventSourcedClubRepository(es),
		personRepo:     domain.NewEventSourcedPersonRepository(es),
		sessionRepo:    domain.NewEventSourcedSessionRepository(es),
		teamRepo:       domain.NewEventSourcedTeamRepository(es),
		teamMemberRepo: domain.NewEventSourcedTeamMemberRepository(es),
		trainingRepo:   domain.NewEventSourcedTrainingRepository(es),
	}
}

func (r *repositories) Account() domain.AccountRepository {
	return r.accountRepo
}

func (r *repositories) Club() domain.ClubRepository {
	return r.clubRepo
}

func (r *repositories) Person() domain.PersonRepository {
	return r.personRepo
}

func (r *repositories) Session() domain.SessionRepository {
	return r.sessionRepo
}

func (r *repositories) Team() domain.TeamRepository {
	return r.teamRepo
}

func (r *repositories) TeamMember() domain.TeamMemberRepository {
	return r.teamMemberRepo
}

func (r *repositories) Training() domain.TrainingRepository {
	return r.trainingRepo
}
