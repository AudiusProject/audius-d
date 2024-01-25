package main

import (
	"github.com/AudiusProject/audius-d/pkg/infra"
	"github.com/spf13/cobra"
)

var (
	infraCmd = &cobra.Command{
		Use:   "infra [command]",
		Short: "Manage audius-d instances",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	infraDestroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return infra.Destroy()
		},
	}

	infraPreviewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Perform a dry-run update for the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return infra.Update(true)
		},
	}

	infraUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update (deploy) the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return infra.Update(false)
		},
	}
)

func init() {
	infraCmd.AddCommand(infraDestroyCmd, infraPreviewCmd, infraUpdateCmd)
}
