package conf

import "fmt"

type IdentityConfig struct {
	BaseServerConfig
	// TODO: rm duplicates
	LogLevel                           string
	MinimumBalance                     int
	MinimumRelayerBalance              int
	MinimumFunderBalance               int
	RelayerPrivateKey                  string
	RelayerPublicKey                   string
	RelayerWallets                     string
	EthFunderAddress                   string
	EthRelayerWallets                  string
	UserVerifierPrivateKey             string
	UserVerifierPublicKey              string
	DbUrl                              string
	RedisHost                          string
	RedisPort                          int
	AaoEndpoint                        string
	AaoAddress                         string
	Web3Provider                       string
	SecondaryWeb3Provider              string
	RegistryAddress                    string
	EntityManagerAddress               string
	OwnerWallet                        string
	EthProviderUrl                     string
	EthTokenAddress                    string
	EthRegistryAddress                 string
	EthOwnerWallet                     string
	SolanaEndpoint                     string
	SolanaTrackListenCountAddress      string
	SolanaAudiusEthRegistryAddress     string
	SolanaValidSigner                  string
	SolanaFeePayerWallets              string
	SolanaSignerPrivateKey             string
	SolanaMintAddress                  string
	SolanaUSDCMintAddress              string
	SolanaClaimableTokenProgramAddress string
	SolanaRewardsManagerProgramId      string
	SolanaRewardsManagerProgramPDA     string
	SolanaRewardsManagerTokenPDA       string
	SolanaAudiusAnchorDataProgramId    string
}

// method to convert an identity config and an associated network config into
// an audius-docker-compose compatible .env file
// the output map of this method is prepared to be written to an override.env
// or exported directly
func (config *IdentityConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	// identity specific mappings
	overrideEnv["logLevel"] = config.LogLevel
	overrideEnv["minimumBalance"] = fmt.Sprintf("%d", config.MinimumBalance)
	overrideEnv["minimumRelayerBalance"] = fmt.Sprintf("%d", config.MinimumRelayerBalance)
	overrideEnv["minimumFunderBalance"] = fmt.Sprintf("%d", config.MinimumFunderBalance)
	overrideEnv["relayerPrivateKey"] = config.RelayerPrivateKey
	overrideEnv["relayerPublicKey"] = config.RelayerPublicKey
	overrideEnv["relayerWallets"] = config.RelayerWallets
	overrideEnv["ethFunderAddress"] = config.EthFunderAddress
	overrideEnv["ethRelayerWallets"] = config.EthRelayerWallets
	overrideEnv["userVerifierPrivateKey"] = config.UserVerifierPrivateKey
	overrideEnv["userVerifierPublicKey"] = config.UserVerifierPublicKey
	overrideEnv["dbUrl"] = config.DbUrl
	overrideEnv["redisHost"] = config.RedisHost
	overrideEnv["redisPort"] = fmt.Sprintf("%d", config.RedisPort)
	overrideEnv["aaoEndpoint"] = config.AaoEndpoint
	overrideEnv["aaoAddress"] = config.AaoAddress
	overrideEnv["web3Provider"] = config.Web3Provider
	overrideEnv["secondaryWeb3Provider"] = config.SecondaryWeb3Provider
	overrideEnv["registryAddress"] = config.RegistryAddress
	overrideEnv["entityManagerAddress"] = config.EntityManagerAddress
	overrideEnv["ownerWallet"] = config.OwnerWallet
	overrideEnv["ethProviderUrl"] = config.EthProviderUrl
	overrideEnv["ethTokenAddress"] = config.EthTokenAddress
	overrideEnv["ethRegistryAddress"] = config.EthRegistryAddress
	overrideEnv["ethOwnerWallet"] = config.EthOwnerWallet
	overrideEnv["solanaEndpoint"] = config.SolanaEndpoint
	overrideEnv["solanaTrackListenCountAddress"] = config.SolanaTrackListenCountAddress
	overrideEnv["solanaAudiusEthRegistryAddress"] = config.SolanaAudiusEthRegistryAddress
	overrideEnv["solanaValidSigner"] = config.SolanaValidSigner
	overrideEnv["solanaFeePayerWallets"] = config.SolanaFeePayerWallets
	overrideEnv["solanaSignerPrivateKey"] = config.SolanaSignerPrivateKey
	overrideEnv["solanaMintAddress"] = config.SolanaMintAddress
	overrideEnv["solanaUSDCMintAddress"] = config.SolanaUSDCMintAddress
	overrideEnv["solanaClaimableTokenProgramAddress"] = config.SolanaClaimableTokenProgramAddress
	overrideEnv["solanaRewardsManagerProgramId"] = config.SolanaRewardsManagerProgramId
	overrideEnv["solanaRewardsManagerProgramPDA"] = config.SolanaRewardsManagerProgramPDA
	overrideEnv["solanaRewardsManagerTokenPDA"] = config.SolanaRewardsManagerTokenPDA
	overrideEnv["solanaAudiusAnchorDataProgramId"] = config.SolanaAudiusAnchorDataProgramId

	return overrideEnv
}
