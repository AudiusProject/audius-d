package acdc

import (
	"math/big"
	"strings"

	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const EM_CONTRACT = "0x1cd8a543596d499b9b6e7a6ec15ecd2b7857fd64"

var AudiusABI abi.ABI

const abiJsonString = `
[
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "_userId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "_signer",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "_entityType",
				"type": "string"
			},
			{
				"indexed": false,
				"name": "_entityId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "_metadata",
				"type": "string"
			},
			{
				"indexed": false,
				"name": "_action",
				"type": "string"
			}
		],
		"name": "ManageEntity",
		"type": "event"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_userId",
				"type": "uint256"
			},
			{
				"name": "_entityType",
				"type": "string"
			},
			{
				"name": "_entityId",
				"type": "uint256"
			},
			{
				"name": "_action",
				"type": "string"
			},
			{
				"name": "_metadata",
				"type": "string"
			},
			{
				"name": "_nonce",
				"type": "bytes32"
			},
			{
				"name": "_subjectSig",
				"type": "bytes"
			}
		],
		"name": "manageEntity",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},

	{
		"constant": false,
		"inputs": [],
		"name": "kill",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	      },
	      {
		"constant": false,
		"inputs": [],
		"name": "renounceOwnership",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	      },
	      {
		"constant": true,
		"inputs": [],
		"name": "owner",
		"outputs": [
		  {
		    "name": "",
		    "type": "address"
		  }
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	      },
	      {
		"constant": true,
		"inputs": [],
		"name": "isOwner",
		"outputs": [
		  {
		    "name": "",
		    "type": "bool"
		  }
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	      },
	      {
		"constant": false,
		"inputs": [
		  {
		    "name": "_registryAddress",
		    "type": "address"
		  }
		],
		"name": "setRegistry",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	      },
	      {
		"constant": false,
		"inputs": [
		  {
		    "name": "newOwner",
		    "type": "address"
		  }
		],
		"name": "transferOwnership",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	      },
	      {
		"constant": true,
		"inputs": [
		  {
		    "name": "",
		    "type": "bytes32"
		  }
		],
		"name": "usedSignatures",
		"outputs": [
		  {
		    "name": "",
		    "type": "bool"
		  }
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	      },
	      {
		"inputs": [
		  {
		    "name": "_registryAddress",
		    "type": "address"
		  },
		  {
		    "name": "_trackStorageRegistryKey",
		    "type": "bytes32"
		  },
		  {
		    "name": "_userFactoryRegistryKey",
		    "type": "bytes32"
		  },
		  {
		    "name": "_networkId",
		    "type": "uint256"
		  }
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	      },
	      {
		"anonymous": false,
		"inputs": [
		  {
		    "indexed": false,
		    "name": "_id",
		    "type": "uint256"
		  },
		  {
		    "indexed": false,
		    "name": "_trackOwnerId",
		    "type": "uint256"
		  },
		  {
		    "indexed": false,
		    "name": "_multihashDigest",
		    "type": "bytes32"
		  },
		  {
		    "indexed": false,
		    "name": "_multihashHashFn",
		    "type": "uint8"
		  },
		  {
		    "indexed": false,
		    "name": "_multihashSize",
		    "type": "uint8"
		  }
		],
		"name": "NewTrack",
		"type": "event"
	      },
	      {
		"anonymous": false,
		"inputs": [
		  {
		    "indexed": false,
		    "name": "_trackId",
		    "type": "uint256"
		  },
		  {
		    "indexed": false,
		    "name": "_trackOwnerId",
		    "type": "uint256"
		  },
		  {
		    "indexed": false,
		    "name": "_multihashDigest",
		    "type": "bytes32"
		  },
		  {
		    "indexed": false,
		    "name": "_multihashHashFn",
		    "type": "uint8"
		  },
		  {
		    "indexed": false,
		    "name": "_multihashSize",
		    "type": "uint8"
		  }
		],
		"name": "UpdateTrack",
		"type": "event"
	      },
	      {
		"anonymous": false,
		"inputs": [
		  {
		    "indexed": false,
		    "name": "_trackId",
		    "type": "uint256"
		  }
		],
		"name": "TrackDeleted",
		"type": "event"
	      },
	      {
		"anonymous": false,
		"inputs": [
		  {
		    "indexed": true,
		    "name": "previousOwner",
		    "type": "address"
		  },
		  {
		    "indexed": true,
		    "name": "newOwner",
		    "type": "address"
		  }
		],
		"name": "OwnershipTransferred",
		"type": "event"
	      },
	      {
		"constant": true,
		"inputs": [
		  {
		    "name": "_id",
		    "type": "uint256"
		  }
		],
		"name": "trackExists",
		"outputs": [
		  {
		    "name": "exists",
		    "type": "bool"
		  }
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	      },
	      {
		"constant": false,
		"inputs": [
		  {
		    "name": "_trackOwnerId",
		    "type": "uint256"
		  },
		  {
		    "name": "_multihashDigest",
		    "type": "bytes32"
		  },
		  {
		    "name": "_multihashHashFn",
		    "type": "uint8"
		  },
		  {
		    "name": "_multihashSize",
		    "type": "uint8"
		  },
		  {
		    "name": "_nonce",
		    "type": "bytes32"
		  },
		  {
		    "name": "_subjectSig",
		    "type": "bytes"
		  }
		],
		"name": "addTrack",
		"outputs": [
		  {
		    "name": "",
		    "type": "uint256"
		  }
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	      },
	      {
		"constant": false,
		"inputs": [
		  {
		    "name": "_trackId",
		    "type": "uint256"
		  },
		  {
		    "name": "_trackOwnerId",
		    "type": "uint256"
		  },
		  {
		    "name": "_multihashDigest",
		    "type": "bytes32"
		  },
		  {
		    "name": "_multihashHashFn",
		    "type": "uint8"
		  },
		  {
		    "name": "_multihashSize",
		    "type": "uint8"
		  },
		  {
		    "name": "_nonce",
		    "type": "bytes32"
		  },
		  {
		    "name": "_subjectSig",
		    "type": "bytes"
		  }
		],
		"name": "updateTrack",
		"outputs": [
		  {
		    "name": "",
		    "type": "bool"
		  }
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	      },
	      {
		"constant": false,
		"inputs": [
		  {
		    "name": "_trackId",
		    "type": "uint256"
		  },
		  {
		    "name": "_nonce",
		    "type": "bytes32"
		  },
		  {
		    "name": "_subjectSig",
		    "type": "bytes"
		  }
		],
		"name": "deleteTrack",
		"outputs": [
		  {
		    "name": "status",
		    "type": "bool"
		  }
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	      },
	      {
		"constant": true,
		"inputs": [
		  {
		    "name": "_caller",
		    "type": "address"
		  },
		  {
		    "name": "_trackId",
		    "type": "uint256"
		  }
		],
		"name": "callerOwnsTrack",
		"outputs": [],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	      },
	      {
		"constant": true,
		"inputs": [
		  {
		    "name": "_id",
		    "type": "uint256"
		  }
		],
		"name": "getTrack",
		"outputs": [
		  {
		    "name": "trackOwnerId",
		    "type": "uint256"
		  },
		  {
		    "name": "multihashDigest",
		    "type": "bytes32"
		  },
		  {
		    "name": "multihashHashFn",
		    "type": "uint8"
		  },
		  {
		    "name": "multihashSize",
		    "type": "uint8"
		  }
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	      }
]
`

// EntityManagerManageEntity represents a ManageEntity event raised by the EntityManager contract.
type EntityManagerManageEntity struct {
	UserId     *big.Int
	Signer     common.Address
	EntityType string
	EntityId   *big.Int
	Metadata   string
	Action     string
	Raw        types.Log // Blockchain specific contextual infos
}

func init() {
	var err error
	AudiusABI, err = abi.JSON(strings.NewReader(abiJsonString))
	if err != nil {
		logger.Error(err)
		return
	}
}
