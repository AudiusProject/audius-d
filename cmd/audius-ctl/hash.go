package main

import (
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
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, arg := range args {
					val, err := strconv.Atoi(arg)
					if err != nil {
						logger.ErrorF("invalid number: %s \n", arg)
						continue
					}
					if hashed, err := hashes.Encode(val); err != nil {
						logger.ErrorF("invalid input: %s \n", arg)
					} else {
						logger.Out(hashed)
					}
				}
				return nil
			},
		},
		&cobra.Command{
			Use: "decode",
			RunE: func(cmd *cobra.Command, args []string) error {
				for _, arg := range args {
					if num, err := hashes.MaybeDecode(arg); err != nil {
						logger.ErrorF("invalid input: %s \n", arg)
					} else {
						logger.Out(num)
					}
				}
				return nil
			},
		},
		&cobra.Command{
			Use: "cid",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) != 1 {
					return logger.Error("filename required")
				}
				file, err := os.Open(args[0])
				if err != nil {
					return logger.Error(err)
				}
				defer file.Close()
				cid, err := hashes.ComputeFileCID(file)
				if err != nil {
					return logger.Error(err)
				}
				logger.Info("cid =", cid)
				return nil
			},
		},
	)
}
