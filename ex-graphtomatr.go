package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func initGraphToMatr(correction bool, graphCode, answerCode int) (e exo, gCode int) {

	e.BasicSetup()
	e.id = graphToMatr
	e.successRequired = 5

	numNodes := 4

	elementSpacing := 100

	// title setup
	e.titleImage = titleGraphToMatrImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup
	e.modifiableGraph = false
	e.modifiableAdjMatr = true
	e.displayGraph = true
	e.displayAdjMatr = true

	if correction {
		e.g.decode(graphCode, numNodes)
		gCode = graphCode
	} else {
		e.g.genConnectedGraph(numNodes, 4, 12, -1, -1)
		gCode = e.g.encode()
	}
	e.g.linkMatrGraph = false
	e.g.clearMatr()
	if correction {
		e.g.decodeMatr(answerCode, numNodes)
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

	// matrix setup
	matrixSize := 5 * matrixCellSize
	e.g.xmatrposition = windowWidth/2 + (windowWidth/2-matrixSize)/2
	e.g.ymatrposition = 2*elementSpacing + yTitleShift

	// question setup
	xshift, yshift := graphToMatrImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(e.g.ysize+4*elementSpacing+yTitleShift))
		screen.DrawImage(
			graphToMatrImage,
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
		return e.g.encodeMatr()
	}

	// return exercise
	return e, gCode
}
