package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type menuInfo struct {
	exoDone  [globalNumExo]bool
	exoTried [globalNumExo]int
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
	for i := 0; i < globalNumExo; i++ {
		sprite := menuExSubimage
		if g.menu.exoDone[i] {
			sprite = menuExDoneSubimage
		}
		if i == selectedEx {
			if g.menu.exoDone[i] {
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
		if g.menu.exoDone[i] {
			screen.DrawImage(
				menuElementsImage.SubImage(menuDoneMarkSubimage).(*ebiten.Image),
				&options,
			)
		}
	}
	g.menu.drawSelectedEx(screen, selectedEx)
}

func (m *menuInfo) drawSelectedEx(screen *ebiten.Image, sel int) {

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

		tries := m.exoTried[sel]
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
		if m.exoDone[sel] {
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
		if !g.menu.exoDone[resID] {
			return resID
		}
	}
	return -1
}
