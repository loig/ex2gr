package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initExistPathMatr() (e exo) {

	e.BasicSetup()
	e.id = existPathMatr
	e.successRequired = 5

	elementSpacing := 100

	// title setup
	e.titleImage = titleExistPathMatrImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift

	// graph setup part 1
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.displayGraph = false
	e.displayAdjMatr = true

	// question setup part 1
	from := rand.Intn(5)
	to := rand.Intn(4)
	if to == from {
		to = 4
	}

	// graph setup part 2
	e.g.genConnectedGraph(5, 7, 10, from, to)
	e.g.linkMatrGraph = false

	matrixSize := 7 * matrixCellSize
	e.g.xmatrposition = (windowWidth - matrixSize) / 2
	e.g.ymatrposition = 2*elementSpacing + yTitleShift

	// question setup part 2
	xshift, yshift := existPathMatrImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(matrixSize+3*elementSpacing+yTitleShift))
		screen.DrawImage(
			existPathMatrImage,
			&options,
		)
		// from label
		options.GeoM.Translate(float64(8*spriteSide), 0)
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
		options.GeoM.Translate(float64(spriteSide+15), 0)
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
	answerSize, _ := ouiImage.Size()
	e.answers.init((windowWidth-3*answerSize)/2, matrixSize+yshift+3*elementSpacing+20+yTitleShift)
	e.answers.addButton(0, 0, ouiImage)
	e.answers.addButton(2*answerSize, 0, nonImage)

	answer := 1
	if e.g.existPath(from, to) {
		answer = 0
	}

	e.checkResult = func() (bool, bool) {
		return e.answers.clics[answer] >= 1, e.answers.clics[0]+e.answers.clics[1] >= 1
	}

	// return exercise
	return e
}
