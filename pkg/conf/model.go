package conf

type ExecutionConfig struct {
	ConfigVersion  string
	CurrentContext string
}

type ContextConfig struct {
	ContextName     string
	ConfigVersion   string
	Network         NetworkConfig
	CreatorNodes    map[string]CreatorConfig
	DiscoveryNodes  map[string]DiscoveryConfig
	IdentityService map[string]IdentityConfig
}

func NewContextConfig() *ContextConfig {
	return &ContextConfig{
		ConfigVersion:   ConfigVersion,
		Network:         NetworkConfig{},
		CreatorNodes:    map[string]CreatorConfig{},
		DiscoveryNodes:  map[string]DiscoveryConfig{},
		IdentityService: map[string]IdentityConfig{},
	}
}

// base structure that all server types need
type BaseServerConfig struct {
	// port that will be exposed via audius-docker-compose
	// i.e. what you would curl in a http://{host}:{port}/health_check
	// defaults to port 80
	InternalHttpPort  uint
	ExternalHttpPort  uint
	InternalHttpsPort uint
	ExternalHttpsPort uint
	Host              string

	// the tag that will be pulled from dockerhub
	// "latest", "stage", "prod", etc may have specific behavior
	// git hashes are also eligible
	Tag string

	OperatorPrivateKey    string
	OperatorWallet        string
	OperatorRewardsWallet string
}

type CreatorConfig struct {
	BaseServerConfig `mapstructure:",squash"`
	// creator specific stuff here
}

func (config *CreatorConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["creatorNodeEndpoint"] = config.Host
	overrideEnv["delegateOwnerWallet"] = config.OperatorWallet
	overrideEnv["delegatePrivateKey"] = config.OperatorPrivateKey
	overrideEnv["spOwnerWallet"] = config.OperatorRewardsWallet

	return overrideEnv
}

type DiscoveryConfig struct {
	BaseServerConfig `mapstructure:",squash"`
}

func (config *DiscoveryConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["audius_discprov_url"] = config.Host
	overrideEnv["audius_delegate_owner_wallet"] = config.OperatorWallet
	overrideEnv["audius_delegate_private_key"] = config.OperatorPrivateKey

	return overrideEnv
}

type IdentityConfig struct {
	BaseServerConfig `mapstructure:",squash"`
	// identity specific stuff here
	SolanaClaimableTokenProgramAddress string
}

func (config *IdentityConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	overrideEnv["solanaClaimableTokenProgramAddress"] = config.SolanaClaimableTokenProgramAddress

	return overrideEnv
}

type NetworkConfig struct {
	// name of the network this/these server(s) belong to
	// analogous to "audius-cli set-network"
	// "dev", "stage", "prod", etc may have specific behavior
	// for a private network set this to any valid string that
	// doesn't have specific behavior
	Name string

	// host that running servers will use to talk to the acdc network
	// example: devnet would have a http://acdc-ganache type string
	AcdcHost string

	// same as AcdcHost but the port if applicable
	AcdcPort uint

	EthMainnetHost string
	EthMainnetPort uint

	SolanaMainnetHost string
	SolanaMainnetPort uint

	Tag string
}

type NodeConfig interface {
	ToOverrideEnv(nc NetworkConfig) []byte
}
