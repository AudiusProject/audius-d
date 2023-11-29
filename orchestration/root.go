package orchestration

import (
	"log"

	"github.com/AudiusProject/audius-d/conf"
	"github.com/spf13/cobra"
)

var (
	UpCmd = &cobra.Command{
		Use:   "up",
		Short: "Uses the currently enabled context to spin up audius nodes.",
		Run: func(cmd *cobra.Command, args []string) {
			RunAudiusWithConfig(readOrCreateContext())
		},
	}
	DownCmd = &cobra.Command{
		Use:   "down",
		Short: "Spin down nodes and network in the current context.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := readOrCreateContext()
			RunDown(ctx)
			DownDevnet(ctx)
		},
	}
	DevnetCmd = &cobra.Command{
		Use:   "devnet [command]",
		Short: "Spin up local ethereum, solana, and acdc chains for development",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := readOrCreateContext()
			StartDevnet(ctx)
		},
	}
	DevnetDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Shut down local ethereum, solana, and acdc chains",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := readOrCreateContext()
			DownDevnet(ctx)
		},
	}
)

func init() {
	DevnetCmd.AddCommand(DevnetDownCmd)
}

func readOrCreateContext() *conf.ContextConfig {
	ctx_config, err := conf.ReadOrCreateContextConfig()
	if err != nil {
		log.Fatal("Failed to retrieve context: ", err)
	}
	return ctx_config
}
