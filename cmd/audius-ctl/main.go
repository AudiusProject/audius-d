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
	var rootCmd = &cobra.Command{
		Use:   "audius-ctl [command]",
		Short: "CLI for provisioning and interacting with audius nodes",
		Run: func(cmd *cobra.Command, args []string) {
			if displayVersion {
				logger.Out(Version)
				return
			}
			cmd.Help()
		},
	}

	rootCmd.Flags().BoolVarP(&displayVersion, "version", "v", false, "Display version info")
	rootCmd.AddCommand(configCmd, devnetCmd, downCmd, infraCmd, registerCmd, restartCmd, sbCmd, testCmd, upCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
