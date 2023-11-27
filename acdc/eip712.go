package acdc

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func preapreEip712(chainId int64, _nonce [32]byte, fields EmArgs) ([]byte, error) {
	var typedData = apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "version",
					Type: "string",
				},
				{
					Name: "chainId",
					Type: "uint256",
				},
				{
					Name: "verifyingContract",
					Type: "address",
				},
			},
			"ManageEntity": []apitypes.Type{
				{
					Name: "userId",
					Type: "uint",
				},
				{
					Name: "entityType",
					Type: "string",
				},
				{
					Name: "entityId",
					Type: "uint",
				},
				{
					Name: "action",
					Type: "string",
				},
				{
					Name: "metadata",
					Type: "string",
				},
				{
					Name: "nonce",
					Type: "bytes32",
				},
			},
		},
		Domain: apitypes.TypedDataDomain{
			Name:              "Entity Manager",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(chainId),
			VerifyingContract: EM_CONTRACT,
		},
		PrimaryType: "ManageEntity",
		Message: map[string]interface{}{
			"userId":     big.NewInt(fields.UserID),
			"entityType": fields.EntityType,
			"entityId":   big.NewInt(fields.EntityID),
			"action":     fields.Action,
			"metadata":   fields.Metadata,
			"nonce":      _nonce[:],
		},
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("eip712domain hash struct: %w", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("primary type hash struct: %w", err)
	}

	// add magic string prefix
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sighash := crypto.Keccak256(rawData)
	return sighash, nil
}
