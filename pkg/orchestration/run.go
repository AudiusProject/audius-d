package orchestration

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/AudiusProject/audius-d/pkg/conf"
)

var (
	devnetIdentityServiceContainerName = "identity-1"
	devnetIdentityServiceConfig        = conf.IdentityConfig{
		BaseServerConfig: conf.BaseServerConfig{
			HttpPort:  7000,
			HttpsPort: 7001,
			Host:      "https://identity-1.audius-d",
		},
		RelayerWallets:         `[{"publicKey": "0xaaaa90Fc2bfa70028D6b444BB9754066d9E2703b", "privateKey": "34efbbc0431c7f481cdba15d65bbc9ef47196b9cf38d5c4b30afa2bcf86fafba"}, {"publicKey": "0xBE718F98a5B5a473186eB6E30888F26E72be0b66", "privateKey": "d3426cd10c4e75207bdc4802c551d21faa89a287546c2c6b3d9a0476f34934d2"}]`,
		SolanaSignerPrivateKey: "d242765e718801781440d77572b9dafcdc9baadf0269eff24cf61510ddbf1003",
		UserVerifierPrivateKey: "ebba299e6163ff3208de4e82ce7db09cf7e434847b5bdab723af96ae7c763a0e",
		SolanaFeePayerWallets:  `[{"privateKey":[170, 161, 84, 122, 118, 210, 128, 213, 96, 185, 143, 218, 54, 254, 217, 204, 157, 175, 137, 71, 202, 108, 51, 242, 21, 50, 56, 77, 54, 116, 103, 56, 251, 64, 77, 100, 199, 88, 103, 189, 42, 163, 67, 251, 101, 204, 7, 59, 70, 109, 113, 50, 209, 154, 55, 164, 227, 108, 203, 146, 121, 148, 85, 119]}]`,
		EthRelayerWallets:      `[{"publicKey": "0xE75dEe171b6472cE30358ede946CcDFfCA70b562", "privateKey": "8a7c63d4aea87647f480e4771ea279f90f8e912fcfe907525bc931f531e564ce"}, {"publicKey": "0xBE718F98a5B5a473186eB6E30888F26E72be0b66", "privateKey": "d3426cd10c4e75207bdc4802c551d21faa89a287546c2c6b3d9a0476f34934d2"}, {"publicKey": "0xaaaa90Fc2bfa70028D6b444BB9754066d9E2703b", "privateKey": "34efbbc0431c7f481cdba15d65bbc9ef47196b9cf38d5c4b30afa2bcf86fafba"}]`,
		RelayerPrivateKey:      "34efbbc0431c7f481cdba15d65bbc9ef47196b9cf38d5c4b30afa2bcf86fafba",
	}
)

func StartDevnet(_ *conf.ContextConfig) {
	startDevnetDocker()
}

func DownDevnet(_ *conf.ContextConfig) {
	downDevnetDocker()
}

func RunAudiusWithConfig(config *conf.ContextConfig, await bool) {
	if config.Network.DeployOn == conf.Devnet {
		startDevnetDocker()
	}

	dashboardVolume := "/dashboard-dist:/dashboard-dist"
	esDataVolume := "/esdata:/esdata"

	// mac local volumes need some extra stuff
	// stick into /var/k8s as if these existed then
	if runtime.GOOS == "darwin" {
		esDataVolume = "/var/k8s/esdata:/esdata"
		dashboardVolume = "/var/k8s/dashboard-dist:/dashboard-dist"
	}

	for cname, cc := range config.CreatorNodes {
		creatorVolumes := []string{"/var/k8s/mediorum:/var/k8s/mediorum", "/var/k8s/creator-node-backend:/var/k8s/creator-node-backend", "/var/k8s/creator-node-db:/var/k8s/creator-node-db", "/var/k8s/bolt:/var/k8s/bolt", dashboardVolume}
		override := cc.ToOverrideEnv(config.Network)
		RunNode(config.Network, cc.BaseServerConfig, override, cname, "creator-node", creatorVolumes)
		if await {
			awaitHealthy(cname, cc.Host, cc.HttpsPort)
		}
	}
	for cname, dc := range config.DiscoveryNodes {
		discoveryVolumes := []string{"/var/k8s/discovery-provider-db:/var/k8s/discovery-provider-db", "/var/k8s/discovery-provider-chain:/var/k8s/discovery-provider-chain", "/var/k8s/bolt:/var/k8s/bolt", esDataVolume, dashboardVolume}
		override := dc.ToOverrideEnv(config.Network)
		RunNode(config.Network, dc.BaseServerConfig, override, cname, "discovery-provider", discoveryVolumes)
		// discovery requires a few extra things
		if config.Network.DeployOn != conf.Devnet {
			audiusCli(cname, "launch-chain")
		}
		if await {
			awaitHealthy(cname, dc.Host, dc.HttpsPort)
		}
	}
	if config.Network.DeployOn == conf.Devnet {
		identityVolumes := []string{"/var/k8s/identity-service-db:/var/lib/postgresql/data"}
		RunNode(
			config.Network,
			devnetIdentityServiceConfig.BaseServerConfig,
			devnetIdentityServiceConfig.ToOverrideEnv(config.Network),
			devnetIdentityServiceContainerName,
			"identity-service",
			identityVolumes,
		)
		if await {
			awaitHealthy(devnetIdentityServiceContainerName, devnetIdentityServiceConfig.Host, devnetIdentityServiceConfig.HttpsPort)
		}
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
	if config.Network.DeployOn == conf.Devnet {
		downDevnetDocker()
	}
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
