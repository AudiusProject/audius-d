package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config [command]",
		Short: "view/modify audius-d configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			dumpCmd.RunE(cmd, args)
			return nil
		},
	}

	dumpOutfile string
	dumpCmd     = &cobra.Command{
		Use:   "dump [-o outfile]",
		Short: "dump current config to stdout or a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx_config, err := conf.ReadOrCreateContextConfig()
			if err != nil {
				return logger.Error("Failed to retrieve context. ", err)
			}
			if dumpOutfile != "" {
				err := conf.WriteConfigToFile(dumpOutfile, ctx_config)
				if err != nil {
					return logger.Error("Failed to write config to file:", err)
				}
			} else {
				str, err := conf.StringifyConfig(ctx_config)
				if err != nil {
					return logger.Error("Failed to dump config:", err)
				}
				logger.Out(str)
			}
			return nil
		},
	}

	editCmd = &cobra.Command{
		Use:               "edit [context]",
		Short:             "edit the current or specified configuration in an external editor",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: contextCompletionFunction,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxName, err := conf.GetCurrentContextName()
			if err != nil {
				return logger.Error(err)
			}
			if len(args) > 0 {
				ctxName = args[0]
			}
			if err := EditConfig(ctxName); err != nil {
				return logger.Error(err)
			}
			return nil
		},
	}

	confFileTemplate string
	createContextCmd = &cobra.Command{
		Use:   "create-context <name>",
		Short: "create an audius-d configuration context, optionally from a template",
		Long: `
		Create an audius-d configuration context.
		Without any flags, creates a bare-bones context with the given name.
		Use '-f [filename]' to specify a template to copy.
		Use '-f -' to read from stdin. Useful for scripts or pipes.
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if confFileTemplate == "-" {
				input, err := io.ReadAll(os.Stdin)
				if err != nil {
					return logger.Error("Error reading from stdin:", err)
				}
				ctx := conf.NewContextConfig()
				if err := conf.ReadConfigFromBytes(input, ctx); err != nil {
					return logger.Error("Could not parse config:", err)
				}
				if err := conf.WriteConfigToContext(args[0], ctx); err != nil {
					return logger.Error("Failed to save config:", err)
				}
			} else {
				err := conf.CreateContextFromTemplate(args[0], confFileTemplate)
				if err != nil {
					return logger.Error("Failed to create context:", err)
				}
			}
			useContextCmd.RunE(cmd, args)
			return nil
		},
	}
	migrateContextCmd = &cobra.Command{
		Use:   "migrate-context <name> <path>",
		Short: "create an audius-d configuration based of an existing audius-docker-compose instance",
		Long: `
		Create an audius-d configuration based of an existing audius-docker-compose instance.
		
		Requires two arguments, the name of the context where the instance will land. 
		A path to an existing audius-docker-compose installation.

		Examples:
		"audius-ctl config migrate-context my-creator-node ~/audius-docker-compose"
		"audius-ctl config migrate-context my-discprov-7 ../audius-docker-compose"
		`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := conf.MigrateAudiusDockerCompose(args[0], args[1]); err != nil {
				return logger.Error("audius-docker-compose migration failed:", err)
			}
			logger.Info("audius-docker-compose migration successful 🎉")
			useContextCmd.RunE(cmd, args)
			return nil
		},
	}
	currentContextCmd = &cobra.Command{
		Use:   "current-context",
		Short: "Show the currently enabled context",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := conf.GetCurrentContextName()
			if err != nil {
				return logger.Error("Failed to retrieve current context:", err)
			}
			logger.Out(ctx)
			return nil
		},
	}
	getContextsCmd = &cobra.Command{
		Use:   "get-contexts",
		Short: "Show all available contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, err := conf.GetContexts()
			if err != nil {
				return logger.Error("Failed to retrieve current context:", err)
			}
			for _, ctx := range ctxs {
				logger.Out(ctx)
			}
			return nil
		},
	}
	useContextCmd = &cobra.Command{
		Use:               "use-context <context>",
		Short:             "Switch to a different context",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: contextCompletionFunction,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := conf.UseContext(args[0])
			if err != nil {
				return logger.Error("Failed to set context:", err)

			}
			logger.Out(args[0])
			logger.Infof("Context set to %s", args[0])
			return nil
		},
	}
	deleteContextCmd = &cobra.Command{
		Use:               "delete-context <context>",
		Short:             "Delete a context",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: contextCompletionFunction,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := conf.DeleteContext(args[0]); err != nil {
				return logger.Error("Failed to delete context:", err)
			}
			logger.Out(args[0])
			logger.Infof("Context %s deleted.", args[0])
			return nil
		},
	}
)

func init() {
	createContextCmd.Flags().StringVarP(&confFileTemplate, "templatefile", "f", "", "'-f <filename>' to copy context from a template file or use '-f -' to read from stdin")
	dumpCmd.Flags().StringVarP(&dumpOutfile, "outfile", "o", "", "-o <outfile")
	configCmd.AddCommand(dumpCmd, createContextCmd, currentContextCmd, getContextsCmd, useContextCmd, deleteContextCmd, editCmd, migrateContextCmd)
}

func contextCompletionFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return getAvailableContextsWithPrefix(toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getAvailableContextsWithPrefix(prefix string) []string {
	ctxs, err := conf.GetContexts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return nil
	}
	matches := make([]string, 0)
	for _, ctx := range ctxs {
		if strings.HasPrefix(ctx, prefix) {
			matches = append(matches, ctx)
		}
	}
	return matches
}

func EditConfig(contextName string) error {
	tempFile, err := os.CreateTemp("", contextName)
	if err != nil {
		return err
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	existingConfig, err := conf.GetContextConfig(contextName)
	if err != nil {
		return err
	}

	if err = conf.WriteConfigToFile(tempFile.Name(), existingConfig); err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		logger.Info("Please set $EDITOR in your shell profile to your preferred text editor.")
		logger.Info("Defaulting to nano")
		editor = "nano"
	}

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return err
	}

	var newConfig conf.ContextConfig
	if err = conf.ReadConfigFromFile(tempFile.Name(), &newConfig); err != nil {
		return err
	}

	if err = conf.WriteConfigToContext(contextName, &newConfig); err != nil {
		return err
	}

	return nil
}
