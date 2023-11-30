package conf

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	RootCmd = &cobra.Command{
		Use:   "config [command]",
		Short: "view/modify audius-d configuration",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			dumpCmd.Run(cmd, args)
		},
	}

	dumpOutfile string
	dumpCmd     = &cobra.Command{
		Use:   "dump [-o outfile]",
		Short: "dump current config to stdout or a file",
		Run: func(cmd *cobra.Command, args []string) {
			ctx_config, err := ReadOrCreateContextConfig()
			if err != nil {
				log.Fatal("Failed to retrieve context. ", err)
			}
			if dumpOutfile != "" {
				err := writeConfigToFile(dumpOutfile, ctx_config)
				if err != nil {
					log.Fatal("Failed to write config to file:", err)
				}
			} else {
				str, err := stringifyConfig(ctx_config)
				if err != nil {
					log.Fatal("Failed to dump config:", err)
				}
				fmt.Println(str)
			}
		},
	}

	setCmd = &cobra.Command{
		Use:   "set <property.name> <value>",
		Short: "modify a configuration value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := setConfigWithViper(args[0], args[1]); err != nil {
				log.Fatal("Failed to set config value: ", err)
			}
		},
	}
	editCmd = &cobra.Command{
		Use:   "edit [context]",
		Short: "edit the current or specified configuration in an external editor",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctxName, err := GetCurrentContextName()
			if err != nil {
				log.Fatal(err)
			}
			if len(args) > 0 {
				ctxName = args[0]
			}
			if err := editConfig(ctxName); err != nil {
				log.Fatal(err)
			}
		},
	}

	confFileTemplate string
	createContextCmd = &cobra.Command{
		Use:   "create-context <name> [options]",
		Short: "create an audius-d configuration context, optionally from a template",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := createContextFromTemplate(args[0], confFileTemplate)
			if err != nil {
				log.Fatal("Failed to create context:", err)
			}
			useContextCmd.Run(cmd, args)
		},
	}
	migrateContextCmd = &cobra.Command{
		Use:   "migrate-context <name> <path>",
		Short: "create an audius-d configuration based of an existing audius-docker-compose instance",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := MigrateAudiusDockerCompose(args[0], args[1]); err != nil {
				log.Fatal("audius-docker-compose migration failed: ", err)
			}
			log.Println("audius-docker-compose migration successful 🎉")
		},
	}
	currentContextCmd = &cobra.Command{
		Use:   "current-context",
		Short: "Show the currently enabled context",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := GetCurrentContextName()
			if err != nil {
				log.Fatal("Failed to retrieve current context: ", err)
			}
			fmt.Println(ctx)
		},
	}
	getContextsCmd = &cobra.Command{
		Use:   "get-contexts",
		Short: "Show all available contexts",
		Run: func(cmd *cobra.Command, args []string) {
			ctxs, err := GetContexts()
			if err != nil {
				log.Fatal("Failed to retrieve current context: ", err)
			}
			for _, ctx := range ctxs {
				fmt.Println(ctx)
			}
		},
	}
	useContextCmd = &cobra.Command{
		Use:   "use-context <context>",
		Short: "Switch to a different context",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := UseContext(args[0])
			if err != nil {
				log.Fatal("Failed to set context: ", err)
			}
			fmt.Printf("Context set to %s\n", args[0])
		},
	}
	deleteContextCmd = &cobra.Command{
		Use:   "delete-context <context>",
		Short: "Delete a context",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := DeleteContext(args[0]); err != nil {
				log.Fatal("Failed to delete context: ", err)
			}
			fmt.Printf("Context %s deleted.\n", args[0])
		},
	}
)

func init() {
	createContextCmd.Flags().StringVarP(&confFileTemplate, "templatefile", "f", "", "-f <config file to build context from>")
	dumpCmd.Flags().StringVarP(&dumpOutfile, "outfile", "o", "", "-o <outfile")
	RootCmd.AddCommand(dumpCmd, createContextCmd, currentContextCmd, getContextsCmd, useContextCmd, deleteContextCmd, setCmd, editCmd, migrateContextCmd)
}

func setConfigWithViper(key string, value string) error {
	v := viper.New()
	cname, err := GetCurrentContextName()
	if err != nil {
		return err
	}
	basedir, err := getContextBaseDir()
	if err != nil {
		return err
	}
	v.SetConfigFile(filepath.Join(basedir, cname))
	v.SetConfigType("toml")
	if err = v.ReadInConfig(); err != nil {
		return err
	}
	if !v.IsSet(key) {
		return fmt.Errorf("key '%s' not found in config", key)
	}
	v.Set(key, value)
	var config ContextConfig
	if err = v.Unmarshal(&config); err != nil {
		return err
	}
	if err = writeConfigToCurrentContext(&config); err != nil {
		return err
	}
	return nil
}

func editConfig(contextName string) error {
	tempFile, err := os.CreateTemp("", contextName)
	if err != nil {
		return err
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	var existingConfig ContextConfig
	if err = readConfigFromContext(contextName, &existingConfig); err != nil {
		return err
	}

	if err = writeConfigToFile(tempFile.Name(), &existingConfig); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		fmt.Println("Please set $EDITOR in your shell profile to your preferred text editor.")
		fmt.Println("Defaulting to nano")
		editor = "nano"
	}

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return err
	}

	var newConfig ContextConfig
	if err = readConfigFromFile(tempFile.Name(), &newConfig); err != nil {
		return err
	}

	if err = writeConfigToContext(contextName, &newConfig); err != nil {
		return err
	}

	return nil
}
