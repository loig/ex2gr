package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initExistArcGraph(correction bool, graphCode, questionCode int) (e exo, gCode, qCode int) {

	e.BasicSetup()
	e.id = existArcGraph
	e.successRequired = 5

	elementSpacing := 100

	// title setup
	e.titleImage = title1Image
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.displayGraph = true
	e.displayAdjMatr = false

	if correction {
		e.g.decode(graphCode, 4)
		gCode = graphCode
	} else {
		e.g.genConnectedGraph(4, 4, 12, -1, -1)
		gCode = e.g.encode()
	}
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
	var from, to int
	if correction {
		from, to = decodeFromToQuestion(questionCode, 4)
		qCode = questionCode
	} else {
		from = rand.Intn(4)
		to = rand.Intn(3)
		if to == from {
			to = 3
		}
		qCode = encodeFromToQuestion(from, to, 4)
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
	answerSize, _ := ouiImage.Size()
	e.answers.init((windowWidth-3*answerSize)/2, e.g.ysize+yshift+4*elementSpacing+20+yTitleShift)
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
	return e, gCode, qCode
}
