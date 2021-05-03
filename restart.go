package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type stateDescription struct {
	numExo              int // exercise number
	descriptionExo      int // exercise binary encoded as an int
	descriptionQuestion int // question binary encoded as an int
	answer              int // answer binary encoded as an int
}

func (s *stateDescription) encode() string {
	return fmt.Sprintf("%x.%x.%x.%x", s.numExo, s.descriptionExo, s.descriptionQuestion, s.answer+1)
}

func (s *stateDescription) decode(code string) {
	parts := strings.Split(code, ".")
	if len(parts) < 4 {
		log.Fatal("Impossible de traiter le code ", code, " (mauvaise structure)")
	}
	numExo, err := strconv.ParseInt(parts[0], 16, 0)
	if err != nil {
		log.Fatal("Impossible de traiter le code ", code, " (mauvais numéro d'exercice)")
	}
	descriptionExo, err := strconv.ParseInt(parts[1], 16, 0)
	if err != nil {
		log.Fatal("Impossible de traiter le code ", code, " (mauvaise description du graphe)")
	}
	question, err := strconv.ParseInt(parts[2], 16, 0)
	if err != nil {
		log.Fatal("Impossible de traiter le code ", code, " (mauvaise description de la question)")
	}
	answer, err := strconv.ParseInt(parts[3], 16, 0)
	if err != nil {
		log.Fatal("Impossible de traiter le code ", code, " (mauvaise description de la réponse)")
	}
	s.numExo = int(numExo)
	s.descriptionExo = int(descriptionExo)
	s.descriptionQuestion = int(question)
	s.answer = int(answer) - 1
}

func (g *graph) encode() int {
	code := 0
	for i := 0; i < len(g.edges); i++ {
		for j := 0; j < len(g.edges[0]); j++ {
			code = 2*code + g.edges[i][j]
		}
	}
	return code
}

func (g *graph) encodeMatr() int {
	code := 0
	for i := 0; i < len(g.adjMatr); i++ {
		for j := 0; j < len(g.adjMatr[0]); j++ {
			code = 2*code + g.adjMatr[i][j]
		}
	}
	return code
}

func (g *graph) decode(code int, numNodes int) {
	g.makeEmptyGraph(numNodes)
	g.linkMatrGraph = true
	for i := len(g.edges) - 1; i >= 0; i-- {
		for j := len(g.edges) - 1; j >= 0; j-- {
			if code%2 != 0 {
				g.addEdge(i, j)
			}
			code = code / 2
		}
	}

	// Prepare draw order
	g.nodesDrawOrder = make([]int, numNodes)
	for i := range g.nodesDrawOrder {
		g.nodesDrawOrder[i] = i
	}
}

func (g *graph) decodeGraph(code int, numNodes int) {
	for i := len(g.edges) - 1; i >= 0; i-- {
		for j := len(g.edges) - 1; j >= 0; j-- {
			g.edges[i][j] = code % 2
			code = code / 2
		}
	}

	// Prepare draw order
	g.nodesDrawOrder = make([]int, numNodes)
	for i := range g.nodesDrawOrder {
		g.nodesDrawOrder[i] = i
	}
}

func (g *graph) decodeMatr(code int, numNodes int) {
	for i := len(g.edges) - 1; i >= 0; i-- {
		for j := len(g.edges) - 1; j >= 0; j-- {
			g.adjMatr[i][j] = code % 2
			code = code / 2
		}
	}
}

func encodeFromToQuestion(from, to, numNodes int) (code int) {
	if numNodes > 16 {
		log.Fatal("impossible de traiter un cas à plus de 16 sommets")
	}
	return to + 16*from
}

func decodeFromToQuestion(code, numNodes int) (from, to int) {
	if numNodes > 16 {
		log.Fatal("impossible de traiter un cas à plus de 16 sommets")
	}
	from = code / 16
	to = code % 16
	return from, to
}
