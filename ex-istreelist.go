package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func initIsTreeList(correction bool, graphCode, answerCode int) (e exo, gCode int) {

	e.BasicSetup()
	e.id = isTreeList
	e.successRequired = 5

	numNodes := 6

	elementSpacing := 100

	// title setup
	e.titleImage = titleIsTreeListImage
	xTitleShift, yTitleShift := e.titleImage.Size()
	e.titleXPosition = (windowWidth - xTitleShift) / 2
	e.titleYPosition = elementSpacing
	e.successCounterY = elementSpacing + yTitleShift
	yTitleShift += spriteSide

	// graph setup
	e.modifiableGraph = false
	e.modifiableAdjMatr = false
	e.displayGraph = false
	e.displayAdjMatr = false
	e.displayList = true

	isTree := true
	if correction {
		e.g.decode(graphCode, numNodes)
		gCode = graphCode
		// todo savoir si c'est un arbre ou pas (probablement grace à la question)
	} else {
		e.g.genTree(numNodes)
		if rand.Intn(2) == 0 {
			e.g.demakeTree()
			isTree = false
		}
		gCode = e.g.encode()
	}
	e.g.linkMatrGraph = false

	listSize := 7 * spriteSide
	e.g.xlistposition = (windowWidth - listSize) / 2
	e.g.ylistposition = 2*elementSpacing + yTitleShift

	// question setup
	xshift, yshift := isTreeListImage.Size()
	e.drawQuestion = func(screen *ebiten.Image) {
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(listSize+3*elementSpacing+yTitleShift))
		screen.DrawImage(
			isTreeListImage,
			&options,
		)
	}

	// answers setup
	e.hasAnswerSheet = true
	answerSize, _ := ouiImage.Size()
	e.answers.init((windowWidth-3*answerSize)/2, listSize+yshift+3*elementSpacing+20+yTitleShift)
	e.answers.addButton(0, 0, ouiImage)
	e.answers.addButton(2*answerSize, 0, nonImage)

	answer := 1
	if isTree {
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
	return e, gCode
}
