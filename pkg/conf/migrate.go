package conf

import (
	"errors"
	"fmt"
	"os"

	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/joho/godotenv"
)

func MigrateAudiusDockerCompose(ctxname, path string) error {
	logger.Infof("migrating audius-docker-compose to context %s", ctxname)

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
	logger.Infof("validating repo at `%s`\n", path)
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
		logger.Info("creator node detected, migrating")
		return "creator-node", nil
	}

	if _, err := os.Stat(discoveryOverride); err == nil {
		logger.Info("discovery provider detected, migrating")
		return "discovery-provider", nil
	}

	return "", errors.New("neither creator or discovery node detected, aborting migration")
}

func readEnv(path, nodeType string) (map[string]string, error) {
	allEnv, err := godotenv.Read(fmt.Sprintf("%s/%s/.env", path, nodeType))
	if err != nil {
		allEnv = make(map[string]string)
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
		Host:           "http://localhost",
		OverrideConfig: make(map[string]string),
		HttpPort:       80,
		HttpsPort:      443,
		Version:        "prerelease",
	}
	net := NetworkConfig{
		DeployOn: Devnet,
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
		case "AUDIUS_DOCKER_COMPOSE_GIT_SHA", "HOST_NAME", "NOTIFICATIONS_TAG", "audius_auto_upgrade_enabled", "autoUpgradeEnabled", "AUDIUS_NATS_ENABLE_JETSTREAM", "IS_V2_ONLY", "MEDIORUM_TAG", "MEDIORUM_PORT", "BACKEND_PORT":
			// dumping ground for unmigrated configs
		case "NETWORK":
			switch v {
			case "dev":
				net.DeployOn = Devnet
				base.Version = "prerelease"
			case "stage":
				net.DeployOn = Testnet
				base.Version = "prerelease"
			case "prod":
				net.DeployOn = Mainnet
				base.Version = "current"
			default:
				net.DeployOn = Devnet
				base.Version = "prerelease"
			}
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
