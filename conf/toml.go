package conf

import "github.com/BurntSushi/toml"

func ReadToml(filepath string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(filepath, &conf)

	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// panics instead of returning error
func ReadTomlUnsafe(filepath string) *Config {
	toml, err := ReadToml(filepath)
	if err != nil {
		panic(err)
	}
	return toml
}
