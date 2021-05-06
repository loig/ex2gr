package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type graph struct {
	xposition      int
	yposition      int
	xsize          int
	ysize          int
	nodes          []node
	nodesDrawOrder []int
	edges          [][]int
	adjMatr        [][]int
	xmatrposition  int
	ymatrposition  int
	successorsList [][]int
	xlistposition  int
	ylistposition  int
	linkMatrGraph  bool
}

const ( // trigonometric order, this is important
	loopBottomRight int = iota
	loopBottomLeft
	loopTopLeft
	loopTopRight
)

type node struct {
	xposition    int
	yposition    int
	loopPosition int
}

const (
	nodeXSize          = spriteSide
	nodeYSize          = spriteSide
	edgeSelectionWidth = 10
	matrixCellSize     = spriteSide
)

func (g *graph) selectListElement(x, y int) (int, int, bool) {
	xx := x - g.xlistposition
	yy := y - g.ylistposition
	if yy >= 0 && yy < spriteSide*len(g.successorsList) {
		yyy := yy / spriteSide
		if xx >= 0 && xx < spriteSide*(len(g.successorsList[yyy])+3) {
			return yyy, xx/spriteSide - 2, yy-(yyy*spriteSide) < spriteSide/2
		}
		if xx >= 0 && xx < spriteSide*4 && len(g.successorsList[yyy]) == 0 {
			return yyy, 1, yy-(yyy*spriteSide) < spriteSide/2
		}
	}
	return -1, -1, false
}

func (g *graph) updateListElement(i, j int, up bool) {
	if i >= 0 && j >= 0 {
		if len(g.successorsList[i]) < len(g.successorsList) {
			// if all nodes are successors, no node can be changed nor added
			if j < len(g.successorsList[i]) {
				// change a node
				g.successorsList[i][j] = getNextOkNode(g.successorsList[i], g.successorsList[i][j], len(g.successorsList))
			} else {
				// change the size of the list
				if up || len(g.successorsList[i]) <= 0 {
					if g.successorsList[i] == nil {
						g.successorsList[i] = make([]int, 0, len(g.successorsList))
					}
					g.successorsList[i] = append(g.successorsList[i], getNextOkNode(g.successorsList[i], -1, len(g.successorsList)))
				} else {
					g.successorsList[i] = g.successorsList[i][:len(g.successorsList[i])-1]
				}
			}
		} else {
			if !up {
				g.successorsList[i] = g.successorsList[i][:len(g.successorsList[i])-1]
			}
		}
	}
}

func getNextOkNode(nodes []int, start int, numNodes int) int {

	res := (start + 1) % numNodes
	for res != start {
		var found bool
		for _, n := range nodes {
			found = n == res
			if found {
				break
			}
		}
		if !found {
			return res
		}
		res = (res + 1) % numNodes
	}
	return res
}

func (g *graph) selectMatrixCell(x, y int) (int, int) {
	xx := x - g.xmatrposition
	yy := y - g.ymatrposition
	if yy >= matrixCellSize && yy < matrixCellSize*(len(g.adjMatr)+1) {
		if xx >= matrixCellSize && xx < matrixCellSize*(len(g.adjMatr[0])+1) {
			return yy/matrixCellSize - 1, xx/matrixCellSize - 1
		}
	}
	return -1, -1
}

func (g *graph) updateMatrixCell(i, j int) {
	g.adjMatr[i][j] = (g.adjMatr[i][j] + 1) % 2
	if g.linkMatrGraph {
		g.edges[i][j] = g.adjMatr[i][j]
	}
}

func (g *graph) selectNode(x, y int) int {
	for i := len(g.nodesDrawOrder) - 1; i >= 0; i-- {
		nodeID := g.nodesDrawOrder[i]
		distToCenterX := (x - g.xposition) - g.nodes[nodeID].xposition
		distToCenterY := (y - g.yposition) - g.nodes[nodeID].yposition
		radius := spriteSide/2 - 3 // hardcoded to correspond to the sprite of the graph node
		if distToCenterX*distToCenterX+distToCenterY*distToCenterY <= radius*radius {
			return nodeID
		}
	}
	return -1
}

func (g *graph) selectEdge(x, y int) (int, int) {
	for i := range g.edges {
		for j, v := range g.edges[i] {
			if v > 0 {
				if i != j {
					if g.isOnNormalEdge(i, j, x-g.xposition, y-g.yposition) {
						return i, j
					}
				} else {
					if g.isOnLoop(i, x-g.xposition, y-g.yposition) {
						return i, i
					}
				}
			}
		}
	}
	return -1, -1
}

