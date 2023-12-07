package main

import (
	"log"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/orchestration"
	"github.com/spf13/cobra"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Uses the currently enabled context to spin up audius nodes.",
		Run: func(cmd *cobra.Command, args []string) {
			orchestration.RunAudiusWithConfig(readOrCreateContext())
		},
	}
	downCmd = &cobra.Command{
		Use:   "down",
		Short: "Spin down nodes and network in the current context.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := readOrCreateContext()
			orchestration.RunDown(ctx)
			orchestration.DownDevnet(ctx)
		},
	}
	devnetCmd = &cobra.Command{
		Use:   "devnet [command]",
		Short: "Spin up local ethereum, solana, and acdc chains for development",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := readOrCreateContext()
			orchestration.StartDevnet(ctx)
		},
	}
	devnetDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Shut down local ethereum, solana, and acdc chains",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := readOrCreateContext()
			orchestration.DownDevnet(ctx)
		},
	}
)

func init() {
	devnetCmd.AddCommand(devnetDownCmd)
}

func readOrCreateContext() *conf.ContextConfig {
	ctx_config, err := conf.ReadOrCreateContextConfig()
	if err != nil {
		log.Fatal("Failed to retrieve context: ", err)
	}
	return ctx_config
}
