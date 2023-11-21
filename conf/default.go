package conf

import (
	_ "embed"
	"sync"

	"github.com/BurntSushi/toml"
)

//go:embed templates/default.devnet.toml
var defaultDevnetRaw []byte

//go:embed templates/default.testnet.toml
var defaultTestnetRaw []byte

//go:embed templates/default.mainnet.toml
var defaultMainnetRaw []byte

// function to parse the embed structs
// and ensure the parsing only happens once
func readDefault(def []byte) ContextConfig {
	return sync.OnceValue(func() ContextConfig {
		conf := &ContextConfig{}
		toml.Unmarshal(def, conf)
		return *conf
	})()
}

func GetDevnetDefaults() ContextConfig {
	return readDefault(defaultDevnetRaw)
}

func GetTestnetDefaults() ContextConfig {
	return readDefault(defaultTestnetRaw)
}

func GetMainnetDefaults() ContextConfig {
	return readDefault(defaultMainnetRaw)
}
