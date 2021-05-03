package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initExistPathGraph(correction bool, graphCode, questionCode, answerCode int) (e exo, gCode, qCode int) {

	e.BasicSetup()
	e.id = existPathGraph
	e.successRequired = 5

	elementSpacing := 100

	// title setup
	e.titleImage = titleExistPathGraphImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup part 1
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.displayGraph = true
	e.displayAdjMatr = false

	// question setup part 1
	var from, to int
	if correction {
		from, to = decodeFromToQuestion(questionCode, 4)
		qCode = questionCode
	} else {
		from = rand.Intn(6)
		to = rand.Intn(5)
		if to == from {
			to = 5
		}
		qCode = encodeFromToQuestion(from, to, 4)
	}

	// graph setup part 2
	if correction {
		e.g.decode(graphCode, 6)
		gCode = graphCode
	} else {
		e.g.genConnectedGraph(6, 8, 12, from, to)
		gCode = e.g.encode()
	}
	e.g.linkMatrGraph = false

	nodeSpacing := 300
	e.g.xposition = (windowWidth - 2*nodeSpacing) / 2
	e.g.yposition = 2*elementSpacing + yTitleShift
	e.g.xsize = 2 * nodeSpacing
	e.g.ysize = nodeSpacing

	e.g.nodes[0].loopPosition = loopTopLeft
	e.g.nodes[1].xposition = nodeSpacing
	e.g.nodes[1].yposition = nodeSpacing / 4
	e.g.nodes[1].loopPosition = loopTopLeft
	e.g.nodes[2].xposition = 2 * nodeSpacing
	e.g.nodes[2].loopPosition = loopTopRight
	e.g.nodes[3].yposition = nodeSpacing
	e.g.nodes[3].loopPosition = loopBottomLeft
	e.g.nodes[4].xposition = nodeSpacing
	e.g.nodes[4].yposition = 3 * nodeSpacing / 4
	e.g.nodes[4].loopPosition = loopBottomLeft
	e.g.nodes[5].xposition = 2 * nodeSpacing
	e.g.nodes[5].yposition = nodeSpacing
	e.g.nodes[5].loopPosition = loopBottomRight

	// question setup part 2
	xshift, yshift := existPathGraphImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(e.g.ysize+4*elementSpacing+yTitleShift))
		screen.DrawImage(
			existPathGraphImage,
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
	e.answers.init((windowWidth-3*answerSize)/2, e.g.ysize+yshift+4*elementSpacing+20+yTitleShift)
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
