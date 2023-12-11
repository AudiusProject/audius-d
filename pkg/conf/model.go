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

	DatabaseUrl string
	CacheUrl    string

	// operations
	Register     bool
	AwaitHealthy bool
	AutoUpgrade  bool
}

type CreatorConfig struct {
	BaseServerConfig `mapstructure:",squash"`
	MediorumEnv      string
}

type DiscoveryConfig struct {
	BaseServerConfig `mapstructure:",squash"`

	CorsAllowAll                   bool
	OpenRestyEnable                bool
	BlockProcessingWindowBlacklist uint
	BlockProcessingWindow          uint
	GetUsersCnodeTtlSec            uint
	UserMetadataServiceUrl         string
	EnableRsyslog                  bool
	GunicornWorkerClass            string
	GunicornWorkers                uint
	SolanaUserBankMinSlot          uint
	SolanaRewardsManagerMinSlot    uint

	// notifications
	NotificationsMaxBlockDiff uint
	// elasticsearch
	ElasticSearchUrl     string
	ElasticSearchEnabled bool
	// relay
	RelayUseAntiAbuseOracle bool
	// comms
	CommsDevMode bool
	// trpc (none present)
}

type IdentityConfig struct {
	BaseServerConfig `mapstructure:",squash"`
	// identity specific stuff here
	SolanaClaimableTokenProgramAddress string
}

type NetworkConfig struct {
	// name of the network this/these server(s) belong to
	// analogous to "audius-cli set-network"
	// "dev", "stage", "prod", etc may have specific behavior
	// for a private network set this to any valid string that
	// doesn't have specific behavior
	Name string

	AcdcRpc          string
	EthMainnetRpc    string
	SolanaMainnetRpc string

	// network singletons
	IdentityServiceUrl string
	AntiAbuseOracleUrl string

	Tag string

	// starts up local containers for acdc, eth, and solana rpcs
	Devnet bool

	// eth mainnet config and addresse
	EthMainnetNetworkId         string
	EthContractsRegistryAddress string

	// acdc config and addresses
	AcdcEntityManagerAddress string

	SolanaRewardsManagerAccount        string
	SolanaRewardsManagerProgramAddress string
	SolanaSignerGroupAddress           string
	SolanaTrackListenCountAddress      string
	SolanaUserBankProgramAddress       string
	SolanaWaudioMint                   string
	SolanaUsdcMint                     string
}

type NodeConfig interface {
	ToOverrideEnv(nc NetworkConfig) []byte
}
