package gui

import (
	"github.com/spf13/cobra"
)

var (
	RunGuiServer = &cobra.Command{
		Use:   "gui",
		Short: "Run gui server",
		Run: func(cmd *cobra.Command, args []string) {
			StartGuiServer()
		},
	}
)
