package conf

import (
	"bytes"
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

	// MkdirAll is idempotent
	// Ensure directory exists before handing it off
	err = os.MkdirAll(confDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return confDir, nil
}

func ReadConfigFromFile(confFilePath string, configTarget interface{}) error {
	if _, err := os.Stat(confFilePath); err != nil {
		return err
	}
	if _, err := toml.DecodeFile(confFilePath, configTarget); err != nil {
		return err
	}
	return nil
}

func WriteConfigToFile(confFilePath string, config interface{}) error {
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

func StringifyConfig(config interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return "", err
	}
	return buf.String(), nil
}
