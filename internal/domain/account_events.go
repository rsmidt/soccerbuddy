package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
)

// ========================================================
// AccountCreatedEvent
// ========================================================

const (
	AccountCreatedEventType    = eventing.EventType("account_created")
	AccountCreatedEventVersion = eventing.EventVersion("v1")
)

var _ eventing.Event = (*AccountCreatedEvent)(nil)

type AccountCreatedEvent struct {
	*eventing.EventBase

	FirstName eventing.EncryptedString `json:"first_name"`
	LastName  eventing.EncryptedString `json:"last_name"`

	Email          eventing.EncryptedString `json:"email"`
	HashedPassword eventing.EncryptedString `json:"hashed_password"`
}

func NewAccountCreatedEvent(id AccountID, firstName, lastName, email string, hashedPassword HashedPassword) *AccountCreatedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), AccountAggregateType, AccountCreatedEventVersion, AccountCreatedEventType)

	return &AccountCreatedEvent{
		EventBase:      base,
		FirstName:      eventing.NewEncryptedString(firstName),
		LastName:       eventing.NewEncryptedString(lastName),
		Email:          eventing.NewEncryptedString(email),
		HashedPassword: eventing.NewEncryptedString(string(hashedPassword)),
	}
}

func (r *AccountCreatedEvent) IsShredded() bool {
	return r.FirstName.IsShredded || r.LastName.IsShredded || r.Email.IsShredded || r.HashedPassword.IsShredded
}

func (r *AccountCreatedEvent) UniqueConstraintsToAdd() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewUniqueConstraint(r.AggregateID(), AccountEmailUniqueConstraint, r.Email.Value),
	}
}

func (r *AccountCreatedEvent) LookupValues() eventing.LookupMap {
	return eventing.LookupMap{
		AccountLookupEmail: eventing.LookupFieldValue(r.Email.Value),
	}
}

func (r *AccountCreatedEvent) DeclareOwners() []eventing.AggregateID {
	return []eventing.AggregateID{r.AggregateID()}
}

func (r *AccountCreatedEvent) AcceptCrypto(transformer eventing.CryptoTransformer) error {
	if err := transformer.Transform(r.AggregateID(), &r.HashedPassword); err != nil {
		return err
	}
	if err := transformer.TransformWithDefault(r.AggregateID(), &r.FirstName, RedactedString); err != nil {
		return err
	}
	if err := transformer.TransformWithDefault(r.AggregateID(), &r.LastName, RedactedString); err != nil {
		return err
	}
	if err := transformer.TransformWithDefault(r.AggregateID(), &r.Email, RedactedString); err != nil {
		return err
	}
	return nil
}

// ========================================================
// AccountLinkedToPersonEvent
// ========================================================

const (
	AccountLinkedToPersonEventType    = eventing.EventType("account_linked_to_person")
	AccountLinkedToPersonEventVersion = eventing.EventVersion("v1")
)

var _ eventing.Event = (*AccountLinkedToPersonEvent)(nil)

type AccountLinkedToPersonEvent struct {
	*eventing.EventBase

	PersonID      PersonID         `json:"person_id"`
	LinkedAs      AccountLink      `json:"linked_as"`
	LinkedBy      *Operator        `json:"linked_by"`
	OwningClubID  ClubID           `json:"owning_club_id"`
	UsedLinkToken *PersonLinkToken `json:"used_link_token"`
}

func NewAccountLinkedToPersonEvent(id AccountID, personID PersonID, linkAs AccountLink, linkedBy *Operator, clubID ClubID, usedToken *PersonLinkToken) *AccountLinkedToPersonEvent {
	// There's unfortunately an import cycle with person, so only stringed types.
	base := eventing.NewEventBase(eventing.AggregateID(id), AccountAggregateType, AccountLinkedToPersonEventVersion, AccountLinkedToPersonEventType)

	return &AccountLinkedToPersonEvent{
		EventBase:     base,
		PersonID:      personID,
		LinkedAs:      linkAs,
		LinkedBy:      linkedBy,
		OwningClubID:  clubID,
		UsedLinkToken: usedToken,
	}
}

func (l *AccountLinkedToPersonEvent) IsShredded() bool {
	return false
}

// ========================================================
// RootAccountCreatedEvent
// ========================================================

const (
	RootAccountCreatedEventType    = eventing.EventType("root_account_created")
	RootAccountCreatedEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event                 = (*RootAccountCreatedEvent)(nil)
	_ eventing.UniqueConstraintAdder = (*RootAccountCreatedEvent)(nil)
	_ eventing.LookupProvider        = (*RootAccountCreatedEvent)(nil)
)

