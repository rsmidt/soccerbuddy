package commands

import (
	"context"
	"errors"
	"github.com/rsmidt/soccerbuddy/internal/app/validation"
	"testing"
)

func TestCommands_CreateClub(t *testing.T) {
	validationCases := []struct {
		name string
		req  CreateClubCommand
		err  error
	}{
		{
			name: "creates club",
			req: CreateClubCommand{
				Name: "Club Name",
			},
			err: nil,
		},
		{
			name: "rejects empty name",
			req: CreateClubCommand{
				Name: "",
			},
			err: validation.Errors{
				validation.NewFieldError("name", validation.ErrRequired),
			},
		},
	}

	for _, tc := range validationCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cmds, _ := setupCommands(t)

			_, err := cmds.CreateClub(context.Background(), tc.req)
			if !errors.Is(err, tc.err) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
