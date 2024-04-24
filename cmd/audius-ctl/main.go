package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	Version        string
	displayVersion bool
	debugLogging   bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "audius-ctl [command]",
		Short: "CLI for provisioning and interacting with audius nodes",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debugLogging {
				logger.SetCliLogLevel(slog.LevelDebug)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if displayVersion {
				fmt.Println(Version)
				return
			}
			cmd.Help()
		},
	}

	rootCmd.Flags().BoolVarP(&displayVersion, "version", "v", false, "Display version info")
	rootCmd.PersistentFlags().BoolVar(&debugLogging, "debug", false, "Print debug logs in console")
	rootCmd.AddCommand(configCmd, devnetCmd, downCmd, infraCmd, jumpCmd, registerCmd, restartCmd, sbCmd, testCmd, upCmd)
	registerCmd.Hidden = true // Hidden as the command is currently only for local devnet registration

	// Handle interrupt/sigterm to mention logfile
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintf(os.Stderr, "Interrupted\n")
		fmt.Fprintf(os.Stderr, "View full debug logs at %s\n", logger.GetLogFilepath())
		os.Exit(1)
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "View full debug logs at %s\n", logger.GetLogFilepath())
		os.Exit(1)
	}
}
