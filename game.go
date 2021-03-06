package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	e                    exo
	readyForNext         bool
	onClicNextButton     int
	selectedNextButton   int
	goToNextEx           answerSheet
	goToNextQ            answerSheet
	exoDone              bool
	succesfulStrike      int
	correctionMode       bool
	exState              stateDescription
	encodedExState       string
	inMenu               bool
	goToMenu             bool
	menu                 MenuInfo
	menuItemSelected     int
	lastMenuItemSelected int
	quitGame             answerSheet
	quitGameSelected     int
	quitGameLastSelected int
	saveFile             string
}

func (g *game) init(code string, saveFile string) {
	g.goToNextEx.init(0, windowHeight-200)
	g.goToNextQ.init(0, windowHeight-200)
	xshift, _ := suivantImage.Size()
	g.goToNextEx.addButton((windowWidth-xshift)/2-200, 0, suivantImage)
	g.goToNextEx.addButton((windowWidth-xshift)/2+200, 0, menuImage)
	xshift, _ = questionImage.Size()
	g.goToNextQ.addButton((windowWidth-xshift)/2, 0, questionImage)
	xshift, _ = quitGameImage.Size()
	g.quitGame.init(0, windowHeight-200)
	g.quitGame.addButton((windowWidth-xshift)/2, 0, quitGameImage)
	g.onClicNextButton = -1
	g.selectedNextButton = -1
	g.menuItemSelected = -1
	g.lastMenuItemSelected = -1
	g.quitGameSelected = -1
	g.quitGameLastSelected = -1
	g.inMenu = true
	if code != "" {
		g.inMenu = false
		g.goToMenu = false
		g.correctionMode = true
		g.exState.decode(code)
		g.initExo(g.exState.numExo)
	}
	g.saveFile = saveFile
	g.menu.init(g.saveFile)
}

func (g *game) reset() {
	g.onClicNextButton = -1
	g.selectedNextButton = -1
	g.readyForNext = false
	g.exoDone = false
	g.menuItemSelected = -1
	g.lastMenuItemSelected = -1
	g.inMenu = false
	g.goToMenu = false
	g.goToNextEx.resetClics()
	g.goToNextQ.resetClics()
	g.quitGame.resetClics()
	g.quitGameSelected = -1
	g.quitGameLastSelected = -1
}

func (g *game) Update() error {
	x, y := ebiten.CursorPosition()

	if g.correctionMode {
		g.e.update(x, y, true)
		return nil
	}

	if g.inMenu {
		g.menuItemSelected = g.checkAboveMenuEx(x, y)
		g.quitGameSelected = g.quitGame.selectButton(x, y)
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if g.menuItemSelected >= 0 && g.menuItemSelected < globalNumExo && g.menuItemSelected == g.lastMenuItemSelected {
				g.menu.ExoTried[g.menuItemSelected]++
				g.initExo(g.menuItemSelected)
				g.reset()
			} else if g.quitGameSelected == g.quitGameLastSelected {
				g.quitGame.clic(g.quitGameSelected)
				if g.quitGame.clics[0] > 0 {
					g.menu.save(g.saveFile)
					return errors.New("Au revoir")
				}
			}
			g.lastMenuItemSelected = -1
			g.menuItemSelected = -1
			return nil
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.lastMenuItemSelected = g.menuItemSelected
			g.quitGameLastSelected = g.quitGameSelected
		}
		return nil
	}

	if !g.exoDone {
		g.e.update(x, y, false)
		if g.e.done {
			g.exoDone = true
			g.exState.answer = g.e.codeAnswer()
			g.encodedExState = g.exState.encode()
			log.Print("Fin de question, le code est??: ", g.encodedExState)
			if g.e.correct {
				g.succesfulStrike++
				g.goToNextQ.aboveText = bravoImage
			} else {
				g.succesfulStrike = 0
				g.menu.ExoTried[g.e.id]++
				g.goToNextQ.aboveText = rateImage
			}
			g.goToNextEx.aboveText = g.goToNextQ.aboveText
		}
		if g.e.quitButton.clics[0] > 0 {
			g.inMenu = true
		}
	}

	if g.exoDone {
		if g.readyForNext {
			// This is where counting of succes and so should be done
			currentID := g.e.id
			if g.succesfulStrike >= g.e.successRequired {
				g.menu.ExoDone[currentID] = true
				if !g.goToMenu {
					currentID = g.getNextUndoneID(currentID)
					if currentID >= 0 {
						g.menu.ExoTried[currentID]++
					} else {
						g.goToMenu = true
					}
				}
				g.succesfulStrike = 0
			}
			if !g.goToMenu {
				g.initExo(currentID)
				g.reset()
			} else {
				g.inMenu = true
			}
			return nil
		}
		usedSheet := &g.goToNextQ
		if g.succesfulStrike >= g.e.successRequired {
			usedSheet = &g.goToNextEx
		}
		g.selectedNextButton = usedSheet.selectButton(x, y)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.onClicNextButton = g.selectedNextButton
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if g.onClicNextButton == g.selectedNextButton {
				usedSheet.clic(g.selectedNextButton)
			}
		}
		g.readyForNext = usedSheet.clics[0] > 0
		if g.succesfulStrike >= g.e.successRequired && len(usedSheet.clics) > 1 {
			g.goToMenu = usedSheet.clics[1] > 0
			g.readyForNext = g.readyForNext || g.goToMenu
		}
		return nil
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	if g.inMenu && !g.correctionMode {
		g.drawMenu(screen, g.menuItemSelected)
		g.quitGame.draw(screen, g.quitGameSelected)
		return
	}

	g.drawSeed(screen)

	g.e.draw(screen, g.correctionMode)

	if !g.correctionMode {
		g.drawSuccessCounter(screen)

		// next exo
		if g.exoDone {
			usedSheet := &g.goToNextQ
			if g.succesfulStrike >= g.e.successRequired {
				usedSheet = &g.goToNextEx
			}
			usedSheet.draw(screen, g.selectedNextButton)
		}
	}
}

func (g *game) Layout(w, h int) (int, int) {
	return windowWidth, windowHeight
}
