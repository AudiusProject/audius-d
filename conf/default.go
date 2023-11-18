package conf

import (
	"embed"
	_ "embed"
	"sync"

	"github.com/BurntSushi/toml"
)

//go:embed templates/default.dev.toml
var defaultDevRaw []byte

//go:embed templates/default.stage.toml
var defaultStageRaw embed.FS

//go:embed templates/default.prod.toml
var defaultProdRaw embed.FS

func GetDevDefaults() Config {
	return sync.OnceValue(func() Config {
		conf := &Config{}
		toml.Unmarshal(defaultDevRaw, conf)
		return *conf
	})()
}

