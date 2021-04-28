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
