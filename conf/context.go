package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type audiusContextKey string

const (
	ConfigVersion = "0.1"
	ContextKey    = audiusContextKey("audius_cobra_context")
)

func ReadOrCreateContextConfig() (*ContextConfig, error) {
	confDir, err := getConfigBaseDir()
	if err != nil {
		return nil, err
	}

	execConfFilePath := filepath.Join(confDir, "audius")
	err = os.MkdirAll(filepath.Join(confDir, "contexts"), os.ModePerm)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(execConfFilePath); os.IsNotExist(err) {
		fmt.Println("No existing config found at ~/.audius, creating new.")
		createExecutionConfig(execConfFilePath)
	}

	var execConf ExecutionConfig
	if err = readExecutionConfig(&execConf); err != nil {
		return nil, err
	}

	contextDir, err := getContextBaseDir()
	if err != nil {
		return nil, err
	}
	contextFilePath := filepath.Join(contextDir, execConf.CurrentContext)
	if _, err = os.Stat(contextFilePath); os.IsNotExist(err) {
		fmt.Printf("Context '%s' not found, creating default.\n", execConf.CurrentContext)
		createDefaultContext(contextFilePath)
		err = UseContext("default")
		if err != nil {
			return nil, err
		}
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
	return contextDir, nil
}

func readExecutionConfig(execConf *ExecutionConfig) error {
	configDir, err := getConfigBaseDir()
	if err != nil {
		return err
	}
	execConfFilePath := filepath.Join(configDir, "audius")
	if err = readConfigFromFile(execConfFilePath, &execConf); err != nil {
		return err
	}
	return nil
}

func GetContext() (string, error) {
	var execConf ExecutionConfig
	if err := readExecutionConfig(&execConf); err != nil {
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

	// verify context to set to actually exists
	if _, err := os.Stat(filepath.Join(ctxDir, ctxName)); os.IsNotExist(err) {
		fmt.Printf("No context named %s\n", ctxName)
		return nil
	}

	var execConf ExecutionConfig
	if err = readExecutionConfig(&execConf); err != nil {
		return err
	}

	// set name to current context
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

func writeConfigToContext(ctxName string, ctxConfig *ContextConfig) error {
	ctxBaseDir, err := getContextBaseDir()
	if err != nil {
		return err
	}
	ctxConfig.ConfigVersion = ConfigVersion
	err = writeConfigToFile(filepath.Join(ctxBaseDir, ctxName), ctxConfig)
	return err
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

func createDefaultContext(contextFilePath string) error {
	conf := ContextConfig{
		ConfigVersion: ConfigVersion,
		Network: NetworkConfig{
			Name: "stage",
		},
		CreatorNodes: map[string]CreatorConfig{},
	}

	file, err := os.OpenFile(contextFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = toml.NewEncoder(file).Encode(conf); err != nil {
		return err
	}
	return nil
}
