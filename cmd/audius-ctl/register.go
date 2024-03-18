package main

import (
	"fmt"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/AudiusProject/audius-d/pkg/orchestration"
	"github.com/AudiusProject/audius-d/pkg/register"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register nodes on ethereum (only works for local devnet)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx_config, err := conf.ReadOrCreateContextConfig()
		if err != nil {
			return logger.Error("Failed to retrieve context:", err)
		}
		for host, nodeConfig := range ctx_config.Nodes {
			if nodeConfig.Type != conf.Creator {
				continue
			}
			privateKey, err := orchestration.NormalizedPrivateKey(host, nodeConfig.PrivateKey)
			if err != nil {
				return logger.Error("Failed to obtain private key:", err)
			}
			err = register.RegisterNode(
				"content-node",
				fmt.Sprintf("https://%s", host),
				"http://localhost:8546",
				register.GanacheAudiusTokenAddress,
				register.GanacheContractRegistryAddress,
				nodeConfig.Wallet,
				privateKey,
			)
			if err != nil {
				return logger.Error("Failed to register node:", err)
			}
		}
		return nil
	},
}
