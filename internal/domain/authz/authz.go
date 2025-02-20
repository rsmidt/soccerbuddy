package authz

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/domain"
)

var (
	// ErrUnauthorized is returned when the action is not allowed on the resource.
	ErrUnauthorized = errors.New("unauthorized")
)

type Resource struct {
	Name string
	ID   string
}

type EntityIDSet map[string]struct{}

type PermissionsSet map[string]struct{}

// Allows checks whether the exact supplied permission is granted.
func (p PermissionsSet) Allows(permission string) bool {
	_, ok := p[permission]
	return ok
}

// Authorizer is an interface that can be implemented by a type to authorize actions on resources.
type Authorizer interface {

	// Authorize checks if the action is allowed on the resource.
	Authorize(ctx context.Context, action string, resource *Resource) error

	// AuthorizedEntities returns the entities that the subject is authorized to perform the action on.
	AuthorizedEntities(ctx context.Context, action, resourceName string) (EntityIDSet, error)

	// Permissions returns the permissions that the subject has on the resource.
	Permissions(ctx context.Context, resource *Resource) (PermissionsSet, error)

	// OptionalActingOperator returns the operator for the subject.
	// It is not required that the principal acts on someone's behalf.
	OptionalActingOperator(ctx context.Context, personID *domain.PersonID) (domain.Operator, error)

	// RequiredActingOperator returns the operator for the subject.
	// It is required that the principal acts on someone's behalf, except for the root account.
	RequiredActingOperator(ctx context.Context, personID *domain.PersonID) (domain.Operator, error)
}

// Basic actions common to most resources.
const (
	ActionView   = "view"
	ActionEdit   = "edit"
	ActionDelete = "delete"
)

const (
	ActionCreateAccount      = "create_account"
	ActionCreateClub         = "create_club"
	ActionCreatePerson       = "create_person"
	ActionCreateTeam         = "create_team"
	ActionListPersons        = "list_persons"
	ActionPersonInitiateLink = "initiate_link"
	ActionScheduleTraining   = "schedule_training"
)

const (
	ResourceSystemName   = "system"
	ResourceUserName     = "user"
	ResourceClubName     = "club"
	ResourceAccountName  = "account"
	ResourcePersonName   = "person"
	ResourceTeamName     = "team"
	ResourceTeamRoleName = "team_role"
	ResourceTrainingName = "training"
)

const (
	RelationOwner  = "owner"
	RelationSystem = "system"
	RelationEditor = "editor"

	// RelationUser describes if a user is either the parent of the person or the person itself.
	RelationUser         = "user"
	RelationPersonSelf   = "self"
	RelationPersonParent = "parent"

	RelationClubPerson = "person"
	RelationClubAdmin  = "admin"

	RelationSystemAdmin      = "admin"
	RelationTeamMember       = "member"
	RelationTeamAdmin        = "admin"
	RelationTeamRoleAssignee = "assignee"

	RelationTrainingTeam        = "team"
	RelationTrainingParticipant = "participant"
)

const (
	SystemMainID = "main"
)

var (
	// SystemResource is the main system resource.
	// All entities without a specific owner are owned by the system.
	SystemResource = &Resource{Name: ResourceSystemName, ID: SystemMainID}

	// NewClubResource creates a new club resource with the given ID.
	NewClubResource = newResourceCreator[domain.ClubID](ResourceClubName)

	// NewAccountResource creates a new account resource with the given ID.
	NewAccountResource = newResourceCreator[domain.AccountID](ResourceAccountName)

	// NewTeamResource creates a new team resource with the given ID.
	NewTeamResource = newResourceCreator[domain.TeamID](ResourceTeamName)

	// NewPersonResource creates a new person resource with the given ID.
	NewPersonResource = newResourceCreator[domain.PersonID](ResourcePersonName)

	// NewTrainingResource creates a new training resource with the given ID.
	NewTrainingResource = newResourceCreator[domain.TrainingID](ResourceTrainingName)
)

type ResourceIdentifier interface {
}

func newResourceCreator[T ~string](name string) func(id T) *Resource {
	return func(id T) *Resource {
		return &Resource{
			Name: name,
			ID:   string(id),
		}
	}
}
