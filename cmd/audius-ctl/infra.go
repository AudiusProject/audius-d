package main

import (
	"context"
	"fmt"

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

	infraCancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel the current in progress update",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return infra.Cancel(ctx)
		},
	}

	infraDestroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return infra.Destroy(ctx)
		},
	}

	infraGetOutputCmd = &cobra.Command{
		Use:   "output <outputName>",
		Short: "Get a specific output from the stack",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			outputName := args[0]
			outputValue, err := infra.GetStackOutput(ctx, outputName)
			if err != nil {
				return err
			}
			fmt.Println(outputValue)
			return nil
		},
	}

	infraPreviewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Perform a dry-run update for the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return infra.Update(ctx, true)
		},
	}

	infraUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update (deploy) the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return infra.Update(ctx, false)
		},
	}
)

func init() {
	infraCmd.AddCommand(infraCancelCmd, infraDestroyCmd, infraGetOutputCmd, infraPreviewCmd, infraUpdateCmd)
}
