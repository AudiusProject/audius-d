package conf

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "config [command]",
		Short: "view/modify audius-d configuration",
		Run: func(cmd *cobra.Command, args []string) {
			dumpCmd.Run(cmd, args)
		},
	}

	dumpOutfile string
	dumpCmd     = &cobra.Command{
		Use:   "dump [-o outfile]",
		Short: "dump current config to stdout or a file",
		Run: func(cmd *cobra.Command, args []string) {
			if dumpOutfile != "" {
				err := writeConfigToFile(dumpOutfile, cmd.Context().Value(ContextKey).(*ContextConfig))
				if err != nil {
					log.Fatal("Failed to write config to file:", err)
				}
			} else {
				str, err := stringifyConfig(cmd.Context().Value(ContextKey).(*ContextConfig))
				if err != nil {
					log.Fatal("Failed to dump config:", err)
				}
				fmt.Println(str)
			}
		},
	}

	confFileTemplate string
	createContextCmd = &cobra.Command{
		Use:   "create-context <name> [options]",
		Short: "view/modify audius-d configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := createContextFromTemplate(args[0], confFileTemplate)
			if err != nil {
				log.Fatal("Failed to create context:", err)
			}
			useContextCmd.Run(cmd, args)
		},
	}
	getContextCmd = &cobra.Command{
		Use:   "get-context",
		Short: "Show the currently enabled context",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, err := GetContext()
			if err != nil {
				log.Fatal("Failed to retrieve current context", err)
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
				log.Fatal("Failed to retrieve current context", err)
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
				log.Fatal("Failed to set context", err)
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
				log.Fatal("Failed to delete context", err)
			}
			fmt.Printf("Context %s deleted.\n", args[0])
		},
	}
)

func init() {
	createContextCmd.Flags().StringVarP(&confFileTemplate, "templatefile", "f", "", "-f <config file to build context from>")
	dumpCmd.Flags().StringVarP(&dumpOutfile, "outfile", "o", "", "-o <outfile")
	RootCmd.AddCommand(dumpCmd, createContextCmd, getContextCmd, getContextsCmd, useContextCmd, deleteContextCmd)
}