func (g *graph) isOnNormalEdge(from, to, x, y int) bool {
	// only check that if mouse is in a rectangle of which the edge is a diagonal
	xfrom := float64(g.nodes[from].xposition)
	yfrom := float64(g.nodes[from].yposition)
	xto := float64(g.nodes[to].xposition)
	yto := float64(g.nodes[to].yposition)
	dx := xto - xfrom
	dy := yto - yfrom
	length := math.Sqrt(dx*dx + dy*dy)
	vx := -dy / length
	vy := dx / length
	shift := float64(spriteSide) / 8
	xfrom = xfrom + shift*vx
	yfrom = yfrom + shift*vy
	xto = xto + shift*vx
	yto = yto + shift*vy
	xx := float64(x)
	yy := float64(y)
	if ((xx < xfrom && xx > xto) || (xx > xfrom && xx < xto) || math.Abs(xfrom-xto) < float64(spriteSide)/4) &&
		((yy < yfrom && yy > yto) || (yy > yfrom && yy < yto) || math.Abs(yfrom-yto) < float64(spriteSide)/4) {
		a := yto - yfrom
		b := xfrom - xto
		c := -(b*yfrom + a*xfrom)
		dist := math.Abs(a*xx+b*yy+c) / math.Sqrt(a*a+b*b)
		return dist < edgeSelectionWidth
	}
	return false
}

func (g *graph) isOnLoop(nodeID, x, y int) bool {
	xCenter := g.nodes[nodeID].xposition
	yCenter := g.nodes[nodeID].yposition
	xShift := spriteSide / 2
	yShift := spriteSide / 2
	switch g.nodes[nodeID].loopPosition {
	case loopBottomRight:
		xCenter += xShift
		yCenter += yShift
	case loopTopRight:
		xCenter += xShift
		yCenter -= yShift
	case loopTopLeft:
		xCenter -= xShift
		yCenter -= yShift
	case loopBottomLeft:
		xCenter -= xShift
		yCenter += yShift
	}
	distToCenterX := x - xCenter
	distToCenterY := y - yCenter
	radius := spriteSide/2 - 3 // hardcoded to correspond to the sprite of the loop
	return distToCenterX*distToCenterX+distToCenterY*distToCenterY <= radius*radius
}

func (g *graph) moveLoop(nodeID int) {
	g.nodes[nodeID].loopPosition = (g.nodes[nodeID].loopPosition + 1) % 4
}

func (g *graph) updateNodePosition(x, y, nodeID int) {
	currentID := g.nodesDrawOrder[len(g.nodesDrawOrder)-1]
	g.nodesDrawOrder[len(g.nodesDrawOrder)-1] = nodeID
	i := len(g.nodesDrawOrder) - 2
	for currentID != nodeID && i >= 0 {
		tmpID := currentID
		currentID = g.nodesDrawOrder[i]
		g.nodesDrawOrder[i] = tmpID
		i--
	}
	g.nodes[nodeID].xposition = x - g.xposition
	if g.nodes[nodeID].xposition > g.xsize {
		g.nodes[nodeID].xposition = g.xsize
	}
	if g.nodes[nodeID].xposition < 0 {
		g.nodes[nodeID].xposition = 0
	}
	g.nodes[nodeID].yposition = y - g.yposition
	if g.nodes[nodeID].yposition > g.ysize {
		g.nodes[nodeID].yposition = g.ysize
	}
	if g.nodes[nodeID].yposition < 0 {
		g.nodes[nodeID].yposition = 0
	}

}

func (g *graph) addEdge(from, to int) {
	//if from != to {
	if from == to && g.edges[from][to] < 1 {
		g.nodes[from].loopPosition = 0
	}
	g.edges[from][to] = 1
	if g.linkMatrGraph {
		g.adjMatr[from][to] = 1
		if g.successorsList[from] == nil {
			g.successorsList[from] = make([]int, 0, len(g.edges))
		}
		g.successorsList[from] = append(g.successorsList[from], to)
	}
	//}
}

