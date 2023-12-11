package acdc

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EmArgs struct {
	UserID     int64
	EntityType string
	EntityID   int64
	Action     string
	Metadata   string
}

func SendEmTx(client *ethclient.Client, privateKey *ecdsa.PrivateKey, fields EmArgs) (*types.Transaction, error) {
	ctx := context.Background()

	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	_nonce := randBytes32()

	// sign it with EIP 712
	sighash, err := prepareEip712(chainId.Int64(), _nonce, fields)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(sighash, privateKey)
	if err != nil {
		return nil, err
	}
	sig[crypto.RecoveryIDOffset] += 27

	// pack manageEntity ABI call
	abiArgs := []interface{}{big.NewInt(fields.UserID), fields.EntityType, big.NewInt(fields.EntityID), fields.Action, fields.Metadata, _nonce, sig}
	packed, err := AudiusABI.Pack("manageEntity", abiArgs...)
	if err != nil {
		return nil, err
	}

	return sendIt(client, privateKey, packed)
}

func randBytes32() [32]byte {
	var _nonce [32]byte
	rand.Read(_nonce[:])
	return _nonce
}

func sendIt(client *ethclient.Client, privateKey *ecdsa.PrivateKey, packed []byte) (*types.Transaction, error) {
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return nil, err
	}

	value := big.NewInt(0)
	gasLimit := uint64(2_100_000)
	gasPrice, _ := client.SuggestGasPrice(context.Background())

	tx := types.NewTransaction(nonce, common.HexToAddress(EM_CONTRACT), value, gasLimit, gasPrice, packed)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
