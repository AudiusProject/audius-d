package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	ConfigVersion = "0.1"
)

func ReadOrCreateContextConfig() (*ContextConfig, error) {
	execConf, err := readOrCreateExecutionConfig()
	if err != nil {
		return nil, err
	}
	contextDir, err := getContextBaseDir()
	if err != nil {
		return nil, err
	}
	contextFilePath := filepath.Join(contextDir, execConf.CurrentContext)
	if _, err = os.Stat(contextFilePath); os.IsNotExist(err) {
		fmt.Printf("Context '%s' not found, using default.\n", execConf.CurrentContext)
		createDefaultContextIfNotExists()
		err = UseContext("default")
		if err != nil {
			return nil, err
		}
		contextFilePath = filepath.Join(contextDir, "default")
	}

	var ctx ContextConfig
	err = readConfigFromFile(contextFilePath, &ctx)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}

func getContextBaseDir() (string, error) {
	confBaseDir, err := getConfigBaseDir()
	if err != nil {
		return "", err
	}
	contextDir := filepath.Join(confBaseDir, "contexts")

	// MkdirAll is idempotent
	// Ensure directory exists before handing it off
	err = os.MkdirAll(contextDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return contextDir, nil
}

func readOrCreateExecutionConfig() (ExecutionConfig, error) {
	var execConf ExecutionConfig
	confDir, err := getConfigBaseDir()
	if err != nil {
		return execConf, err
	}

	execConfFilePath := filepath.Join(confDir, "audius")
	if _, err := os.Stat(execConfFilePath); os.IsNotExist(err) {
		fmt.Println("No existing config found at ~/.audius, creating new.")
		if err = createExecutionConfig(execConfFilePath); err != nil {
			return execConf, err
		}

	}

	if err = readConfigFromFile(execConfFilePath, &execConf); err != nil {
		fmt.Printf("Failed to read execution config: %s\nAttempting to recreate...\n", err)
		if err = createExecutionConfig(execConfFilePath); err != nil {
			return execConf, err
		}
		if err = readConfigFromFile(execConfFilePath, &execConf); err != nil {
			return execConf, err
		}
	}
	return execConf, nil
}

func GetCurrentContextName() (string, error) {
	execConf, err := readOrCreateExecutionConfig()
	if err != nil {
		return "", err
	}
	return execConf.CurrentContext, nil
}

func GetContexts() ([]string, error) {
	ctxDir, err := getContextBaseDir()
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(ctxDir)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, file := range files {
		if !file.IsDir() {
			ret = append(ret, file.Name())
		}
	}
	return ret, nil
}

func UseContext(ctxName string) error {
	ctxDir, err := getContextBaseDir()
	if err != nil {
		return err
	}
	confBaseDir, err := getConfigBaseDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(ctxDir, ctxName)); os.IsNotExist(err) {
		fmt.Printf("No context named %s\n", ctxName)
		return nil
	}

	execConf, err := readOrCreateExecutionConfig()
	if err != nil {
		return err
	}

	// set new name and rewrite file
	execConf.CurrentContext = ctxName

	execConfFilePath := filepath.Join(confBaseDir, "audius")
	if err = writeConfigToFile(execConfFilePath, &execConf); err != nil {
		return err
	}
	return nil
}

func DeleteContext(ctxName string) error {
	ctxDir, err := getContextBaseDir()
	if err != nil {
		return err
	}
	ctxFilepath := filepath.Join(ctxDir, ctxName)
	if _, err := os.Stat(ctxFilepath); os.IsNotExist(err) {
		fmt.Printf("No context named %s\n", ctxName)
		return nil
	}
	if err := os.Remove(ctxFilepath); err != nil {
		return err
	}
	return nil
}

func readConfigFromContext(contextName string, configTarget *ContextConfig) error {
	contextBaseDir, err := getContextBaseDir()
	if err != nil {
		return err
	}
	err = readConfigFromFile(filepath.Join(contextBaseDir, contextName), configTarget)
	if err != nil {
		return err
	}
	return nil
}

func writeConfigToContext(ctxName string, ctxConfig *ContextConfig) error {
	ctxBaseDir, err := getContextBaseDir()
	if err != nil {
		return err
	}
	ctxConfig.ConfigVersion = ConfigVersion
	err = writeConfigToFile(filepath.Join(ctxBaseDir, ctxName), ctxConfig)
	return err
}

func writeConfigToCurrentContext(ctxConfig *ContextConfig) error {
	ctxName, err := GetCurrentContextName()
	if err != nil {
		return err
	}
	return writeConfigToContext(ctxName, ctxConfig)
}

func createContextFromTemplate(name string, templateFilePath string) error {
	var ctxConfig ContextConfig
	if err := readConfigFromFile(templateFilePath, &ctxConfig); err != nil {
		return err
	}
	if err := writeConfigToContext(name, &ctxConfig); err != nil {
		return err
	}
	return nil
}

func createExecutionConfig(confFilePath string) error {
	execConfig := ExecutionConfig{
		ConfigVersion:  ConfigVersion,
		CurrentContext: "default",
	}
	err := writeConfigToFile(confFilePath, &execConfig)
	return err
}

func createDefaultContextIfNotExists() error {
	contextDir, err := getContextBaseDir()
	if err != nil {
		return err
	}

	var conf ContextConfig
	if err = readConfigFromContext("default", &conf); err == nil {
		return nil
	}

	conf = ContextConfig{
		ConfigVersion: ConfigVersion,
		Network: NetworkConfig{
			Name: "stage",
		},
	}

	file, err := os.OpenFile(filepath.Join(contextDir, "default"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = toml.NewEncoder(file).Encode(conf); err != nil {
		return err
	}
	return nil
}
