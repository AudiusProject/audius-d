package hashid

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/speps/go-hashids/v2"
	"github.com/spf13/cobra"
)

var (
	hasher  *hashids.HashID
	HashCmd *cobra.Command
)

func init() {
	initHasher()
	initCmd()
}

func initHasher() {
	hd := hashids.NewData()
	hd.Salt = "azowernasdfoia"
	hd.MinLength = 5
	hasher, _ = hashids.NewWithData(hd)
}

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
	)
}

func Encode(id int) (string, error) {
	return hasher.Encode([]int{id})
}

func Decode(hid string) (int, error) {
	ids, err := hasher.DecodeWithError(hid)
	if err != nil {
		return 0, err
	}
	if len(ids) < 1 {
		return 0, errors.New("invalid hash")
	}
	return ids[0], nil
}

func MaybeDecode(hid string) (int, error) {
	id, err := Decode(hid)
	if err != nil {
		id, err = strconv.Atoi(hid)
	}
	return id, err
}
