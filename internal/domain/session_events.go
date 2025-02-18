package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"net"
	"time"
)

// ========================================================
// SessionCreatedEvent
// ========================================================

const (
	SessionCreatedEventType    = eventing.EventType("session_created")
	SessionCreatedEventVersion = eventing.EventVersion("v1")
)

var (
	_ eventing.Event                 = (*SessionCreatedEvent)(nil)
	_ eventing.LookupProvider        = (*SessionCreatedEvent)(nil)
	_ eventing.UniqueConstraintAdder = (*SessionCreatedEvent)(nil)
)

type SessionCreatedEvent struct {
	*eventing.EventBase

	Token      SessionToken  `json:"token"`
	AccountID  AccountID     `json:"account_id"`
	UserAgent  string        `json:"user_agent"`
	IPAddress  net.IP        `json:"ip_address"`
	ValidUntil time.Time     `json:"valid_until"`
	Role       PrincipalRole `json:"role"`
}

func NewSessionCreatedEvent(
	id SessionID,
	token SessionToken,
	accountID AccountID,
	userAgent string,
	ipAddress net.IP,
	validUntil time.Time,
	role PrincipalRole,
) *SessionCreatedEvent {
	base := eventing.NewEventBase(eventing.AggregateID(id), SessionAggregateType, SessionCreatedEventVersion, SessionCreatedEventType)

	return &SessionCreatedEvent{
		EventBase:  base,
		Token:      token,
		AccountID:  accountID,
		UserAgent:  userAgent,
		IPAddress:  ipAddress,
		ValidUntil: validUntil,
		Role:       role,
	}
}

func (c *SessionCreatedEvent) IsShredded() bool {
	return false
}

func (c *SessionCreatedEvent) LookupValues() eventing.LookupMap {
	return eventing.LookupMap{
		SessionLookupToken: eventing.LookupFieldValue(c.Token),
	}
}

func (c *SessionCreatedEvent) UniqueConstraintsToAdd() []eventing.UniqueConstraint {
	return []eventing.UniqueConstraint{
		eventing.NewUniqueConstraint(c.AggregateID(), SessionIDUniqueConstraint, string(c.Token)),
	}
}
