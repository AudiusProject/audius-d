package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// global var
var ctx = context.Background()

// connects to docker host on machine
func ConnectDocker() *client.Client {
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		exitWithError("Error connecting to docker", err)
	}
	return docker
}

// downs the existing container, creates, and starts the new one
func Run(dc *client.Client, imageName, imageTag, containerName string, containerConfig *container.Config, mounts []mount.Mount) {
	ensureDirectory("/tmp/dind")

	image := imageName + ":" + imageTag

	// removing running container if name exists
	runningContainer := FindRunningContainer(dc, containerName)
	if runningContainer != nil {

		runningContainerId := runningContainer.ID
		err := dc.ContainerRemove(ctx, runningContainerId, types.ContainerRemoveOptions{Force: true})

		if err != nil {
			exitWithError("Could not remove container"+image+":", err)
		}
	}

	// pull latest image
	if !localImage {
		reader, err := dc.ImagePull(ctx, image, types.ImagePullOptions{})
		if err != nil {
			exitWithError("Error pulling image:", err)
		}
		defer reader.Close()
	}

	// create container
	var conf *container.Config
	if containerConfig != nil {
		conf = containerConfig
		conf.Image = image
	} else {
		conf = &container.Config{
			Image: image,
		}
	}

	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/tmp/dind",
		Target: "/var/lib/docker",
	})

	// config all node types require
	hostConf := &container.HostConfig{
		Privileged: true,
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%d", port)): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", port),
				},
			},
			nat.Port(fmt.Sprintf("%d", tlsPort)): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", tlsPort),
				},
			},
		},
		Mounts: mounts,
	}

	resp, err := dc.ContainerCreate(ctx, conf, hostConf, nil, nil, "creator-node")

	if err != nil {
		exitWithError("Creating creator-node container failed:", err)
	}

	// run it
	if err := dc.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		exitWithError(err)
	}

	fmt.Printf("%s %s started \n", image, containerName)

	execResp, err := dc.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
		Cmd: []string{"cd creator-node && docker compose up"},
	})

	if err != nil {
		exitWithError("Error exec failed:", err)
	}

	fmt.Printf("%s", execResp.ID)
}

// finds a container id given an image name and tag
func FindRunningContainer(dc *client.Client, containerName string) *types.Container {
	containers, err := dc.ContainerList(ctx, types.ContainerListOptions{All: true})

	if err != nil {
		exitWithError("Error could not find container for ", err)
	}

	for _, container := range containers {
		runningContainerName := container.Names[0][1:]
		if runningContainerName == containerName {
			return &container
		}
	}
	return nil
}

func ensureDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			exitWithError("Failed to create directory:", err)
		}
	}
}

// downs all known node types
func DownAll(dc *client.Client) {
	creator := FindRunningContainer(dc, "creator-node")
	discovery := FindRunningContainer(dc, "discovery-provider")

	if creator != nil {
		fmt.Println("removing creator-node")
		dc.ContainerRemove(ctx, creator.ID, types.ContainerRemoveOptions{Force: true})
	}

	if discovery != nil {
		fmt.Println("removing discovery-provider")
		dc.ContainerRemove(ctx, discovery.ID, types.ContainerRemoveOptions{Force: true})
	}
}
