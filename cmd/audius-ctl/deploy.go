package main

import (
	"github.com/AudiusProject/audius-d/pkg/deploy"
	"github.com/spf13/cobra"
)

var (
	deployCmd = &cobra.Command{
		Use:   "deploy [command]",
		Short: "deploy current-context",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			deploy.DeployContext()
			return nil
		},
	}
)

func init() {
	testCmd.AddCommand(testContextCmd)
}
