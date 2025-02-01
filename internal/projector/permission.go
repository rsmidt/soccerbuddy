package projector

import (
	"context"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
)

const (
	PermissionProjectorName eventing.ProjectionName = "permission_projector"
)

type permissionProjector struct {
	relationStore authz.RelationStore
}

func NewPermissionProjector(relationStore authz.RelationStore) eventing.Projector {
	return &permissionProjector{
		relationStore: relationStore,
	}
}

func (a *permissionProjector) Init(ctx context.Context) error {
	return nil
}

func (a *permissionProjector) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(domain.AccountAggregateType).
		Events(
			domain.AccountCreatedEventType,
			domain.RootAccountCreatedEventType,
			domain.AccountLinkedToPersonEventType,
			domain.AccountRegisteredEventType,
		).Finish().
		WithAggregate(domain.TeamMemberAggregateType).
		Events(domain.PersonInvitedToTeamEventType).Finish().
		WithAggregate(domain.PersonAggregateType).
		Events(domain.PersonCreatedEventType).Finish().
		WithAggregate(domain.TeamAggregateType).
		Events(domain.TeamCreatedEventType, domain.TeamDeletedEventType).Finish().
		WithAggregate(domain.ClubAggregateType).
		Events(domain.ClubCreatedEventType, domain.ClubAdminAddedEventType).Finish().
		WithAggregate(domain.TrainingAggregateType).
		Events(domain.TrainingScheduledEventType, domain.PersonsNominatedForTrainingEventType).Finish().
		MustBuild()
}

func (a *permissionProjector) Projection() eventing.ProjectionName {
	return PermissionProjectorName
}

func (a *permissionProjector) Project(ctx context.Context, events ...*eventing.JournalEvent) error {
	ctx, span := tracing.Tracer.Start(ctx, "projector.Permission.Project")
	defer span.End()

	var err error
	for _, event := range events {
		switch e := event.Event.(type) {
		case *domain.AccountCreatedEvent:
			err = a.createAccountPermissions(ctx, event, e)
		case *domain.AccountLinkedToPersonEvent:
			err = a.createLinkedToPersonPermissions(ctx, event, e)
		case *domain.RootAccountCreatedEvent:
			err = a.createRootAccountPermissions(ctx, event, e)
		case *domain.AccountRegisteredEvent:
			err = a.createAccountRegisteredPermissions(ctx, event, e)
		case *domain.PersonInvitedToTeamEvent:
			err = a.createTeamMemberPermissions(ctx, event, e)
		case *domain.PersonCreatedEvent:
			err = a.createPersonPermissions(ctx, event, e)
		case *domain.TeamCreatedEvent:
			err = a.createTeamPermissions(ctx, event, e)
		case *domain.TeamDeletedEvent:
			err = a.deleteTeamPermissions(ctx, event, e)
		case *domain.ClubCreatedEvent:
			err = a.createClubPermissions(ctx, event, e)
		case *domain.ClubAdminAddedEvent:
			err = a.createClubAdminPermissions(ctx, event, e)
		case *domain.TrainingScheduledEvent:
			err = a.createTrainingPermissions(ctx, event, e)
		case *domain.PersonsNominatedForTrainingEvent:
			err = a.createPersonsNominatedForTrainingPermissions(ctx, event, e)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *permissionProjector) createAccountPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountCreatedEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate the system to the account.
		Entity(authz.ResourceAccountName, event.AggregateID().Deref()).
		Subject(authz.ResourceSystemName, authz.SystemMainID).
		Relate(authz.RelationSystem).And().
		// Relate the user to the account as owner.
		Entity(authz.ResourceAccountName, event.AggregateID().Deref()).
		Subject(authz.ResourceUserName, event.AggregateID().Deref()).
		Relate(authz.RelationOwner).And().
		Build()
	return a.relationStore.AddRelations(ctx, relations)
}

func (a *permissionProjector) createRootAccountPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.RootAccountCreatedEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate system to the account.
		Entity(authz.ResourceAccountName, event.AggregateID().Deref()).
		Subject(authz.ResourceSystemName, authz.SystemMainID).
		Relate(authz.RelationSystem).And().
		// Relate the user to the system as an admin.
		Entity(authz.ResourceSystemName, authz.SystemMainID).
		Subject(authz.ResourceUserName, event.AggregateID().Deref()).
		Relate(authz.RelationSystemAdmin).Build()
	return a.relationStore.AddRelations(ctx, relations)
}

func (a *permissionProjector) createAccountRegisteredPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountRegisteredEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate the system to the account.
		Entity(authz.ResourceAccountName, event.AggregateID().Deref()).
		Subject(authz.ResourceSystemName, authz.SystemMainID).
		Relate(authz.RelationSystem).
		// Relate the user to the account as owner.
		Entity(authz.ResourceAccountName, event.AggregateID().Deref()).
		Subject(authz.ResourceUserName, event.AggregateID().Deref()).
		Relate(authz.RelationOwner).And().
		Build()
	return a.relationStore.AddRelations(ctx, relations)
}

