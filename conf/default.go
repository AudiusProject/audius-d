package conf

func NewContextConfig() *ContextConfig {
	return &ContextConfig{
		ConfigVersion:   ConfigVersion,
		Network:         NetworkConfig{},
		CreatorNodes:    map[string]CreatorConfig{},
		DiscoveryNodes:  map[string]DiscoveryConfig{},
		IdentityService: map[string]IdentityConfig{},
	}
}

func NewCreatorConfig() *CreatorConfig {
	base := NewBaseServerConfig("creator-node")
	return &CreatorConfig{
		BaseServerConfig: *base,
	}
}

func NewDiscoveryConfig() *DiscoveryConfig {
	base := NewBaseServerConfig("discovery-provider")
	return &DiscoveryConfig{
		BaseServerConfig: *base,
	}
}

func NewIdentityConfig() *IdentityConfig {
	base := NewBaseServerConfig("identity-service")
	return &IdentityConfig{
		BaseServerConfig: *base,
	}
}

func NewNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		Name: "stage",
		// dind internal tag
		// audius-cli set-tag
		Tag: "latest",
	}
}

func NewBaseServerConfig(nodetype string) *BaseServerConfig {
	base := &BaseServerConfig{
		Host: "http://localhost",
		// audius-docker-compose image tag, external
		Tag: "latest",
	}
	switch nodetype {
	case "creator-node":
		base.ExternalHttpPort = 80
		base.InternalHttpPort = 80
		base.ExternalHttpsPort = 443
		base.InternalHttpsPort = 443
	case "discovery-provider":
		base.ExternalHttpPort = 5000
		base.InternalHttpPort = 5000
		base.ExternalHttpsPort = 5001
		base.InternalHttpsPort = 5001
	case "identity-service":
		base.ExternalHttpPort = 7000
		base.InternalHttpPort = 7000
		base.ExternalHttpsPort = 7001
		base.InternalHttpsPort = 7001
	}
	return &BaseServerConfig{}
}
