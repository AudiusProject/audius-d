package conf

type DiscoveryConfig struct {
	BaseServerConfig
	AudiusWeb3Host                           string
	AudiusWeb3EthProviderUrl                 string
	AudiusContractsRegistry                  string
	AudiusContractsEntityManagerAddress      string
	AudiusEthContractsRegistry               string
	AudiusEthContractsToken                  string
	AudiusSolanaEndpoint                     string
	AudiusSolanaTrackListenCountAddress      string
	AudiusSolanaSignerGroupAddress           string
	AudiusSolanaUserBankProgramAddress       string
	AudiusSolanaWaudioMint                   string
	AudiusSolanaUsdcMint                     string
	AudiusSolanaRewardsManagerProgramAddress string
	AudiusSolanaRewardsManagerAccount        string
	AudiusSolanaAnchorDataProgramId          string
	AudiusSolanaAnchorAdminStoragePublicKey  string
	AudiusDiscprovDevMode                    string
	AudiusDiscprovLoglevel                   string
}

// method to convert a discovery config and an associated network config into
// an audius-docker-compose compatible .env file
// the output map of this method is prepared to be written to an override.env
// or exported directly
func (config *DiscoveryConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	// discovery specific mappings
	overrideEnv["audius_web3_host"] = config.AudiusWeb3Host
	overrideEnv["audius_web3_eth_provider_url"] = config.AudiusWeb3EthProviderUrl
	overrideEnv["audius_contracts_registry"] = config.AudiusContractsRegistry
	overrideEnv["audius_contracts_entity_manager_address"] = config.AudiusContractsEntityManagerAddress
	overrideEnv["audius_eth_contracts_registry"] = config.AudiusEthContractsRegistry
	overrideEnv["audius_eth_contracts_token"] = config.AudiusEthContractsToken
	overrideEnv["audius_solana_endpoint"] = config.AudiusSolanaEndpoint
	overrideEnv["audius_solana_track_listen_count_address"] = config.AudiusSolanaTrackListenCountAddress
	overrideEnv["audius_solana_signer_group_address"] = config.AudiusSolanaSignerGroupAddress
	overrideEnv["audius_solana_user_bank_program_address"] = config.AudiusSolanaUserBankProgramAddress
	overrideEnv["audius_solana_waudio_mint"] = config.AudiusSolanaWaudioMint
	overrideEnv["audius_solana_usdc_mint"] = config.AudiusSolanaUsdcMint
	overrideEnv["audius_solana_rewards_manager_program_address"] = config.AudiusSolanaRewardsManagerProgramAddress
	overrideEnv["audius_solana_rewards_manager_account"] = config.AudiusSolanaRewardsManagerAccount
	overrideEnv["audius_solana_anchor_data_program_id"] = config.AudiusSolanaAnchorDataProgramId
	overrideEnv["audius_solana_anchor_admin_storage_public_key"] = config.AudiusSolanaAnchorAdminStoragePublicKey
	overrideEnv["audius_discprov_dev_mode"] = config.AudiusDiscprovDevMode
	overrideEnv["audius_discprov_loglevel"] = config.AudiusDiscprovLoglevel

	return overrideEnv
}
