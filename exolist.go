package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	existArcGraph int = iota
	existArcList
	existArcMatr
	existPathGraph
	existPathList
	existPathMatr
	existCycleGraph
	existCycleList
	existCycleMatr
	listToGraph
	listToMatr
	matrToGraph
	matrToList
	graphToList
	graphToMatr
	isTreeGraph
	isTreeList
	isTreeMatr
	globalNumExo
)

func (g *game) initExo(exNum int) {

	g.exState.numExo = exNum

	if !g.correctionMode {
		g.exState.descriptionExo = 0
		g.exState.descriptionQuestion = 0
		g.exState.answer = -1
	}

	switch exNum {
	case existArcGraph:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistArcGraph(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case existArcMatr:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistArcMatr(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case existPathGraph:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistPathGraph(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case existPathMatr:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistPathMatr(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case graphToMatr:
		g.e, g.exState.descriptionExo = initGraphToMatr(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	case matrToGraph:
		g.e, g.exState.descriptionExo = initMatrToGraph(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	case isTreeGraph:
		g.e, g.exState.descriptionExo = initIsTreeGraph(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	case isTreeMatr:
		g.e, g.exState.descriptionExo = initIsTreeMatr(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	case existCycleGraph:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistCycleGraph(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case existCycleMatr:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistCycleMatr(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case existArcList:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistArcList(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case existPathList:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistPathList(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case existCycleList:
		g.e, g.exState.descriptionExo, g.exState.descriptionQuestion = initExistCycleList(g.correctionMode, g.exState.descriptionExo, g.exState.descriptionQuestion, g.exState.answer)
	case isTreeList:
		g.e, g.exState.descriptionExo = initIsTreeList(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	case listToGraph:
		g.e, g.exState.descriptionExo = initListToGraph(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	case graphToList:
		g.e, g.exState.descriptionExo = initGraphToList(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	case listToMatr:
		g.e, g.exState.descriptionExo = initListToMatr(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	case matrToList:
		g.e, g.exState.descriptionExo = initMatrToList(g.correctionMode, g.exState.descriptionExo, g.exState.answer)
	}

	log.Print(g.exState.encode())

}

func getExTitlePerNum(exNum int) *ebiten.Image {

	switch exNum {
	case existArcGraph:
		return title1Image
	case existArcMatr:
		return title2Image
	case existPathGraph:
		return titleExistPathGraphImage
	case existPathMatr:
		return titleExistPathMatrImage
	case graphToMatr:
		return titleGraphToMatrImage
	case matrToGraph:
		return titleMatrToGraphImage
	case isTreeGraph:
		return titleIsTreeGraphImage
	case isTreeMatr:
		return titleIsTreeMatrImage
	case existCycleGraph:
		return titleExistCycleGraphImage
	case existCycleMatr:
		return titleExistCycleMatrImage
	case existArcList:
		return titleExistArcListImage
	case existPathList:
		return titleExistPathListImage
	case existCycleList:
		return titleExistCycleListImage
	case isTreeList:
		return titleIsTreeListImage
	case listToGraph:
		return titleListToGraphImage
	case graphToList:
		return titleGraphToListImage
	case listToMatr:
		return titleListToMatrImage
	case matrToList:
		return titleMatrToListImage
	}

	return nil

}
