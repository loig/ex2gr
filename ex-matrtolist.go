package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func initMatrToList(correction bool, graphCode, answerCode int) (e exo, gCode int) {

	e.BasicSetup()
	e.id = matrToList
	e.successRequired = 5

	numNodes := 4

	elementSpacing := 100

	// title setup
	e.titleImage = titleMatrToListImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.modifiableList = true
	e.displayGraph = false
	e.displayAdjMatr = true
	e.displayList = true

	if correction {
		e.g.decode(graphCode, numNodes)
		gCode = graphCode
	} else {
		e.g.genConnectedGraph(numNodes, 4, 8, -1, -1)
		gCode = e.g.encode()
	}
	e.g.linkMatrGraph = false
	e.g.clearList()
	if correction {
		e.g.decodeList(answerCode, numNodes)
	}

	// matrix setup
	matrixSize := (numNodes + 1) * matrixCellSize
	e.g.xmatrposition = (windowWidth/2 - matrixSize) / 2
	e.g.ymatrposition = 2*elementSpacing + yTitleShift

	// list setup
	listSize := (numNodes + 2) * spriteSide
	e.g.xlistposition = windowWidth/2 + (windowWidth/2-listSize)/2
	e.g.ylistposition = 2*elementSpacing + yTitleShift

	// question setup
	xshift, yshift := matrToListImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(matrixSize+3*elementSpacing+yTitleShift))
		screen.DrawImage(
			matrToListImage,
			&options,
		)
	}

	// answers setup
	e.hasAnswerSheet = true
	answerSize, _ := finiImage.Size()
	e.answers.init((windowWidth-answerSize)/2, matrixSize+yshift+3*elementSpacing+20+yTitleShift)
	e.answers.addButton(0, 0, finiImage)

	e.checkResult = func() (bool, bool) {
		return e.g.checkListMatrEquality(), e.answers.clics[0] >= 1
	}

	e.codeAnswer = func() int {
		return e.g.encodeList()
	}

	// return exercise
	return e, gCode
}
