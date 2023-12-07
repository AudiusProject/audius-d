//go:build !mac
// +build !mac

package main

import (
	"github.com/spf13/cobra"
)

var sbCmd = &cobra.Command{
	Use:   "statusbar",
	Short: "Run mac status bar",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
