package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
)

// ========================================================
// PersonCreatedEvent
// ========================================================

const (
	PersonCreatedEventType    = eventing.EventType("person_created")
	PersonCreatedEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event          = (*PersonCreatedEvent)(nil)
	_ eventing.EncryptedEvent = (*PersonCreatedEvent)(nil)
)

type PersonCreatedEvent struct {
	*eventing.EventBase

	FirstName eventing.EncryptedString `json:"firstname"`
	LastName  eventing.EncryptedString `json:"lastname"`
	Birthdate eventing.EncryptedString `json:"birthdate"`

	Creator      Operator `json:"creator"`
	OwningClubID ClubID   `json:"owning_club_id"`
}

func NewPersonCreatedEvent(id PersonID, firstName, lastName string, birthdate time.Time, creator Operator, owningClubID ClubID) *PersonCreatedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), PersonAggregateType, PersonCreatedEventVersion, PersonCreatedEventType)

	return &PersonCreatedEvent{
		EventBase:    base,
		FirstName:    eventing.NewEncryptedString(firstName),
		LastName:     eventing.NewEncryptedString(lastName),
		Birthdate:    eventing.NewEncryptedString(birthdate.Format(time.RFC3339)),
		Creator:      creator,
		OwningClubID: owningClubID,
	}
}

func (p *PersonCreatedEvent) IsShredded() bool {
	return p.FirstName.IsShredded || p.LastName.IsShredded || p.Birthdate.IsShredded
}

func (p *PersonCreatedEvent) DeclareOwners() []eventing.AggregateID {
	return []eventing.AggregateID{p.AggregateID()}
}

func (p *PersonCreatedEvent) AcceptCrypto(transformer eventing.CryptoTransformer) error {
	if err := transformer.TransformWithDefault(p.AggregateID(), &p.FirstName, RedactedString); err != nil {
		return err
	}
	if err := transformer.TransformWithDefault(p.AggregateID(), &p.LastName, RedactedString); err != nil {
		return err
	}
	if err := transformer.TransformWithDefault(p.AggregateID(), &p.Birthdate, time.Time{}.Format(time.RFC3339)); err != nil {
		return err
	}
	return nil
}

// ========================================================
// PersonLinkClaimedEvent
// ========================================================

const (
	PersonLinkClaimedEventType    = eventing.EventType("person_link_claimed")
	PersonLinkClaimedEventVersion = eventing.EventVersion("v1")
)

var _ eventing.Event = (*PersonLinkClaimedEvent)(nil)

type PersonLinkClaimedEvent struct {
	*eventing.EventBase

	AccountID AccountID       `json:"account_id"`
	LinkedAs  AccountLink     `json:"linked_as"`
	UsedToken PersonLinkToken `json:"used_token"`
}

func NewPersonLinkClaimedEvent(id PersonID, accountID AccountID, as AccountLink, usedToken PersonLinkToken) *PersonLinkClaimedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), PersonAggregateType, PersonLinkClaimedEventVersion, PersonLinkClaimedEventType)

	return &PersonLinkClaimedEvent{
		EventBase: base,
		AccountID: accountID,
		LinkedAs:  as,
		UsedToken: usedToken,
	}
}

func (l *PersonLinkClaimedEvent) IsShredded() bool {
	return false
}

// ========================================================
// PersonLinkInitiatedEvent
// ========================================================

const (
	PersonLinkInitiatedEventType    = eventing.EventType("person_link_initiated")
	PersonLinkInitiatedEventVersion = eventing.EventVersion("v1")
)

var _ eventing.Event = (*PersonLinkInitiatedEvent)(nil)

type PersonLinkInitiatedEvent struct {
	*eventing.EventBase

	LinkAs    AccountLink     `json:"link_as"`
	Token     PersonLinkToken `json:"token"`
	ExpiresAt time.Time       `json:"expires_at"`
	InvitedBy Operator        `json:"invited_by"`
}

func NewPersonLinkInitiatedEvent(id PersonID, invitedBy Operator, linkAs AccountLink, token PersonLinkToken, expiresAt time.Time) *PersonLinkInitiatedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), PersonAggregateType, PersonLinkInitiatedEventVersion, PersonLinkInitiatedEventType)

	return &PersonLinkInitiatedEvent{
		EventBase: base,
		LinkAs:    linkAs,
		Token:     token,
		ExpiresAt: expiresAt,
		InvitedBy: invitedBy,
	}
}

func (l *PersonLinkInitiatedEvent) IsShredded() bool {
	return false
}
