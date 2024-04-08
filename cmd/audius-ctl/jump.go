package main

import (
	"github.com/AudiusProject/audius-d/pkg/orchestration"
	"github.com/spf13/cobra"
)

var (
	jumpCmd = &cobra.Command{
		Use:               "jump <host>",
		Short:             "Open a shell into the audius-d container running on a host.",
		ValidArgsFunction: hostsCompletionFunction,
		Args:              cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return orchestration.ShellIntoNode(args[0])
		},
	}
)
