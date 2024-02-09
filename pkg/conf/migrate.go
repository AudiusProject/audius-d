package conf

import (
	"errors"
	"fmt"
	"net/url"
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
	if err := envToContextConfig(nodeType, env, configContext); err != nil {
		return err
	}

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
func determineNodeType(adcpath string) (NodeType, error) {
	creatorOverride := fmt.Sprintf("%s/creator-node/override.env", adcpath)
	discoveryOverride := fmt.Sprintf("%s/discovery-provider/override.env", adcpath)

	if _, err := os.Stat(creatorOverride); err == nil {
		logger.Info("creator node detected, migrating")
		return Creator, nil
	}

	if _, err := os.Stat(discoveryOverride); err == nil {
		logger.Info("discovery provider detected, migrating")
		return Discovery, nil
	}

	return "", errors.New("neither creator or discovery node detected, aborting migration")
}

func readEnv(path string, nodeType NodeType) (map[string]string, error) {
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

func envToContextConfig(nodeType NodeType, env map[string]string, ctx *ContextConfig) error {
	node := NewNodeConfig(nodeType)
	net := NetworkConfig{
		DeployOn: Devnet,
	}
	host := ""
	for k, v := range env {
		switch k {
		case "delegatePrivateKey", "audius_delegate_private_key":
			node.PrivateKey = v
		case "delegateOwnerWallet", "audius_delegate_owner_wallet":
			node.Wallet = v
			// Allow "spOwnerWallet" to override
			if node.RewardsWallet == "" {
				node.RewardsWallet = v
			}
		case "spOwnerWallet":
			node.RewardsWallet = v
		case "creatorNodeEndpoint", "audius_discprov_url":
			u, err := url.Parse(v)
			if err != nil {
				return err
			}
			host = u.Host
		case "AUDIUS_DOCKER_COMPOSE_GIT_SHA", "HOST_NAME", "NOTIFICATIONS_TAG", "audius_auto_upgrade_enabled", "autoUpgradeEnabled", "AUDIUS_NATS_ENABLE_JETSTREAM", "IS_V2_ONLY", "MEDIORUM_TAG", "MEDIORUM_PORT", "BACKEND_PORT":
			// dumping ground for unmigrated configs
		case "NETWORK":
			switch v {
			case "dev":
				net.DeployOn = Devnet
				node.Version = "prerelease"
			case "stage":
				net.DeployOn = Testnet
				node.Version = "prerelease"
			case "prod":
				net.DeployOn = Mainnet
				node.Version = "current"
			default:
				net.DeployOn = Devnet
				node.Version = "prerelease"
			}
		default:
			node.OverrideConfig[k] = v
		}
	}
	ctx.Nodes[host] = node
	ctx.Network = net
	return nil
}
