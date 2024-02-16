package orchestration

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/joho/godotenv"
)

// deploys a server node generically
func runNode(
	host string,
	config conf.NodeConfig,
	nconf conf.NetworkConfig,
	override map[string]string,
	internalVolumes []string,
	audiusdTag string,
) error {
	dockerClient, err := getDockerClient(host)
	if err != nil {
		return logger.Error(err)
	}
	defer dockerClient.Close()

	if isContainerRunning(dockerClient, host) {
		logger.Infof("Audius container already running on %s", host)
		return nil
	} else if isContainerNameInUse(dockerClient, host) {
		logger.Infof("stopped container exists on %s, removing and starting with current config", host)
		if err := removeContainerByName(dockerClient, host); err != nil {
			return logger.Error(err)
		}
	}

	if audiusdTag == "" {
		audiusdTag = "default"
	}

	containerConfig := &container.Config{
		Image: fmt.Sprintf("audius/audius-d:%s", audiusdTag),
	}
	hostConfig := &container.HostConfig{
		Privileged: true,
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeVolume,
				Source: fmt.Sprintf("audius-d-%s", host),
				Target: "/var/lib/docker",
			},
		},
	}
	for _, vol := range internalVolumes {
		splitVols := strings.Split(vol, ":")
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: splitVols[0],
			Target: splitVols[1],
		})
	}

	var port uint = 80
	var tlsPort uint = 443
	if config.HttpPort != 0 {
		port = config.HttpPort
	}
	if config.HttpsPort != 0 {
		tlsPort = config.HttpsPort
	}
	httpPorts := fmt.Sprintf("%d:%d", port, port)
	httpsPorts := fmt.Sprintf("%d:%d", tlsPort, tlsPort)
	portSet, portBindings, err := nat.ParsePortSpecs([]string{httpPorts, httpsPorts})
	if err != nil {
		return logger.Error(err)
	}
	containerConfig.ExposedPorts = portSet
	hostConfig.PortBindings = portBindings

	if nconf.DeployOn == conf.Devnet {
		hostConfig.NetworkMode = "deployments_devnet"
		hostConfig.ExtraHosts = []string{
			"creator-1.devnet.audius-d:172.100.0.1",
			"discovery-1.devnet.audius-d:172.100.0.1",
			"identity.devnet.audius-d:172.100.0.1",
			"eth-ganache.devnet.audius-d:172.100.0.1",
			"acdc-ganache.devnet.audius-d:172.100.0.1",
			"solana-test-validator.devnet.audius-d:172.100.0.1",
		}
		containerConfig.Env = []string{"HOST_DOCKER_INTERNAL=172.100.0.1"}
	}

	// pull audius-d
	pullResp, err := dockerClient.ImagePull(context.Background(), containerConfig.Image, types.ImagePullOptions{})
	if err != nil {
		return logger.Error("Failed to pull image:", err)
	}
	defer pullResp.Close()
	scanner := bufio.NewScanner(pullResp)
	for scanner.Scan() {
		logger.Info(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return logger.Error("Error ImagePull output:", err)
	}

	// create wrapper container
	createResponse, err := dockerClient.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		nil,
		nil,
		host,
	)
	if err != nil {
		return logger.Error("Failed to create container:", err)
	}
	if err := dockerClient.ContainerStart(
		context.Background(),
		createResponse.ID,
		container.StartOptions{},
	); err != nil {
		return logger.Error("Failed to start container:", err)
	}

	// generate the override.env file locally
	// WARNING: NOT THREAD SAFE
	localOverridePath := "./override.env"
	if err := godotenv.Write(override, localOverridePath); err != nil {
		return logger.Error(err)
	}

	// copy the override.env file to the server and then into the wrapper container
	var adcDir string
	switch config.Type {
	case conf.Creator:
		adcDir = "creator-node"
	case conf.Discovery:
		adcDir = "discovery-provider"
	case conf.Identity:
		adcDir = "identity-service"
	}
	overrideFile, err := os.Open(localOverridePath)
	if err != nil {
		return logger.Error(err)
	}
	defer overrideFile.Close()
	tarReader, err := archive.Tar(overrideFile.Name(), archive.Gzip)
	if err != nil {
		return logger.Error(err)
	}
	if err := dockerClient.CopyToContainer(
		context.Background(),
		createResponse.ID,
		fmt.Sprintf("/root/audius-docker-compose/%s", adcDir),
		tarReader,
		types.CopyToContainerOptions{},
	); err != nil {
		return logger.Error(err)
	}
	if err := os.Remove(localOverridePath); err != nil {
		return logger.Error(err)
	}

	// Configure branch
	var branch string
	switch config.Version {
	case "", "current":
		branch = "main"
	case "prerelease":
		branch = "stage"
	default:
		branch = config.Version
	}
	if err := audiusCli(dockerClient, host, "pull-reset", branch); err != nil {
		return logger.Error(err)
	}

	// auto update hourly, starting 55 minutes from now (for randomness + prevent updates during CI)
	currentTime := time.Now()
	fiveMinutesAgo := currentTime.Add(-5 * time.Minute)
	if err := audiusCli(dockerClient, host, "auto-upgrade", fmt.Sprintf("%d * * * *", fiveMinutesAgo.Minute())); err != nil {
		return logger.Error(err)
	}
	if err := dockerExec(dockerClient, host, "crond"); err != nil {
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
	audiusCli(dockerClient, host, "set-network", network)

	if err := audiusCli(dockerClient, host, "launch", "-y", adcDir); err != nil {
		return logger.Error(err)
	}

	return nil
}

