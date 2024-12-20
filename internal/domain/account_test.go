package domain

import (
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount_InitAsRoot(t *testing.T) {
	accID := idgen.New[AccountID]()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		email         string
		password      HashedPassword
		firstName     string
		lastName      string
		expectedError error
	}{
		{
			name:          "InitAsRoot success",
			initialEvents: createInitialEvents(),
			emittedEvents: []eventing.Event{
				NewRootAccountCreatedEvent(accID, "root@example.com", "hashedpassword", "Jane", "Doe"),
			},
			email:         "root@example.com",
			password:      "hashedpassword",
			firstName:     "Jane",
			lastName:      "Doe",
			expectedError: nil,
		},
		{
			name: "InitAsRoot already initialized",
			initialEvents: createInitialEvents(
				NewRootAccountCreatedEvent(accID, "root@example.com", "hashedpassword", "Jane", "Doe"),
			),
			email:         "root@example.com",
			password:      "hashedpassword",
			firstName:     "Jane",
			lastName:      "Doe",
			expectedError: NewInvalidAggregateStateError(NewAccount(accID).Aggregate(), int(AccountStateUnspecified), int(AccountStateActive)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account := NewAccount(accID)
			account.Reduce(tt.initialEvents)
			err := account.InitAsRoot(tt.email, tt.password, tt.firstName, tt.lastName)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, account.Changes().Events())
		})
	}
}

func TestAccount_Init(t *testing.T) {
	accID := idgen.New[AccountID]()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		firstName     string
		lastName      string
		email         string
		password      HashedPassword
		expectedError error
	}{
		{
			name:          "Init success",
			initialEvents: createInitialEvents(),
			emittedEvents: []eventing.Event{
				NewAccountCreatedEvent(accID, "John", "Doe", "john@example.com", "hashedpassword"),
			},
			firstName:     "John",
			lastName:      "Doe",
			email:         "john@example.com",
			password:      "hashedpassword",
			expectedError: nil,
		},
		{
			name: "Init already initialized",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "john@example.com", "hashedpassword"),
			),
			firstName:     "John",
			lastName:      "Doe",
			email:         "john@example.com",
			password:      "hashedpassword",
			expectedError: NewInvalidAggregateStateError(NewAccount(accID).Aggregate(), int(AccountStateUnspecified), int(AccountStateActive)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account := NewAccount(accID)
			account.Reduce(tt.initialEvents)
			err := account.Init(tt.firstName, tt.lastName, tt.email, tt.password)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, account.Changes().Events())
		})
	}
}

func TestAccount_Link(t *testing.T) {
	accID := idgen.New[AccountID]()
	per1ID := idgen.New[PersonID]()
	per2ID := idgen.New[PersonID]()
	club1ID := idgen.New[ClubID]()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		personID      PersonID
		linkAs        AccountLink
		linkedBy      *Operator
		clubID        ClubID
		expectedError error
	}{
		{
			name: "Link success",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "password", "password"),
			),
			emittedEvents: []eventing.Event{
				NewAccountLinkedToPersonEvent(accID, per1ID, AccountLinkParent, nil, club1ID),
			},
			personID:      per1ID,
			linkAs:        AccountLinkParent,
			linkedBy:      nil,
			clubID:        club1ID,
			expectedError: nil,
		},
		{
			name: "Multiple links to different persons succeed",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "test@example.com", "password"),
				NewAccountLinkedToPersonEvent(accID, per1ID, AccountLinkParent, nil, club1ID),
			),
			emittedEvents: []eventing.Event{
				NewAccountLinkedToPersonEvent(accID, per2ID, AccountLinkParent, nil, club1ID),
			},
			personID:      per2ID,
			linkAs:        AccountLinkParent,
			linkedBy:      nil,
			clubID:        club1ID,
			expectedError: nil,
		},
		{
			name: "Multiple links to same person are forbidden",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "test@example.com", "password"),
				NewAccountLinkedToPersonEvent(accID, per1ID, AccountLinkParent, nil, club1ID),
			),
			personID:      per1ID,
			linkAs:        AccountLinkParent,
			linkedBy:      nil,
			clubID:        club1ID,
			expectedError: ErrAccountAlreadyLinkedToPerson,
		},
		{
			name: "Multiple self links fail",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "test@example.com", "password"),
				NewAccountLinkedToPersonEvent(accID, per1ID, AccountLinkSelf, nil, club1ID),
			),
			personID:      "person2",
			linkAs:        AccountLinkSelf,
			linkedBy:      nil,
			clubID:        club1ID,
			expectedError: ErrAccountAlreadyHasSelfLink,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account := NewAccount(accID)
			account.Reduce(tt.initialEvents)
			err := account.Link(tt.personID, tt.linkAs, tt.linkedBy, tt.clubID)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, account.Changes().Events())
		})
	}
}

