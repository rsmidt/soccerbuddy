package domain

import (
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"net"
	"time"
)

type (
	// SessionID is the unique identifier for a session, distributed to the client.
	SessionID string

	SessionToken string
)

const (
	SessionAggregateType = eventing.AggregateType("session")

	// SessionLookupToken is the key used to look up a session by the session token.
	SessionLookupToken = "session_token"

	SessionIDUniqueConstraint = "session_id"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type SessionState int

const (
	SessionStateUnspecified SessionState = iota
	SessionStateActive
)

type Session struct {
	eventing.BaseWriter

	State SessionState

	ID        SessionID
	AccountID AccountID
	Role      PrincipalRole
	Token     SessionToken
}

func NewSession(id SessionID) *Session {
	return &Session{
		BaseWriter: *eventing.NewBaseWriter(eventing.AggregateID(id), SessionAggregateType, eventing.VersionMatcherExact),
		ID:         id,
	}
}

func (s *Session) Query() eventing.JournalQuery {
	var builder eventing.JournalQueryBuilder
	return builder.
		WithAggregate(s.Aggregate().AggregateType).
		AggregateID(s.Aggregate().AggregateID).
		Finish().MustBuild()
}

func (s *Session) Reduce(events []*eventing.JournalEvent) {
	for _, event := range events {
		switch e := event.Event.(type) {
		case *SessionCreatedEvent:
			s.State = SessionStateActive
			s.ID = SessionID(e.AggregateID())
			s.Role = e.Role
			s.AccountID = e.AccountID
			s.Token = e.Token
		}
	}
	s.BaseWriter.Reduce(events)
}

func (s *Session) Init(
	token SessionToken,
	accountID AccountID,
	userAgent string,
	ipAddress net.IP,
	validUntil time.Time,
	role PrincipalRole,
) error {
	if s.State != SessionStateUnspecified {
		return NewInvalidAggregateStateError(s.Aggregate(), int(SessionStateUnspecified), int(s.State))
	}
	event := NewSessionCreatedEvent(s.ID, token, accountID, userAgent, ipAddress, validUntil, role)
	s.Append(event)
	return nil
}
