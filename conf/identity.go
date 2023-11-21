package conf

type IdentityConfig struct {
	BaseServerConfig
	// identity specific stuff here
}

// method to convert an identity config and an associated network config into
// an audius-docker-compose compatible .env file
// the output map of this method is prepared to be written to an override.env
// or exported directly
func (id *IdentityConfig) ToOverrideEnv(nc NetworkConfig) map[string]string {
	overrideEnv := make(map[string]string)

	// identity specific mappings
	overrideEnv["testKey"] = "value"

	return overrideEnv
}
