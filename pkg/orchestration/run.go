package orchestration

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/AudiusProject/audius-d/pkg/register"
)

func StartDevnet(_ *conf.ContextConfig) {
	startDevnetDocker()
}

func DownDevnet(_ *conf.ContextConfig) {
	downDevnetDocker()
}

func RunAudiusWithConfig(config *conf.ContextConfig, await bool, audiusdTagOverride string) {
	if config.Network.DeployOn == conf.Devnet {
		startDevnetDocker()
		for _, cc := range config.CreatorNodes {
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
				logger.Warnf("Failed to register creator node: %s", err)
			}
		}
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
		if err := RunNode(
			config.Network,
			cc.BaseServerConfig,
			override,
			cname,
			"creator-node",
			creatorVolumes,
			audiusdTagOverride,
		); err != nil {
			logger.Warnf("Failure starting node %s: %s\nSkipping", cname, err.Error())
		}
	}
	for cname, dc := range config.DiscoveryNodes {
		discoveryVolumes := []string{"/var/k8s/discovery-provider-db:/var/k8s/discovery-provider-db", "/var/k8s/discovery-provider-chain:/var/k8s/discovery-provider-chain", "/var/k8s/bolt:/var/k8s/bolt", esDataVolume, dashboardVolume}
		override := dc.ToOverrideEnv(config.Network)
		if err := RunNode(
			config.Network,
			dc.BaseServerConfig,
			override,
			cname,
			"discovery-provider",
			discoveryVolumes,
			audiusdTagOverride,
		); err != nil {
			logger.Warnf("Failure starting node %s: %s\nSkipping", cname, err.Error())
		}
		// discovery requires a few extra things
		if config.Network.DeployOn != conf.Devnet {
			audiusCli(cname, "launch-chain")
		}
	}
	for cname, id := range config.IdentityService {
		identityVolumes := []string{"/var/k8s/identity-service-db:/var/lib/postgresql/data"}
		override := id.ToOverrideEnv(config.Network)
		if err := RunNode(
			config.Network,
			id.BaseServerConfig,
			override,
			cname,
			"identity-service",
			identityVolumes,
			audiusdTagOverride,
		); err != nil {
			logger.Warnf("Failure starting node %s: %s\nSkipping", cname, err.Error())
		}
	}

	if await {
		awaitHealthy(config)
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
