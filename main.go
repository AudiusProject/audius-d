package main

import (
	_ "embed"

	"github.com/AudiusProject/audius-d/conf"
	"github.com/AudiusProject/audius-d/migration"
	"github.com/AudiusProject/audius-d/orchestration"
	"github.com/AudiusProject/audius-d/register"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "audius-d [command]",
		Short: "CLI for provisioning and interacting with audius nodes",
		Run: func(cmd *cobra.Command, args []string) {
			orchestration.UpCmd.Run(cmd, args)
		},
	}
	rootCmd.AddCommand(orchestration.UpCmd, orchestration.DownCmd, orchestration.DevnetCmd)
	rootCmd.AddCommand(register.RootCmd)
	rootCmd.AddCommand(conf.RootCmd)
	rootCmd.AddCommand(migration.MigrateCmd)
	rootCmd.Execute()
}
