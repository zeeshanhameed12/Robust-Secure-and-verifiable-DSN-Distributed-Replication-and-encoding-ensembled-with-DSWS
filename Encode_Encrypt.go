package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/klauspost/reedsolomon"
)



// Helper function to encrypt data using AES-GCM
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Function to encrypt each shard

func encoding_EncryptingData(data []byte, dataShards, parityShards int, privateKey []byte) ([][]byte, *Metadata, error) {

	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return nil, nil, err
	}

	shards, err := enc.Split(data)
	if err != nil {
		return nil, nil, err
	}

	shardSize := len(data) / dataShards

	// err = enc.Encode(shards)
	// if err != nil {
	//     return nil, nil, err
	// }
	shardsData := make([][]byte, dataShards+parityShards)
	for i := 0; i < dataShards; i++ {
		start := i * shardSize
		end := start + shardSize
		if end > len(data) {
			end = len(data)
		}
		shardsData[i] = data[start:end]
	}

	// Encode shards
	err = enc.Encode(shards)
	if err != nil {
		return nil, nil, err
	}

	metadata := &Metadata{
		FileSize:     int64(len(data)),
		DataShards:   int64( dataShards),
		ParityShards: int64( parityShards),
		ShardHashes:  make([]string, len(shards)),
		IpfsHashes:   make([]string, len(shards)),
		ShardOrder:   make([]int, len(shards)),
	}

	for i, shard := range shards {
		
		hash := sha256.Sum256(shard)
		// fmt.Println("shard hash",hash)
		shardHash := fmt.Sprintf("%x", hash)
		metadata.ShardHashes[i] = shardHash

		encryptedShard, err := encrypt(shard, privateKey) //ecrypt each chunk 
		if err != nil {
			return nil, nil, err
		}
		shards[i] = encryptedShard
	}

	return shards, metadata, nil
}