func (g *graph) removeEdge(from, to int) {
	g.edges[from][to] = 0
	if g.linkMatrGraph {
		g.adjMatr[from][to] = 0
		toID := 0
		for toID < len(g.successorsList[from]) && g.successorsList[from][toID] != to {
			toID++
		}
		if toID < len(g.successorsList[from]) {
			for toID < len(g.successorsList[from])-1 {
				g.successorsList[from][toID] = g.successorsList[from][toID+1]
				toID++
			}
			g.successorsList[from] = g.successorsList[from][:len(g.successorsList[from])-1]
		}
	}
}

func (g *graph) draw(screen *ebiten.Image, selectedNodes []int, selectedEdge []int, selectedCell []int) {
	g.drawGraph(screen, selectedNodes, selectedEdge)
	g.drawMatrix(screen, selectedCell)
}

func (g *graph) drawList(screen *ebiten.Image, selectedElement []int, up bool, modifiable bool) {

	for i := 0; i < len(g.successorsList); i++ {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(g.xlistposition), float64(g.ylistposition+i*spriteSide))
		xLabel := i % 10
		yLabel := i / 10
		if selectedElement[0] == i {
			screen.DrawImage(
				graphElementsImage.SubImage(nodeSelectedSubimage).(*ebiten.Image),
				&options,
			)
		}
		labelSubimage := image.Rect(
			xLabel*spriteSide, (yLabel+1)*spriteSide,
			(xLabel+1)*spriteSide, (yLabel+2)*spriteSide,
		)
		screen.DrawImage(
			graphElementsImage.SubImage(labelSubimage).(*ebiten.Image),
			&options,
		)
		options.GeoM.Translate(float64(spriteSide), 0)
		screen.DrawImage(
			graphElementsImage.SubImage(twoDotsSubimage).(*ebiten.Image),
			&options,
		)
		if g.successorsList[i] == nil || len(g.successorsList[i]) == 0 {
			options.GeoM.Translate(float64(spriteSide), 0)
			screen.DrawImage(
				graphElementsImage.SubImage(emptyListSubimage).(*ebiten.Image),
				&options,
			)
			if modifiable {
				options.GeoM.Translate(float64(spriteSide), 0)
				if selectedElement[0] == i && selectedElement[1] >= 1 && up {
					screen.DrawImage(
						graphElementsImage.SubImage(moreSelectedListSubimage).(*ebiten.Image),
						&options,
					)
				} else {
					screen.DrawImage(
						graphElementsImage.SubImage(moreListSubimage).(*ebiten.Image),
						&options,
					)
				}
			}
		} else {
			screen.DrawImage(
				graphElementsImage.SubImage(openListSubimage).(*ebiten.Image),
				&options,
			)
			for jID, j := range g.successorsList[i] {
				options.GeoM.Translate(float64(spriteSide), 0)
				xLabel := j % 10
				yLabel := j / 10
				if selectedElement[0] == i && selectedElement[1] == jID {
					screen.DrawImage(
						graphElementsImage.SubImage(nodeSelectedSubimage).(*ebiten.Image),
						&options,
					)
				}
				labelSubimage := image.Rect(
					xLabel*spriteSide, (yLabel+1)*spriteSide,
					(xLabel+1)*spriteSide, (yLabel+2)*spriteSide,
				)
				screen.DrawImage(
					graphElementsImage.SubImage(labelSubimage).(*ebiten.Image),
					&options,
				)
				if jID < len(g.successorsList[i])-1 {
					screen.DrawImage(
						graphElementsImage.SubImage(sepListSubimage).(*ebiten.Image),
						&options,
					)
				}
			}
			options.GeoM.Translate(float64(spriteSide), 0)
			screen.DrawImage(
				graphElementsImage.SubImage(closeListSubimage).(*ebiten.Image),
				&options,
			)
			if modifiable {
				if selectedElement[0] == i && selectedElement[1] >= len(g.successorsList[i]) {
					if up {
						if len(g.successorsList[i]) < len(g.successorsList) {
							screen.DrawImage(
								graphElementsImage.SubImage(moreSelectedListSubimage).(*ebiten.Image),
								&options,
							)
						}
						if len(g.successorsList[i]) > 0 {
							screen.DrawImage(
								graphElementsImage.SubImage(lessListSubimage).(*ebiten.Image),
								&options,
							)
						}
					} else {
						if len(g.successorsList[i]) < len(g.successorsList) {
							screen.DrawImage(
								graphElementsImage.SubImage(moreListSubimage).(*ebiten.Image),
								&options,
							)
						}
						if len(g.successorsList[i]) > 0 {
							screen.DrawImage(
								graphElementsImage.SubImage(lessSelectedListSubimage).(*ebiten.Image),
								&options,
							)
						}
					}
				} else {
					if len(g.successorsList[i]) < len(g.successorsList) {
						screen.DrawImage(
							graphElementsImage.SubImage(moreListSubimage).(*ebiten.Image),
							&options,
						)
					}
					if len(g.successorsList[i]) > 0 {
						screen.DrawImage(
							graphElementsImage.SubImage(lessListSubimage).(*ebiten.Image),
							&options,
						)
					}
				}
			}
		}

	}

}

