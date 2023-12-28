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
	// use an existing override .env file
	// instead of creating one on the fly
	// based on toml
	OverrideEnvPath string

	OperatorPrivateKey    string
	OperatorWallet        string
	OperatorRewardsWallet string
	EthOwnerWallet        string

	// operations
	Register     bool
	AwaitHealthy bool
	AutoUpgrade  string

	// Stores any as-yet unstructured configuration
	// (for compatibility with audius-docker-compose migrations)
	OverrideConfig map[string]string
}

type CreatorConfig struct {
	BaseServerConfig `mapstructure:",squash"`
}

type DiscoveryConfig struct {
	BaseServerConfig `mapstructure:",squash"`
}

type IdentityConfig struct {
	BaseServerConfig       `mapstructure:",squash"`
	UserVerifierPrivateKey string
	SolanaFeePayerWallets  string
	EthRelayerWallets      string
	RelayerPrivateKey      string
	SolanaSignerPrivateKey string
	RelayerWallets         string
}

type NetworkConfig struct {
	AudiusComposeNetwork string

	AcdcRpc          string
	EthMainnetRpc    string
	SolanaMainnetRpc string

	// network singletons
	IdentityServiceUrl     string
	AntiAbuseOracleUrl     string
	AntiAbuseOracleAddress string

	Tag string

	// starts up local containers for acdc, eth, and solana rpcs
	Devnet bool

	HostDockerInternal string
}

type NodeConfig interface {
	ToOverrideEnv(nc NetworkConfig) []byte
}
