package subcommands

import (
	"context"

	"github.com/eser/aya.is/services/pkg/api/adapters/appcontext"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/spf13/cobra"
)

func CmdProfilesList() *cobra.Command {
	profilesListCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "list",
		Short: "Lists profiles",
		Long:  "Lists all available profiles registered on the site",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execProfilesList(cmd.Context())
		},
	}

	return profilesListCmd
}

func execProfilesList(ctx context.Context) error {
	appContext := appcontext.New()

	err := appContext.Init(ctx)
	if err != nil {
		return err //nolint:wrapcheck
	}

	profileList, err := appContext.ProfileService.List(ctx, "en", cursors.NewCursor(0, nil))
	if err != nil {
		panic(err)
	}

	for _, record := range profileList.Data {
		appContext.Logger.InfoContext(ctx, "profile entry", "profile", record)
	}

	return nil
}
