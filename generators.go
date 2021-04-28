package main

import "math/rand"

// Build a connected graph and returns it
// the graph and the adjacency matrix are linked
// don't care about positions
func (g *graph) genConnectedGraph(numNodes, minEdges, maxEdges int) {
	if minEdges < numNodes-1 {
		minEdges = numNodes - 1
	}

	g.makeEmptyGraph(numNodes)
	g.linkMatrGraph = true

	// Ensure that the graph is connected
	notConnectedNodes := make([]int, numNodes)
	for i := range notConnectedNodes {
		notConnectedNodes[i] = i
	}
	rand.Shuffle(numNodes, func(i, j int) {
		notConnectedNodes[i], notConnectedNodes[j] = notConnectedNodes[j], notConnectedNodes[i]
	})

	g.addEdge(notConnectedNodes[0], notConnectedNodes[1])
	currentNode := rand.Intn(2)
	nodeOrder := rand.Intn(2)
	if nodeOrder == 0 {
		g.addEdge(notConnectedNodes[currentNode], notConnectedNodes[2])
	} else {
		g.addEdge(notConnectedNodes[2], notConnectedNodes[currentNode])
	}
	currentNode = rand.Intn(3)
	nodeOrder = rand.Intn(2)
	if nodeOrder == 0 {
		g.addEdge(notConnectedNodes[currentNode], notConnectedNodes[3])
	} else {
		g.addEdge(notConnectedNodes[3], notConnectedNodes[currentNode])
	}

	// Add a few more edges if needed
	edgesAdded := numNodes - 1
	edgesNeeded := rand.Intn(maxEdges-minEdges+1) + minEdges - edgesAdded
	edgesPossible := numNodes*numNodes - edgesAdded
	if edgesNeeded > edgesPossible {
		edgesNeeded = edgesPossible
	}
	for edgesNeeded > 0 {
		nextEdgeNumber := rand.Intn(edgesPossible) + 1
	edgesLoop:
		for i := range g.edges {
			for j, v := range g.edges[i] {
				if v == 0 {
					nextEdgeNumber--
					if nextEdgeNumber == 0 {
						g.addEdge(i, j)
						break edgesLoop
					}
				}
			}
		}
		edgesAdded++
		edgesPossible--
		edgesNeeded--
	}

	// Prepare draw order
	g.nodesDrawOrder = make([]int, numNodes)
	for i := range g.nodesDrawOrder {
		g.nodesDrawOrder[i] = i
	}

}

func (g *graph) makeEmptyGraph(numNodes int) {
	g.nodes = make([]node, numNodes)
	g.edges = make([][]int, numNodes)
	for i := range g.edges {
		g.edges[i] = make([]int, numNodes)
	}
	g.adjMatr = make([][]int, numNodes)
	for i := range g.adjMatr {
		g.adjMatr[i] = make([]int, numNodes)
	}
}
