package main

import (
	"fmt"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/AudiusProject/audius-d/pkg/test"
	"github.com/spf13/cobra"
)

var (
	testCmd = &cobra.Command{
		Use:   "test [command]",
		Short: "test audius-d connectivity",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	testContextCmd = &cobra.Command{
		Use:          "context",
		Short:        "Test the health of the current context",
		SilenceUsage: true, // do not print --help text on failed node health
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxConfig, err := conf.ReadOrCreateContextConfig()
			if err != nil {
				return logger.Error("Failed to retrieve context. ", err)
			}

			responses, err := test.CheckNodeHealth(ctxConfig)
			if err != nil {
				return err
			}

			var encounteredError bool
			for _, response := range responses {
				if response.Error != nil {
					fmt.Printf("%-50s Error: %v\n", response.Host, response.Error)
					encounteredError = true
				} else {
					fmt.Printf("%-50s %t\n", response.Host, response.Result)
				}
			}

			if encounteredError {
				return fmt.Errorf("\none or more health checks failed")
			}

			return nil
		},
	}
)

func init() {
	testCmd.AddCommand(testContextCmd)
}
