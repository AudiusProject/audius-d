package main

import (
	"github.com/AudiusProject/audius-d/pkg/statusbar"
	"github.com/spf13/cobra"
)

var sbCmd = &cobra.Command{
	Use:   "statusbar",
	Short: "Run mac status bar",
	Run: func(cmd *cobra.Command, args []string) {
		statusbar.RunStatusBar()
	},
}
