package conf

type Config struct {
	Network          NetworkConfig
	CreatorNodes     map[string]CreatorConfig
	DiscoveryNodes   map[string]DiscoveryConfig
	IdentityServices map[string]IdentityConfig
}

// base structure that all server types need
type BaseServerConfig struct {
	// port that will be exposed via audius-docker-compose
	// i.e. what you would curl in a http://{host}:{port}/health_check
	// defaults to port 80
	HttpPort uint
	// port that will be exposed via audius-docker-compose
	// i.e. what you would curl in a http://{host}:{port}/health_check
	// defaults to port 443
	HttpsPort uint
	Host      string

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

	// most servers have a db
	DatabaseUrl string

	// most servers have a cache
	CacheUrl  string
	CachePort uint
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
	AcdcRpc                  string
	AcdcNetworkId            uint
	AcdcPort                 uint
	AcdcEntityManagerAddress string
	AcdcRegistryAddress      string

	EthMainnetRpc             string
	EthMainnetPort            uint
	EthMainnetNetworkId       uint
	EthMainnetRegistryAddress string
	EthMainnetTokenAddress    string

	SolanaMainnetRpc  string
	SolanaMainnetPort uint

	Tag string
}

type NodeConfig interface {
	ToOverrideEnv(nc NetworkConfig) []byte
}
