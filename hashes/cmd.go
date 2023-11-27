package hashes

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func initCmd() {
	HashCmd = &cobra.Command{
		Use: "hash",
	}
	HashCmd.AddCommand(
		&cobra.Command{
			Use: "encode",
			Run: func(cmd *cobra.Command, args []string) {
				for _, arg := range args {
					val, err := strconv.Atoi(arg)
					if err != nil {
						fmt.Printf("invalid number: %s \n", arg)
						continue
					}
					if hashed, err := Encode(val); err != nil {
						fmt.Printf("invalid number: %s \n", arg)
					} else {
						fmt.Printf("%s => %s \n", arg, hashed)
					}
				}
			},
		},
		&cobra.Command{
			Use: "decode",
			Run: func(cmd *cobra.Command, args []string) {
				for _, arg := range args {
					if num, err := MaybeDecode(arg); err != nil {
						fmt.Printf("invalid number: %s \n", arg)
					} else {
						fmt.Printf("%s => %d \n", arg, num)
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
				cid, err := ComputeFileCID(file)
				if err != nil {
					return err
				}
				fmt.Println("cid =", cid)
				return nil
			},
		},
	)
}