type RootAccountCreatedEvent struct {
	*eventing.EventBase

	Email          string         `json:"email"`
	HashedPassword HashedPassword `json:"hashed_password"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
}

func NewRootAccountCreatedEvent(id AccountID, email string, hashedPassword HashedPassword, firstName, lastName string) *RootAccountCreatedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), AccountAggregateType, RootAccountCreatedEventVersion, RootAccountCreatedEventType)

	return &RootAccountCreatedEvent{
		EventBase:      base,
		Email:          email,
		HashedPassword: hashedPassword,
		FirstName:      firstName,
		LastName:       lastName,
	}
}

func (r *RootAccountCreatedEvent) IsShredded() bool {
	return false
}

func (r *RootAccountCreatedEvent) UniqueConstraintsToAdd() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewUniqueConstraint(r.AggregateID(), AccountEmailUniqueConstraint, r.Email),
	}
}

func (r *RootAccountCreatedEvent) LookupValues() eventing.LookupMap {
	return eventing.LookupMap{
		AccountLookupEmail: eventing.LookupFieldValue(r.Email),
	}
}

// ========================================================
// MobileDeviceAttachedToAccountEvent
// ========================================================

const (
	MobileDeviceAttachedToAccountEventType    = eventing.EventType("mobile_device_attached_to_account")
	MobileDeviceAttachedToAccountEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event = (*MobileDeviceAttachedToAccountEvent)(nil)
)

type MobileDeviceAttachedToAccountEvent struct {
	*eventing.EventBase

	InstallationID          InstallationID          `json:"installation_id"`
	NotificationDeviceToken NotificationDeviceToken `json:"notification_device_token"`
}

func NewMobileDeviceAttachedToAccountEvent(id AccountID, installationID InstallationID, notificationDeviceToken NotificationDeviceToken) *MobileDeviceAttachedToAccountEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), AccountAggregateType, MobileDeviceAttachedToAccountEventVersion, MobileDeviceAttachedToAccountEventType)

	return &MobileDeviceAttachedToAccountEvent{
		EventBase:               base,
		InstallationID:          installationID,
		NotificationDeviceToken: notificationDeviceToken,
	}
}

func (r *MobileDeviceAttachedToAccountEvent) IsShredded() bool {
	return false
}

// ========================================================
// AccountNotificationDeviceTokenChangedEvent
// ========================================================

const (
	AccountNotificationDeviceTokenChangedEventType    = eventing.EventType("account_notification_device_token_changed")
	AccountNotificationDeviceTokenChangedEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event = (*AccountNotificationDeviceTokenChangedEvent)(nil)
)

type AccountNotificationDeviceTokenChangedEvent struct {
	*eventing.EventBase

	InstallationID          InstallationID          `json:"installation_id"`
	NotificationDeviceToken NotificationDeviceToken `json:"notification_device_token"`
}

func NewAccountNotificationDeviceTokenChangedEvent(id AccountID, installationID InstallationID, notificationDeviceToken NotificationDeviceToken) *AccountNotificationDeviceTokenChangedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), AccountAggregateType, AccountNotificationDeviceTokenChangedEventVersion, AccountNotificationDeviceTokenChangedEventType)

	return &AccountNotificationDeviceTokenChangedEvent{
		EventBase:               base,
		InstallationID:          installationID,
		NotificationDeviceToken: notificationDeviceToken,
	}
}

func (r *AccountNotificationDeviceTokenChangedEvent) IsShredded() bool {
	return false
}

// ========================================================
// AccountRegisteredEvent
// ========================================================

const (
	AccountRegisteredEventType    = eventing.EventType("account_registered")
	AccountRegisteredEventVersion = eventing.EventVersion("v1")
)

var _ eventing.Event = (*AccountRegisteredEvent)(nil)

type AccountRegisteredEvent struct {
	*eventing.EventBase

	FirstName eventing.EncryptedString `json:"first_name"`
	LastName  eventing.EncryptedString `json:"last_name"`

	Email          eventing.EncryptedString `json:"email"`
	HashedPassword eventing.EncryptedString `json:"hashed_password"`
	UsedLinkToken  PersonLinkToken          `json:"link_token"`
}

func NewAccountRegisteredEvent(id AccountID, firstName, lastName, email string, hashedPassword HashedPassword, usedLinkToken PersonLinkToken) *AccountRegisteredEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), AccountAggregateType, AccountRegisteredEventVersion, AccountRegisteredEventType)

	return &AccountRegisteredEvent{
		EventBase:      base,
		FirstName:      eventing.NewEncryptedString(firstName),
		LastName:       eventing.NewEncryptedString(lastName),
		Email:          eventing.NewEncryptedString(email),
		HashedPassword: eventing.NewEncryptedString(string(hashedPassword)),
		UsedLinkToken:  usedLinkToken,
	}
}

func (r *AccountRegisteredEvent) IsShredded() bool {
	return r.FirstName.IsShredded || r.LastName.IsShredded || r.Email.IsShredded || r.HashedPassword.IsShredded
}

func (r *AccountRegisteredEvent) UniqueConstraintsToAdd() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewUniqueConstraint(r.AggregateID(), AccountEmailUniqueConstraint, r.Email.Value),
		eventing.NewUniqueConstraint(r.AggregateID(), AccountUsedLinkTokenUniqueConstraint, string(r.UsedLinkToken)),
	}
}

func (r *AccountRegisteredEvent) LookupValues() eventing.LookupMap {
	return eventing.LookupMap{
		AccountLookupEmail: eventing.LookupFieldValue(r.Email.Value),
	}
}

func (r *AccountRegisteredEvent) DeclareOwners() []eventing.AggregateID {
	return []eventing.AggregateID{r.AggregateID()}
}

func (r *AccountRegisteredEvent) AcceptCrypto(transformer eventing.CryptoTransformer) error {
	if err := transformer.Transform(r.AggregateID(), &r.HashedPassword); err != nil {
		return err
	}
	if err := transformer.TransformWithDefault(r.AggregateID(), &r.FirstName, RedactedString); err != nil {
		return err
	}
	if err := transformer.TransformWithDefault(r.AggregateID(), &r.LastName, RedactedString); err != nil {
		return err
	}
	if err := transformer.TransformWithDefault(r.AggregateID(), &r.Email, RedactedString); err != nil {
		return err
	}
	return nil
}
