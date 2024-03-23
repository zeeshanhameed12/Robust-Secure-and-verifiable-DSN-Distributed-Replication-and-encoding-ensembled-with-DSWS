package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

type ChunkResult struct {
	Index int
	Data  []byte
	Valid bool
}

func getFile() {
	fmt.Println("Provide the name of the file you want to retrieve:")
	var fileName string
	fmt.Scan(&fileName)

	metadata, err := LoadMetadataFromContract(fileName) //Loading the metadata from blockchain
	if err != nil {
		panic(err)
	}

	retreiveddata := make([][]byte, len(metadata.IpfsHashes))
	fmt.Println("########Retreiving########")
	api := shell.NewShell("localhost:5001")

	ctx, _ := context.WithCancel(context.Background())

	done := make(chan bool)
	startt := time.Now()
	//go func() {
	results, validChunks := downloadChunks(ctx, metadata, api) // Downloading the chunks in parallal
	if validChunks < int64(metadata.DataShards) {
		fmt.Println("Not enough valid chunks to reconstruct the file")
		return
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Index < results[j].Index
	})
	//fmt.Println("length of doenloaded results:", len(results))
	//fmt.Println("number of valid chunks", validChunks)
	for i, chunk := range results {
		retreiveddata[i] = chunk.Data
	}

	decodedData, err := decodingData(retreiveddata, metadata)
	if err != nil {
		panic(err)
	}

	// Save the decoded file
	err = ioutil.WriteFile(fileName, decodedData, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("File retrieved successfully!")
	close(done)
	//}()

	//<-done // This will block until the file is retrieved or an error occurs
	elapsed := time.Since(startt)
	//totalDuration += elapsed
	fmt.Println(elapsed)
	//time.Sleep(2 * time.Second)
}

func downloadChunks(ctx context.Context, metadata *Metadata, api *shell.Shell) ([]ChunkResult, int64) {
	results := make([]ChunkResult, len(metadata.IpfsHashes))
	validChunks := int64(0)
	mux := &sync.Mutex{}
	done := make(chan bool)
	stop := false // Use a boolean to control the loop

	for i, hash := range metadata.IpfsHashes {
		// Use the stop variable to check if we should continue launching goroutines
		mux.Lock()
		if stop {
			mux.Unlock()
			break // Correctly exit the for loop
		}
		mux.Unlock()

		go func(i int, hash string) {
			select {
			case <-ctx.Done():
				return
			default:
				fileContents, err := api.Cat(hash)
				if err != nil {
					return
				}
				defer fileContents.Close()

				data, err := io.ReadAll(fileContents)
				if err != nil {
					return
				}

				decryptedChunk := decriptChunk(data) // decrypt the chunk
				hash := sha256.Sum256(decryptedChunk)
				//fmt.Printf(" downloaded shared %v\n uploaded shard hash %v\n", fmt.Sprintf("%x", sha256.Sum256(decryptedChunk)), metadata.ShardHashes[i])
				if fmt.Sprintf("%x", hash) == metadata.ShardHashes[i] { // validate the chunks
					mux.Lock()
					results[i] = ChunkResult{Index: i, Data: decryptedChunk, Valid: true}
					validChunks++
					if validChunks == metadata.DataShards {  // download only minimum required vaid chunks to recover the file (number of data chunks)
						stop = true // stop the loop
						done <- true // Use done channel to signal completion
					}
					mux.Unlock()
				} else {
					results[i] = ChunkResult{Index: i}
				}
			}
		}(i, hash)
	}

	<-done // Wait here until done signal is received
	return results, validChunks
}
