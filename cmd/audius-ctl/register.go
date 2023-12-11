package main

import (
	"log"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/register"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register nodes on ethereum (only works for local devnet)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx_config, err := conf.ReadOrCreateContextConfig()
		if err != nil {
			log.Fatal("Failed to retrieve context: ", err)
		}
		for _, cc := range ctx_config.CreatorNodes {
			register.RegisterNode(
				"content-node",
				cc.Host,
				ctx_config.Network.EthMainnetRpc,
				"0xdcB2fC9469808630DD0744b0adf97C0003fC29B2", // hardcoded ganache address
				"0xABbfF712977dB51f9f212B85e8A4904c818C2b63", // "
				cc.OperatorWallet,
				cc.OperatorPrivateKey,
			)
		}
	},
}
