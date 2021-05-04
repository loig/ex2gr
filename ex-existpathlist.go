package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initExistPathList(correction bool, graphCode, questionCode, answerCode int) (e exo, gCode, qCode int) {

	e.BasicSetup()
	e.id = existPathList
	e.successRequired = 5

	numNodes := 5

	elementSpacing := 100

	// title setup
	e.titleImage = titleExistPathListImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift

	// graph setup part 1
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.displayGraph = false
	e.displayAdjMatr = false
	e.displayList = true

	// question setup part 1
	var from, to int
	if correction {
		from, to = decodeFromToQuestion(questionCode, numNodes)
		qCode = questionCode
	} else {
		from = rand.Intn(numNodes)
		to = rand.Intn(numNodes - 1)
		if to == from {
			to = numNodes - 1
		}
		qCode = encodeFromToQuestion(from, to, numNodes)
	}

	// graph setup part 2
	if correction {
		e.g.decode(graphCode, numNodes)
		gCode = graphCode
	} else {
		e.g.genConnectedGraph(numNodes, 7, 10, from, to)
		gCode = e.g.encode()
	}
	e.g.linkMatrGraph = false

	listSize := 7 * spriteSide
	e.g.xlistposition = (windowWidth - listSize) / 2
	e.g.ylistposition = 2*elementSpacing + yTitleShift

	// question setup part 2
	xshift, yshift := existPathListImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(listSize+2*elementSpacing+yTitleShift))
		screen.DrawImage(
			existPathListImage,
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
	e.answers.init((windowWidth-3*answerSize)/2, listSize+yshift+2*elementSpacing+40+yTitleShift)
	e.answers.addButton(0, 0, ouiImage)
	e.answers.addButton(2*answerSize, 0, nonImage)

	answer := 1
	if e.g.existPath(from, to) {
		answer = 0
	}

	if correction {
		if answerCode < len(e.answers.clics) && answerCode >= 0 {
			e.answers.clics[answerCode] = 1
		}
	}

	e.checkResult = func() (bool, bool) {
		return e.answers.clics[answer] >= 1, e.answers.clics[0]+e.answers.clics[1] >= 1
	}

	e.codeAnswer = func() int {
		if e.answers.clics[0] >= 1 {
			return 0
		}
		return 1
	}

	// return exercise
	return e, gCode, qCode
}
