package conf

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func getConfigBaseDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	confDir := filepath.Join(usr.HomeDir, ".audius")
	return confDir, nil
}

func readConfigFromFile(confFilePath string, configTarget interface{}) error {
	if _, err := os.Stat(confFilePath); err != nil {
		return err
	}
	if _, err := toml.DecodeFile(confFilePath, configTarget); err != nil {
		return err
	}
	return nil
}

func writeConfigToFile(confFilePath string, config interface{}) error {
	file, err := os.OpenFile(confFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = toml.NewEncoder(file).Encode(config); err != nil {
		return err
	}
	return nil
}
