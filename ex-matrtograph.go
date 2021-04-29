package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initMatrToGraph() (e exo) {

	e.BasicSetup()
	e.id = matrToGraph
	e.successRequired = 5

	elementSpacing := 100

	// title setup
	e.titleImage = titleMatrToGraphImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup
	e.modifiableGraph = true
	e.modifiableAdjMatr = false
	e.displayGraph = true
	e.displayAdjMatr = true

	e.g.genConnectedGraph(4, 4, 12)
	e.g.linkMatrGraph = false
	e.g.clearGraph()

	nodeSpacing := 300
	e.g.xposition = windowWidth/2 + (windowWidth/2-nodeSpacing)/2
	e.g.yposition = 2*elementSpacing + yTitleShift
	e.g.xsize = nodeSpacing
	e.g.ysize = nodeSpacing

	e.g.nodes[1].xposition = nodeSpacing
	e.g.nodes[2].yposition = nodeSpacing
	e.g.nodes[3].xposition = nodeSpacing
	e.g.nodes[3].yposition = nodeSpacing

	// matrix setup
	matrixSize := 5 * matrixCellSize
	e.g.xmatrposition = (windowWidth/2 - matrixSize) / 2
	e.g.ymatrposition = 2*elementSpacing + yTitleShift

	// question setup
	from := rand.Intn(4)
	to := rand.Intn(3)
	if to == from {
		to = 3
	}
	xshift, yshift := matrToGraphImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(e.g.ysize+4*elementSpacing+yTitleShift))
		screen.DrawImage(
			matrToGraphImage,
			&options,
		)
	}

	// answers setup
	e.hasAnswerSheet = true
	answerSize, _ := finiImage.Size()
	e.answers.init((windowWidth-answerSize)/2, e.g.ysize+yshift+4*elementSpacing+20+yTitleShift)
	e.answers.addButton(0, 0, finiImage)

	e.checkResult = func() (bool, bool) {
		return e.g.checkGraphMatrEquality(), e.answers.clics[0] >= 1
	}

	// return exercise
	return e
}
