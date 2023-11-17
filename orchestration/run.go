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
	for _, cc := range config.CreatorNodes {
		runNodeDocker("creator-node", config.Network.Name, cc.Tag, false)
	}
	for _, dc := range config.DiscoveryNodes {
		runNodeDocker("discovery-provider", config.Network.Name, dc.Tag, false)
	}
	if config.IdentityService.Tag != "" {
		runNodeDocker("identity-service", config.Network.Name, config.IdentityService.Tag, false)
	}
}

func RunDown(_ *conf.ContextConfig) {
	runDownDocker()
	downDevnetDocker()
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
