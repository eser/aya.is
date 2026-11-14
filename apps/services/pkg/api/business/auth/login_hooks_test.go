//nolint:testpackage // exercises the unexported runLoginHooks
package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/eser/aya.is/services/pkg/api/business/users"
)

var errHook = errors.New("hook failed")

func TestLoginHooksRunInOrderAndAbortOnError(t *testing.T) {
	t.Parallel()

	service := &Service{}            //nolint:exhaustruct // only the hooks matter here
	user := &users.User{ID: "user1"} //nolint:exhaustruct

	var calls []string

	service.AddLoginHook(func(_ context.Context, hookUser *users.User) error {
		calls = append(calls, "first:"+hookUser.ID)

		return nil
	})
	service.AddLoginHook(func(_ context.Context, _ *users.User) error {
		calls = append(calls, "second")

		return errHook
	})
	service.AddLoginHook(func(_ context.Context, _ *users.User) error {
		calls = append(calls, "third")

		return nil
	})

	err := service.runLoginHooks(context.Background(), user)
	if !errors.Is(err, ErrFailedToHandleCallback) || !errors.Is(err, errHook) {
		t.Fatalf("a failing hook must abort the login with ErrFailedToHandleCallback: %v", err)
	}

	if len(calls) != 2 || calls[0] != "first:user1" || calls[1] != "second" {
		t.Fatalf("hooks ran as %v, want the third skipped after the failure", calls)
	}
}

func TestLoginWithoutHooksSucceeds(t *testing.T) {
	t.Parallel()

	err := (&Service{}).runLoginHooks(context.Background(), &users.User{}) //nolint:exhaustruct
	if err != nil {
		t.Fatal(err)
	}
}
