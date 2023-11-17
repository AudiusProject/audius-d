package main

import (
	"fmt"

	"github.com/AudiusProject/audius-d/conf"
)

func DeployCreator(name string, config conf.CreatorConfig) error {
	tag := config.Tag

	imageTag := fmt.Sprintf("audius/audius-docker-compose:%s", tag)
	volumes := []string{}
	httpPort := config.HttpPort
	httpsPort := config.HttpsPort

	return nil
}
