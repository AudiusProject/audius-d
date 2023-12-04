package conf

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func MigrateAudiusDockerCompose(ctxname, path string) error {
	log.Printf("migrating audius-docker-compose to context %s", ctxname)

	if err := assertRepoPath(path); err != nil {
		return err
	}

	nodeType, err := determineNodeType(path)
	if err != nil {
		return err
	}

	env, err := readOverrideEnv(path, nodeType)
	if err != nil {
		return err
	}

	configContext := NewContextConfig()
	envToContextConfig(nodeType, env, configContext)

	writeConfigToContext(ctxname, configContext)

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

func envToContextConfig(nodeType string, env map[string]string, ctx *ContextConfig) {
	if nodeType == "creator-node" {
		creator := NewCreatorConfig()

		creator.OperatorPrivateKey = env["delegatePrivateKey"]
		creator.OperatorWallet = env["delegateOwnerWallet"]
		creator.OperatorRewardsWallet = env["spOwnerWallet"]

		ctx.CreatorNodes["creator-node"] = *creator
	}
	if nodeType == "discovery-provider" {
		discovery := NewDiscoveryConfig()

		discovery.OperatorPrivateKey = env["audius_delegate_private_key"]
		discovery.OperatorWallet = env["audius_delegate_owner_wallet"]
		discovery.OperatorRewardsWallet = env["audius_delegate_owner_wallet"]

		ctx.DiscoveryNodes["discovery-provider"] = *discovery
	}

	ctx.Network = *NewNetworkConfig()
}