func (g *graph) drawMatrix(screen *ebiten.Image, selectedCell []int) {

	optionsVert := ebiten.DrawImageOptions{}
	optionsVert.GeoM.Translate(float64(g.xmatrposition), float64(g.ymatrposition))
	screen.DrawImage(
		graphElementsImage.SubImage(matrixLeftSubimage).(*ebiten.Image),
		&optionsVert,
	)
	for i := range g.adjMatr {
		optionsVert.GeoM.Translate(0, float64(matrixCellSize))
		screen.DrawImage(
			graphElementsImage.SubImage(matrixLeftSubimage).(*ebiten.Image),
			&optionsVert,
		)
		if selectedCell[0] == i {
			screen.DrawImage(
				graphElementsImage.SubImage(nodeSelectedSubimage).(*ebiten.Image),
				&optionsVert,
			)
		}
		xLabel := i % 10
		yLabel := i / 10
		labelSubimage := image.Rect(
			xLabel*spriteSide, (yLabel+1)*spriteSide,
			(xLabel+1)*spriteSide, (yLabel+2)*spriteSide,
		)
		screen.DrawImage(
			graphElementsImage.SubImage(labelSubimage).(*ebiten.Image),
			&optionsVert,
		)
	}

	optionsHor := ebiten.DrawImageOptions{}
	optionsHor.GeoM.Translate(float64(g.xmatrposition), float64(g.ymatrposition))
	screen.DrawImage(
		graphElementsImage.SubImage(matrixTopSubimage).(*ebiten.Image),
		&optionsHor,
	)
	if len(g.adjMatr) > 0 {
		for j := range g.adjMatr[0] {
			optionsHor.GeoM.Translate(float64(matrixCellSize), 0)
			screen.DrawImage(
				graphElementsImage.SubImage(matrixTopSubimage).(*ebiten.Image),
				&optionsHor,
			)
			if selectedCell[1] == j {
				screen.DrawImage(
					graphElementsImage.SubImage(nodeSelectedSubimage).(*ebiten.Image),
					&optionsHor,
				)
			}
			xLabel := j % 10
			yLabel := j / 10
			labelSubimage := image.Rect(
				xLabel*spriteSide, (yLabel+1)*spriteSide,
				(xLabel+1)*spriteSide, (yLabel+2)*spriteSide,
			)
			screen.DrawImage(
				graphElementsImage.SubImage(labelSubimage).(*ebiten.Image),
				&optionsHor,
			)
		}
	}

	for i := range g.adjMatr {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(g.xmatrposition), float64(g.ymatrposition+(i+1)*matrixCellSize))
		for j, v := range g.adjMatr[i] {
			options.GeoM.Translate(float64(matrixCellSize), 0)
			if selectedCell[0] == i && selectedCell[1] == j {
				screen.DrawImage(
					graphElementsImage.SubImage(nodeSelectedSubimage).(*ebiten.Image),
					&options,
				)
			}
			valueSubimage := image.Rect(
				(6+v)*spriteSide, 3*spriteSide,
				(7+v)*spriteSide, 4*spriteSide,
			)
			screen.DrawImage(
				graphElementsImage.SubImage(valueSubimage).(*ebiten.Image),
				&options,
			)
		}
	}

}

func (g *graph) drawGraphLayout(screen *ebiten.Image) {

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(g.xposition-spriteSide), float64(g.yposition-spriteSide))
	screen.DrawImage(
		graphElementsImage.SubImage(graphLayoutTopLeftSubimage).(*ebiten.Image),
		&options,
	)

	options = ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(g.xposition+g.xsize), float64(g.yposition-spriteSide))
	screen.DrawImage(
		graphElementsImage.SubImage(graphLayoutTopRightSubimage).(*ebiten.Image),
		&options,
	)

	options = ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(g.xposition+g.xsize), float64(g.yposition+g.ysize))
	screen.DrawImage(
		graphElementsImage.SubImage(graphLayoutBottomRightSubimage).(*ebiten.Image),
		&options,
	)

	options = ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(g.xposition-spriteSide), float64(g.yposition+g.ysize))
	screen.DrawImage(
		graphElementsImage.SubImage(graphLayoutBottomLeftSubimage).(*ebiten.Image),
		&options,
	)
}

