package main

func (g *graph) existPath(fromID, toID int) bool {
	nexts := make([]int, 1, len(g.edges))
	marked := make([]bool, len(g.edges))
	nextID := 0
	nexts[0] = fromID
	marked[fromID] = true
	for nextID < len(nexts) {
		currentID := nexts[nextID]
		nextID++
		if currentID == toID {
			return true
		}
		for potentialID := 0; potentialID < len(g.edges[currentID]); potentialID++ {
			if g.edges[currentID][potentialID] > 0 && !marked[potentialID] {
				nexts = append(nexts, potentialID)
				marked[potentialID] = true
			}
		}
	}
	return false
}

func (g *graph) existCycle(atID int) bool {
	if atID < len(g.edges) {
		for j := 0; j < len(g.edges[atID]); j++ {
			if g.edges[atID][j] > 0 {
				if j == atID || g.existPath(j, atID) {
					return true
				}
			}
		}
	}
	return false
}
