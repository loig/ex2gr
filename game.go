package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	e                  exo
	readyForNext       bool
	onClicNextButton   int
	selectedNextButton int
	goToNext           answerSheet
	exoDone            bool
	succesfulStrike    int
	correctionMode     bool
	exState            stateDescription
}

func (g *game) init(code string) {
	g.goToNext.init(0, windowHeight-200)
	xshift, _ := suivantImage.Size()
	g.goToNext.addButton((windowWidth-xshift)/2, 0, suivantImage)
	g.onClicNextButton = -1
	g.selectedNextButton = -1
	if code != "" {
		g.correctionMode = true
		g.exState.decode(code)
		g.initExo(g.exState.numExo)
	} else {
		g.initExo(listToGraph)
	}
}

func (g *game) reset() {
	g.onClicNextButton = -1
	g.selectedNextButton = -1
	g.readyForNext = false
	g.exoDone = false
	g.goToNext.resetClics()
}

func (g *game) Update() error {
	x, y := ebiten.CursorPosition()

	if g.correctionMode {
		g.e.update(x, y, true)
		return nil
	}

	if !g.exoDone {
		g.e.update(x, y, false)
		if g.e.done {
			g.exoDone = true
			g.exState.answer = g.e.codeAnswer()
			log.Print(g.exState.encode())
			if g.e.correct {
				g.succesfulStrike++
				g.goToNext.aboveText = bravoImage
			} else {
				g.succesfulStrike = 0
				g.goToNext.aboveText = rateImage
			}
		}
	}

	if g.exoDone {
		if g.readyForNext {
			// This is where counting of succes and so should be done
			currentID := g.e.id
			if g.succesfulStrike >= g.e.successRequired {
				currentID++
				g.succesfulStrike = 0
			}
			g.initExo(currentID)
			g.reset()
		}
		g.selectedNextButton = g.goToNext.selectButton(x, y)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.onClicNextButton = g.selectedNextButton
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if g.onClicNextButton == g.selectedNextButton {
				g.goToNext.clic(g.selectedNextButton)
			}
		}
		g.readyForNext = g.goToNext.clics[0] > 0
		return nil
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	g.e.draw(screen, g.correctionMode)

	if !g.correctionMode {
		g.drawSuccessCounter(screen)

		// next exo
		if g.exoDone {
			g.goToNext.draw(screen, g.selectedNextButton)
		}
	}
}

func (g *game) Layout(w, h int) (int, int) {
	return windowWidth, windowHeight
}
