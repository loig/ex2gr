package main

import "github.com/hajimehoshi/ebiten/v2"

type answerSheet struct {
	xposition int
	yposition int
	aboveText *ebiten.Image
	buttons   []button
	clics     []int
}

type button struct {
	xposition int
	yposition int
	xsize     int
	ysize     int
	content   *ebiten.Image
}

func (a *answerSheet) init(x, y int) {
	a.xposition = x
	a.yposition = y
	a.buttons = make([]button, 0)
	a.clics = make([]int, 0)
}

func (a *answerSheet) addButton(x, y int, content *ebiten.Image) {
	xsize, ysize := content.Size()
	a.buttons = append(a.buttons, button{
		xposition: x,
		yposition: y,
		xsize:     xsize,
		ysize:     ysize,
		content:   content,
	})
	a.clics = append(a.clics, 0)
}

func (a *answerSheet) selectButton(x, y int) int {
	x = x - a.xposition
	y = y - a.yposition
	if a.aboveText != nil {
		_, yshift := a.aboveText.Size()
		y -= yshift
	}
	for i, b := range a.buttons {
		if b.isOver(x, y) {
			return i
		}
	}
	return -1
}

func (b *button) isOver(x, y int) bool {
	return b.xposition <= x && b.xposition+b.xsize >= x &&
		b.yposition <= y && b.yposition+b.ysize >= y
}

func (a *answerSheet) clic(buttonID int) {
	if buttonID >= 0 && buttonID < len(a.buttons) {
		a.clics[buttonID]++
	}
}

func (a *answerSheet) resetClics() {
	for i := range a.clics {
		a.clics[i] = 0
	}
}

func (a *answerSheet) draw(screen *ebiten.Image, buttonSelected int) {
	var xshift int
	yshift := 0
	if a.aboveText != nil {
		xshift, yshift = a.aboveText.Size()
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64((windowWidth-xshift)/2), float64(a.yposition-spriteSide/2))
		screen.DrawImage(a.aboveText, &options)
	}
	for i, b := range a.buttons {
		b.draw(screen, a.xposition, a.yposition+yshift, (i == buttonSelected) || (a.clics[i] > 0))
	}
}

func (b *button) draw(screen *ebiten.Image, xshift, yshift int, selected bool) {
	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(xshift+b.xposition), float64(yshift+b.yposition))
	screen.DrawImage(
		b.content,
		&options,
	)

	// left
	subimage := buttonLeftSubimage
	if selected {
		subimage = buttonLeftSelectedSubimage
	}
	screen.DrawImage(
		graphElementsImage.SubImage(subimage).(*ebiten.Image),
		&options,
	)

	// middle
	length := spriteSide
	subimage = buttonCenterSubimage
	if selected {
		subimage = buttonCenterSelectedSubimage
	}
	for length+spriteSide < b.xsize {
		options.GeoM.Translate(float64(spriteSide), 0)
		screen.DrawImage(
			graphElementsImage.SubImage(subimage).(*ebiten.Image),
			&options,
		)
		length += spriteSide
	}

	// right
	subimage = buttonRightSubimage
	if selected {
		subimage = buttonRightSelectedSubimage
	}
	options.GeoM.Translate(float64(b.xsize-length), 0)
	screen.DrawImage(
		graphElementsImage.SubImage(subimage).(*ebiten.Image),
		&options,
	)

}
