package conf

import (
	"log"

	"github.com/AudiusProject/audius-d/utils"
)

type CreatorConfig struct {
	BaseServerConfig
	IdentityService string

	// TODO: /up UI
}

// method to convert a creator config and an associated network config into
// an audius-docker-compose compatible .env file
// the output map of this method is prepared to be written to an override.env
// or exported directly
func (cc *CreatorConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	// TODO: handle error
	ownerWallet, err := utils.GenerateAddress(cc.OperatorPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// creator specific mappings
	overrideEnv["creatorNodeEndpoint"] = cc.Host
	overrideEnv["delegatePrivateKey"] = cc.OperatorPrivateKey
	overrideEnv["delegateOwnerWallet"] = *ownerWallet
	overrideEnv["identityService"] = cc.IdentityService

	return overrideEnv
}
