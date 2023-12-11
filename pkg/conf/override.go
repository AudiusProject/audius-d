package conf

import "fmt"

/** Mappings of toml config to override .env for each node type. */

func (config *DiscoveryConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)
	overrideEnv["audius_delegate_owner_wallet"] = config.OperatorWallet
	overrideEnv["audius_delegate_private_key"] = config.OperatorPrivateKey
	overrideEnv["audius_discprov_url"] = config.Host

	// network config
	overrideEnv["audius_contracts_registry"] = nc.EthContractsRegistryAddress
	overrideEnv["audius_eth_contracts_registry"] = nc.EthContractsRegistryAddress
	overrideEnv["audius_contracts_entity_manager_address"] = nc.AcdcEntityManagerAddress
	overrideEnv["audius_contracts_nethermind_entity_manager_address"] = nc.AcdcEntityManagerAddress
	overrideEnv["AUDIUS_IS_STAGING"] = fmt.Sprintf("%t", nc.Name == "stage")
	overrideEnv["audius_discprov_env"] = nc.Name
	overrideEnv["audius_discprov_identity_service_url"] = nc.IdentityServiceUrl
	overrideEnv["audius_solana_endpoint"] = nc.SolanaMainnetRpc
	overrideEnv["audius_solana_rewards_manager_account"] = nc.SolanaRewardsManagerAccount
	overrideEnv["audius_solana_rewards_manager_program_address"] = nc.SolanaRewardsManagerProgramAddress
	overrideEnv["audius_solana_signer_group_address"] = nc.SolanaSignerGroupAddress
	overrideEnv["audius_solana_track_listen_count_address"] = nc.SolanaTrackListenCountAddress
	overrideEnv["audius_solana_user_bank_program_address"] = nc.SolanaUserBankProgramAddress
	overrideEnv["audius_solana_waudio_mint"] = nc.SolanaWaudioMint
	overrideEnv["audius_solana_usdc_mint"] = nc.SolanaUsdcMint
	overrideEnv["audius_web3_eth_provider_url"] = nc.EthMainnetRpc
	overrideEnv["audius_web3_localhost"] = nc.AcdcRpc
	overrideEnv["audius_web3_host"] = nc.AcdcRpc
	overrideEnv["audius_web3_port"] = fmt.Sprintf("%d", 443)

	// discovery config
	overrideEnv["audius_redis_url"] = config.CacheUrl
	overrideEnv["audius_db_url"] = config.DatabaseUrl
	overrideEnv["audius_db_url_read_replica"] = config.DatabaseUrl
	overrideEnv["audius_cors_allow_all"] = fmt.Sprintf("%t", config.CorsAllowAll)
	overrideEnv["audius_openresty_enable"] = fmt.Sprintf("%t", config.OpenRestyEnable)
	overrideEnv["audius_discprov_blacklist_block_processing_window"] = fmt.Sprintf("%d", config.BlockProcessingWindowBlacklist)
	overrideEnv["audius_discprov_block_processing_window"] = fmt.Sprintf("%d", config.BlockProcessingWindow)
	overrideEnv["audius_discprov_get_users_cnode_ttl_sec"] = fmt.Sprintf("%d", config.GetUsersCnodeTtlSec)
	overrideEnv["audius_discprov_user_metadata_service_url"] = config.UserMetadataServiceUrl
	overrideEnv["audius_enable_rsyslog"] = fmt.Sprintf("%t", config.EnableRsyslog)
	overrideEnv["audius_gunicorn_worker_class"] = config.GunicornWorkerClass
	overrideEnv["audius_gunicorn_workers"] = fmt.Sprintf("%d", config.GunicornWorkers)
	overrideEnv["audius_solana_rewards_manager_min_slot"] = fmt.Sprintf("%d", config.SolanaRewardsManagerMinSlot)
	overrideEnv["audius_solana_user_bank_min_slot"] = fmt.Sprintf("%d", config.SolanaUserBankMinSlot)

	// notifications
	overrideEnv["audius_discprov_notifications_max_block_diff"] = fmt.Sprintf("%d", config.NotificationsMaxBlockDiff)

	// elasticsearch
	overrideEnv["audius_elasticsearch_url"] = config.ElasticSearchUrl
	overrideEnv["audius_elasticsearch_search_enabled"] = fmt.Sprintf("%t", config.ElasticSearchEnabled)

	// relay
	overrideEnv["audius_aao_endpoint"] = nc.AntiAbuseOracleUrl
	overrideEnv["audius_use_aao"] = fmt.Sprintf("%t", config.RelayUseAntiAbuseOracle)

	return overrideEnv
}

func (config *CreatorConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["creatorNodeEndpoint"] = config.Host
	overrideEnv["delegateOwnerWallet"] = config.OperatorWallet
	overrideEnv["delegatePrivateKey"] = config.OperatorPrivateKey
	overrideEnv["spOwnerWallet"] = config.OperatorRewardsWallet
	overrideEnv["MEDIORUM_ENV"] = nc.Name
	overrideEnv["ethProviderUrl"] = nc.EthMainnetRpc

	overrideEnv["ethRegistryAddress"] = nc.EthContractsRegistryAddress

	overrideEnv["entityManagerAddress"] = nc.AcdcEntityManagerAddress

	return overrideEnv
}

func (config *IdentityConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["solanaClaimableTokenProgramAddress"] = config.SolanaClaimableTokenProgramAddress

	overrideEnv["ethRegistryAddress"] = nc.EthContractsRegistryAddress
	overrideEnv["entityManagerAddress"] = nc.AcdcEntityManagerAddress

	return overrideEnv
}
