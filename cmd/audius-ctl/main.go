package main

import (
	"os"

	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	Version        string
	displayVersion bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "audius-ctl [command]",
		Short: "CLI for provisioning and interacting with audius nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if displayVersion {
				logger.Out(Version)
			} else {
				upCmd.Run(cmd, args)
			}
			return nil
		},
	}
	rootCmd.Flags().BoolVar(&displayVersion, "version", false, "--version")
	rootCmd.AddCommand(upCmd, downCmd, devnetCmd, registerCmd, configCmd, guiCmd, sbCmd, emCmd, hashCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
