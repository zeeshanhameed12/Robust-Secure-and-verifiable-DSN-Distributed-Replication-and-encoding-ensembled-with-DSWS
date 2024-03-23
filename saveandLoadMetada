package main

import (

	//"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
)

// Metadata structure
type Metadata struct {
	FileName     string
	FileSize     int64
	DataShards   *big.Int
	ParityShards *big.Int
	ShardHashes  []string
	ShardOrder   []*big.Int
	IpfsHashes   []string
}

// type Metadata struct {
// 	FileName     string
// 	FileSize     int64
// 	DataShards   *big.Int
// 	ParityShards *big.Int
// 	ShardHashes  []string
// 	ShardOrder   []*big.Int
// 	IpfsHashes   []string
// }

func main() {
	// Example usage
	//contractAddressHex := "0x2bA13060F3962269B0567C352BDF0E9BBB672bC4"

	// metadata := Metadata{
	// 	FileName:     "File15",
	// 	FileSize:     1,
	// 	DataShards:   big.NewInt(1),
	// 	ParityShards: big.NewInt(2),
	// 	ShardHashes:  []string{"Hash1", "Hash2", "Hash3"},
	// 	ShardOrder:   []*big.Int{big.NewInt(3), big.NewInt(5), big.NewInt(6), big.NewInt(7)},
	// 	IpfsHashes:   []string{"ipfsHash1", "ipfsHash2", "ipfsHash3"},
	// }

	// txHash, err := SaveMetadataToContract(contractAddressHex, privateKeyHex, metadata)
	// if err != nil {
	// 	log.Fatalf("Error: %s", err)
	// }

	// fmt.Printf("Metadata saved; transaction hash: %s\n", txHash)

	fileName := "File15"

	metadata, err := LoadMetadataFromContract(contractAddressHex, privateKeyHex, fileName)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	fmt.Println("Retrieved metadata:")
	fmt.Printf("File Name: %s\n", metadata.FileName)
	fmt.Printf("File Size: %d\n", metadata.FileSize)
	fmt.Printf("Data Shards: %d\n", metadata.DataShards)
	fmt.Printf("Parity Shards: %d\n", metadata.ParityShards)
	fmt.Println("Shard Hashes:", metadata.ShardHashes)
	fmt.Println("Shard Order:", metadata.ShardOrder)
	fmt.Println("IPFS Hashes:", metadata.IpfsHashes)

}