func (a *permissionProjector) createTeamMemberPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonInvitedToTeamEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate the person to the team as a team member.
		Entity(authz.ResourceTeamName, string(e.TeamID)).
		Subject(authz.ResourcePersonName, string(e.PersonID)).
		Relate(authz.RelationTeamMember).And().
		// Relate the person to the team role.
		Entity(authz.ResourceTeamRoleName, string(e.AssignedRole)).
		Subject(authz.ResourcePersonName, string(e.PersonID)).
		Relate(authz.RelationTeamRoleAssignee).
		Build()
	return a.relationStore.AddRelations(ctx, relations)
}

func (a *permissionProjector) createPersonPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonCreatedEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate the club to the person as the owner.
		Entity(authz.ResourcePersonName, event.AggregateID().Deref()).
		Subject(authz.ResourceClubName, string(e.OwningClubID)).
		Relate(authz.RelationOwner).And().
		// Relate the person the club as a domain.
		Entity(authz.ResourceClubName, string(e.OwningClubID)).
		Subject(authz.ResourcePersonName, event.AggregateID().Deref()).
		Relate(authz.RelationClubPerson).Build()
	return a.relationStore.AddRelations(ctx, relations)
}

func (a *permissionProjector) createTeamPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.TeamCreatedEvent) error {
	var builder authz.RelationBuilder
	b := builder.
		// Relate the club to the team as the owner.
		Entity(authz.ResourceTeamName, event.AggregateID().Deref()).
		Subject(authz.ResourceClubName, string(e.OwningClubID)).
		Relate(authz.RelationOwner).And().
		// Relate the team role "trainer" to the team as an editor.
		Entity(authz.ResourceTeamName, event.AggregateID().Deref()).
		Subject(authz.ResourceTeamRoleName, string(domain.TeamMemberRoleCoach)).
		Relate(authz.RelationEditor)
	if e.CreatedBy.OnBehalfOf != nil {
		// If acted on behalf of a person, relate the person as the admin of the team.
		b = b.And().
			Entity(authz.ResourceTeamName, event.AggregateID().Deref()).
			Subject(authz.ResourcePersonName, string(*e.CreatedBy.OnBehalfOf)).
			Relate(authz.RelationTeamAdmin)
	}
	return a.relationStore.AddRelations(ctx, b.Build())
}

func (a *permissionProjector) deleteTeamPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.TeamDeletedEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate the club to the team as the owner.
		Entity(authz.ResourceTeamName, event.AggregateID().Deref()).
		Subject(authz.ResourceClubName, string(e.OwningClubID)).
		Relate(authz.RelationOwner).Build()
	return a.relationStore.RemoveRelations(ctx, relations)
}

func (a *permissionProjector) createClubPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.ClubCreatedEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate the system to the club as the owner.
		Entity(authz.ResourceClubName, e.AggregateID().Deref()).
		Subject(authz.ResourceSystemName, authz.SystemMainID).
		Relate(authz.RelationSystem).
		Build()
	return a.relationStore.AddRelations(ctx, relations)
}

func (a *permissionProjector) createLinkedToPersonPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.AccountLinkedToPersonEvent) error {
	var builder authz.RelationBuilder
	b := builder.
		// Relate the user to the person.
		Entity(authz.ResourcePersonName, string(e.PersonID)).
		Subject(authz.ResourceUserName, event.AggregateID().Deref())
	if e.LinkedAs == domain.AccountLinkParent {
		b = b.Relate(authz.RelationPersonParent)
	} else if e.LinkedAs == domain.AccountLinkSelf {
		b = b.Relate(authz.RelationPersonSelf)
	}
	return a.relationStore.AddRelations(ctx, b.Build())
}

func (a *permissionProjector) createTrainingPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.TrainingScheduledEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate the team to the training as the owner.
		Entity(authz.ResourceTrainingName, e.AggregateID().Deref()).
		Subject(authz.ResourceTeamName, string(e.TeamID)).
		Relate(authz.RelationTrainingTeam).And().
		// Relate the club to the training as the owner.
		Entity(authz.ResourceTrainingName, e.AggregateID().Deref()).
		Subject(authz.ResourceClubName, string(e.OwningClubID)).
		Relate(authz.RelationOwner).
		Build()
	return a.relationStore.AddRelations(ctx, relations)
}

func (a *permissionProjector) createPersonsNominatedForTrainingPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.PersonsNominatedForTrainingEvent) error {
	builder := &authz.RelationBuilder{}
	for _, player := range e.NominatedPlayers {
		builder = builder.
			Entity(authz.ResourceTrainingName, e.AggregateID().Deref()).
			Subject(authz.ResourcePersonName, string(player)).
			Relate(authz.RelationTrainingParticipant).And()
	}
	for _, staff := range e.NominatedStaff {
		builder = builder.
			Entity(authz.ResourceTrainingName, e.AggregateID().Deref()).
			Subject(authz.ResourcePersonName, string(staff)).
			Relate(authz.RelationTrainingParticipant).And()
	}
	return a.relationStore.AddRelations(ctx, builder.Build())
}

func (a *permissionProjector) createClubAdminPermissions(ctx context.Context, event *eventing.JournalEvent, e *domain.ClubAdminAddedEvent) error {
	var builder authz.RelationBuilder
	relations := builder.
		// Relate the user to the club as an admin.
		Entity(authz.ResourceClubName, e.AggregateID().Deref()).
		Subject(authz.ResourceUserName, string(e.AddedUserID)).
		Relate(authz.RelationClubAdmin).
		Build()
	return a.relationStore.AddRelations(ctx, relations)
}
