package conf

type DiscoveryConfig struct {
	BaseServerConfig
	// discovery specific stuff here
}

// method to convert a discovery config and an associated network config into
// an audius-docker-compose compatible .env file
// the output map of this method is prepared to be written to an override.env
// or exported directly
func (dp *DiscoveryConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	// discovery specific mappings
	overrideEnv["testKey"] = "value"

	return overrideEnv
}
