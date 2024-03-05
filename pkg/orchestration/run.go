package orchestration

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/AudiusProject/audius-d/pkg/register"
)

func StartDevnet(_ *conf.ContextConfig) error {
	return startDevnetDocker()
}

func DownDevnet(_ *conf.ContextConfig) error {
	return downDevnetDocker()
}

func RunAudiusNodes(nodes map[string]conf.NodeConfig, network conf.NetworkConfig, await bool, audiusdTagOverride string) {
	// Handle devnet-specific setup
	if network.DeployOn == conf.Devnet {
		if err := startDevnetDocker(); err != nil {
			logger.Warnf("Failed to start devnet: %s", err.Error())
		}

		// register all content nodes
		for host, nodeConfig := range nodes {
			if nodeConfig.Type != conf.Creator {
				continue
			}
			err := register.RegisterNode(
				"content-node",
				fmt.Sprintf("https://%s", host),
				"http://localhost:8546",
				register.GanacheAudiusTokenAddress,
				register.GanacheContractRegistryAddress,
				nodeConfig.Wallet,
				nodeConfig.PrivateKey,
			)
			if err != nil {
				logger.Warnf("Failed to register creator node: %s", err)
			}
		}
	}

	for host, nodeConfig := range nodes {
		if err := runNode(
			host,
			nodeConfig,
			network,
			audiusdTagOverride,
		); err != nil {
			logger.Warnf("Error encountered starting node %s: %s", host, err.Error())
		}
	}

	if await {
		awaitHealthy(nodes)
	}
}

func RunDownNodes(nodes map[string]conf.NodeConfig) {
	for host := range nodes {
		if err := downDockerNode(host); err != nil {
			logger.Warnf("Error encountered spinning down %s: %s", host, err.Error())
		} else {
			logger.Infof("Node %s spun down.", host)
		}
	}
}

func execLocal(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func execRemote(host string, stdout io.Writer, stderr io.Writer, command string, args ...string) error {
	cmd := exec.Command("ssh", append([]string{host, command}, args...)...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
