package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
func Run(dc *client.Client, imageName, imageTag, containerName string, containerConfig *container.Config) {
	ensureDirectory("/tmp/dind")

	image := imageName + ":" + imageTag

	// removing running container if name exists
	runningContainerId := FindRunningContainer(dc, containerName)
	if runningContainerId != nil {
		err := dc.ContainerRemove(ctx, *runningContainerId, types.ContainerRemoveOptions{})

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

	hostConf := &container.HostConfig{Privileged: true}

	resp, err := dc.ContainerCreate(ctx, conf, hostConf, nil, nil, "creator-node")

	if err != nil {
		exitWithError("Creating creator-node container failed:", err)
	}

	// run it
	if err := dc.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		exitWithError(err)
	}

	fmt.Printf("%s %s started \n", image, containerName)
}

// finds a container id given an image name and tag
func FindRunningContainer(dc *client.Client, containerName string) *string {
	containers, err := dc.ContainerList(ctx, types.ContainerListOptions{All: true})

	if err != nil {
		exitWithError("Error could not find container for ", err)
	}

	for _, container := range containers {
		runningContainerName := container.Names[0][1:]
		if runningContainerName == containerName {
			return &container.ID
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
		dc.ContainerRemove(ctx, *creator, types.ContainerRemoveOptions{Force: true})
	}

	if discovery != nil {
		fmt.Println("removing discovery-provider")
		dc.ContainerRemove(ctx, *discovery, types.ContainerRemoveOptions{Force: true})
	}
}
