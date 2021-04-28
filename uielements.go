package main

import "github.com/hajimehoshi/ebiten/v2"

func (g *game) drawSuccessCounter(screen *ebiten.Image) {

	requested := g.e.successRequired
	achieved := g.succesfulStrike

	x := (windowWidth - g.e.successRequired*spriteSide) / 2

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(x), float64(g.e.successCounterY))

	for i := 0; i < requested; i++ {
		subImage := undoneQuestionSubimage
		if i < achieved {
			subImage = doneQuestionSubimage
		}
		screen.DrawImage(
			graphElementsImage.SubImage(subImage).(*ebiten.Image),
			&options,
		)
		options.GeoM.Translate(float64(spriteSide), 0)
	}

}
