package main

import (
	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/AudiusProject/audius-d/pkg/register"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register nodes on ethereum (only works for local devnet)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx_config, err := conf.ReadOrCreateContextConfig()
		if err != nil {
			return logger.Error("Failed to retrieve context: ", err)
		}
		for _, cc := range ctx_config.CreatorNodes {
			err := register.RegisterNode(
				"content-node",
				cc.Host,
				"http://localhost:8546",
				register.GanacheAudiusTokenAddress,
				register.GanacheContractRegistryAddress,
				cc.OperatorWallet,
				cc.OperatorPrivateKey,
			)
			if err != nil {
				return logger.Error("Failed to register node:", err)
			}
		}
		return nil
	},
}
