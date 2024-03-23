package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	// Replace the below path with the actual path to your generated Go bindings
	//"your/package/path/to/generated/bindings"
)

func main() {
	ganacheURL := "http://localhost:8545"// Default Ganache CLI URL
	privateKeyHex := "5345d475be2a8561ad0b61fb6e8b8e2b13e6cde1f602a5a5bab28a1d29151de9"

	client, err := ethclient.Dial(ganacheURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337)) // Ganache default chain ID is 1337
	if err != nil {
		log.Fatalf("Failed to create an authorized transactor: %v", err)
	}

	// The address that will own the deployed contract.
	auth.From = crypto.PubkeyToAddress(privateKey.PublicKey)

	// Deploy the contract
	address, tx, _, err := DeployMain(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy the contract: %v", err)
	}

	fmt.Printf("Contract deployed! Address: %s, Tx: %s\n", address.Hex(), tx.Hash().Hex())
}
