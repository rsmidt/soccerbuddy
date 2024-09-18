package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPerson_Init(t *testing.T) {
	personID := idgen.New[PersonID]()
	clubID := idgen.New[ClubID]()
	creator := NewOperator(idgen.New[AccountID](), nil)
	birthdate1 := time.Now()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		firstName     string
		lastName      string
		birthdate     time.Time
		creator       Operator
		owningClubID  ClubID
	}{
		{
			name:          "Succeeds if person is initialized correctly",
			initialEvents: createInitialEvents(),
			emittedEvents: []eventing.Event{
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate1, creator, clubID),
			},
			firstName:    "John",
			lastName:     "Doe",
			birthdate:    birthdate1,
			creator:      creator,
			owningClubID: clubID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			person := NewPerson(personID)
			person.Reduce(tt.initialEvents)
			person.Init(tt.firstName, tt.lastName, tt.birthdate, tt.creator, tt.owningClubID)
			assert.Equal(t, tt.emittedEvents, person.Changes().Events())
		})
	}
}

func TestPerson_Reduce(t *testing.T) {
	personID := idgen.New[PersonID]()
	clubID := idgen.New[ClubID]()
	creator := NewOperator(idgen.New[AccountID](), nil)
	birthdate1 := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		expectedState PersonState
		expectedName  string
	}{
		{
			name: "Succeeds if person state is updated correctly",
			initialEvents: createInitialEvents(
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate1, creator, clubID),
			),
			expectedState: PersonStateActive,
			expectedName:  "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			person := NewPerson(personID)
			person.Reduce(tt.initialEvents)
			assert.Equal(t, tt.expectedState, person.State)
			assert.Equal(t, tt.expectedName, person.Firstname+" "+person.Lastname)
		})
	}
}

func TestPerson_InitiateNewLink(t *testing.T) {
	personID := idgen.New[PersonID]()
	clubID := idgen.New[ClubID]()
	creator := NewOperator(idgen.New[AccountID](), nil)
	token := PersonLinkToken("token")
	expiresAt := time.Now().Add(24 * time.Hour)
	birthdate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		linkAs        AccountLink
		token         PersonLinkToken
		expiresAt     time.Time
		expectedError error
	}{
		{
			name: "Succeeds if new link is initiated correctly",
			initialEvents: createInitialEvents(
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate, creator, clubID),
			),
			emittedEvents: []eventing.Event{
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkParent, token, expiresAt),
			},
			linkAs:        AccountLinkParent,
			token:         token,
			expiresAt:     expiresAt,
			expectedError: nil,
		},
		{
			name: "Fails if too many pending links",
			initialEvents: createInitialEvents(
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate, creator, clubID),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkParent, "token1", expiresAt),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkParent, "token2", expiresAt),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkParent, "token3", expiresAt),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkParent, "token4", expiresAt),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkParent, "token5", expiresAt),
			),
			linkAs:        AccountLinkParent,
			token:         token,
			expiresAt:     expiresAt,
			expectedError: ErrPersonHasTooManyPendingLinks,
		},
		{
			name: "Fails if person already self linked",
			initialEvents: createInitialEvents(
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate, creator, clubID),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkSelf, "token1", expiresAt),
				NewPersonLinkClaimedEvent(personID, "account1", AccountLinkSelf, "token1"),
			),
			linkAs:        AccountLinkSelf,
			token:         token,
			expiresAt:     expiresAt,
			expectedError: ErrPersonAlreadySelfLinked,
		},
		{
			name: "Fails if token is reused",
			initialEvents: createInitialEvents(
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate, creator, clubID),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkSelf, token, expiresAt),
			),
			linkAs:        AccountLinkSelf,
			token:         token,
			expiresAt:     expiresAt,
			expectedError: ErrPersonInvalidLinkToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			person := NewPerson(personID)
			person.Reduce(tt.initialEvents)
			err := person.InitiateNewLink(creator, tt.linkAs, tt.token, tt.expiresAt)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, person.Changes().Events())
		})
	}
}

func TestPerson_Claim(t *testing.T) {
	personID := idgen.New[PersonID]()
	clubID := idgen.New[ClubID]()
	creator := NewOperator(idgen.New[AccountID](), nil)
	token := PersonLinkToken("token")
	expiresAt := time.Now().Add(24 * time.Hour)
	accountID := idgen.New[AccountID]()
	birthdate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		token         PersonLinkToken
		accountID     AccountID
		expectedError error
	}{
		{
			name: "Succeeds if link is claimed correctly",
			initialEvents: createInitialEvents(
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate, creator, clubID),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkParent, token, expiresAt),
			),
			emittedEvents: []eventing.Event{
				NewPersonLinkClaimedEvent(personID, accountID, AccountLinkParent, token),
			},
			token:         token,
			accountID:     accountID,
			expectedError: nil,
		},
		{
			name: "Fails if token is invalid",
			initialEvents: createInitialEvents(
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate, creator, clubID),
			),
			token:         token,
			accountID:     accountID,
			expectedError: ErrPersonInvalidLinkToken,
		},
		{
			name: "Fails if token is expired",
			initialEvents: createInitialEvents(
				NewPersonCreatedEvent(personID, "John", "Doe", birthdate, creator, clubID),
				NewPersonLinkInitiatedEvent(personID, creator, AccountLinkParent, token, time.Now().Add(-24*time.Hour)),
			),
			token:         token,
			accountID:     accountID,
			expectedError: ErrPersonLinkTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			person := NewPerson(personID)
			person.Reduce(tt.initialEvents)
			err := person.Claim(tt.token, tt.accountID)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, person.Changes().Events())
		})
	}
}
