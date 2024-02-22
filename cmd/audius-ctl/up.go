package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/AudiusProject/audius-d/pkg/orchestration"
	"github.com/spf13/cobra"
)

var (
	awaitHealthy = false
	audiusdTag   = "default"
	upCmd        = &cobra.Command{
		Use:   "up [hosts]",
		Short: "Spin up the audius nodes specified in your config, optionally specifying which hosts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var nodesToRunUp map[string]conf.NodeConfig
			var err error
			ctx := readOrCreateContext()
			if len(args) == 0 {
				nodesToRunUp = ctx.Nodes
			} else {
				nodesToRunUp, err = filterNodesFromContext(args, ctx)
				if err != nil {
					return err
				}
			}
			orchestration.RunAudiusNodes(nodesToRunUp, ctx.Network, awaitHealthy, audiusdTag)
			return nil
		},
	}

	downAll   = false
	downForce = false
	downCmd   = &cobra.Command{
		Use:   "down [hosts] [--all/-a, --force/-f]",
		Short: "Spin down nodes and network in the current context.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if downAll && len(args) > 0 {
				return logger.Error("Cannot combine specific nodes with flag --all/-a.")
			} else if !downAll && len(args) == 0 {
				return logger.Error("Must specify which nodes to take down or use --all/-a.")
			}

			ctx := readOrCreateContext()
			var nodesToRunDown map[string]conf.NodeConfig
			var err error
			if downAll {
				nodesToRunDown = ctx.Nodes
			} else {
				nodesToRunDown, err = filterNodesFromContext(args, ctx)
				if err != nil {
					return err
				}
			}

			infoString := "This will run down the following nodes:"
			for host := range nodesToRunDown {
				infoString += fmt.Sprintf("\n%s", host)
			}
			logger.Info(infoString)

			if !downForce && !askForConfirmation("Are you sure you want to continue?") {
				return logger.Error("Aborted")
			}

			orchestration.RunDownNodes(nodesToRunDown)
			return nil
		},
	}
	devnetCmd = &cobra.Command{
		Use:   "devnet [command]",
		Short: "Spin up local ethereum, solana, and acdc chains for development",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := readOrCreateContext()
			return orchestration.StartDevnet(ctx)
		},
	}
	devnetDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Shut down local ethereum, solana, and acdc chains",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := readOrCreateContext()
			return orchestration.DownDevnet(ctx)
		},
	}
)

func init() {
	upCmd.Flags().BoolVarP(&awaitHealthy, "await-healthy", "w", false, "Wait for services to become healthy before returning.")
	upCmd.Flags().StringVar(&audiusdTag, "audius-d-version", "default", "(Development) override docker image tag to use (audius/audius-d:<version>)")
	downCmd.Flags().BoolVarP(&downAll, "all", "a", false, "Take down all nodes in the current configuration.")
	downCmd.Flags().BoolVarP(&downForce, "force", "f", false, "Do not ask for confirmation.")
	devnetCmd.AddCommand(devnetDownCmd)
}

func readOrCreateContext() *conf.ContextConfig {
	ctx_config, err := conf.ReadOrCreateContextConfig()
	if err != nil {
		logger.Error("Failed to retrieve context: ", err)
		return nil
	}
	return ctx_config
}

func filterNodesFromContext(desired []string, ctx *conf.ContextConfig) (map[string]conf.NodeConfig, error) {
	result := make(map[string]conf.NodeConfig)
	for _, desiredHost := range desired {
		matched := false
		for host, config := range ctx.Nodes {
			if desiredHost == host {
				matched = true
				result[host] = config
			}
		}
		if !matched {
			return nil, logger.Errorf("Node %s does not exist. Check your configuration (`audius-ctl config`)", desiredHost)
		}
	}
	return result, nil

}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Error encountered reading from stdin")
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
