package main

import (
	"fmt"
	
)



type Metadata struct {
	FileName     string
	FileSize     int64
	DataShards   int64
	ParityShards int64
	ShardHashes  []string
	ShardOrder   []int
	IpfsHashes   []string
}
type Workflow struct {
	//SecFlag bool    // Security flag
	//hashOfFile string // Security level
	fileSize float32
	CPU      int // CPU utilization
	RAM      int // Memory utilization
	Location string

	//Best int // Best position
}

// structure for unmarshelling the the data for getting free available space frome nodes
type peers struct {
	Name          string
	Peer          string
	Value         string
	Expire        int64
	Valid         bool
	Weight        int64
	Partitionable bool
	Received_at   int64
}

type listOfPeer struct {
	peer     string
	space    float32
	CPU      int
	RAM      int
	Location string
}

var PeerData = make([]listOfPeer, 0)

const constprivateKeyHex = "6b9ed19b8fa8cb1966e4e4e1d312a14d936b010bac4e325db3c2a83e4da60620" // Private Key for AES

func switchMenu() {

	var input int
	fmt.Print("Please select the available options below.\n 1 : For file storage \n 2 : For retreiving files\n")
	fmt.Scan(&input)

	switch input {
	// case 1:
	// 	getPostID()
	// case 2:
	// 	getFreeSpace()
	case 1:
		sendFiletoPeer()
	case 2:

		getFile()
	default:
		fmt.Println("Invalid command")
	}
}

func main() {

	for {
		switchMenu()
	}
	

}