func (g *graph) drawGraph(screen *ebiten.Image, selectedNodes []int, selectedEdge []int) {

	g.drawGraphLayout(screen)

	for i := range g.edges {
		for j, v := range g.edges[i] {
			if v > 0 {
				g.drawEdge(i, g.nodes[j].xposition, g.nodes[j].yposition, screen, i == selectedEdge[0] && j == selectedEdge[1], i == j, j)
			}
		}
	}

	for _, i := range g.nodesDrawOrder {
		g.drawNode(i, screen, i == selectedNodes[0] || i == selectedNodes[1])
	}

}

func (g *graph) drawNode(i int, screen *ebiten.Image, selected bool) {
	xLeft := float64(g.nodes[i].xposition+g.xposition) - float64(nodeXSize)/2
	yTop := float64(g.nodes[i].yposition+g.yposition) - float64(nodeYSize)/2
	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(xLeft, yTop)
	subimage := nodeSubimage
	if selected {
		subimage = nodeSelectedSubimage
	}
	screen.DrawImage(
		graphElementsImage.SubImage(subimage).(*ebiten.Image),
		&options,
	)
	xLabel := i % 10
	yLabel := i / 10
	labelSubimage := image.Rect(
		xLabel*spriteSide, (yLabel+1)*spriteSide,
		(xLabel+1)*spriteSide, (yLabel+2)*spriteSide,
	)
	screen.DrawImage(
		graphElementsImage.SubImage(labelSubimage).(*ebiten.Image),
		&options,
	)
	//ebitenutil.DebugPrintAt(screen, g.nodes[i].label, g.nodes[i].xposition+g.xposition-nodeXSize/4, g.nodes[i].yposition+g.yposition-nodeYSize/2)
}

func (g *graph) drawEdge(from, xto, yto int, screen *ebiten.Image, selected bool, isLoop bool, toNode int) {
	if isLoop {
		g.drawLoopEdge(from, screen, selected)
		return
	}
	g.drawNormalEdge(from, xto, yto, screen, selected, toNode)
}

func (g *graph) drawNormalEdge(from, xto, yto int, screen *ebiten.Image, selected bool, toNode int) {
	xFrom := float64(g.nodes[from].xposition + g.xposition)
	yFrom := float64(g.nodes[from].yposition + g.yposition)
	if toNode >= 0 {
		xto = g.nodes[toNode].xposition
		yto = g.nodes[toNode].yposition
	}
	xTo := float64(xto + g.xposition)
	yTo := float64(yto + g.yposition)
	dx := xTo - xFrom
	dy := yTo - yFrom
	angle := math.Atan(dy / (dx))
	if xTo < xFrom {
		angle += math.Pi
	}
	angle -= math.Pi / 2
	length := math.Sqrt(dx*dx + dy*dy)
	vx := -dy / length
	vy := dx / length
	xTrans := xFrom + float64(spriteSide/2)*vx
	yTrans := yFrom + float64(spriteSide/2)*vy
	options := ebiten.DrawImageOptions{}
	options.GeoM.Rotate(angle)
	options.GeoM.Translate(xTrans, yTrans)
	subimage := edgeSubimage
	if selected {
		subimage = edgeSelectedSubimage
	}
	if toNode >= 0 {
		length -= float64(spriteSide)/2 - 5
	}
	if length >= float64(spriteSide) {
		i := float64(spriteSide)
		for i+float64(spriteSide) < length {
			screen.DrawImage(
				graphElementsImage.SubImage(subimage).(*ebiten.Image),
				&options,
			)
			i += float64(spriteSide)
			options.GeoM.Translate(float64(spriteSide)*vy, float64(-spriteSide)*vx)
		}
		trans := length - i
		screen.DrawImage(
			graphElementsImage.SubImage(subimage).(*ebiten.Image),
			&options,
		)
		options.GeoM.Translate(trans*vy, -trans*vx)
	}
	toSubimage := edgeToSubimage
	if selected {
		toSubimage = edgeToSelectedSubimage
	}
	screen.DrawImage(
		graphElementsImage.SubImage(toSubimage).(*ebiten.Image),
		&options,
	)
	/*
		col := color.RGBA{255, 255, 255, 255}
		if selected {
			col = color.RGBA{255, 0, 0, 255}
		}
		ebitenutil.DrawLine(screen, xFrom, yFrom, xTo, yTo, col)
	*/
}

