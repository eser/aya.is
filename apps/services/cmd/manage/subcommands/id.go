package subcommands

import (
	"context"
	"fmt"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/spf13/cobra"
)

func CmdID() *cobra.Command {
	var flagCount int

	idCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "id",
		Short: "Generates id",
		Long:  "Generates id",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execID(cmd.Context(), flagCount)
		},
	}

	idCmd.Flags().IntVarP(&flagCount, "count", "n", 1, "count of ids will be generated")

	return idCmd
}

func execID(_ context.Context, count int) error {
	for range count {
		id := lib.IDsGenerateUnique()

		fmt.Println(id) //nolint:forbidigo
	}

	return nil
}
