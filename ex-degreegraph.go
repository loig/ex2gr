package main

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initDegreeGraph(correction bool, graphCode, questionCode, answerCode int) (e exo, gCode, qCode int) {

	e.BasicSetup()
	e.id = degreeGraph
	e.successRequired = 5

	numNodes := 4

	elementSpacing := 100

	// title setup
	e.titleImage = titleDegreeGraphImage
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
		e.g.decode(graphCode, numNodes)
		gCode = graphCode
	} else {
		e.g.genConnectedGraph(numNodes, 4, 12, -1, -1)
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
	var from int
	if correction {
		from, _ = decodeFromToQuestion(questionCode, numNodes)
		qCode = questionCode
	} else {
		from = rand.Intn(numNodes)
		qCode = encodeFromToQuestion(from, from, numNodes)
	}
	xshift, yshift := degreeGraphImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(e.g.ysize+4*elementSpacing+yTitleShift))
		screen.DrawImage(
			degreeGraphImage,
			&options,
		)
		// from label
		options.GeoM.Translate(float64(8*spriteSide+5), 8)
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
	}

	// answers setup
	e.hasAnswerSheet = true
	answerSize := menuSpriteSide / 2
	e.answers.init((windowWidth-15*answerSize)/2, e.g.ysize+yshift+4*elementSpacing+20+yTitleShift)
	for i := 0; i < 10; i++ {
		x := (i % 6) * (menuSpriteSide / 2)
		y := (i / 6) * (menuSpriteSide / 2)
		e.answers.addButton(i*3*answerSize/2+answerSize/4, 0, menuElementsImage.SubImage(image.Rect(x, y+menuSpriteSide, x+menuSpriteSide/2, y+menuSpriteSide+menuSpriteSide/2)).(*ebiten.Image))
	}

	answer := len(e.g.successorsList[from])
	for i := 0; i < len(e.g.edges); i++ {
		if i != from {
			if e.g.edges[i][from] > 0 {
				answer++
			}
		}
	}

	if correction {
		if answerCode < len(e.answers.clics) && answerCode >= 0 {
			e.answers.clics[answerCode] = 1
		}
	}

	e.checkResult = func() (bool, bool) {
		return e.answers.clics[answer] >= 1, *(e.answers.numClics) >= 1
	}

	e.codeAnswer = func() int {
		for answerID, clics := range e.answers.clics {
			if clics > 0 {
				return answerID
			}
		}
		return 0
	}

	// return exercise
	return e, gCode, qCode
}
