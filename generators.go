package main

import "math/rand"

// Build a connected graph and returns it
// the graph and the adjacency matrix are linked
// don't care about positions
// forbids edge (from, to)
func (g *graph) genConnectedGraph(numNodes, minEdges, maxEdges, noFrom, noTo int) {
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

	if notConnectedNodes[0] != noFrom || notConnectedNodes[1] != noTo {
		g.addEdge(notConnectedNodes[0], notConnectedNodes[1])
	} else {
		g.addEdge(notConnectedNodes[1], notConnectedNodes[0])
	}
	for i := 2; i < numNodes; i++ {
		currentNode := rand.Intn(i)
		nodeOrder := rand.Intn(2)
		if (nodeOrder == 0 && (notConnectedNodes[currentNode] != noFrom || notConnectedNodes[i] != noTo)) ||
			(notConnectedNodes[i] == noFrom && notConnectedNodes[currentNode] == noTo) {
			g.addEdge(notConnectedNodes[currentNode], notConnectedNodes[i])
		} else {
			g.addEdge(notConnectedNodes[i], notConnectedNodes[currentNode])
		}
	}

	// Add a few more edges if needed
	edgesAdded := numNodes - 1
	edgesNeeded := rand.Intn(maxEdges-minEdges+1) + minEdges - edgesAdded
	edgesPossible := numNodes*numNodes - edgesAdded
	if noFrom >= 0 && noTo >= 0 {
		edgesPossible--
	}
	if edgesNeeded > edgesPossible {
		edgesNeeded = edgesPossible
	}
	for edgesNeeded > 0 {
		nextEdgeNumber := rand.Intn(edgesPossible) + 1
	edgesLoop:
		for i := range g.edges {
			for j, v := range g.edges[i] {
				if v == 0 && (i != noFrom || j != noTo) {
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
