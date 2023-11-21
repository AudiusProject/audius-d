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

	// the one key to rule them all
	OperatorPrivateKey string

	// currently only eligible on devnet
	// will automatically register the node if set to true
	Register bool

	// will query `http://{host}:{port}/health_check` until
	// a 2XX response is received,
	AwaitHealthy bool
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
