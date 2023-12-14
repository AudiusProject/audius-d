package main

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			if displayVersion {
				logger.Info(Version)
			} else {
				upCmd.Run(cmd, args)
			}
		},
	}
	rootCmd.Flags().BoolVar(&displayVersion, "version", false, "--version")
	rootCmd.AddCommand(upCmd, downCmd, devnetCmd, registerCmd, configCmd, guiCmd, sbCmd, emCmd, hashCmd)
	rootCmd.Execute()
}
