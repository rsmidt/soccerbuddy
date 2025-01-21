package commands

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/domain/authz"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	insecrand "math/rand"
	"time"
)

type CreateAccountCommand struct {
	FirstName string
	LastName  string

	Email    string
	Password string
}

func (c CreateAccountCommand) Validate() error {
	var errs validation.Errors
	if c.FirstName == "" {
		errs = append(errs, validation.NewFieldError("first_name", validation.ErrRequired))
	}
	if c.LastName == "" {
		errs = append(errs, validation.NewFieldError("last_name", validation.ErrRequired))
	}
	if c.Email == "" {
		errs = append(errs, validation.NewFieldError("email", validation.ErrRequired))
	}
	if c.Password == "" {
		errs = append(errs, validation.NewFieldError("password", validation.ErrRequired))
	} else if len(c.Password) < 8 {
		errs = append(errs, validation.NewMinLengthError("password", 8))
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) CreateAccount(ctx context.Context, cmd CreateAccountCommand) (domain.AccountID, error) {
	ctx, span := tracing.Tracer.Start(ctx, "commands.CreateAccount")
	defer span.End()

	if err := c.authorizer.Authorize(ctx, authz.ActionCreateAccount, authz.SystemResource); err != nil {
		return "", err
	}
	if err := cmd.Validate(); err != nil {
		return "", err
	}

	exists, err := c.repos.Account().ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", validation.NewExistsError("email")
	}

	hashedPW, err := domain.Argon2idHashPassword(cmd.Password)
	if err != nil {
		return "", err
	}
	id := idgen.New[domain.AccountID]()
	account, err := c.repos.Account().FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if err := account.Init(cmd.FirstName, cmd.LastName, cmd.Email, hashedPW); err != nil {
		return "", err
	}

	if err := c.repos.Account().Save(ctx, account); err != nil {
		return "", err
	}
	return id, nil
}

type LoginAccountCommand struct {
	Email     string
	Password  string
	UserAgent string
	IPAddress string
}

