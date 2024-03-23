

    func DSWS(wList []Workflow, nodesInAllRegions [][]listOfPeer) [][]NodeListWithCost {
        //var wg sync.WaitGroup
        var regionalBest [][]NodeListWithCost
        //dist := time.Now()
        for _, nodeGroup := range nodes {
            //wg.Add(1)
    
            //go func(workflows []Workflow, peers []listOfPeer) {
            //defer wg.Done()
            nodeMat := PSO(wList, nodeGroup)
    
            //fmt.Println("GA execution", nodeMat)
    
            for p := range nodeMat {
                sort.Slice(nodeMat[p], func(i, j int) bool {
                    return nodeMat[p][i].cost < nodeMat[p][j].cost
                })
            }
    
            regionalBest = append(regionalBest, nodeMat...)
        }

        var homeRegionNodes []NodeListWithCost
        var homeRegionIndex int
        var found bool
        var listOfAllNodes []NodeListWithCost
        // Find the home region
        for i, region := range regionalBest {
        
            listOfAllNodes = append(listOfAllNodes, region...)
        
            if len(region) > 0 && region[0].location == loc {
                homeRegionNodes = region
                homeRegionIndex = i
                found = true
                //break
            }
        }
        sort.Slice(listOfAllNodes, func(i, j int) bool {
            return listOfAllNodes[i].cost < listOfAllNodes[j].cost
        })
           selectedMap := make(map[string]bool)
        
            for _, node := range homeRegionNodes {
                selectedMap[node.peer] = true
            }
            var globalBest []NodeListWithCost
            for _, node := range listOfAllNodes {
                if _, exists := selectedMap[node.peer]; !exists {
                    globalBest = append(remainingNodes, node)
                }
            }

    
        return regionalBest, globalBest
    }