package utils

import "github.com/AudiusProject/audius-protocol/mediorum/ethcontracts"

func GenerateAddress(pkey string) (*string, error) {
	pk, err := ethcontracts.ParsePrivateKeyHex(pkey)
	if err != nil {
		return nil, err
	}

	addr := ethcontracts.ComputeAddressFromPrivateKey(pk)
	return &addr, nil
}
