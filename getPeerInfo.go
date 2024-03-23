package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

func getFreeSpace() []listOfPeer {

	response,
		err := http.Get("http://127.0.0.1:9094/monitor/metrics/freespace")
	if err != nil {
		fmt.Println("API ttp://127.0.0.1:9094 is not responding\nSwitching to ttp://127.0.0.1:9098")
		// os.Exit(1)
		response,
			err := http.Get("http://127.0.0.1:9098/monitor/metrics/freespace")
		if err != nil {
			fmt.Println("API ttp://127.0.0.1:9098 is not responding\nSwitching to ttp://127.0.0.1:9099")
			response,
				err := http.Get("http://127.0.0.1:9099/monitor/metrics/freespace")
			if err != nil {
				fmt.Println("API ttp://127.0.0.1:9099 is not responding")
				os.Exit(1)
			}
			data := getInfo(response)
			return data
		}
		dtat := getInfo(response)
		return dtat
	}
	data := getInfo(response)
	return data
}

func getInfo(response *http.Response) []listOfPeer {
	var data = make([]listOfPeer, 0)
	bytes,
		err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	prs := []peers{}

	json.Unmarshal(bytes, &prs)

	//fmt.Println("It is working")
	for _, info := range prs {
		var peerlist = listOfPeer{

			peer:  info.Peer,
			space: float32(info.Weight),

			//CPU:   rand.Intn(10) + 1,
		}

		//fmt.Println(peerlist)

		// Arbitarary assign CPU to nodes

		if peerlist.peer == "12D3KooWABUagR9bzQhLzPj21fx9zCnKEoCvJVZcm9VbKMjSHz6h" {
			peerlist.CPU = 9
			peerlist.RAM = 128
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWAJNg7iYXNabxGgMYbQjSWuYUHoKPoYiL5QGYSDRyeQ3W" {
			peerlist.CPU = 8

			peerlist.RAM = 16
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWAh9bEeHPUuUx3EbdKEzyxNRtWtkeWEUhG1d98MUtpb5X" {
			peerlist.CPU = 8

			peerlist.RAM = 32
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWBKYmsvV3JPjZzujmEaDGjZgebKyDWXnt3vwpELobA3Pd" {
			peerlist.CPU = 11
			peerlist.RAM = 64
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWCduUwcqYZBQEK9uVZ8jrS63R6UeF1JuiW2WkkbqgLCZ7" {
			peerlist.CPU = 7
			peerlist.RAM = 128
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWE2AXWhk3eBMy4CEtvSZro9RkToGMmLH1YXYbk9uRoCq2" {
			peerlist.CPU = 9
			peerlist.RAM = 16
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWFuVghmRqxtm8oK4VPYJ48uv27THuyxbZiXQh4kHcf7SK" {
			peerlist.CPU = 9
			peerlist.RAM = 256
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWK7Xd92XUaAkuEPmoQ7mmfVbh2gWKkURdmgS5zHvRBKJ5" {
			peerlist.CPU = 8
			peerlist.RAM = 128
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWKG5zVDvoJP1KD2mTKaFw29tYSzwxxZpz9XGaJLsEbuyx" {
			peerlist.CPU = 11
			peerlist.RAM = 64
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWLjbzzH1ar9gcxjAR2e9DyvCgZmtbuhsy8yL9GvL4R7bq" {
			peerlist.CPU = 11
			peerlist.RAM = 256
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWLxnnQM2Ln2huEt5p2NEoVvsZw7witfhsNkHUSZxBb3NF" {
			peerlist.CPU = 11
			peerlist.RAM = 120
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWMB4E9RtMRiXxdrLf8sy9KrnVDdzLLMkYxeUADnfVYH6f" {
			peerlist.CPU = 11
			peerlist.RAM = 64
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWN7XXQ7BVYCDa439uLRE1NTvekwoRzouAgysQfn3qdE17" {
			peerlist.CPU = 11
			peerlist.RAM = 80
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWRXtjozGvAFeA8JSMdJ26Ei9i2kYSNZK37mp3QwZ8R4io" {
			peerlist.CPU = 11
			peerlist.RAM = 110
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWSc6HTGXPzM3qYgHYsuu7BE9NgZqHJnnmyQsYBysW2cHF" {
			peerlist.CPU = 11
			peerlist.RAM = 110
			peerlist.Location = "Africa"

		}
		// if peerlist.peer == "12D3KooWNAtW79NsLAGQK6UmnPL3TbiVP5SVQWTdjBUDmCnVpbMj" {
		// 	peerlist.CPU = 11
		// 	peerlist.RAM = 110
		// 	peerlist.Location = "Africa"

		// }
		// if peerlist.peer == "12D3KooWPhCAgpjrmUwqFnCxkVTwdoRFuQVH6vpgott1BSXcHdB9" {
		// 	peerlist.CPU = 11
		// 	peerlist.RAM = 110
		// 	peerlist.Location = "Africa"

		// }
		// if peerlist.peer == "12D3KooWQ9P7GPHCbN1Xx6669KSXw6TT77J6ewfdDcPEsah6NbFo" {
		// 	peerlist.CPU = 11
		// 	peerlist.RAM = 110
		// 	peerlist.Location = "Africa"

		// }
		// if peerlist.peer == "12D3KooWQfd8r78DTBEagXkkhs6gS2ajVQRU8i56vQPmeAax4tyn" {
		// 	peerlist.CPU = 11
		// 	peerlist.RAM = 110
		// 	peerlist.Location = "Africa"

		// }
		// if peerlist.peer == "12D3KooWRUMd9PxeDX3mMRViJvciDdkyWbGCKC84w46AEJF9kqeB" {
		// 	peerlist.CPU = 11
		// 	peerlist.RAM = 110
		// 	peerlist.Location = "Africa"

		// }
		// if peerlist.peer == "12D3KooWSjC8TCtzaQVLWCdqf1QWgcmywAtvCZ4w8oD9G8qb3RVg" {
		// 	peerlist.CPU = 11
		// 	peerlist.RAM = 110
		// 	peerlist.Location = "Africa"

		// }
		// if peerlist.peer == "12D3KooWSy6eZ1v7YDUbgWbfwhMhdaLKa6szKUUrS8RrpATFH2hT" {
		// 	peerlist.CPU = 11
		// 	peerlist.RAM = 110
		// 	peerlist.Location = "Africa"

		// }

		data = append(data, peerlist)
		//fmt.Println("All the Nodes with availabale free disk space")

	}

	// fmt.Println("All the Nodes with availabale free disk space")
	// for i, dt := range data {
	// 	fmt.Println(i, dt)
	// }

	sort.Slice(data, func(i, j int) bool {
		return data[i].space > data[j].space // use ">" if you want descending order
	})
	// for _, fs := range PeerData {

	// 	fsp = append(fsp, fs.space)
	// }

	//fmt.Println("lenght of string ",string(responseDat))
	//ioutil.WriteFile()
	//print(responseDat)

	return data
}
