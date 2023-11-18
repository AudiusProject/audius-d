package conf

import (
	_ "embed"
	"sync"

	"github.com/BurntSushi/toml"
)

//go:embed templates/default.dev.toml
var defaultDevRaw []byte

//go:embed templates/default.stage.toml
var defaultStageRaw []byte

//go:embed templates/default.prod.toml
var defaultProdRaw []byte

func readDefault(def []byte) Config {
	return sync.OnceValue(func() Config {
		conf := &Config{}
		toml.Unmarshal(def, conf)
		return *conf
	})()
}

func GetDevDefaults() Config {
	return readDefault(defaultDevRaw)
}

func GetStageDefaults() Config {
	return readDefault(defaultStageRaw)
}

func GetProdDefaults() Config {
	return readDefault(defaultProdRaw)
}
