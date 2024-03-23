package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func sendFiletoPeer() {
	// Connect to the IPFS API
	api := shell.NewShell("localhost:5001")
	var filePath string
	var ReqCPU, ReqRAM int
	var loc string

	fmt.Println("Enter Required RAM:")
	fmt.Scan(&ReqRAM)
	fmt.Println("Enter Required CPU:")
	fmt.Scan(&ReqCPU)
	fmt.Println("Enter full path of the file you want to save:")
	fmt.Scan(&filePath)
	fmt.Println("Enter the region e.g. Europe,Asia...:")
	fmt.Scan(&loc)
	// fmt.Printf("File size: %d bytes\n", fileSize)

	var totalDuration time.Duration
	for i := 0; i < 1; i++ {
		encodingTime := time.Now()
		data, err := ioutil.ReadFile(string(filePath))
		if err != nil {
			panic(err)
		}

		NodesInfo := getFreeSpace() // getting the available storage space
		groupedNodes, err := splitNodesByLocation(NodesInfo)  // split the stoage network into multiple regions
		if err != nil {
			fmt.Println("Error splitting nodes by location:", err)
			return
		}

		privateKeyBytes, err := hex.DecodeString(constprivateKeyHex)  // decoding the private key
		if err != nil {
			log.Fatal(err)
		}
		

		dataShards, parityShards := 3, 2
		encodedData, metadata, err := encoding_EncryptingData(data, dataShards, parityShards, privateKeyBytes)  // Encoding and encrypting the data
		//Encoding_time := time.Since(encodingTime)
		//fmt.Println("Encoding time", Encoding_time)
		metadata.FileName = filePath
		if err != nil {
			log.Fatal(err)
		}
		size := len(encodedData[0])
		wList := []Workflow{
			{fileSize: float32(size), CPU: ReqCPU, RAM: ReqRAM, Location: loc},
		}
		//psoTime := time.Now()
		allResults, _ := executePSO(wList, groupedNodes) // PSO execution
		//fmt.Println("returned time", psotime)
		//fmt.Println("PSO time", time.Since(psoTime))
		//dist := time.Now()

		var homeRegionNodes []NodeListWithCost
		var homeRegionIndex int
		var found bool
		var listOfAllNodes []NodeListWithCost
		// Find the home region
		for i, region := range allResults {

			listOfAllNodes = append(listOfAllNodes, region...)

			if len(region) > 0 && region[0].location == loc {
				homeRegionNodes = region
				homeRegionIndex = i
				found = true
				//break
			}
		}

		if !found {
			log.Fatalf("Home region %s not found", loc)
			return
		}
		totalShards := dataShards + parityShards
		fmt.Println("Required number of nodes:", totalShards)
		fmt.Printf("Nodes in home (%v) region: %v \n", loc, len(homeRegionNodes))

		sort.Slice(listOfAllNodes, func(i, j int) bool {
			return listOfAllNodes[i].cost < listOfAllNodes[j].cost
		})

		//fmt.Println("all nodes////////:", listOfAllNodes)
		selectedMap := make(map[string]bool)
		if len(homeRegionNodes) >= totalShards {
			//If the home region has enough nodes

			//Finding the remaining nodes other that those to whome we have sent the chunks
			for _, node := range homeRegionNodes {
				selectedMap[node.peer] = true
			}
			var remainingNodes []NodeListWithCost
			for _, node := range listOfAllNodes {
				if _, exists := selectedMap[node.peer]; !exists {
					remainingNodes = append(remainingNodes, node)
				}
			}

			// fmt.Println("Global nodes excluding home region nodes:", remainingNodes)
			// fmt.Println(" distributing to the global best nodes:")
			// distributeData(encodedData, remainingNodes, metadata, api)
			// fmt.Println("Distribution time", time.Since(dist))
			//fmt.Println("Home region nodes:", homeRegionNodes)
			distributeData(encodedData, homeRegionNodes, remainingNodes, metadata, api) // distribute chunks in parallal
			//distributionTime := time.Since(dist)
			//fmt.Println("distribution time", time.Since(encodingTime))
			//fmt.Println("Home region has enogh service providers.\n Uploadin time ", distributionTime + Encoding_time + psotime)
			return
		}
		// if home region does not have enough storage nodes
		requiredAdditionalNodes := totalShards - len(homeRegionNodes)
		allRegionsCount := len(allResults)
		selectedNodes := homeRegionNodes

		for offset := 1; requiredAdditionalNodes > 0; offset++ {
			nextRegionIndex := (homeRegionIndex + offset) % allRegionsCount
			prevRegionIndex := (homeRegionIndex - offset + allRegionsCount) % allRegionsCount

			if nextRegionIndex != homeRegionIndex && len(allResults[nextRegionIndex]) > 0 {
				nodesToTake := min(requiredAdditionalNodes, len(allResults[nextRegionIndex]))
				selectedNodes = append(selectedNodes, allResults[nextRegionIndex][:nodesToTake]...)
				requiredAdditionalNodes -= nodesToTake
				fmt.Printf("%v Chunks will be stored in %v regions\n", nodesToTake, allResults[nextRegionIndex][0].location)
			}

			if requiredAdditionalNodes > 0 && prevRegionIndex != homeRegionIndex && len(allResults[prevRegionIndex]) > 0 {
				nodesToTake := min(requiredAdditionalNodes, len(allResults[prevRegionIndex]))
				selectedNodes = append(selectedNodes, allResults[prevRegionIndex][:nodesToTake]...)
				requiredAdditionalNodes -= nodesToTake
				fmt.Printf("%v Chunks will be stored in %v regions\n", nodesToTake, allResults[prevRegionIndex][0].location)
			}

			if nextRegionIndex == homeRegionIndex && prevRegionIndex == homeRegionIndex {
				log.Println("Not enough nodes available even in the neighboring regions")
				return
			}
		}

		//selectedMap := make(map[string]bool)
		for _, node := range selectedNodes {
			selectedMap[node.peer] = true
		}
		var remainingNodes []NodeListWithCost
		for _, node := range listOfAllNodes {
			if _, exists := selectedMap[node.peer]; !exists {
				remainingNodes = append(remainingNodes, node)
			}
		}
		//t := 1.719 // seconds
		//ET := time.Duration(t * float64(time.Second))
		//distributionTime := time.Since(psoTime)
		duration := time.Since(encodingTime)
		//fmt.Println("distribution time", time.Since(encodingTime))
		//distributeData(encodedData, selectedNodes, metadata, api)
		fmt.Println(" distributing to global best nodes:")
		distributeData(encodedData, selectedNodes, remainingNodes, metadata, api) // distribute chunks in parallal
		totalDuration += duration

	}
	fmt.Println("Uploading time", totalDuration/5)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

