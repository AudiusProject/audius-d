package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	ConfigVersion = "0.1"
)

func ReadOrCreateContextConfig() (*ContextConfig, error) {
	execConf, err := readOrCreateExecutionConfig()
	if err != nil {
		return nil, err
	}
	contextDir, err := GetContextBaseDir()
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
	err = ReadConfigFromFile(contextFilePath, &ctx)
	if err != nil {
		return nil, err
	}

	return &ctx, nil
}

func GetContextBaseDir() (string, error) {
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

	if err = ReadConfigFromFile(execConfFilePath, &execConf); err != nil {
		fmt.Printf("Failed to read execution config: %s\nAttempting to recreate...\n", err)
		if err = createExecutionConfig(execConfFilePath); err != nil {
			return execConf, err
		}
		if err = ReadConfigFromFile(execConfFilePath, &execConf); err != nil {
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
	ctxDir, err := GetContextBaseDir()
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

func ReadContext(ctxName string) (*ContextConfig, error) {
	ctxDir, err := GetContextBaseDir()
	if err != nil {
		return nil, err
	}
	_, err = ioutil.ReadDir(ctxDir)
	if err != nil {
		return nil, err
	}
	contextFilePath := filepath.Join(ctxDir, ctxName)
	var ctx ContextConfig
	err = ReadConfigFromFile(contextFilePath, &ctx)
	if err != nil {
		return nil, err
	}
	return &ctx, nil
}

func UseContext(ctxName string) error {
	ctxDir, err := GetContextBaseDir()
	if err != nil {
		return err
	}
	confBaseDir, err := getConfigBaseDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(ctxDir, ctxName)); os.IsNotExist(err) {
		return fmt.Errorf("No context named %s\n", ctxName)
	}

	execConf, err := readOrCreateExecutionConfig()
	if err != nil {
		return err
	}

	// set new name and rewrite file
	execConf.CurrentContext = ctxName

	execConfFilePath := filepath.Join(confBaseDir, "audius")
	if err = WriteConfigToFile(execConfFilePath, &execConf); err != nil {
		return err
	}
	return nil
}

func DeleteContext(ctxName string) error {
	ctxDir, err := GetContextBaseDir()
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

func ReadConfigFromContext(contextName string, configTarget *ContextConfig) error {
	contextBaseDir, err := GetContextBaseDir()
	if err != nil {
		return err
	}
	err = ReadConfigFromFile(filepath.Join(contextBaseDir, contextName), configTarget)
	if err != nil {
		return err
	}
	return nil
}

func WriteConfigToContext(ctxName string, ctxConfig *ContextConfig) error {
	ctxBaseDir, err := GetContextBaseDir()
	if err != nil {
		return err
	}
	ctxConfig.ConfigVersion = ConfigVersion
	err = WriteConfigToFile(filepath.Join(ctxBaseDir, ctxName), ctxConfig)
	return err
}

func WriteConfigToCurrentContext(ctxConfig *ContextConfig) error {
	ctxName, err := GetCurrentContextName()
	if err != nil {
		return err
	}
	return WriteConfigToContext(ctxName, ctxConfig)
}

func CreateContextFromTemplate(name string, templateFilePath string) error {
	var ctxConfig ContextConfig
	if err := ReadConfigFromFile(templateFilePath, &ctxConfig); err != nil {
		return err
	}
	if err := WriteConfigToContext(name, &ctxConfig); err != nil {
		return err
	}
	return nil
}

func createExecutionConfig(confFilePath string) error {
	execConfig := ExecutionConfig{
		ConfigVersion:  ConfigVersion,
		CurrentContext: "default",
	}
	err := WriteConfigToFile(confFilePath, &execConfig)
	return err
}

func createDefaultContextIfNotExists() error {
	contextDir, err := GetContextBaseDir()
	if err != nil {
		return err
	}

	var conf ContextConfig
	if err = ReadConfigFromContext("default", &conf); err == nil {
		return nil
	}

	conf = ContextConfig{
		ConfigVersion: ConfigVersion,
		Network:       NetworkConfig{},
	}

	if err = WriteConfigToFile(filepath.Join(contextDir, "default"), &conf); err != nil {
		return err
	}

	return nil
}
