package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// const (
// 	infuraURL  = "https://sepolia.infura.io/v3/a2b0e9a36e9a481b8a5356585850b404"
// 	//ganacheurl = "http://localhost:8545"
// 	// walletPath = "./wallet/UTC--2024-03-06T18-26-46.040477236Z--4a02341be8e888617faedd979d25adc80061cece"
// 	// walletPassword = "account1"
// 	contractAddressHex = "0x7e4D9F2AD3AaBC0dC7DEAFAE3564257Cf2daa80e"
// 	privateKeyHex      = "f58aefca586e6b0791b83109e00d885f5abc6750305dbe635aceeb487b6ad47b"
// )

//var client, err = ethclient.Dial(ganacheurl)

// SaveMetadataToContract stores metadata in the blockchain and returns the transaction hash.
func SaveMetadataToContract( metadata *Metadata) (string, error) {
	
	ganacheurl := "http://localhost:8545"
	// walletPath = "./wallet/UTC--2024-03-06T18-26-46.040477236Z--4a02341be8e888617faedd979d25adc80061cece"
	// walletPassword = "account1"
	contractAddressHex := "0x8842b262eC9B7A6961F62470fd01CAC3E172cF93"
	privateKeyHex      := "5345d475be2a8561ad0b61fb6e8b8e2b13e6cde1f602a5a5bab28a1d29151de9"
	client, err := ethclient.Dial(ganacheurl)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337)) // Adjust the chain ID accordingly
	if err != nil {
		return "", fmt.Errorf("failed to create authorized transactor: %v", err)
	}

	// Set up nonce for the transaction
	nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	contractAddress := common.HexToAddress(contractAddressHex)
	metadataStorageInstance, err := NewMain(contractAddress, client)
	if err != nil {
		return "", fmt.Errorf("failed to instantiate a MetadataStorage contract: %v", err)
	}

	// Save metadata using the contract

	tx, err := metadataStorageInstance.SaveMetadata(auth, metadata.FileName, metadata.FileSize, big.NewInt(int64( metadata.DataShards)),big.NewInt(int64( metadata.ParityShards)), metadata.ShardHashes,convertIntSliceToBigIntSlice(metadata.ShardOrder), metadata.IpfsHashes)
	if err != nil {
		return "", fmt.Errorf("failed to save metadata: %v", err)
	}

	return tx.Hash().Hex(), nil
}

// LoadMetadataFromContract retrieves metadata from the blockchain.
func LoadMetadataFromContract( fileName string) (*Metadata, error) {
	// Connect to the Ethereum client
	ganacheURL := "http://localhost:8545"
	// walletPath = "./wallet/UTC--2024-03-06T18-26-46.040477236Z--4a02341be8e888617faedd979d25adc80061cece"
	// walletPassword = "account1"
	contractAddressHex := "0x8842b262eC9B7A6961F62470fd01CAC3E172cF93"
	privateKeyHex      := "5345d475be2a8561ad0b61fb6e8b8e2b13e6cde1f602a5a5bab28a1d29151de9"
	client, err := ethclient.Dial(ganacheURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Convert the private key from hex to an ECDSA private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	// Get the public address of the private key
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Instantiate the contract
	contractAddress := common.HexToAddress(contractAddressHex)
	contract, err := NewMain(contractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the contract: %v", err)
	}

	// Prepare the call options
	opts := &bind.CallOpts{
		Context: context.Background(),
		From:    publicAddress,
	}

	// Call the loadMetadata function of the contract
	metadataResult, err := contract.LoadMetadata(opts, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve metadata: %v", err)
	}

	//Convert the result to our Metadata struct
	metadata := &Metadata{
		FileName:     metadataResult.FileName,
		FileSize:     metadataResult.FileSize,
		DataShards:   metadataResult.DataShards.Int64(),
		ParityShards: metadataResult.ParityShards.Int64(),
		ShardHashes:  metadataResult.ShardHashes,
		IpfsHashes:   metadataResult.IpfsHashes,
	}

	//Convert shardOrder from []*big.Int to []int64
	for _, order := range metadataResult.ShardOrder {
		metadata.ShardOrder = append(metadata.ShardOrder, int(order.Int64()))
	}

	return metadata, err
}

func convertIntSliceToBigIntSlice(intSlice []int) []*big.Int {
    var bigIntSlice []*big.Int
    for _, num := range intSlice {
        // Convert int to *big.Int and append to the slice
        bigIntSlice = append(bigIntSlice, big.NewInt(int64(num)))
    }
    return bigIntSlice
}
