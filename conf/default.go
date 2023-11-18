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

func readDefault(def []byte) Config {
	return sync.OnceValue(func() Config {
		conf := &Config{}
		toml.Unmarshal(def, conf)
		return *conf
	})()
}

func GetDevnetDefaults() Config {
	return readDefault(defaultDevnetRaw)
}

func GetTestnetDefaults() Config {
	return readDefault(defaultTestnetRaw)
}

func GetMainnetDefaults() Config {
	return readDefault(defaultMainnetRaw)
}
