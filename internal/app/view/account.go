package view

import (
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"time"
)

type GetMe struct {
	ID            domain.AccountID
	Email         string
	FirstName     string
	LastName      string
	LinkedPersons []*GetMeLinkedPerson
	IsSuper       bool
}

type GetMeOperator struct {
	FullName string
	// If the operator that performed this account is this account owner.
	IsMe bool
}

type GetMeLinkedPerson struct {
	ID              domain.PersonID
	FirstName       string
	LastName        string
	LinkedAs        domain.AccountLink
	LinkedAt        time.Time
	LinkedBy        *GetMeOperator
	TeamMemberships []*GetMeLinkedPersonTeamMemberships
	OwningClubID    domain.ClubID
}

type GetMeLinkedPersonTeamMemberships struct {
	ID           domain.TeamID
	Name         string
	OwningClubID domain.ClubID
	JoinedAt     time.Time
	Role         domain.TeamMemberRole
}
