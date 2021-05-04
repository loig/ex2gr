package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type exo struct {
	g                 graph
	modifiableGraph   bool
	modifiableAdjMatr bool
	displayGraph      bool
	displayAdjMatr    bool
	displayList       bool
	selectedNode      int
	nodeFrom          int
	nodeAbove         int
	edgeAboveI        int
	edgeAboveJ        int
	selectedCellI     int
	selectedCellJ     int
	hasAnswerSheet    bool
	answers           answerSheet
	onClicButton      int
	selectedButton    int
	checkResult       func() (correct bool, finished bool)
	drawQuestion      func(*ebiten.Image)
	done              bool
	correct           bool
	titleXPosition    int
	titleYPosition    int
	titleImage        *ebiten.Image
	id                int
	successRequired   int
	successCounterY   int
	codeAnswer        func() int
}

func (e *exo) BasicSetup() {
	e.selectedNode = -1
	e.nodeFrom = -1
	e.nodeAbove = -1
	e.edgeAboveI = -1
	e.edgeAboveJ = -1
	e.selectedCellI = -1
	e.selectedCellJ = -1
	e.selectedButton = -1
}

func (e *exo) update(x, y int, correction bool) {

	if !correction && !e.done {
		e.correct, e.done = e.checkResult()
	}

	if !correction && e.done {
		return
	}

	if !correction && e.hasAnswerSheet {
		e.selectedButton = e.answers.selectButton(x, y)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.onClicButton = e.selectedButton
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if e.onClicButton == e.selectedButton {
				e.answers.clic(e.selectedButton)
			}
		}
	}

	if e.displayGraph {
		// graph
		nodeID := e.g.selectNode(x, y)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			if e.nodeFrom < 0 {
				e.selectedNode = nodeID
			}
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			e.selectedNode = -1
		}

		if e.selectedNode >= 0 {
			e.g.updateNodePosition(x, y, e.selectedNode)
		}

		if e.selectedNode >= 0 {
			e.nodeAbove = e.selectedNode
		} else {
			e.nodeAbove = nodeID
		}

		edgeI, edgeJ := e.g.selectEdge(x, y)
		if !correction && e.modifiableGraph {

			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if nodeID >= 0 {
					if e.selectedNode < 0 {
						e.nodeFrom = nodeID
					}
				} else if edgeI >= 0 {
					e.g.removeEdge(edgeI, edgeJ)
				}
			}

			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				if e.nodeFrom >= 0 {
					if nodeID >= 0 {
						e.g.addEdge(e.nodeFrom, nodeID)
					}
				}
				e.nodeFrom = -1
			}

		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) &&
			!inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if nodeID < 0 {
				if edgeI >= 0 && edgeI == edgeJ {
					e.g.moveLoop(edgeI)
				}
			}
		}

		if e.selectedNode >= 0 || e.nodeFrom >= 0 || nodeID >= 0 {
			e.edgeAboveI = -1
			e.edgeAboveJ = -1
		} else {
			e.edgeAboveI = edgeI
			e.edgeAboveJ = edgeJ
		}

	}

	if e.displayAdjMatr {
		// matrix
		e.selectedCellI, e.selectedCellJ = e.g.selectMatrixCell(x, y)

		if !correction && e.modifiableAdjMatr {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if e.selectedCellI >= 0 && e.selectedCellJ >= 0 {
					e.g.updateMatrixCell(e.selectedCellI, e.selectedCellJ)
				}
			}
		}
	}

}

func (e *exo) draw(screen *ebiten.Image, correction bool) {

	// title
	e.drawTitle(screen)

	// graph
	if !correction && e.modifiableGraph && e.displayGraph {
		if e.nodeFrom >= 0 {
			x, y := ebiten.CursorPosition()
			e.g.drawEdge(e.nodeFrom, x-e.g.xposition, y-e.g.yposition, screen, true, e.nodeAbove == e.nodeFrom, e.nodeAbove)
		}
	}

	if e.displayAdjMatr {
		e.g.drawMatrix(screen, []int{e.selectedCellI, e.selectedCellJ})
	}

	if e.displayGraph {
		e.g.drawGraph(screen, []int{e.nodeAbove, e.nodeFrom}, []int{e.edgeAboveI, e.edgeAboveJ})
	}

	if e.displayList {
		e.g.drawList(screen)
	}

	// question
	e.drawQuestion(screen)

	// answers
	if e.hasAnswerSheet {
		e.answers.draw(screen, e.selectedButton)
	}

}

func (e *exo) drawTitle(screen *ebiten.Image) {

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(e.titleXPosition), float64(e.titleYPosition))
	screen.DrawImage(
		e.titleImage,
		&options,
	)
}
