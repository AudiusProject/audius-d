package orchestration

import (
	"os"
	"os/exec"

	"github.com/AudiusProject/audius-d/conf"
)

func StartDevnet(_ *conf.ContextConfig) {
	startDevnetDocker()
}

func DownDevnet(_ *conf.ContextConfig) {
	downDevnetDocker()
}

func RunAudiusWithConfig(config *conf.ContextConfig) {
	for cname, cc := range config.CreatorNodes {
		creatorVolumes := []string{"/var/k8s/mediorum:/var/k8s/mediorum", "/var/k8s/creator-node-backend:/var/k8s/creator-node-backend", "/var/k8s/creator-node-db:/var/k8s/creator-node-db"}
		override := cc.ToOverrideEnv(config.Network)
		RunNode(*config, cc.BaseServerConfig, override, cname, "creator-node", creatorVolumes)
	}
	for cname, dc := range config.DiscoveryNodes {
		discoveryVolumes := []string{"/var/k8s/discovery-provider-db:/var/k8s/discovery-provider-db", "/var/k8s/discovery-provider-chain:/var/k8s/discovery-provider-chain"}
		override := dc.ToOverrideEnv(config.Network)
		RunNode(*config, dc.BaseServerConfig, override, cname, "discovery-provider", discoveryVolumes)
	}
	for cname, id := range config.IdentityService {
		identityVolumes := []string{"/var/k8s/identity-service-db:/var/lib/postgresql/data"}
		override := id.ToOverrideEnv(config.Network)
		RunNode(*config, id.BaseServerConfig, override, cname, "identity-service", identityVolumes)
	}
}

func RunDown(config *conf.ContextConfig) {
	// easiest way
	cnames := []string{"rm", "-f"}

	for cname, _ := range config.CreatorNodes {
		cnames = append(cnames, cname)
	}
	for cname, _ := range config.DiscoveryNodes {
		cnames = append(cnames, cname)
	}
	for cname, _ := range config.IdentityService {
		cnames = append(cnames, cname)
	}
	runCommand("docker", cnames...)
	downDevnetDocker()
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
