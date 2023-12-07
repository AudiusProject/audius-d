package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/AudiusProject/audius-d/hashes"
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
						fmt.Println(hashed)
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
						fmt.Println(num)
					}
				}
			},
		},
		&cobra.Command{
			Use: "cid",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) != 1 {
					log.Fatal("filename required")
				}
				file, err := os.Open(args[0])
				if err != nil {
					return err
				}
				defer file.Close()
				cid, err := hashes.ComputeFileCID(file)
				if err != nil {
					return err
				}
				fmt.Println("cid =", cid)
				return nil
			},
		},
	)
}
