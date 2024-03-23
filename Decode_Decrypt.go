package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"log"

	"github.com/klauspost/reedsolomon"
)

// decrypt decrypts data using AES-GCM.
// The key argument should be the AES key, either 16, 24, or 32 bytes
// to select AES-128, AES-192, or AES-256.
func decrypt(encryptedData []byte, key []byte) ([]byte, error) {
	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// The nonce size should be the same as the one used during encryption
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// Extract the nonce from the encrypted data
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt the data
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}



func decriptChunk(echunk []byte) []byte {

	privateKeyBytes, err := hex.DecodeString(constprivateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Derived Key for decryption: ", key)
	decryptedShard, err := decrypt(echunk, privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to decrypt shard here %d: ", err)
	}

	return decryptedShard

}

func decodingData(encodedData [][]byte, metadata *Metadata) ([]byte, error) { // reconstruct the file

	// Create Reed-Solomon decoder
	dec, err := reedsolomon.New(int(metadata.DataShards), int(metadata.ParityShards))
	if err != nil {
		return nil, err
	}

	err = dec.Reconstruct(encodedData)
	if err != nil {
		return nil, err
	}

	// Join data shards
	data := bytes.Join(encodedData[:metadata.DataShards], nil)

	// Truncate to original file size
	return data[:metadata.FileSize], nil
}
