package main

import (
	"github.com/eser/aya.is/services/cmd/manage/subcommands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "manage",
		Short: "aya.is/services CLI for site management",
		Long:  "aya.is/services CLI provides various functionalities for site management including reporting and administration.", //nolint:lll
	}

	rootCmd.AddCommand(subcommands.CmdID())
	rootCmd.AddCommand(subcommands.CmdReady())
	rootCmd.AddCommand(subcommands.CmdProfiles())
	rootCmd.AddCommand(subcommands.CmdScrape())

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
