package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	windowWidth  = 1024
	windowHeight = 1448
)

type game struct {
	e                  exo
	readyForNext       bool
	onClicNextButton   int
	selectedNextButton int
	goToNext           answerSheet
	exoDone            bool
}

type exo struct {
	g                 graph
	modifiableGraph   bool
	modifiableAdjMatr bool
	displayGraph      bool
	displayAdjMatr    bool
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
}

func (e *exo) drawTitle(screen *ebiten.Image) {

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(e.titleXPosition), float64(e.titleYPosition))
	screen.DrawImage(
		e.titleImage,
		&options,
	)
}

func (g *game) init() {
	g.goToNext.init(0, windowHeight-200)
	xshift, _ := suivantImage.Size()
	g.goToNext.addButton((windowWidth-xshift)/2, 0, suivantImage)
	g.onClicNextButton = -1
	g.selectedNextButton = -1
	g.e = initEx1()
}

func (g *game) reset() {
	g.onClicNextButton = -1
	g.selectedNextButton = -1
	g.readyForNext = false
	g.exoDone = false
	g.goToNext.resetClics()
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

func (g *game) Update() error {
	x, y := ebiten.CursorPosition()
	if !g.exoDone {
		g.e.update(x, y)
		if g.e.done {
			g.exoDone = true
			if g.e.correct {
				g.goToNext.aboveText = bravoImage
			} else {
				g.goToNext.aboveText = rateImage
			}
		}
	}

	if g.exoDone {
		if g.readyForNext {
			// This is where counting of succes and so should be done
			g.e = initEx1()
			g.reset()
		}
		g.selectedNextButton = g.goToNext.selectButton(x, y)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.onClicNextButton = g.selectedNextButton
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if g.onClicNextButton == g.selectedNextButton {
				g.goToNext.clic(g.selectedNextButton)
			}
		}
		g.readyForNext = g.goToNext.clics[0] > 0
		return nil
	}
	return nil
}

func (e *exo) update(x, y int) {

	if !e.done {
		e.correct, e.done = e.checkResult()
	}

	if e.done {
		return
	}

	if e.hasAnswerSheet {
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
		if e.modifiableGraph {

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

		if e.modifiableAdjMatr {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if e.selectedCellI >= 0 && e.selectedCellJ >= 0 {
					e.g.updateMatrixCell(e.selectedCellI, e.selectedCellJ)
				}
			}
		}
	}

}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	g.e.draw(screen)

	// next exo
	if g.exoDone {
		g.goToNext.draw(screen, g.selectedNextButton)
	}
}

func (e *exo) draw(screen *ebiten.Image) {

	// title
	e.drawTitle(screen)

	// graph
	if e.modifiableGraph && e.displayGraph {
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

	// question
	e.drawQuestion(screen)

	// answers
	if e.hasAnswerSheet {
		e.answers.draw(screen, e.selectedButton)
	}

}

func (g *game) Layout(w, h int) (int, int) {
	return windowWidth, windowHeight
}

func main() {

	ebiten.SetWindowTitle("ex2gr : exercices de graphes")
	ebiten.SetWindowResizable(true)
	ebiten.MaximizeWindow()

	loadAssets()

	g := game{}
	g.init()

	err := ebiten.RunGame(&g)
	log.Print(err)

}
