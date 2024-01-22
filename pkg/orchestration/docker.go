package orchestration

import (
	"fmt"
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
func RunNode(
	nconf conf.NetworkConfig,
	serverConfig conf.BaseServerConfig,
	override OverrideEnv,
	containerName string,
	nodeType string,
	internalVolumes []string,
	audiusdTag string,
) error {
	if isContainerRunning(containerName) {
		logger.Infof("container %s already running", containerName)
		return nil
	}

	if isContainerNameInUse(containerName) {
		logger.Infof("container %s already exists, removing and starting with current config", containerName)
		if err := removeContainerByName(containerName); err != nil {
			return logger.Error(err)
		}
	}

	// naive check for now, populate this with existing dotenv during migration step
	useProvidedOverrideEnv := serverConfig.OverrideEnvPath != ""
	providedOverrideEnvVolume := fmt.Sprintf("%s:/root/audius-docker-compose/%s/override.env", serverConfig.OverrideEnvPath, nodeType)
	if useProvidedOverrideEnv {
		internalVolumes = append(internalVolumes, providedOverrideEnvVolume)
	}

	if audiusdTag == "" {
		audiusdTag = "default"
	}
	imageTag := fmt.Sprintf("audius/audius-d:%s", audiusdTag)
	externalVolume := fmt.Sprintf("audius-d-%s", containerName)
	httpPorts := fmt.Sprintf("-p %d:%d", serverConfig.HttpPort, serverConfig.HttpPort)
	httpsPorts := fmt.Sprintf("-p %d:%d", serverConfig.HttpsPort, serverConfig.HttpsPort)
	formattedInternalVolumes := " -v " + strings.Join(internalVolumes, " -v ")

	devnetAddendum := ""
	if nconf.DeployOn == conf.Devnet {
		dockerNetwork := "--network deployments_devnet"
		hostDockerInternal := "-e HOST_DOCKER_INTERNAL=172.100.0.1"
		extraHosts := "--add-host creator-1.devnet.audius-d:172.100.0.1 --add-host discovery-1.devnet.audius-d:172.100.0.1 --add-host identity.devnet.audius-d:172.100.0.1 --add-host eth-ganache.devnet.audius-d:172.100.0.1 --add-host acdc-ganache.devnet.audius-d:172.100.0.1 --add-host solana-test-validator.devnet.audius-d:172.100.0.1"
		devnetAddendum = fmt.Sprintf("%s %s %s", dockerNetwork, hostDockerInternal, extraHosts)
	}

	// assemble wrapper command and run
	// todo: handle https port
	upCmd := fmt.Sprintf("docker run --privileged %s -d -v %s:/var/lib/docker %s %s --name %s %s %s",
		devnetAddendum, externalVolume, httpPorts, httpsPorts, containerName, formattedInternalVolumes, imageTag)
	if err := Sh(upCmd); err != nil {
		return logger.Error(err)
	}

	// generate override based on toml if not provided an existing one
	if !useProvidedOverrideEnv {
		localOverridePath := fmt.Sprintf("./%s-override.env", containerName)
		if err := godotenv.Write(override, localOverridePath); err != nil {
			return logger.Error(err)
		}

		envCmd := fmt.Sprintf("docker cp %s %s:/root/audius-docker-compose/%s/override.env", localOverridePath, containerName, nodeType)
		if err := Sh(envCmd); err != nil {
			return logger.Error(err)
		}

		cmd := fmt.Sprintf(`docker exec %s sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done"`, containerName)
		if err := runCommand("/bin/sh", "-c", cmd); err != nil {
			return logger.Error(err)
		}

		if err := os.Remove(localOverridePath); err != nil {
			return logger.Error(err)
		}

	}

	// Configure branch
	var branch string
	switch serverConfig.Version {
	case "", "current":
		branch = "main"
	case "prerelease":
		branch = "stage"
	default:
		branch = serverConfig.Version
	}
	if err := audiusCli(containerName, "pull-reset", branch); err != nil {
		return logger.Error(err)
	}

	// auto update hourly, starting 55 minutes from now (for randomness + prevent updates during CI)
	currentTime := time.Now()
	fiveMinutesAgo := currentTime.Add(-5 * time.Minute)
	if err := audiusCli(containerName, "auto-upgrade", fmt.Sprintf("%d * * * *", fiveMinutesAgo.Minute())); err != nil {
		return logger.Error(err)
	}
	if err := runCommand("docker", "exec", containerName, "crond"); err != nil {
		return logger.Error(err)
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
		return logger.Error(err)
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
		logger.Warn(err.Error())
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

func isContainerNameInUse(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		logger.Warn(err.Error())
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
		return logger.Error(err)
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
	audCli := []string{"exec", container, "audius-cli"}
	cmds := append(audCli, args...)
	err := runCommand("docker", cmds...)
	if err != nil {
		return logger.Error("Error with audius-cli:", err)
	}
	return nil
}
