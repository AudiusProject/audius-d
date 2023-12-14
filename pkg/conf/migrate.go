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

	env, err := readEnv(path, nodeType)
	if err != nil {
		return err
	}

	configContext := NewContextConfig()
	envToContextConfig(nodeType, env, configContext)

	WriteConfigToContext(ctxname, configContext)

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

func readEnv(path, nodeType string) (map[string]string, error) {
	allEnv := make(map[string]string)
	hiddenEnv, err := godotenv.Read(fmt.Sprintf("%s/%s/.env", path, nodeType))
	if err != nil {
		hiddenEnv = make(map[string]string)
	}
	if network := hiddenEnv["NETWORK"]; network == "stage" || network == "prod" {
		networkEnv, err := godotenv.Read(fmt.Sprintf("%s/%s/%s.env", path, nodeType, network))
		if err != nil {
			networkEnv = make(map[string]string)
		}
		for k, v := range networkEnv {
			allEnv[k] = v
		}
	}

	orenv, err := godotenv.Read(fmt.Sprintf("%s/%s/override.env", path, nodeType))
	if err != nil {
		return nil, err
	}
	for k, v := range orenv {
		allEnv[k] = v
	}
	return allEnv, nil
}

func envToContextConfig(nodeType string, env map[string]string, ctx *ContextConfig) {
	base := BaseServerConfig{
		Host:              "http://localhost",
		Tag:               "latest",
		OverrideConfig:    make(map[string]string),
		ExternalHttpPort:  80,
		InternalHttpPort:  80,
		ExternalHttpsPort: 443,
		InternalHttpsPort: 443,
	}
	net := NetworkConfig{
		Name: "stage",
		Tag:  "latest",
	}
	discoveryConf := DiscoveryConfig{}
	creatorConf := CreatorConfig{}
	for k, v := range env {
		switch k {
		case "delegatePrivateKey", "audius_delegate_private_key":
			base.OperatorPrivateKey = v
		case "delegateOwnerWallet", "audius_delegate_owner_wallet":
			base.OperatorWallet = v
			// Allow "spOwnerWallet" to override
			if base.OperatorRewardsWallet == "" {
				base.OperatorRewardsWallet = v
			}
		case "spOwnerWallet":
			base.OperatorRewardsWallet = v
		case "creatorNodeEndpoint", "audius_discprov_url":
			base.Host = v
		case "autoUpgradeEnabled", "audius_auto_upgrade_enabled":
			base.AutoUpgrade = "*/15 * * * *"
		case "NETWORK":
			net.Name = v
		default:
			base.OverrideConfig[k] = v
		}
	}
	if nodeType == "creator-node" {
		creatorConf.BaseServerConfig = base
		ctx.CreatorNodes["creator-node"] = creatorConf
	}
	if nodeType == "discovery-provider" {
		discoveryConf.BaseServerConfig = base
		ctx.DiscoveryNodes["discovery-provider"] = discoveryConf
	}
	ctx.Network = net
}
