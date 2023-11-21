package conf

type CreatorConfig struct {
	BaseServerConfig
	DirTemplate        string
	DbUrlTemplate      string
	HostNameTemplate   string
	IdentityService    string
	Web3EthProviderUrl string
	EthProviderUrl     string
	EthRegistryAddress string
}

// method to convert a creator config and an associated network config into
// an audius-docker-compose compatible .env file
// the output map of this method is prepared to be written to an override.env
// or exported directly
func (config *CreatorConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["dirTemplate"] = config.DirTemplate
	overrideEnv["dbUrlTemplate"] = config.DbUrlTemplate
	overrideEnv["hostNameTemplate"] = config.HostNameTemplate
	overrideEnv["identityService"] = config.IdentityService
	overrideEnv["web3EthProviderUrl"] = config.Web3EthProviderUrl
	overrideEnv["ethProviderUrl"] = config.EthProviderUrl
	overrideEnv["ethRegistryAddress"] = config.EthRegistryAddress

	return overrideEnv
}
