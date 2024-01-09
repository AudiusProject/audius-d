package orchestration

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/joho/godotenv"
)

type OverrideEnv = map[string]string

// deploys a server node generically
func RunNode(nconf conf.NetworkConfig, serverConfig conf.BaseServerConfig, override OverrideEnv, containerName string, nodeType string, internalVolumes []string) error {
	if isContainerRunning(containerName) {
		logger.Infof("container %s already running\n", containerName)
		return nil
	}

	if isContainerNameInUse(containerName) {
		logger.Infof("container %s already exists, removing and starting with current config\n", containerName)
		if err := removeContainerByName(containerName); err != nil {
			return err
		}
	}

	// naive check for now, populate this with existing dotenv during migration step
	useProvidedOverrideEnv := serverConfig.OverrideEnvPath != ""
	providedOverrideEnvVolume := fmt.Sprintf("%s:/root/audius-docker-compose/%s/override.env", serverConfig.OverrideEnvPath, nodeType)
	if useProvidedOverrideEnv {
		internalVolumes = append(internalVolumes, providedOverrideEnvVolume)
	}

	version := nconf.ADCTagOverride
	if version == "" {
		version = "latest"
	}
	imageTag := fmt.Sprintf("audius/audius-docker-compose:%s", version)
	externalVolume := fmt.Sprintf("audius-d-%s", containerName)
	httpPorts := fmt.Sprintf("-p %d:%d", serverConfig.HttpPort, serverConfig.HttpPort)
	httpsPorts := fmt.Sprintf("-p %d:%d", serverConfig.HttpsPort, serverConfig.HttpsPort)
	formattedInternalVolumes := " -v " + strings.Join(internalVolumes, " -v ")

	// devnet networking
	dockerNetwork := "--network deployments_devnet"
	hostDockerInternal := "-e HOST_DOCKER_INTERNAL=172.100.0.1"
	extraHosts := "--add-host creator-1.audius-d:172.100.0.1 --add-host discovery-1.audius-d:172.100.0.1 --add-host identity-1.audius-d:172.100.0.1 --add-host eth-ganache.audius-d:172.100.0.1 --add-host acdc-ganache.audius-d:172.100.0.1 --add-host solana-test-validator.audius-d:172.100.0.1"

	// assemble wrapper command and run
	// todo: handle https port
	upCmd := fmt.Sprintf("docker run --privileged %s %s %s -d -v %s:/var/lib/docker %s %s --name %s %s %s",
		dockerNetwork, hostDockerInternal, extraHosts, externalVolume, httpPorts, httpsPorts, containerName, formattedInternalVolumes, imageTag)
	if err := Sh(upCmd); err != nil {
		return err
	}

	// generate override based on toml if not provided an existing one
	if !useProvidedOverrideEnv {
		localOverridePath := fmt.Sprintf("./%s-override.env", containerName)
		if err := godotenv.Write(override, localOverridePath); err != nil {
			return err
		}

		envCmd := fmt.Sprintf("docker cp %s %s:/root/audius-docker-compose/%s/override.env", localOverridePath, containerName, nodeType)
		if err := Sh(envCmd); err != nil {
			return err
		}

		cmd := fmt.Sprintf(`docker exec %s sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done"`, containerName)
		if err := runCommand("/bin/sh", "-c", cmd); err != nil {
			return err
		}

		if err := os.Remove(localOverridePath); err != nil {
			return err
		}

	}

	if version == "latest" {
		minute := rand.Intn(60)
		if err := audiusCli(containerName, "auto-upgrade", fmt.Sprintf("*/%d * * * *", minute)); err != nil {
			// don't propogate audiusCli error, service can still start if auto-upgrade fails
			logger.Infof("Auto upgrade failed: %v\nSkipping", err)
		}
	}

	// set network
	var network string
	switch nconf.DeployOn {
	case conf.Devnet:
		network = "dev"
	case conf.Testnet:
		network = "stage"
	case conf.Mainnet:
		network = "prod"
	default:
		network = "dev"
	}
	audiusCli(containerName, "set-network", network)

	// assemble inner command and run
	startCmd := fmt.Sprintf(`docker exec %s sh -c "cd %s && docker compose up -d"`, containerName, nodeType)
	if err := Sh(startCmd); err != nil {
		return err
	}

	return nil
}

func Sh(cmd string) error {
	logger.Info(cmd)
	return runCommand("/bin/sh", "-c", cmd)
}

func isContainerRunning(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-q", "-f", "name="+containerName)
	output, err := cmd.Output()
	if err != nil {
		logger.Error(err)
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

func isContainerNameInUse(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		logger.Error(err)
		return false
	}

	// Split the output into individual container names
	containers := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Check if the given container name exists in the list
	for _, name := range containers {
		if name == containerName {
			return true
		}
	}
	return false
}

func removeContainerByName(containerName string) error {
	cmd := exec.Command("docker", "rm", "-f", containerName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func startDevnetDocker() {
	logger.Info("Creating docker network")
	runCommand("docker", "network", "create", "--subnet=172.100.0.0/16", "--gateway=172.100.0.1", "deployments_devnet")
	logger.Info("Starting local eth, sol, and acdc chains")
	runCommand("docker", "compose", "-f", "./deployments/docker-compose.devnet.yml", "up", "-d")
	time.Sleep(5 * time.Second)
}

func downDevnetDocker() {
	runCommand("docker", "compose", "-f", "./deployments/docker-compose.devnet.yml", "down")
	runCommand("docker", "network", "rm", "deployments_devnet")
}

func audiusCli(container string, args ...string) error {
	audCli := []string{"exec", container, ".venv/bin/python3", "audius-cli"}
	cmds := append(audCli, args...)
	err := runCommand("docker", cmds...)
	if err != nil {
		return logger.Error("Error with audius-cli:", err)
	}
	return nil
}
