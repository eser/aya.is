package subcommands

import (
	"github.com/spf13/cobra"
)

func CmdProfiles() *cobra.Command {
	profilesCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "profiles",
		Short: "Manages profiles",
		Long:  "Manages profiles registered on the site",
	}

	profilesCmd.AddCommand(CmdProfilesList())
	profilesCmd.AddCommand(CmdProfilesImport())

	return profilesCmd
}
