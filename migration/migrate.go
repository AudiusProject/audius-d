package migration

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func MigrateAudiusDockerCompose(path string) error {
	if err := assertRepoPath(path); err != nil {
		return err
	}

	nodeType, err := determineNodeType(path)
	if err != nil {
		return err
	}

	env, err := readOverrideEnv(path, nodeType)
	spew.Dump(env)

	return nil
}

// checks that the audius-docker-compose repo is at the path
// provided to the cmd
func assertRepoPath(path string) error {
	log.Printf("validating repo at `%s`\n", path)
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return nil
}

// determines the audius-docker-compose node type
// based on the existence of an override.env file
func determineNodeType(adcpath string) (string, error) {
	creatorOverride := fmt.Sprintf("%s/creator-node/override.env", adcpath)
	discoveryOverride := fmt.Sprintf("%s/discovery-provider/override.env", adcpath)

	if _, err := os.Stat(creatorOverride); err == nil {
		log.Println("creator node detected, migrating")
		return "creator-node", nil
	}

	if _, err := os.Stat(discoveryOverride); err == nil {
		log.Println("discovery provider detected, migrating")
		return "discovery-provider", nil
	}

	return "", errors.New("neither creator or discovery node detected, aborting migration")
}

func readOverrideEnv(path, nodeType string) (map[string]string, error) {
	orpath := fmt.Sprintf("%s/%s/override.env", path, nodeType)
	return godotenv.Read(orpath)
}

func writeMigratedContextConfig() error {
	return nil
}
