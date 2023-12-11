package orchestration

import (
	"os"
	"os/exec"

	"github.com/AudiusProject/audius-d/pkg/conf"
)

func StartDevnet(_ *conf.ContextConfig) {
	startDevnetDocker()
}

func DownDevnet(_ *conf.ContextConfig) {
	downDevnetDocker()
}

func RunAudiusWithConfig(config *conf.ContextConfig) {
	// stand up devnet should it be required
	if config.Network.Devnet {
		startDevnetDocker()
		registerDevnetNodes(config)
	}

	for cname, cc := range config.CreatorNodes {
		creatorVolumes := []string{"/var/k8s/mediorum:/var/k8s/mediorum", "/var/k8s/creator-node-backend:/var/k8s/creator-node-backend", "/var/k8s/creator-node-db:/var/k8s/creator-node-db"}
		override := cc.ToOverrideEnv(config.Network)
		RunNode(config.Network, cc.BaseServerConfig, override, cname, "creator-node", creatorVolumes)
	}
	for cname, dc := range config.DiscoveryNodes {
		discoveryVolumes := []string{"/var/k8s/discovery-provider-db:/var/k8s/discovery-provider-db", "/var/k8s/discovery-provider-chain:/var/k8s/discovery-provider-chain", "/var/k8s/bolt:/var/k8s/bolt"}
		override := dc.ToOverrideEnv(config.Network)
		RunNode(config.Network, dc.BaseServerConfig, override, cname, "discovery-provider", discoveryVolumes)
		// discovery requires a few extra things
		if !config.Network.Devnet {
			audiusCli(cname, "launch-chain")
		}
	}
	for cname, id := range config.IdentityService {
		identityVolumes := []string{"/var/k8s/identity-service-db:/var/lib/postgresql/data"}
		override := id.ToOverrideEnv(config.Network)
		RunNode(config.Network, id.BaseServerConfig, override, cname, "identity-service", identityVolumes)
	}
}

func RunDown(config *conf.ContextConfig) {
	// easiest way
	cnames := []string{"rm", "-f"}

	for cname := range config.CreatorNodes {
		cnames = append(cnames, cname)
	}
	for cname := range config.DiscoveryNodes {
		cnames = append(cnames, cname)
	}
	for cname := range config.IdentityService {
		cnames = append(cnames, cname)
	}
	runCommand("docker", cnames...)
	if config.Network.Devnet {
		downDevnetDocker()
	}
}

func registerDevnetNodes(config *conf.ContextConfig) {
	// for _, cc := range config.CreatorNodes {
	// 	register.RegisterNode(
	// 		"content-node",
	// 		cc.Host,
	// 		config.Network.EthMainnetRpc,
	// 		"0xdcB2fC9469808630DD0744b0adf97C0003fC29B2", // hardcoded ganache address
	// 		"0xABbfF712977dB51f9f212B85e8A4904c818C2b63", // "
	// 		cc.OperatorWallet,
	// 		cc.OperatorPrivateKey,
	// 	)
	// }
	// for _, dc := range config.DiscoveryNodes {
	// 	register.RegisterNode(
	// 		"content-node",
	// 		dc.Host,
	// 		config.Network.EthMainnetRpc,
	// 		"0xdcB2fC9469808630DD0744b0adf97C0003fC29B2", // hardcoded ganache address
	// 		"0xABbfF712977dB51f9f212B85e8A4904c818C2b63", // "
	// 		dc.OperatorWallet,
	// 		dc.OperatorPrivateKey,
	// 	)
	// }
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
