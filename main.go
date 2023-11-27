package main

import (
	"context"
	_ "embed"
	"log"

	"github.com/AudiusProject/audius-d/conf"
	"github.com/AudiusProject/audius-d/em"
	"github.com/AudiusProject/audius-d/hashes"
	"github.com/AudiusProject/audius-d/orchestration"
	"github.com/AudiusProject/audius-d/register"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "audius-d [command]",
		Short: "CLI for provisioning and interacting with audius nodes",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetContext(context.WithValue(context.Background(), conf.ContextKey, readOrCreateContext()))
		},
		Run: func(cmd *cobra.Command, args []string) {
			orchestration.UpCmd.Run(cmd, args)
		},
	}
	rootCmd.AddCommand(orchestration.UpCmd, orchestration.DownCmd, orchestration.DevnetCmd)
	rootCmd.AddCommand(register.RootCmd)
	rootCmd.AddCommand(conf.RootCmd)
	rootCmd.AddCommand(em.EmCmd)
	rootCmd.AddCommand(hashes.HashCmd)
	rootCmd.Execute()
}

func readOrCreateContext() *conf.ContextConfig {
	ctx_config, err := conf.ReadOrCreateContextConfig()
	if err != nil {
		log.Fatal("Failed to retrieve context:", err)
	}
	return ctx_config
}
