package main

import "log"

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
	matrToGraph
	graphToList
	graphToMatr
	isTreeGraph
	isTreeList
	isTreeMatr
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
	}

	log.Print(g.exState.encode())

}
