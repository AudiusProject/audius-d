package orchestration

import (
	"github.com/AudiusProject/audius-d/conf"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var (
	UpCmd = &cobra.Command{
		Use:   "up",
		Short: "Uses the currently enabled context to spin up audius nodes.",
		Run: func(cmd *cobra.Command, args []string) {
			audius := figure.NewColorFigure("Audius", "starwars", "purple", true)
			audius.Print()
			RunAudiusWithConfig(cmd.Context().Value(conf.ContextKey).(*conf.ContextConfig))
		},
	}
	DownCmd = &cobra.Command{
		Use:   "down",
		Short: "Spin down nodes and network in the current context.",
		Run: func(cmd *cobra.Command, args []string) {
			RunDown(cmd.Context().Value(conf.ContextKey).(*conf.ContextConfig))
			DownDevnet(cmd.Context().Value(conf.ContextKey).(*conf.ContextConfig))
		},
	}
	DevnetCmd = &cobra.Command{
		Use:   "devnet [command]",
		Short: "Spin up local ethereum, solana, and acdc chains for development",
		Run: func(cmd *cobra.Command, args []string) {
			StartDevnet(cmd.Context().Value(conf.ContextKey).(*conf.ContextConfig))
		},
	}
	DevnetDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Shut down local ethereum, solana, and acdc chains",
		Run: func(cmd *cobra.Command, args []string) {
			DownDevnet(cmd.Context().Value(conf.ContextKey).(*conf.ContextConfig))
		},
	}
)

func init() {
	DevnetCmd.AddCommand(DevnetDownCmd)
}
