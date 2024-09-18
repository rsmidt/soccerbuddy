package commands

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rsmidt/soccerbuddy/gen/eventregistry"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"github.com/rsmidt/soccerbuddy/internal/domain/account"
	"github.com/rsmidt/soccerbuddy/internal/eventing/pg"
	"github.com/rsmidt/soccerbuddy/internal/postgres"
	"log/slog"
	"testing"
)

func setupCommands(t *testing.T) (*Commands, *pgxpool.Pool) {
	pool, cleanup := postgres.GetTestPool()
	t.Cleanup(cleanup)
	es := pg.NewEventStore(pool, eventregistry.Default, pg.NewEventCrypto(pool), slog.Default())
	commands := NewCommands(es)
	return commands, pool
}

func TestCommands_CreateRootAccount(t *testing.T) {
	t.Run("rejects if root account already exists", func(t *testing.T) {
		t.Parallel()
		commands, _ := setupCommands(t)

		successfulReq := CreateRootAccountCommand{
			Email:    "test@testerino.com",
			Password: "password",
		}
		if err := commands.CreateRootAccount(context.Background(), successfulReq); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		failingReq := CreateRootAccountCommand{
			Email:    "another@test.com",
			Password: "password",
		}
		if err := commands.CreateRootAccount(context.Background(), failingReq); err == nil {
			t.Fatalf("expected error, got nil")
		} else if !errors.Is(err, account.ErrRootAccountAlreadyInitialized) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	validationCases := []struct {
		name string
		req  CreateRootAccountCommand
		err  error
	}{
		{
			name: "creates root account",
			req: CreateRootAccountCommand{
				Email:    "valid@example.com",
				Password: "password",
			},
		},
		{
			name: "rejects empty email",
			req: CreateRootAccountCommand{
				Email:    "",
				Password: "password",
			},
			err: validation.Errors{
				validation.NewFieldError("email", validation.ErrRequired),
			},
		},
		{
			name: "rejects empty password",
			req: CreateRootAccountCommand{
				Email:    "email@email.com",
				Password: "",
			},
			err: validation.Errors{
				validation.NewFieldError("password", validation.ErrRequired),
			},
		},
		{
			name: "rejects too short password",
			req: CreateRootAccountCommand{
				Email:    "email@email.com",
				Password: "short",
			},
			err: validation.Errors{
				validation.NewMinLengthError("password", 8),
			},
		},
	}
	for _, tc := range validationCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cmds, _ := setupCommands(t)

			err := cmds.CreateRootAccount(context.Background(), tc.req)
			if !errors.Is(err, tc.err) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
