package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initEx1() (e exo) {

	e.BasicSetup()

	elementSpacing := 100

	// title setup
	e.titleImage = title1Image
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing

	// graph setup
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.displayGraph = true
	e.displayAdjMatr = false

	e.g.genConnectedGraph(4, 4, 12)
	e.g.linkMatrGraph = false

	nodeSpacing := 300
	e.g.xposition = (windowWidth - nodeSpacing) / 2
	e.g.yposition = 2*elementSpacing + yTitleShift
	e.g.xsize = nodeSpacing
	e.g.ysize = nodeSpacing

	e.g.nodes[0].loopPosition = loopTopLeft
	e.g.nodes[1].xposition = nodeSpacing
	e.g.nodes[1].loopPosition = loopTopRight
	e.g.nodes[2].yposition = nodeSpacing
	e.g.nodes[2].loopPosition = loopBottomLeft
	e.g.nodes[3].xposition = nodeSpacing
	e.g.nodes[3].yposition = nodeSpacing
	e.g.nodes[3].loopPosition = loopBottomRight

	// question setup
	from := rand.Intn(4)
	to := rand.Intn(3)
	if to == from {
		to = 3
	}
	xshift, yshift := ex1Image.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(e.g.ysize+4*elementSpacing+yTitleShift))
		screen.DrawImage(
			ex1Image,
			&options,
		)
		// from label
		options.GeoM.Translate(float64(3*spriteSide), 0)
		xLabel := from % 10
		yLabel := from / 10
		labelSubimage := image.Rect(
			xLabel*spriteSide, (yLabel+1)*spriteSide,
			(xLabel+1)*spriteSide, (yLabel+2)*spriteSide,
		)
		screen.DrawImage(
			graphElementsImage.SubImage(labelSubimage).(*ebiten.Image),
			&options,
		)
		// to label
		options.GeoM.Translate(float64(spriteSide), 0)
		xLabel = to % 10
		yLabel = to / 10
		labelSubimage = image.Rect(
			xLabel*spriteSide, (yLabel+1)*spriteSide,
			(xLabel+1)*spriteSide, (yLabel+2)*spriteSide,
		)
		screen.DrawImage(
			graphElementsImage.SubImage(labelSubimage).(*ebiten.Image),
			&options,
		)
	}

	// answers setup
	e.hasAnswerSheet = true
	answerSize := 128
	e.answers.init((windowWidth-3*answerSize)/2, e.g.ysize+yshift+5*elementSpacing+yTitleShift)
	e.answers.addButton(0, 0, ouiImage)
	e.answers.addButton(2*answerSize, 0, nonImage)

	answer := 1
	if e.g.edges[from][to] > 0 {
		answer = 0
	}

	e.checkResult = func() (bool, bool) {
		return e.answers.clics[answer] >= 1, e.answers.clics[0]+e.answers.clics[1] >= 1
	}

	// return exercise
	return e
}

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
