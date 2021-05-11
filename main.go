package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 1024
	windowHeight = 1448
)

func main() {

	rand.Seed(time.Now().UnixNano())

	restartPoint := flag.String("seed", "", "Utiliser pour reprendre sur un exercice donné.")
	flag.Parse()
	saveFile := flag.String("sauvegarde", "ex2grSave.json", "Utiliser pour changer le fichier de sauvegarde")

	ebiten.SetWindowTitle("ex2gr : exercices de graphes")
	ebiten.SetWindowResizable(true)
	ebiten.MaximizeWindow()

	loadAssets()

	g := game{}
	g.init(*restartPoint, *saveFile)

	err := ebiten.RunGame(&g)
	log.Print(err)

}
