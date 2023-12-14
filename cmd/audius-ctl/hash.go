package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AudiusProject/audius-d/pkg/hashes"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/spf13/cobra"
)

var hashCmd *cobra.Command

func init() {
	hashCmd = &cobra.Command{
		Use: "hash",
	}
	hashCmd.AddCommand(
		&cobra.Command{
			Use: "encode",
			Run: func(cmd *cobra.Command, args []string) {
				for _, arg := range args {
					val, err := strconv.Atoi(arg)
					if err != nil {
						fmt.Fprintf(os.Stderr, "invalid number: %s \n", arg)
						continue
					}
					if hashed, err := hashes.Encode(val); err != nil {
						fmt.Fprintf(os.Stderr, "invalid input: %s \n", arg)
					} else {
						logger.Info(hashed)
					}
				}
			},
		},
		&cobra.Command{
			Use: "decode",
			Run: func(cmd *cobra.Command, args []string) {
				for _, arg := range args {
					if num, err := hashes.MaybeDecode(arg); err != nil {
						fmt.Fprintf(os.Stderr, "invalid input: %s \n", arg)
					} else {
						logger.Info("", num)
					}
				}
			},
		},
		&cobra.Command{
			Use: "cid",
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) != 1 {
					logger.Error("filename required")
					return
				}
				file, err := os.Open(args[0])
				if err != nil {
					logger.Error(err)
					return
				}
				defer file.Close()
				cid, err := hashes.ComputeFileCID(file)
				if err != nil {
					logger.Error(err)
					return
				}
				logger.Info("cid =", cid)
			},
		},
	)
}
