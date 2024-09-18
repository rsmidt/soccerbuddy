package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSession_Init(t *testing.T) {
	sessionID := idgen.New[SessionID]()
	accountID := idgen.New[AccountID]()
	token := SessionToken("token")
	userAgent := "Mozilla/5.0"
	ipAddress := "192.168.1.1"
	validUntil := time.Now().Add(24 * time.Hour)
	role := PrincipalRoleRegular

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		token         SessionToken
		accountID     AccountID
		userAgent     string
		ipAddress     string
		validUntil    time.Time
		role          PrincipalRole
		expectedError error
	}{
		{
			name:          "Succeeds if session is initialized correctly",
			initialEvents: createInitialEvents(),
			emittedEvents: []eventing.Event{
				NewSessionCreatedEvent(sessionID, token, accountID, userAgent, ipAddress, validUntil, role),
			},
			token:         token,
			accountID:     accountID,
			userAgent:     userAgent,
			ipAddress:     ipAddress,
			validUntil:    validUntil,
			role:          role,
			expectedError: nil,
		},
		{
			name: "Fails if session is already initialized",
			initialEvents: createInitialEvents(
				NewSessionCreatedEvent(sessionID, token, accountID, userAgent, ipAddress, validUntil, role),
			),
			token:         token,
			accountID:     accountID,
			userAgent:     userAgent,
			ipAddress:     ipAddress,
			validUntil:    validUntil,
			role:          role,
			expectedError: NewInvalidAggregateStateError(NewSession(sessionID).Aggregate(), int(SessionStateUnspecified), int(SessionStateActive)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			session := NewSession(sessionID)
			session.Reduce(tt.initialEvents)
			err := session.Init(tt.token, tt.accountID, tt.userAgent, tt.ipAddress, tt.validUntil, tt.role)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, session.Changes().Events())
		})
	}
}

func TestSession_Reduce(t *testing.T) {
	sessionID := idgen.New[SessionID]()
	accountID := idgen.New[AccountID]()
	token := SessionToken("token")
	userAgent := "Mozilla/5.0"
	ipAddress := "192.168.1.1"
	validUntil := time.Now().Add(24 * time.Hour)
	role := PrincipalRoleRegular

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		expectedState SessionState
		expectedToken SessionToken
		expectedRole  PrincipalRole
		expectedID    SessionID
		expectedAccID AccountID
	}{
		{
			name: "Succeeds if session state is updated correctly",
			initialEvents: createInitialEvents(
				NewSessionCreatedEvent(sessionID, token, accountID, userAgent, ipAddress, validUntil, role),
			),
			expectedState: SessionStateActive,
			expectedToken: token,
			expectedRole:  role,
			expectedID:    sessionID,
			expectedAccID: accountID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			session := NewSession(sessionID)
			session.Reduce(tt.initialEvents)
			assert.Equal(t, tt.expectedState, session.State)
			assert.Equal(t, tt.expectedToken, session.Token)
			assert.Equal(t, tt.expectedRole, session.Role)
			assert.Equal(t, tt.expectedID, session.ID)
			assert.Equal(t, tt.expectedAccID, session.AccountID)
		})
	}
}
