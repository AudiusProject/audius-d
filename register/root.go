package register

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

//go:embed ABIs/ERC20Detailed.json
var erc20ABIFile string

//go:embed ABIs/Registry.json
var registryABIFile string

//go:embed ABIs/ServiceProviderFactory.json
var spfABIFile string

var RootCmd *cobra.Command

func init() {
	RootCmd = &cobra.Command{
		Use:   "register",
		Short: "Register nodes on ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("got here")
			//conf := cmd.ctx.Value(config.ContextKey)
			//RegisterNode()
		},
	}
}

func RegisterNode(registrationNodeType string, nodeEndpoint string, ethProviderUrl string, tokenAddress string, contractRegistryAddress string, ownerWallet string, privateKey string) {
	client, err := ethclient.Dial(ethProviderUrl)
	if err != nil {
		log.Fatal("Failed to dial ethereum client:", err)
	}
	delegateOwnerWallet := common.HexToAddress(ownerWallet)
	pKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal("Failed to encode private key:", err)
	}
	ethRegistryAddress := common.HexToAddress(contractRegistryAddress)
	tokenABI := getContractABI(erc20ABIFile)
	serviceProviderFactoryABI := getContractABI(spfABIFile)

	var tokenDecimals uint8
	tokenDecimalsData, err := tokenABI.Pack("decimals")
	if err != nil {
		log.Fatal("Failed to pack tokenABI for token decimals:", err)
	}
	ethTokenAddress := common.HexToAddress(tokenAddress)
	tokenDecimalsResult, err := client.CallContract(
		context.Background(),
		ethereum.CallMsg{
			To:   &ethTokenAddress,
			Data: tokenDecimalsData,
		},
		nil,
	)
	if err != nil {
		log.Fatal("Failed to retrieve token decimals:", err)
	}
	if err = tokenABI.UnpackIntoInterface(&tokenDecimals, "decimals", tokenDecimalsResult); err != nil {
		log.Fatal("Failed to unpack token decimals result:", err)
	}

	coeff := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(tokenDecimals)), nil)
	stakedTokensAmount := new(big.Int).Mul(big.NewInt(200000), coeff)

	tokenApprovalData, err := tokenABI.Pack(
		"approve",
		getContractAddress(client, ethRegistryAddress, "StakingProxy"),
		stakedTokensAmount,
	)
	err = client.SendTransaction(
		context.Background(),
		getSignedTx(client, tokenApprovalData, delegateOwnerWallet, ethTokenAddress, pKey),
	)
	if err != nil {
		log.Fatal("Failed to approve tokens:", err)
	}

	var bytes32NodeType [32]byte
	copy(bytes32NodeType[:], []byte(registrationNodeType))

	spfAddress := getContractAddress(client, ethRegistryAddress, "ServiceProviderFactory")
	spfRegisterData, err := serviceProviderFactoryABI.Pack(
		"register",
		bytes32NodeType,
		nodeEndpoint,
		stakedTokensAmount,
		delegateOwnerWallet,
	)
	if err != nil {
		log.Fatal("Failed to pack serviceProviderFactoryABI:", err)
	}
	err = client.SendTransaction(
		context.Background(),
		getSignedTx(client, spfRegisterData, delegateOwnerWallet, spfAddress, pKey),
	)
	if err != nil {
		log.Fatal("Failed to register node transaction:", err)
	}
}

func getContractABI(abiFile string) abi.ABI {
	resultABI, err := abi.JSON(strings.NewReader(abiFile))
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create contract ABI from file '%s':", abiFile), err)
	}

	return resultABI
}

func getContractAddress(client *ethclient.Client, ethRegistryAddress common.Address, contractName string) common.Address {
	registryABI := getContractABI(registryABIFile)

	var bytes32Key [32]byte
	copy(bytes32Key[:], []byte(contractName))

	// The actual method is getContract(bytes32), but it's overloaded and go-ethereum is dumb.
	data, err := registryABI.Pack("getContract0", bytes32Key)
	if err != nil {
		log.Fatal("Failed to pack registryABI:", err)
	}

	msg := ethereum.CallMsg{
		To:   &ethRegistryAddress,
		Data: data,
	}

	resultBytes, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatal("Failed to retrieve contract address:", err)
	}

	var contractAddr common.Address
	if err = registryABI.UnpackIntoInterface(&contractAddr, "getContract0", resultBytes); err != nil {
		log.Fatal("Failed to unpack result:", err)
	}

	return contractAddr
}

func getSignedTx(client *ethclient.Client, txData []byte, from common.Address, to common.Address, privateKey *ecdsa.PrivateKey) *types.Transaction {
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("Failed to get chain id:", err)
	}
	gasLimit, err := client.EstimateGas(
		context.Background(),
		ethereum.CallMsg{
			From: from,
			To:   &to,
			Data: txData,
		},
	)
	if err != nil {
		log.Fatal("Failed to estimate gas limit:", err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Failed to suggest gas price:", err)
	}
	tx := types.NewTransaction(nonce, to, big.NewInt(0), gasLimit, gasPrice, txData)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("Failed to sign tx:", err)
	}
	return signedTx
}
