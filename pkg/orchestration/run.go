package orchestration

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

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

func NormalizedPrivateKey(host, privateKeyConfigValue string) (string, error) {
	privateKey := privateKeyConfigValue
	if strings.HasPrefix(privateKeyConfigValue, "/") {
		// get key value from file on host
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		if err := execOnHost(host, outBuf, errBuf, "cat", privateKeyConfigValue); err != nil {
			return "", logger.Error(errBuf.String(), err)
		}
		privateKey = strings.TrimSpace(outBuf.String())
	}
	privateKey = strings.TrimPrefix(privateKey, "0x")
	if len(privateKey) != 64 {
		return "", logger.Error("Invalid private key")
	}
	return privateKey, nil
}

func execLocal(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func execOnHost(host string, stdout io.Writer, stderr io.Writer, command string, args ...string) error {
	var cmd *exec.Cmd
	isLocalhost, err := resolvesToLocalhost(host)
	if err != nil {
		return logger.Error("Error determining origin of host:", err)
	} else if isLocalhost {
		cmd = exec.Command(command, args...)
	} else {
		cmd = exec.Command("ssh", append([]string{host, command}, args...)...)
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func resolvesToLocalhost(host string) (bool, error) {
	ips, err := net.LookupHost(host)
	if err != nil {
		return false, logger.Errorf("Cannot resolve host %s: %s", host, err.Error())
	}

	for _, ip := range ips {
		if ip == "127.0.0.1" || ip == "::1" {
			return true, nil
		}
	}
	return false, nil
}