func isContainerRunning(dockerClient *client.Client, containerName string) bool {
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		logger.Warn(err.Error())
		return false
	}
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName && c.State == "running" {
				return true
			}
		}
	}
	return false
}

func isContainerNameInUse(dockerClient *client.Client, containerName string) bool {
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		logger.Warn(err.Error())
		return false
	}
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				return true
			}
		}
	}
	return false
}

func removeContainerByName(dockerClient *client.Client, containerName string) error {
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return err
	}
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				err := dockerClient.ContainerRemove(context.Background(), c.ID, container.RemoveOptions{Force: true})
				return err
			}
		}
	}
	logger.Warnf("Container %s does not exist.", containerName)
	return nil
}

func startDevnetDocker() error {
	logger.Info("Starting local eth, sol, and acdc chains")
	err := execLocal("docker", "compose", "-f", "./deployments/docker-compose.devnet.yml", "up", "-d")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	return nil
}

func downDockerNode(host string) error {
	dockerClient, err := getDockerClient(host)
	if err != nil {
		return logger.Error(err)
	}
	defer dockerClient.Close()

	if err := audiusCli(dockerClient, host, "down"); err != nil {
		logger.Warnf("Failed to spin down internal services on host %s: %s", host, err.Error())
	}
	if err := removeContainerByName(dockerClient, host); err != nil {
		return logger.Error(err)
	}
	return nil
}

func downDevnetDocker() error {
	if err := execLocal("docker", "compose", "-f", "./deployments/docker-compose.devnet.yml", "down"); err != nil {
		return err
	}
	return nil
}

func audiusCli(dockerClient *client.Client, host string, args ...string) error {
	cmds := []string{"audius-cli"}
	cmds = append(cmds, args...)
	return dockerExec(dockerClient, host, cmds...)
}

func dockerExec(dockerClient *client.Client, host string, cmds ...string) error {
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmds,
	}
	resp, err := dockerClient.ContainerExecCreate(context.Background(), host, execConfig)
	if err != nil {
		return err
	}
	execResp, err := dockerClient.ContainerExecAttach(context.Background(), resp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer execResp.Close()
	scanner := bufio.NewScanner(execResp.Reader)
	for scanner.Scan() {
		logger.Info(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return logger.Error("Error reading command output:", err)
	}

	return nil
}

func getDockerClient(host string) (*client.Client, error) {
	isLocalhost, err := resolvesToLocalhost(host)
	if err != nil {
		return nil, err
	} else if isLocalhost {
		return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	} else {
		hostScheme := fmt.Sprintf("ssh://%s", host)
		helper, err := connhelper.GetConnectionHelper(hostScheme)
		if err != nil {
			return nil, err
		}
		return client.NewClientWithOpts(
			client.WithHTTPClient(
				&http.Client{
					Transport: &http.Transport{
						DialContext: helper.Dialer,
					},
				},
			),
			client.WithHost(helper.Host),
			client.WithDialContext(helper.Dialer),
			client.WithAPIVersionNegotiation(),
		)
	}
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
