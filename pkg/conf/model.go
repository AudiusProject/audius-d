package conf

type ExecutionConfig struct {
	ConfigVersion  string
	CurrentContext string
}

type ContextConfig struct {
	Network NetworkConfig         `yaml:"network"`
	Nodes   map[string]NodeConfig `yaml:"nodes"`
}

func NewContextConfig() *ContextConfig {
	return &ContextConfig{
		Network: NetworkConfig{
			DeployOn: Devnet,
		},
		Nodes: make(map[string]NodeConfig),
	}
}

// base structure that all server types need
type NodeConfig struct {
	Type NodeType `yaml:"type"`

	// i.e. https://{host}:{httpsPort}/health_check
	HttpPort  uint `yaml:"httpPort,omitempty"`
	HttpsPort uint `yaml:"httpsPort,omitempty"`

	// A string of additional port bindings to allow exposing docker-in-docker containers to the host
	// e.g. "5433:5432,9201:9200" would expose the postgres and elastic search dind containers
	//      on the host ports 5433 and 9201 respectively
	HostPorts string `yaml:"hostPorts,omitempty"`

	// One of "current", "prerelease", or an audius-docker-compose git branch
	// "current" corresponds to main adc branch
	// "prelease" corresponds to stage
	// empty string defaults to "current" behavior
	// anything else is the adc branch (for dev purposes)
	Version string `yaml:"version,omitempty"`

	PrivateKey    string `yaml:"privateKey"`
	Wallet        string `yaml:"wallet"`
	RewardsWallet string `yaml:"rewardsWallet"`
	EthWallet     string `yaml:"ethWallet,omitempty"`

	// content storage
	StorageUrl         string `yaml:"storageUrl,omitempty"`
	StorageCredentials string `yaml:"storageCredentials,omitempty"`

	// Stores any as-yet unstructured configuration
	// (for compatibility with audius-docker-compose migrations)
	OverrideConfig map[string]string `yaml:"overrideConfig,omitempty"`
}

func NewNodeConfig(nodeType NodeType) NodeConfig {
	return NodeConfig{
		Type:      nodeType,
		HttpPort:  80,
		HttpsPort: 443,
		Version:   "current",
	}
}

type NetworkType string

const (
	Devnet  NetworkType = "devnet"
	Testnet NetworkType = "testnet"
	Mainnet NetworkType = "mainnet"
)

type NodeType string

const (
	Creator   NodeType = "creator"
	Discovery NodeType = "discovery"
	Identity  NodeType = "identity"
)

type NetworkConfig struct {
	// Network that the node should be configured to deploy on.
	// Choose "devnet", "testnet", or "mainnet"
	// "devnet" will automatically spin up local chains and identity service
	DeployOn NetworkType `yaml:"deployOn"`

	// Optional Infrastructure API credentials
	Infra *Infra `yaml:"infra,omitempty"`
}

type Infra struct {
	CloudflareAPIKey string `yaml:"cloudflareAPIKey,omitempty"`
	CloudflareZoneId string `yaml:"cloudflareZoneId,omitempty"`
	CloudflareTLD    string `yaml:"cloudflareTld,omitempty"`

	AWSAccessKeyID     string `yaml:"awsAccessKeyID,omitempty"`
	AWSSecretAccessKey string `yaml:"awsSecretAccessKey,omitempty"`
	AWSRegion          string `yaml:"awsRegion,omitempty"`
}
