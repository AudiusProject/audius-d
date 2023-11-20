package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/AudiusProject/audius-d/conf"
	"github.com/joho/godotenv"
)

func Deploy(config conf.Config) error {
	for name, creatorConf := range config.CreatorNodes {
		DeployCreator(config, name, creatorConf)
	}
	return nil
}

func DeployCreator(config conf.Config, creatorName string, creatorConfig conf.CreatorConfig) error {
	// assemble and format config
	containerName := creatorName
	imageTag := fmt.Sprintf("audius/audius-docker-compose:%s", config.Network.Tag)
	externalVolume := fmt.Sprintf("audius-d-%s", containerName)
	internalVolumes := []string{"/var/k8s/mediorum:/var/k8s/mediorum", "/var/k8s/creator-node-backend:/var/k8s/creator-node-backend", "/var/k8s/creator-node-db:/var/k8s/creator-node-db"}
	httpPort := creatorConfig.HttpPort
	httpsPort := creatorConfig.HttpsPort
	formattedInternalVolumes := " -v " + strings.Join(internalVolumes, " -v ")

	// assemble wrapper command and run
	upCmd := fmt.Sprintf("docker run --privileged -d -v %s:/var/lib/docker -p %d:80 -p %d:443 --name %s %s %s", externalVolume, httpPort, httpsPort, containerName, formattedInternalVolumes, imageTag)
	if err := Sh(upCmd); err != nil {
		return err
	}

	// initialize override.env file
	overrides := creatorConfig.ToOverrideEnv(config.Network)
	localOverridePath := fmt.Sprintf("./%s-override.env", containerName)
	if err := godotenv.Write(overrides, localOverridePath); err != nil {
		return err
	}

	envCmd := fmt.Sprintf("docker cp %s %s:/root/audius-docker-compose/creator-node/override.env", localOverridePath, containerName)
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

	// assemble inner command and run
	startCmd := fmt.Sprintf(`docker exec %s sh -c "cd creator-node && docker compose up -d"`, containerName)
	if err := Sh(startCmd); err != nil {
		return err
	}

	return nil
}

func Sh(cmd string) error {
	fmt.Println(cmd)
	return runCommand("/bin/sh", "-c", cmd)
}
