package commands

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"testing"
)

func TestCommands_CreateAccount(t *testing.T) {
	validationCases := []struct {
		name string
		req  CreateAccountCommand
		err  error
	}{
		{
			name: "creates account",
			req: CreateAccountCommand{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "valid@example.com",
				Password:  "password",
			},
		},
		{
			name: "rejects empty first name",
			req: CreateAccountCommand{
				FirstName: "",
				LastName:  "Doe",
				Email:     "valid@example.com",
				Password:  "password",
			},
			err: validation.Errors{
				validation.NewFieldError("first_name", validation.ErrRequired),
			},
		},
		{
			name: "rejects empty last name",
			req: CreateAccountCommand{
				FirstName: "John",
				LastName:  "",
				Email:     "valid@example.com",
				Password:  "password",
			},
			err: validation.Errors{
				validation.NewFieldError("last_name", validation.ErrRequired),
			},
		},
		{
			name: "rejects empty email",
			req: CreateAccountCommand{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "",
				Password:  "password",
			},
			err: validation.Errors{
				validation.NewFieldError("email", validation.ErrRequired),
			},
		},
	}
	for _, tc := range validationCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cmds, _ := setupCommands(t)

			_, err := cmds.CreateAccount(context.Background(), tc.req)
			if !errors.Is(err, tc.err) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
