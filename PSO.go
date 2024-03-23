package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type NodeListWithCost struct {
	index    int
	peer     string
	location string
	cost     float32
}

type Particle struct {
	position     float32 // Assume position is represented as a float for simplicity
	velocity     float32
	bestPosition float32
	bestCost     float32
}

func pso(wList []Workflow, nList []listOfPeer) [][]NodeListWithCost {
	nodeMat := make([][]NodeListWithCost, len(wList))
	inertia := float32(0.3)
	cognitiveConst := float32(1.3)
	socialConst := float32(1.5)
	randCognitive := rand.Float32()
	randSocial := rand.Float32()

	particles := make([]Particle, len(nList))
	globalBest := float32(^uint(0) >> 1) // Initialized with max possible value

	// Initializing particles
	for i := range particles {
		particles[i] = Particle{
			position:     rand.Float32(), // Initial position should be set according to your problem's context
			velocity:     0,              // Starting with zero velocity
			bestPosition: 0,              // To be updated during the iterations
			bestCost:     float32(^uint(0) >> 1),
		}
	}

	iterations := 1000
	for iter := 0; iter < iterations; iter++ {
		for i, w := range wList {
			var list []NodeListWithCost

			for j, n := range nList {
				if n.CPU >= w.CPU && n.space >= w.fileSize && float32(n.RAM) >= float32(w.RAM) {
					p := particles[j]

					// Assuming the cost function as the objective to minimize
					c := cost(w, n)

					if c < p.bestCost {
						p.bestCost = c
						p.bestPosition = p.position
					}

					if c < globalBest {
						globalBest = c
					}

					// Update velocity and position
					p.velocity = inertia*p.velocity +
						cognitiveConst*randCognitive*(p.bestPosition-p.position) +
						socialConst*randSocial*(globalBest-p.position)

					p.position += p.velocity

					particles[j] = p // Update the particle in the slice
					list = append(list, NodeListWithCost{j, n.peer, n.Location, c})
				}
			}
			nodeMat[i] = list
		}
	}

	return nodeMat
}

// func pso(wList []Workflow, nList []listOfPeer) [][]NodeListWithCost {
// 	nodeMat := make([][]NodeListWithCost, len(wList))

// 	for i, w := range wList {
// 		var list []NodeListWithCost

// 		for j, n := range nList {
// 			if n.CPU >= w.CPU && n.space >= w.fileSize && float32(n.RAM) >= float32(w.RAM) {
// 				c := cost(w, n)
// 				list = append(list, NodeListWithCost{j, n.peer, n.Location ,c})
// 			}
// 		}

// 		nodeMat[i] = list
// 	}

// 	return nodeMat
// }

func cost(w Workflow, n listOfPeer) float32 {
	return float32(w.CPU/n.CPU) + (w.fileSize / n.space) + (float32(w.RAM) / float32(n.RAM))
}

func splitNodesByLocation(nodes []listOfPeer) ([][]listOfPeer, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes available")
	}

	var As_Nodes, Af_Nodes, Eu_Nodes []listOfPeer

	for _, node := range nodes {
		switch node.Location {
		case "Asia":
			As_Nodes = append(As_Nodes, node)
		case "Africa":
			Af_Nodes = append(Af_Nodes, node)
		case "Europe": // Added "Europe" assuming Eu_Nodes means European nodes
			Eu_Nodes = append(Eu_Nodes, node)
		default:
			return nil, fmt.Errorf("unknown location: %s", node.Location)
		}
	}

	return [][]listOfPeer{As_Nodes, Af_Nodes, Eu_Nodes}, nil
}

func executePSO(wList []Workflow, nodes [][]listOfPeer) ([][]NodeListWithCost, time.Duration) {
	var wg sync.WaitGroup
	results := make(chan [][]NodeListWithCost, len(nodes))
	dist := time.Now()
	for _, nodeGroup := range nodes {
		wg.Add(1)

		go func(workflows []Workflow, peers []listOfPeer) {
			defer wg.Done()
			nodeMat := pso(workflows, peers)
			for p := range nodeMat {
				sort.Slice(nodeMat[p], func(i, j int) bool {
					return nodeMat[p][i].cost < nodeMat[p][j].cost
				})
			}
			results <- nodeMat
		}(wList, nodeGroup)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	execTime :=time.Since(dist)
	fmt.Println("PSO execution time :", execTime)
	var allResults [][]NodeListWithCost
	
	for result := range results {
		allResults = append(allResults, result...)
	
	}

	return allResults, execTime
}

// var list = make([]NodeListWithCost, 0)

// func pso(wList []Workflow, nList []listOfPeer) [][]NodeListWithCost {
// 	// weight := 0.3
// 	// cp := 1.5
// 	// cg := 1.3

// 	nodeMat := make([][]NodeListWithCost, len(wList))
// 	//fmt.Printf("Wlict%v",wList)
// 	// Initialize the swarm and particle properties

// 	for i := 0; i < len(wList); i++ {
// 		//nodeList := make([]float64, 0)
// 		// currPos := 0
// 		// swarmBest := 0
// 		// wList[i].Best = 0
// 		// wList[i].Vel = 0

// 		for j := 0; j < len(nList); j++ {
// 			// rp := rand.Float64()
// 			// rg := rand.Float64()
// 			if nList[j].CPU >= wList[i].CPU && nList[j].space >= wList[i].fileSize && float32(nList[j].RAM) >= float32(wList[i].RAM) && wList[i].Location == nList[j].Location {
// 				c := cost(wList[i], nList[j])
// 				var indexCost = NodeListWithCost{
// 					cost:  c,
// 					index: j,
// 					peer:  nList[j].peer,
// 				}
// 				list = append(list, indexCost)
// 				//nodeList = append(nodeList, c)

// 			}

// 		}

// 		//sort.Slice(nodeList, func(i, j int) bool { return nodeList[i] < nodeList[j] })

// 		nodeMat[i] = list
// 		list = nil

// 	}

// 	return nodeMat

// }

// // Cost Function

// func cost(w Workflow, n listOfPeer) float32 {
// 	// Calculate the cost of assigning workflow w to node n
// 	// (you can implement this function according to your requirements)
// 	c := float32(w.CPU/n.CPU) + (w.fileSize / n.space) + (float32(w.RAM) / float32(n.RAM))
// 	return c
// }
