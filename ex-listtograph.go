package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func initListToGraph(correction bool, graphCode, answerCode int) (e exo, gCode int) {

	e.BasicSetup()
	e.id = listToGraph
	e.successRequired = 5

	numNodes := 4

	elementSpacing := 100

	// title setup
	e.titleImage = titleListToGraphImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup
	e.modifiableGraph = true
	e.modifiableAdjMatr = false
	e.displayGraph = true
	e.displayAdjMatr = false
	e.displayList = true

	if correction {
		e.g.decode(graphCode, numNodes)
		gCode = graphCode
	} else {
		e.g.genConnectedGraph(numNodes, 4, 8, -1, -1)
		gCode = e.g.encode()
	}
	e.g.linkMatrGraph = false
	e.g.clearGraph()
	if correction {
		e.g.decodeGraph(answerCode, numNodes)
	}

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
	listSize := (numNodes + 1) * spriteSide
	e.g.xlistposition = (windowWidth/2 - listSize) / 2
	e.g.ylistposition = 2*elementSpacing + yTitleShift

	// question setup
	xshift, yshift := listToGraphImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(e.g.ysize+4*elementSpacing+yTitleShift))
		screen.DrawImage(
			listToGraphImage,
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

	e.codeAnswer = func() int {
		return e.g.encode()
	}

	// return exercise
	return e, gCode
}
