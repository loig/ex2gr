package main

import (
	"encoding/json"
	"image"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type MenuInfo struct {
	ExoDone  [globalNumExo]bool
	ExoTried [globalNumExo]int
}

func (m *MenuInfo) init(saveFile string) {

	b, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Panic("Impossible de récupérer les données du fichier ", saveFile)
	}
}

func (m *MenuInfo) save(saveFile string) {
	b, err := json.Marshal(*m)
	if err != nil {
		log.Panic("Impossible de préparer les données pour écriture dans le fichier ", saveFile)
	}

	err = ioutil.WriteFile(saveFile, b, 0644)
	if err != nil {
		log.Panic("Impossible d'écrire les données dans le fichier ", saveFile)
	}
}

const (
	menuColumns    = 6
	menuSpace      = 10
	menuWidth      = menuColumns*menuSpriteSide + (menuColumns-1)*menuSpace
	menuHeight     = (globalNumExo/menuColumns+1)*(menuSpriteSide+menuSpace) - menuSpace
	menuXPos       = (windowWidth - menuWidth) / 2
	menuYPos       = 300
	menuTitleSpace = 50
)

func (g *game) drawMenu(screen *ebiten.Image, selectedEx int) {

	screen.DrawImage(generalTitleImage, &ebiten.DrawImageOptions{})

	for i := 0; i < globalNumExo; i++ {
		sprite := menuExSubimage
		if g.menu.ExoDone[i] {
			sprite = menuExDoneSubimage
		}
		if i == selectedEx {
			if g.menu.ExoDone[i] {
				sprite = menuExDoneSelectedSubimage
			} else {
				sprite = menuExSelectedSubimage
			}
		}
		options := ebiten.DrawImageOptions{}
		xshift := menuXPos + (i%menuColumns)*(menuSpriteSide+menuSpace)
		yshift := menuYPos + (i/menuColumns)*(menuSpriteSide+menuSpace)
		options.GeoM.Translate(float64(xshift), float64(yshift))
		screen.DrawImage(
			menuElementsImage.SubImage(sprite).(*ebiten.Image),
			&options,
		)
		if g.menu.ExoDone[i] {
			screen.DrawImage(
				menuElementsImage.SubImage(menuDoneMarkSubimage).(*ebiten.Image),
				&options,
			)
		}
	}
	g.menu.drawSelectedEx(screen, selectedEx)
}

func (m *MenuInfo) drawSelectedEx(screen *ebiten.Image, sel int) {

	if sel >= 0 && sel < globalNumExo {
		var xsize, ysize int

		titleImage := getExTitlePerNum(sel)
		if titleImage != nil {
			xsize, ysize = titleImage.Size()
			options := ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64((windowWidth-xsize)/2), float64(menuYPos+menuHeight+menuTitleSpace))
			screen.DrawImage(
				titleImage,
				&options,
			)
		}

		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xsize)/2+menuTitleSpace), float64(menuYPos+menuHeight+ysize+menuTitleSpace))
		screen.DrawImage(
			menuElementsImage.SubImage(menuTriesText).(*ebiten.Image),
			&options,
		)

		tries := m.ExoTried[sel]
		digits := make([]int, 0)
		if tries == 0 {
			digits = append(digits, 0)
		}
		for tries > 0 {
			digits = append(digits, tries%10)
			tries = tries / 10
		}
		options.GeoM.Translate(float64(menuSpriteSide)+10, 0)
		for i := len(digits) - 1; i >= 0; i-- {
			digitSubimage := image.Rect((digits[i]%6)*spriteSide, menuSpriteSide+(digits[i]/6)*spriteSide, (digits[i]%6+1)*spriteSide, menuSpriteSide+(digits[i]/6+1)*spriteSide)
			screen.DrawImage(
				menuElementsImage.SubImage(digitSubimage).(*ebiten.Image),
				&options,
			)
			options.GeoM.Translate(float64(spriteSide/3), 0)
		}

		options = ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xsize)/2+menuTitleSpace), float64(menuYPos+menuHeight+ysize+menuTitleSpace+menuSpriteSide/2))
		screen.DrawImage(
			menuElementsImage.SubImage(menuDoneText).(*ebiten.Image),
			&options,
		)
		options.GeoM.Translate(float64(menuSpriteSide), 0)
		sprite := nonImage
		if m.ExoDone[sel] {
			sprite = ouiImage
		}
		screen.DrawImage(
			sprite,
			&options,
		)
	}
}

func (g *game) checkAboveMenuEx(x, y int) int {
	if x > menuXPos && x < menuWidth+menuXPos {
		if y > menuYPos && y < menuHeight+menuYPos {
			return (x-menuXPos)/(menuSpriteSide+menuSpace) + ((y-menuYPos)/(menuSpriteSide+menuSpace))*menuColumns
		}
	}
	return -1
}

func (g *game) getNextUndoneID(lastID int) int {
	resID := lastID
	for i := 0; i < globalNumExo; i++ {
		resID = (resID + 1) % globalNumExo
		if !g.menu.ExoDone[resID] {
			return resID
		}
	}
	return -1
}
