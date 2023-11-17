package main

import (
	"context"
	_ "embed"
	"log"

	"github.com/AudiusProject/audius-d/conf"
	"github.com/AudiusProject/audius-d/orchestration"
	"github.com/AudiusProject/audius-d/register"
	"github.com/spf13/cobra"
)

//go:embed sample.audius.conf
var confExample string

// with the intent of reducing configuration,
// the latest audius-docker-compose sha (from stage branch) is set at build time via ci.
// this bakes the (tested) image dependency in, so we know that the built binary will always work.
var imageTag string

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
	rootCmd.Execute()
}

func readOrCreateContext() *conf.ContextConfig {
	ctx_config, err := conf.ReadOrCreateContextConfig()
	if err != nil {
		log.Fatal("Failed to retrieve context:", err)
	}
	return ctx_config
}
