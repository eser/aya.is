package subcommands

import (
	"context"

	"github.com/eser/aya.is/services/pkg/api/adapters/appcontext"
	"github.com/spf13/cobra"
)

func CmdReady() *cobra.Command {
	readyCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "ready",
		Short: "Checks the readiness of the site",
		Long:  "Checks the readiness of the site",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execReady(cmd.Context())
		},
	}

	return readyCmd
}

func execReady(ctx context.Context) error {
	appContext := appcontext.New()

	err := appContext.Init(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	appContext.Logger.InfoContext(ctx, "readiness check passed")

	return nil
}
