package main

import (
	"github.com/AudiusProject/audius-d/pkg/gui"
	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Run gui server",
	Run: func(cmd *cobra.Command, args []string) {
		gui.StartGuiServer()
	},
}
