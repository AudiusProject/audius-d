package conf

/** Mappings of toml config to override .env for each node type. */

func (config *DiscoveryConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)
	overrideEnv["audius_delegate_owner_wallet"] = config.OperatorWallet
	overrideEnv["audius_delegate_private_key"] = config.OperatorPrivateKey
	overrideEnv["audius_discprov_url"] = config.Host

	return overrideEnv
}

func (config *CreatorConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)
	overrideEnv["creatorNodeEndpoint"] = config.Host
	overrideEnv["delegateOwnerWallet"] = config.OperatorWallet
	overrideEnv["delegatePrivateKey"] = config.OperatorPrivateKey
	overrideEnv["spOwnerWallet"] = config.OperatorRewardsWallet
	overrideEnv["ethOwnerWallet"] = config.OperatorRewardsWallet

	return overrideEnv
}

func (config *IdentityConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	return overrideEnv
}
