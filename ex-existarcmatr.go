package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initExistArcMatr() (e exo) {

	e.BasicSetup()
	e.id = existArcMatr
	e.successRequired = 5

	elementSpacing := 100

	// title setup
	e.titleImage = title2Image
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.displayGraph = false
	e.displayAdjMatr = true

	e.g.genConnectedGraph(4, 4, 12)
	e.g.linkMatrGraph = false

	matrixSize := 5 * matrixCellSize
	e.g.xmatrposition = (windowWidth - matrixSize) / 2
	e.g.ymatrposition = 2*elementSpacing + yTitleShift

	// question setup
	from := rand.Intn(4)
	to := rand.Intn(3)
	if to == from {
		to = 3
	}
	xshift, yshift := ex2Image.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(matrixSize+3*elementSpacing+yTitleShift))
		screen.DrawImage(
			ex2Image,
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
	e.answers.init((windowWidth-3*answerSize)/2, matrixSize+yshift+4*elementSpacing+yTitleShift)
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
