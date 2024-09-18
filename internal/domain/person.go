package domain

import (
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/core"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
)

type (
	PersonID        string
	PersonLinkToken string
)

const (
	PersonAggregateType = eventing.AggregateType("person")

	PersonMaxPendingLinks = 5
)

var (
	PersonMaxBirthdate = time.Now()
	PersonMinBirthdate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	ErrOwningClubNotFound           = errors.New("owning club not found")
	ErrPersonNotFound               = errors.New("person not found")
	ErrPersonAlreadySelfLinked      = errors.New("person already self linked")
	ErrPersonHasTooManyPendingLinks = errors.New("too many pending links for person")
	ErrPersonInvalidLinkToken       = errors.New("invalid link token")
	ErrPersonLinkTokenExpired       = errors.New("link token has expired")
)

type PersonState int

const (
	PersonStateUnspecified PersonState = iota
	PersonStateActive
)

type PersonLinkedAccount struct {
	AccountID AccountID
	LinkedAs  AccountLink
	LinkedAt  time.Time
}

type PendingLink struct {
	LinkAs    AccountLink
	Token     PersonLinkToken
	ExpiresAt time.Time
}

type Person struct {
	eventing.BaseWriter

	ID             PersonID
	Firstname      string
	Lastname       string
	Birthdate      time.Time
	State          PersonState
	OwningClubID   ClubID
	LinkedAccounts []PersonLinkedAccount
	PendingLinks   map[PersonLinkToken]PendingLink
	Creator        Operator
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewPerson(id PersonID) *Person {
	return &Person{
		BaseWriter:   *eventing.NewBaseWriter(eventing.AggregateID(id), PersonAggregateType, eventing.VersionMatcherExact),
		ID:           id,
		PendingLinks: map[PersonLinkToken]PendingLink{},
	}
}

func (p *Person) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.WithAggregate(PersonAggregateType).
		AggregateID(eventing.AggregateID(p.ID)).
		Finish().MustBuild()
}

func (p *Person) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *PersonCreatedEvent:
			p.State = PersonStateActive
			p.Firstname = e.FirstName.Value
			p.Lastname = e.LastName.Value
			p.Birthdate = core.Must2(time.Parse(time.RFC3339, e.Birthdate.Value))
			p.OwningClubID = e.OwningClubID
			p.CreatedAt = event.InsertedAt()
			p.UpdatedAt = event.InsertedAt()
		case *PersonLinkInitiatedEvent:
			p.PendingLinks[e.Token] = PendingLink{
				LinkAs:    e.LinkAs,
				ExpiresAt: e.ExpiresAt,
				Token:     e.Token,
			}
		case *PersonLinkClaimedEvent:
			// Remove the pending link.
			delete(p.PendingLinks, e.UsedToken)
			p.LinkedAccounts = append(p.LinkedAccounts, PersonLinkedAccount{
				AccountID: e.AccountID,
				LinkedAs:  e.LinkedAs,
				LinkedAt:  event.InsertedAt(),
			})
		}
	}
	p.BaseWriter.Reduce(events)
}

func (p *Person) Init(firstName, lastName string, birthdate time.Time, creator Operator, owningClubID ClubID) {
	p.Firstname = firstName
	p.Lastname = lastName
	p.Birthdate = birthdate
	p.OwningClubID = owningClubID
	p.State = PersonStateActive
	event := NewPersonCreatedEvent(p.ID, firstName, lastName, birthdate, creator, owningClubID)
	p.Append(event)
}

func (p *Person) InitiateNewLink(operator Operator, linkAs AccountLink, token PersonLinkToken, expiresAt time.Time) error {
	// Only allow links for active persons.
	if p.State != PersonStateActive {
		return NewInvalidAggregateStateError(p.Aggregate(), int(PersonStateActive), int(p.State))
	}

	// Don't allow more than N pending links at the same time.
	if len(p.PendingLinks) >= PersonMaxPendingLinks {
		return ErrPersonHasTooManyPendingLinks
	}

	// Determine if there's an existing self link. We only allow one.
	if linkAs == AccountLinkSelf {
		var hasExistingSelfLink bool
		for _, link := range p.LinkedAccounts {
			if link.LinkedAs == AccountLinkSelf {
				hasExistingSelfLink = true
				break
			}
		}
		if hasExistingSelfLink {
			return ErrPersonAlreadySelfLinked
		}
	}

	// If there's already a pending link with the same token, don't allow it.
	if _, ok := p.PendingLinks[token]; ok {
		return ErrPersonInvalidLinkToken
	}

	event := NewPersonLinkInitiatedEvent(p.ID, operator, linkAs, token, expiresAt)
	p.Append(event)
	return nil
}

func (p *Person) Claim(token PersonLinkToken, id AccountID) error {
	// Only allow links for active persons.
	if p.State != PersonStateActive {
		return NewInvalidAggregateStateError(p.Aggregate(), int(PersonStateActive), int(p.State))
	}

	pl, ok := p.PendingLinks[token]
	if !ok {
		return ErrPersonInvalidLinkToken
	} else if pl.ExpiresAt.Before(time.Now()) {
		return ErrPersonLinkTokenExpired
	}
	event := NewPersonLinkClaimedEvent(p.ID, id, pl.LinkAs, token)
	p.Append(event)
	return nil
}

func (p *Person) FindPendingLink(token PersonLinkToken) (PendingLink, error) {
	for _, link := range p.PendingLinks {
		if link.Token == token {
			return link, nil
		}
	}
	return PendingLink{}, ErrPersonInvalidLinkToken
}
