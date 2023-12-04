package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "audius-d [command]",
		Short: "CLI for provisioning and interacting with audius nodes",
		Run: func(cmd *cobra.Command, args []string) {
			upCmd.Run(cmd, args)
		},
	}
	rootCmd.AddCommand(upCmd, downCmd, devnetCmd, registerCmd, configCmd, guiCmd)
	rootCmd.Execute()
}
