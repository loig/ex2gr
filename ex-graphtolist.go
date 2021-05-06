package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func initGraphToList(correction bool, graphCode, answerCode int) (e exo, gCode int) {

	e.BasicSetup()
	e.id = graphToList
	e.successRequired = 5

	numNodes := 4

	elementSpacing := 100

	// title setup
	e.titleImage = titleGraphToListImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.modifiableList = true
	e.displayGraph = true
	e.displayAdjMatr = false
	e.displayList = true

	if correction {
		e.g.decode(graphCode, numNodes)
		gCode = graphCode
	} else {
		e.g.genConnectedGraph(numNodes, 4, 12, -1, -1)
		gCode = e.g.encode()
	}
	e.g.linkMatrGraph = false
	e.g.clearList()
	if correction {
		e.g.decodeList(answerCode, numNodes)
	}

	nodeSpacing := 300
	e.g.xposition = (windowWidth/2 - nodeSpacing) / 2
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

	// list setup
	listSize := (numNodes + 2) * spriteSide
	e.g.xlistposition = windowWidth/2 + (windowWidth/2-listSize)/2
	e.g.ylistposition = 2*elementSpacing + yTitleShift

	// question setup
	xshift, yshift := graphToListImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(e.g.ysize+4*elementSpacing+yTitleShift))
		screen.DrawImage(
			graphToListImage,
			&options,
		)
	}

	// answers setup
	e.hasAnswerSheet = true
	answerSize, _ := finiImage.Size()
	e.answers.init((windowWidth-answerSize)/2, e.g.ysize+yshift+4*elementSpacing+20+yTitleShift)
	e.answers.addButton(0, 0, finiImage)

	e.checkResult = func() (bool, bool) {
		return e.g.checkGraphListEquality(), e.answers.clics[0] >= 1
	}

	e.codeAnswer = func() int {
		return e.g.encodeList()
	}

	// return exercise
	return e, gCode
}