func (c LoginAccountCommand) Validate() error {
	var errs validation.Errors
	if c.Email == "" {
		errs = append(errs, validation.NewFieldError("email", validation.ErrRequired))
	}
	if c.Password == "" {
		errs = append(errs, validation.NewFieldError("password", validation.ErrRequired))
	}
	if c.IPAddress == "" {
		errs = append(errs, validation.NewFieldError("ip_address", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type LoginAccountResult struct {
	Token     domain.SessionToken
	ExpiresAt time.Time
}

func (c *Commands) Login(ctx context.Context, cmd LoginAccountCommand) (*LoginAccountResult, error) {
	ctx, span := tracing.Tracer.Start(ctx, "commands.Login")
	defer span.End()

	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	account, err := c.repos.Account().FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if ok, err := account.VerifyPassword(cmd.Password, domain.Argon2idVerifyPassword); err != nil {
		return nil, err
	} else if !ok {
		// Add some random delay to prevent timing attacks.
		time.Sleep(time.Duration(50+insecrand.Intn(50)) * time.Millisecond)
		return nil, domain.ErrWrongCredentials
	}

	id := idgen.New[domain.SessionID]()
	session, err := c.repos.Session().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	var role domain.PrincipalRole
	if account.IsRoot {
		role = domain.PrincipalRoleRoot
	} else {
		role = domain.PrincipalRoleRegular
	}
	if err := session.Init(token, account.ID, cmd.UserAgent, cmd.IPAddress, expiresAt, role); err != nil {
		return nil, err
	}
	if err := c.repos.Session().Save(ctx, session); err != nil {
		return nil, err
	}
	return &LoginAccountResult{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

type AttachMobileDeviceCommand struct {
	InstallationID          domain.InstallationID
	NotificationDeviceToken domain.NotificationDeviceToken
}

func (c *AttachMobileDeviceCommand) Validate() error {
	var errs validation.Errors

	if c.InstallationID == "" {
		errs = append(errs, validation.NewFieldError("installation_id", validation.ErrRequired))
	}
	if c.NotificationDeviceToken == "" {
		errs = append(errs, validation.NewFieldError("notification_device_token", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (c *Commands) AttachMobileDevice(ctx context.Context, cmd *AttachMobileDeviceCommand) error {
	ctx, span := tracing.Tracer.Start(ctx, "commands.AttachMobileDevice")
	defer span.End()

	principal, ok := domain.PrincipalFromContext(ctx)
	if !ok {
		return domain.ErrUnauthenticated
	}
	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := c.authorizer.Authorize(ctx, authz.ActionEdit, authz.NewAccountResource(principal.AccountID)); err != nil {
		return err
	}

	account, err := c.repos.Account().FindByID(ctx, principal.AccountID)
	if err != nil {
		return err
	}
	if err := account.AttachMobileDevice(cmd.InstallationID, cmd.NotificationDeviceToken); err != nil {
		return err
	}
	if err := c.repos.Account().Save(ctx, account); err != nil {
		return err
	}
	return nil
}

func generateSessionToken() (domain.SessionToken, error) {
	id, err := randomString(16)
	if err != nil {
		return "", err
	}
	return domain.SessionToken(id), nil
}

func randomString(nBytes int) (string, error) {
	randomBytes := make([]byte, nBytes)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}

type RegisterAccountCommand struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	LinkToken string
	UserAgent string
	IPAddress string
}

func (c *RegisterAccountCommand) Validate() error {
	var errs validation.Errors
	if c.FirstName == "" {
		errs = append(errs, validation.NewFieldError("first_name", validation.ErrRequired))
	}
	if c.LastName == "" {
		errs = append(errs, validation.NewFieldError("last_name", validation.ErrRequired))
	}
	if c.Email == "" {
		errs = append(errs, validation.NewFieldError("email", validation.ErrRequired))
	}
	if c.Password == "" {
		errs = append(errs, validation.NewFieldError("password", validation.ErrRequired))
	}
	if c.LinkToken == "" {
		errs = append(errs, validation.NewFieldError("link_token", validation.ErrRequired))
	}
	if c.UserAgent == "" {
		errs = append(errs, validation.NewFieldError("user_agent", validation.ErrRequired))
	}
	if c.IPAddress == "" {
		errs = append(errs, validation.NewFieldError("ip_address", validation.ErrRequired))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type RegisterAccountResult struct {
	AccountID    domain.AccountID
	SessionToken domain.SessionToken
	ExpiresAt    time.Time
}

func (c *Commands) RegisterAccount(ctx context.Context, cmd *RegisterAccountCommand) (*RegisterAccountResult, error) {
	ctx, span := tracing.Tracer.Start(ctx, "commands.RegisterAccount")
	defer span.End()

	linkToken := domain.PersonLinkToken(cmd.LinkToken)
	pers, err := c.getPersonProjectionByPendingToken(ctx, linkToken)
	if len(pers) == 0 {
		return nil, domain.ErrPersonInvalidLinkToken
	}

	exists, err := c.repos.Account().ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, validation.NewExistsError("email")
	}
	hashedPW, err := domain.Argon2idHashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}
	id := idgen.New[domain.AccountID]()
	account, err := c.repos.Account().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := account.Register(cmd.FirstName, cmd.LastName, cmd.Email, hashedPW, domain.PersonLinkToken(cmd.LinkToken)); err != nil {
		return nil, err
	}
	if err := c.repos.Account().Save(ctx, account); err != nil {
		return nil, err
	}
	// TODO: No need to do the full login flow here.
	res, err := c.Login(ctx, LoginAccountCommand{
		Email:    cmd.Email,
		Password: cmd.Password,
		// TODO: Add UserAgent and IPAddress.
		UserAgent: cmd.UserAgent,
		IPAddress: cmd.IPAddress,
	})
	if err != nil {
		// TODO: Rollback account creation?.
		return nil, err
	}
	return &RegisterAccountResult{
		AccountID:    id,
		SessionToken: res.Token,
		ExpiresAt:    res.ExpiresAt,
	}, nil
}
