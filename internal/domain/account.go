package domain

import (
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

const (
	AccountAggregateType = eventing.AggregateType("account")

	AccountEmailUniqueConstraint = "account_email"
	AccountLookupEmail           = "account_email"
)

var (
	ErrAccountNotFound               = errors.New("account not found")
	ErrRootAccountAlreadyInitialized = errors.New("root account already initialized")
	ErrWrongCredentials              = errors.New("wrong credentials")
	ErrAccountAlreadyLinkedToPerson  = errors.New("already linked to person")
	ErrAccountAlreadyHasSelfLink     = errors.New("already has self link")
)

type (
	AccountID               string
	HashedPassword          string
	InstallationID          string
	NotificationDeviceToken string
)

type AccountState int

const (
	AccountStateUnspecified AccountState = iota
	AccountStateActive
)

type AccountLink string

const (
	AccountLinkParent AccountLink = "parent"
	AccountLinkSelf   AccountLink = "self"
)

type Account struct {
	eventing.BaseWriter

	State AccountState

	ID            AccountID
	FirstName     string
	LastName      string
	LinkedPersons []*AccountLinkedPerson
	Password      HashedPassword

	AppInstallations map[InstallationID]*AppInstallation

	// IsRoot specifies if this the base service account.
	IsRoot bool
}

type AppInstallation struct {
	ID                      InstallationID
	NotificationDeviceToken NotificationDeviceToken
}

type AccountLinkedPerson struct {
	ID       PersonID
	LinkedAs AccountLink
	// Only set if someone other than themselves linked the person.
	LinkedBy *Operator
}

func NewAccount(accountID AccountID) *Account {
	return &Account{
		BaseWriter: *eventing.NewBaseWriter(eventing.AggregateID(accountID), AccountAggregateType, eventing.VersionMatcherExact),
		ID:         accountID,
	}
}

func (a *Account) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.WithAggregate(AccountAggregateType).
		AggregateID(eventing.AggregateID(a.ID)).
		Finish().MustBuild()
}

func (a *Account) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *RootAccountCreatedEvent:
			a.State = AccountStateActive
			a.Password = e.HashedPassword
			a.IsRoot = true
			a.AppInstallations = make(map[InstallationID]*AppInstallation)
		case *AccountCreatedEvent:
			a.State = AccountStateActive
			a.FirstName = e.FirstName.Value
			a.LastName = e.LastName.Value
			a.Password = HashedPassword(e.HashedPassword.Value)
			a.AppInstallations = make(map[InstallationID]*AppInstallation)
		case *AccountLinkedToPersonEvent:
			a.LinkedPersons = append(a.LinkedPersons, &AccountLinkedPerson{
				ID:       e.PersonID,
				LinkedAs: e.LinkedAs,
				LinkedBy: e.LinkedBy,
			})
		case *MobileDeviceAttachedToAccountEvent:
			a.AppInstallations[e.InstallationID] = &AppInstallation{
				ID:                      e.InstallationID,
				NotificationDeviceToken: e.NotificationDeviceToken,
			}
		case *AccountNotificationDeviceTokenChangedEvent:
			a.AppInstallations[e.InstallationID].NotificationDeviceToken = e.NotificationDeviceToken
		}
	}
	a.BaseWriter.Reduce(events)
}

func (a *Account) InitAsRoot(email string, password HashedPassword, firstName, lastName string) error {
	if a.State != AccountStateUnspecified {
		return NewInvalidAggregateStateError(a.Aggregate(), int(AccountStateUnspecified), int(a.State))
	}
	accountID := AccountID(a.Aggregate().AggregateID)
	a.Append(NewRootAccountCreatedEvent(accountID, email, password, firstName, lastName))
	return nil
}

func (a *Account) Init(firstName, lastName, email string, password HashedPassword) error {
	if a.State != AccountStateUnspecified {
		return NewInvalidAggregateStateError(a.Aggregate(), int(AccountStateUnspecified), int(a.State))
	}
	a.Append(NewAccountCreatedEvent(a.ID, firstName, lastName, email, password))
	return nil
}

func (a *Account) Link(id PersonID, linkAs AccountLink, linkedBy *Operator, clubID ClubID) error {
	if a.State != AccountStateActive {
		return NewInvalidAggregateStateError(a.Aggregate(), int(AccountStateActive), int(a.State))
	}
	var existingLinkForPerson *AccountLinkedPerson
	for _, link := range a.LinkedPersons {
		if link.ID == id {
			existingLinkForPerson = link
			break
		}
	}
	if existingLinkForPerson != nil {
		return ErrAccountAlreadyLinkedToPerson
	}
	var existingSelfLink *AccountLinkedPerson
	for _, link := range a.LinkedPersons {
		if link.LinkedAs == AccountLinkSelf {
			existingSelfLink = link
			break
		}
	}
	if linkAs == AccountLinkSelf && existingSelfLink != nil {
		return ErrAccountAlreadyHasSelfLink
	}
	event := NewAccountLinkedToPersonEvent(a.ID, id, linkAs, linkedBy, clubID)
	a.Append(event)
	return nil
}

func (a *Account) VerifyPassword(password string, verifier PasswordVerifier) (bool, error) {
	if a.State != AccountStateActive {
		return false, NewInvalidAggregateStateError(a.Aggregate(), int(AccountStateActive), int(a.State))
	}
	return verifier.Verify(password, a.Password)
}

func (a *Account) AttachMobileDevice(id InstallationID, token NotificationDeviceToken) error {
	if a.State != AccountStateActive {
		return NewInvalidAggregateStateError(a.Aggregate(), int(AccountStateActive), int(a.State))
	}
	if existingInstallation, ok := a.AppInstallations[id]; ok {
		if existingInstallation.NotificationDeviceToken != token {
			a.Append(NewAccountNotificationDeviceTokenChangedEvent(a.ID, id, token))
			return nil
		}
		// Existing installation, but no token change. Nothing to do.
		return nil
	}
	a.Append(NewMobileDeviceAttachedToAccountEvent(a.ID, id, token))
	return nil
}