func TestAccount_VerifyPassword(t *testing.T) {
	accID := idgen.New[AccountID]()

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		password      string
		expectedError error
		expectedValid bool
	}{
		{
			name: "Argon2idVerifyPassword success",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "password", "password"),
			),
			password:      "password",
			expectedError: nil,
			expectedValid: true,
		},
		{
			name: "Argon2idVerifyPassword wrong password",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "password", "password"),
			),
			password:      "wrongpassword",
			expectedError: nil,
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account := NewAccount(accID)
			account.Reduce(tt.initialEvents)
			valid, err := account.VerifyPassword(tt.password, plainPasswordVerifier)
			assert.Equal(t, tt.expectedValid, valid)
			assert.Equal(t, tt.expectedError, err)
			assert.Empty(t, account.Changes().Events())
		})
	}
}

var plainPasswordVerifier PasswordVerifierFunc = func(password string, hashed HashedPassword) (bool, error) {
	return password == string(hashed), nil
}

func TestAccount_AttachMobileDevice(t *testing.T) {
	accID := idgen.New[AccountID]()
	installationID := InstallationID("install1")
	token1 := NotificationDeviceToken("token1")
	token2 := NotificationDeviceToken("token2")

	tests := []struct {
		name          string
		initialEvents []*eventing.JournalEvent
		emittedEvents []eventing.Event
		installID     InstallationID
		token         NotificationDeviceToken
		expectedError error
	}{
		{
			name: "Initial device attachment success",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "john@example.com", "password"),
			),
			emittedEvents: []eventing.Event{
				NewMobileDeviceAttachedToAccountEvent(accID, installationID, token1),
			},
			installID:     installationID,
			token:         token1,
			expectedError: nil,
		},
		{
			name: "Update existing installation token",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "john@example.com", "password"),
				NewMobileDeviceAttachedToAccountEvent(accID, installationID, token1),
			),
			emittedEvents: []eventing.Event{
				NewAccountNotificationDeviceTokenChangedEvent(accID, installationID, token2),
			},
			installID:     installationID,
			token:         token2,
			expectedError: nil,
		},
		{
			name: "No events when token hasn't changed",
			initialEvents: createInitialEvents(
				NewAccountCreatedEvent(accID, "John", "Doe", "john@example.com", "password"),
				NewMobileDeviceAttachedToAccountEvent(accID, installationID, token1),
			),
			installID:     installationID,
			token:         token1,
			expectedError: nil,
		},
		{
			name:          "Fail when account not initialized",
			initialEvents: createInitialEvents(),
			installID:     installationID,
			token:         token1,
			expectedError: NewInvalidAggregateStateError(NewAccount(accID).Aggregate(), int(AccountStateActive), int(AccountStateUnspecified)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account := NewAccount(accID)
			account.Reduce(tt.initialEvents)
			err := account.AttachMobileDevice(tt.installID, tt.token)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.emittedEvents, account.Changes().Events())
		})
	}
}
