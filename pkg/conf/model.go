package conf

type ExecutionConfig struct {
	ConfigVersion  string
	CurrentContext string
}

type ContextConfig struct {
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
	HttpPort  uint
	HttpsPort uint
	Host      string

	// use an existing override .env file
	// instead of creating one on the fly
	// based on toml
	OverrideEnvPath string

	OperatorPrivateKey    string
	OperatorWallet        string
	OperatorRewardsWallet string
	EthOwnerWallet        string

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

type NetworkType string

const (
	Devnet  NetworkType = "devnet"
	Testnet NetworkType = "testnet"
	Mainnet NetworkType = "mainnet"
)

type NetworkConfig struct {
	// Network that the node should be configured to deploy on.
	// Choose "devnet", "testnet", or "mainnet"
	// "devnet" will automatically spin up local chains and identity service
	DeployOn           NetworkType
	ADCTagOverride     string
	HostDockerInternal string
}

type NodeConfig interface {
	ToOverrideEnv(nc NetworkConfig) []byte
}
