package view

import (
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"time"
)

type PersonOverviewTeam struct {
	ID       domain.TeamID
	Name     string
	Role     domain.TeamMemberRole
	JoinedAt time.Time
}

type PersonOverviewLinkedAccount struct {
	FullName  string
	LinkedAs  domain.AccountLink
	LinkedAt  time.Time
	InvitedBy *Operator
	InvitedAt *time.Time
	LinkedBy  *Operator
}

type PersonOverviewPendingAccountLink struct {
	LinkedAs  domain.AccountLink
	InvitedBy Operator
	InvitedAt time.Time
	ExpiresAt time.Time
}

type PersonOverview struct {
	ID                  domain.PersonID
	FirstName           string
	LastName            string
	Birthdate           time.Time
	CreatedAt           time.Time
	CreatedBy           Operator
	LinkedAccounts      []*PersonOverviewLinkedAccount
	PendingAccountLinks []*PersonOverviewPendingAccountLink
	Teams               []*PersonOverviewTeam
}
