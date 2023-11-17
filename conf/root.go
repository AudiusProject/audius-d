package conf

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "config [command]",
	Short: "view/modify audius-d configuration",
	Run: func(cmd *cobra.Command, args []string) {
		spew.Dump(cmd.Context().Value(ContextKey).(*ContextConfig))
	},
}
var confFileTemplate string
var createContextCmd = &cobra.Command{
	Use:   "create-context <name> [options]",
	Short: "view/modify audius-d configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createContextFromTemplate(args[0], confFileTemplate)
	},
}
var getContextCmd = &cobra.Command{
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
var setContextCmd = &cobra.Command{
	Use:   "set-context <context>",
	Short: "Switch to a different context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := SetContext(args[0])
		if err != nil {
			log.Fatal("Failed to set context", err)
		}
		fmt.Printf("Context set to %s\n", args[0])
	},
}

func init() {
	createContextCmd.Flags().StringVarP(&confFileTemplate, "templatefile", "f", "", "-f <config file to build context from>")
	RootCmd.AddCommand(createContextCmd, getContextCmd, setContextCmd)
}
