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
				upCmd.RunE(cmd, args)
			}
			return nil
		},
	}
	rootCmd.Flags().BoolVar(&displayVersion, "version", false, "Display version info")
	rootCmd.AddCommand(configCmd, devnetCmd, downCmd, emCmd, guiCmd, hashCmd, registerCmd, sbCmd, testCmd, upCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
