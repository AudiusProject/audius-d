package conf

/** Mappings of toml config to override .env for each node type. */

func (config *DiscoveryConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)
	for k, v := range config.OverrideConfig {
		overrideEnv[k] = v
	}
	overrideEnv["audius_delegate_owner_wallet"] = config.OperatorWallet
	overrideEnv["audius_delegate_private_key"] = config.OperatorPrivateKey
	overrideEnv["audius_discprov_url"] = config.Host
	if nc.Devnet {
		overrideEnv["comms_sandbox"] = "true"
	}

	// Everything else we don't yet capture in audius-d models
	for k, v := range config.OverrideConfig {
		overrideEnv[k] = v
	}

	// Remove empty configs
	result := make(map[string]string)
	for k, v := range overrideEnv {
		if v != "" {
			result[k] = v
		}
	}
	return result
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

	overrideEnv["userVerifierPrivateKey"] = config.UserVerifierPrivateKey
	overrideEnv["solanaFeePayerWallets"] = config.SolanaFeePayerWallets
	overrideEnv["ethRelayerWallets"] = config.EthRelayerWallets
	overrideEnv["relayerPrivateKey"] = config.RelayerPrivateKey
	overrideEnv["solanaSignerPrivateKey"] = config.SolanaSignerPrivateKey
	overrideEnv["relayerWallets"] = config.RelayerWallets

	return overrideEnv
}
