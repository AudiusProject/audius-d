package main

import (
	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/AudiusProject/audius-d/pkg/orchestration"
	"github.com/spf13/cobra"
)

var (
	awaitHealthy = false
	upCmd        = &cobra.Command{
		Use:   "up",
		Short: "Uses the currently enabled context to spin up audius nodes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			orchestration.RunAudiusWithConfig(readOrCreateContext(), awaitHealthy)
			return nil
		},
	}
	downCmd = &cobra.Command{
		Use:   "down",
		Short: "Spin down nodes and network in the current context.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := readOrCreateContext()
			orchestration.RunDown(ctx)
			orchestration.DownDevnet(ctx)
			return nil
		},
	}
	devnetCmd = &cobra.Command{
		Use:   "devnet [command]",
		Short: "Spin up local ethereum, solana, and acdc chains for development",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := readOrCreateContext()
			orchestration.StartDevnet(ctx)
			return nil
		},
	}
	devnetDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Shut down local ethereum, solana, and acdc chains",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := readOrCreateContext()
			orchestration.DownDevnet(ctx)
			return nil
		},
	}
)

func init() {
	upCmd.Flags().BoolVar(&awaitHealthy, "await-healthy", false, "Wait for services to become healthy before returning.")
	devnetCmd.AddCommand(devnetDownCmd)
}

func readOrCreateContext() *conf.ContextConfig {
	ctx_config, err := conf.ReadOrCreateContextConfig()
	if err != nil {
		logger.Error("Failed to retrieve context: ", err)
		return nil
	}
	return ctx_config
}
