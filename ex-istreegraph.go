package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initIsTreeGraph(correction bool, graphCode int) (e exo, gCode int) {

	e.BasicSetup()
	e.id = isTreeGraph
	e.successRequired = 5

	elementSpacing := 100

	// title setup
	e.titleImage = titleIsTreeGraphImage
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

	isTree := true
	if correction {
		e.g.decode(graphCode, 6)
		gCode = graphCode
		// todo savoir si c'est un arbre (probablement via le transfert de la question)
	} else {
		e.g.genTree(6)
		if rand.Intn(2) == 0 {
			e.g.demakeTree()
			isTree = false
		}
		gCode = e.g.encode()
	}
	e.g.linkMatrGraph = false

	e.g.xposition = (windowWidth - 400) / 2
	e.g.yposition = 2*elementSpacing + yTitleShift
	e.g.xsize = 400
	e.g.ysize = 400

	e.g.nodes[0].yposition = 100
	e.g.nodes[1].xposition = 200
	e.g.nodes[2].xposition = 400
	e.g.nodes[2].yposition = 100
	e.g.nodes[3].yposition = 300
	e.g.nodes[4].xposition = 200
	e.g.nodes[4].yposition = 400
	e.g.nodes[5].xposition = 400
	e.g.nodes[5].yposition = 300

	// question setup
	xshift, yshift := isTreeGraphImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(e.g.ysize+3*elementSpacing+yTitleShift))
		screen.DrawImage(
			isTreeGraphImage,
			&options,
		)
	}

	// answers setup
	e.hasAnswerSheet = true
	answerSize, _ := ouiImage.Size()
	e.answers.init((windowWidth-3*answerSize)/2, e.g.ysize+yshift+3*elementSpacing+20+yTitleShift)
	e.answers.addButton(0, 0, ouiImage)
	e.answers.addButton(2*answerSize, 0, nonImage)

	answer := 1
	if isTree {
		answer = 0
	}

	e.checkResult = func() (bool, bool) {
		return e.answers.clics[answer] >= 1, e.answers.clics[0]+e.answers.clics[1] >= 1
	}

	// return exercise
	return e, gCode
}