func (g *graph) drawLoopEdge(from int, screen *ebiten.Image, selected bool) {
	xLeft := float64(g.nodes[from].xposition + g.xposition) //- float64(nodeXSize)/2 + float64(spriteSide/2+10)
	yTop := float64(g.nodes[from].yposition + g.yposition)  //- float64(nodeYSize)/2 - float64(spriteSide/2)
	options := ebiten.DrawImageOptions{}
	options.GeoM.Rotate(float64(g.nodes[from].loopPosition) * math.Pi / 2)
	options.GeoM.Translate(xLeft, yTop)
	subimage := loopSubimage
	if selected {
		subimage = loopSelectedSubimage
	}
	screen.DrawImage(
		graphElementsImage.SubImage(subimage).(*ebiten.Image),
		&options,
	)
}

func (g *graph) drawSelectedNode(i int, screen *ebiten.Image) {
	g.drawNode(i, screen, true)
}

func (g *graph) drawSelectedEdge(from, to int, screen *ebiten.Image) {
	g.drawEdge(from, g.nodes[to].xposition, g.nodes[to].yposition, screen, true, from == to, to)
}

func (g *graph) drawSelectedMatrixCell(i, j int, screen *ebiten.Image) {
	g.drawMatrixCell(i+1, j+1, screen)
	g.drawMatrixCell(0, j+1, screen)
	g.drawMatrixCell(i+1, 0, screen)
}

func (g *graph) drawMatrixCell(i, j int, screen *ebiten.Image) {
	col := color.RGBA{255, 0, 0, 255}
	xLeft := float64(j*matrixCellSize + g.xmatrposition)
	xRight := float64((j+1)*matrixCellSize + g.xmatrposition)
	yTop := float64(i*matrixCellSize + g.ymatrposition)
	yBottom := float64((i+1)*matrixCellSize + g.ymatrposition)
	ebitenutil.DrawLine(screen, xLeft, yTop, xLeft, yBottom, col)
	ebitenutil.DrawLine(screen, xLeft, yBottom, xRight, yBottom, col)
	ebitenutil.DrawLine(screen, xRight, yBottom, xRight, yTop, col)
	ebitenutil.DrawLine(screen, xRight, yTop, xLeft, yTop, col)
}

func (g *graph) checkGraphMatrEquality() bool {

	for i := 0; i < len(g.edges); i++ {
		for j := 0; j < len(g.edges[i]); j++ {
			if i >= len(g.adjMatr) || j > len(g.adjMatr[i]) || g.adjMatr[i][j] != g.edges[i][j] {
				return false
			}
		}
	}

	return true
}

func (g *graph) checkGraphListEquality() bool {

	for i := 0; i < len(g.edges); i++ {
		for j := 0; j < len(g.edges[i]); j++ {
			mustFound := g.edges[i][j] > 0
			found := false
			for jID := 0; jID < len(g.successorsList[i]); jID++ {
				if g.successorsList[i][jID] == j {
					found = true
					break
				}
			}
			if (!found && mustFound) || (found && !mustFound) {
				return false
			}
		}
	}

	return true
}

func (g *graph) checkListMatrEquality() bool {

	for i := 0; i < len(g.adjMatr); i++ {
		for j := 0; j < len(g.adjMatr[i]); j++ {
			mustFound := g.adjMatr[i][j] > 0
			found := false
			for jID := 0; jID < len(g.successorsList[i]); jID++ {
				if g.successorsList[i][jID] == j {
					found = true
					break
				}
			}
			if (!found && mustFound) || (found && !mustFound) {
				return false
			}
		}
	}

	return true
}

func (g *graph) clearMatr() {
	for i := 0; i < len(g.adjMatr); i++ {
		for j := 0; j < len(g.adjMatr[i]); j++ {
			g.adjMatr[i][j] = 0
		}
	}
}

func (g *graph) clearGraph() {
	for i := 0; i < len(g.edges); i++ {
		for j := 0; j < len(g.edges[i]); j++ {
			g.edges[i][j] = 0
		}
	}
}

func (g *graph) clearList() {
	for i := 0; i < len(g.successorsList); i++ {
		g.successorsList[i] = nil
	}
}
