package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AudiusProject/audius-d/pkg/infra"
	"github.com/spf13/cobra"
)

var (
	infraCmd = &cobra.Command{
		Use:   "infra [command]",
		Short: "Manage audius-d instances",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	infraCancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel the current in progress update",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return infra.Cancel(ctx)
		},
	}

	infraDestroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Are you sure you want to destroy? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			response = strings.TrimSpace(response)
			if strings.ToLower(response) != "y" {
				fmt.Println("Destroy canceled.")
				return nil
			}
			ctx := context.Background()
			return infra.Destroy(ctx)
		},
	}

	infraGetOutputCmd = &cobra.Command{
		Use:   "output <outputName>",
		Short: "Get a specific output from the stack",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			outputName := args[0]
			outputValue, err := infra.GetStackOutput(ctx, outputName)
			if err != nil {
				return err
			}
			fmt.Println(outputValue)
			return nil
		},
	}

	infraPreviewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Perform a dry-run update for the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return infra.Update(ctx, true)
		},
	}

	infraSshCmd = &cobra.Command{
		Use:   "ssh -- <command>",
		Short: "SSH into the EC2 instance and execute commands",
		Long:  `Use this command to SSH into the EC2 instance and execute commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			publicIP, err := infra.GetStackOutput(ctx, "instancePublicIp")
			if err != nil {
				return err
			}
			privateKeyFilePath, err := infra.GetStackOutput(ctx, "instancePrivateKeyFilePath")
			if err != nil {
				return err
			}
			command := "echo 'Please specify a command to run on the remote server.'"
			if len(args) > 0 {
				command = strings.Join(args, " ")
			}

			output, err := infra.ExecuteSSHCommand(privateKeyFilePath, publicIP, command)
			if err != nil {
				fmt.Printf("Error executing SSH command: %v\n", err)
				return err
			}

			fmt.Println(output)
			return nil
		},
	}

	infraUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update (deploy) the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			return infra.Update(ctx, false)
		},
	}
)

func init() {
	infraCmd.AddCommand(infraCancelCmd, infraDestroyCmd, infraGetOutputCmd, infraPreviewCmd, infraSshCmd, infraUpdateCmd)
}
