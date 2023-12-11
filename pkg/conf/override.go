package conf

/** Mappings of toml config to override .env for each node type. */

func (config *CreatorConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["creatorNodeEndpoint"] = config.Host
	overrideEnv["delegateOwnerWallet"] = config.OperatorWallet
	overrideEnv["delegatePrivateKey"] = config.OperatorPrivateKey
	overrideEnv["spOwnerWallet"] = config.OperatorRewardsWallet
	overrideEnv["MEDIORUM_ENV"] = config.MediorumEnv
	overrideEnv["ethProviderUrl"] = nc.EthMainnetRpc

	overrideEnv["ethRegistryAddress"] = nc.EthContractsRegistryAddress

	overrideEnv["entityManagerAddress"] = nc.AcdcEntityManagerAddress

	return overrideEnv
}

func (config *DiscoveryConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["audius_discprov_url"] = config.Host
	overrideEnv["audius_delegate_owner_wallet"] = config.OperatorWallet
	overrideEnv["audius_delegate_private_key"] = config.OperatorPrivateKey

	overrideEnv["audius_contracts_registry"] = nc.EthContractsRegistryAddress
	overrideEnv["audius_contracts_entity_manager_address"] = nc.AcdcEntityManagerAddress

	return overrideEnv
}

func (config *IdentityConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["solanaClaimableTokenProgramAddress"] = config.SolanaClaimableTokenProgramAddress

	overrideEnv["ethRegistryAddress"] = nc.EthContractsRegistryAddress
	overrideEnv["entityManagerAddress"] = nc.AcdcEntityManagerAddress

	return overrideEnv
}
