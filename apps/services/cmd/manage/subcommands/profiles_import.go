package subcommands

import (
	"context"

	"github.com/eser/aya.is/services/pkg/api/adapters/appcontext"
	"github.com/spf13/cobra"
)

func CmdProfilesImport() *cobra.Command {
	profilesImportCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "import",
		Short: "Imports profile data",
		Long:  "Imports profile related-data such as external posts on remote systems",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execProfilesImport(cmd.Context())
		},
	}

	return profilesImportCmd
}

func execProfilesImport(ctx context.Context) error {
	appContext := appcontext.New()

	err := appContext.Init(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	err = appContext.ProfileService.Import(ctx, appContext.Arcade)
	if err != nil {
		panic(err)
	}

	appContext.Logger.InfoContext(ctx, "profile imports completed")

	return nil
}
