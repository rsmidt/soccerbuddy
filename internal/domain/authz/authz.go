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

// Authorizer is an interface that can be implemented by a type to authorize actions on resources.
type Authorizer interface {

	// Authorize checks if the action is allowed on the resource.
	Authorize(ctx context.Context, action string, resource *Resource) error

	// AuthorizedEntities returns the entities that the subject is authorized to perform the action on.
	AuthorizedEntities(ctx context.Context, action, resourceName string) (EntityIDSet, error)

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
)

const (
	ResourceSystemName   = "system"
	ResourceUserName     = "user"
	ResourceClubName     = "club"
	ResourceAccountName  = "account"
	ResourcePersonName   = "person"
	ResourceTeamName     = "team"
	ResourceTeamRoleName = "team_role"
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
)

const (
	SystemMainID = "main"

	RoleTrainer = "trainer"
)

var (
	// SystemResource is the main system resource.
	// All entities without a specific owner are owned by the system.
	SystemResource = &Resource{Name: ResourceSystemName, ID: SystemMainID}

	// NewClubResource creates a new club resource with the given ID.
	NewClubResource = newResourceCreator(ResourceClubName)

	// NewAccountResource creates a new account resource with the given ID.
	NewAccountResource = newResourceCreator(ResourceAccountName)

	// NewTeamResource creates a new team resource with the given ID.
	NewTeamResource = newResourceCreator(ResourceTeamName)

	// NewPersonResource creates a new person resource with the given ID.
	NewPersonResource = newResourceCreator(ResourcePersonName)
)

type ResourceIdentifier interface {
}

func newResourceCreator(name string) func(id string) *Resource {
	return func(id string) *Resource {
		return &Resource{
			Name: name,
			ID:   id,
		}
	}
}
