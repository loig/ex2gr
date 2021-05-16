package main

import (
	"fmt"
	"image"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
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

func (g *graph) encodeList() int {
	code := 0
	for i := 0; i < len(g.successorsList); i++ {
		for j := 0; j < len(g.successorsList); j++ {
			found := false
			for _, jj := range g.successorsList[i] {
				if jj == j {
					found = true
					break
				}
			}
			if found {
				code = 2*code + 1
			} else {
				code = 2 * code
			}
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

func (g *graph) decodeList(code int, numNodes int) {
	for i := len(g.successorsList) - 1; i >= 0; i-- {
		for j := len(g.successorsList) - 1; j >= 0; j-- {
			if code%2 == 1 {
				if g.successorsList[i] == nil {
					g.successorsList[i] = make([]int, 0, len(g.successorsList))
				}
				g.successorsList[i] = append(g.successorsList[i], j)
			}
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

func (g *game) drawSeed(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(5, 0)
	for _, c := range g.encodedExState {
		if c <= 'f' && c >= 'a' {
			options.GeoM.Translate(0, float64(-spriteSide/8))
			x := int(c - 97)
			screen.DrawImage(graphElementsImage.SubImage(image.Rect(x*spriteSide, spriteSide, (x+1)*spriteSide, 2*spriteSide)).(*ebiten.Image), options)
			options.GeoM.Translate(0, float64(+spriteSide/8))
		} else if c <= '9' && c >= '0' {
			options.GeoM.Translate(float64(spriteSide/4), float64(-spriteSide/4))
			v := int(c - 48)
			x := v % 6
			y := v / 6
			screen.DrawImage(menuElementsImage.SubImage(image.Rect(x*spriteSide, y*spriteSide+menuSpriteSide, (x+1)*spriteSide, (y+1)*spriteSide+menuSpriteSide)).(*ebiten.Image), options)
			options.GeoM.Translate(-float64(spriteSide/4), float64(spriteSide/4))
		} else if c == '.' {
			screen.DrawImage(graphElementsImage.SubImage(undoneQuestionSubimage).(*ebiten.Image), options)
		}
		options.GeoM.Translate(float64(spriteSide/2-10), 0)
	}
}
