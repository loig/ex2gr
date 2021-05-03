package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 1024
	windowHeight = 1448
)

// ajouter flag pour lire la seed

func main() {

	restartPoint := flag.String("seed", "", "Utiliser pour reprendre sur un exercice donné.")
	flag.Parse()

	ebiten.SetWindowTitle("ex2gr : exercices de graphes")
	ebiten.SetWindowResizable(true)
	ebiten.MaximizeWindow()

	loadAssets()

	g := game{}
	g.init(*restartPoint)

	err := ebiten.RunGame(&g)
	log.Print(err)

}
