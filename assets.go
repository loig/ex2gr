package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/sprites.png
var graphAssets []byte
var graphElementsImage *ebiten.Image

//go:embed assets/oui.png
var ouiAsset []byte
var ouiImage *ebiten.Image

//go:embed assets/non.png
var nonAsset []byte
var nonImage *ebiten.Image

//go:embed assets/suivant.png
var suivantAsset []byte
var suivantImage *ebiten.Image

//go:embed assets/bravo.png
var bravoAsset []byte
var bravoImage *ebiten.Image

//go:embed assets/rate.png
var rateAsset []byte
var rateImage *ebiten.Image

//go:embed assets/ex1.png
var ex1Asset []byte
var ex1Image *ebiten.Image // should maybe be created only when needed

//go:embed assets/titreex1.png
var title1Asset []byte
var title1Image *ebiten.Image // should maybe be created only when needed

//go:embed assets/ex2.png
var ex2Asset []byte
var ex2Image *ebiten.Image // should maybe be created only when needed

//go:embed assets/titreex2.png
var title2Asset []byte
var title2Image *ebiten.Image // should maybe be created only when needed

//go:embed assets/ex-existpathmatr.png
var existPathMatrAsset []byte
var existPathMatrImage *ebiten.Image // should maybe be created only when needed

//go:embed assets/extitle-existpathmatr.png
var titleExistPathMatrAsset []byte
var titleExistPathMatrImage *ebiten.Image // should maybe be created only when needed

const spriteSide int = 64

func loadAssets() {
	var err error
	graphAssetsDecoded, _, err := image.Decode(bytes.NewReader(graphAssets))
	if err != nil {
		log.Fatal(err)
	}
	graphElementsImage = ebiten.NewImageFromImage(graphAssetsDecoded)

	ouiAssetDecoded, _, err := image.Decode(bytes.NewReader(ouiAsset))
	if err != nil {
		log.Fatal(err)
	}
	ouiImage = ebiten.NewImageFromImage(ouiAssetDecoded)

	nonAssetDecoded, _, err := image.Decode(bytes.NewReader(nonAsset))
	if err != nil {
		log.Fatal(err)
	}
	nonImage = ebiten.NewImageFromImage(nonAssetDecoded)

	suivantAssetDecoded, _, err := image.Decode(bytes.NewReader(suivantAsset))
	if err != nil {
		log.Fatal(err)
	}
	suivantImage = ebiten.NewImageFromImage(suivantAssetDecoded)

	bravoAssetDecoded, _, err := image.Decode(bytes.NewReader(bravoAsset))
	if err != nil {
		log.Fatal(err)
	}
	bravoImage = ebiten.NewImageFromImage(bravoAssetDecoded)

	rateAssetDecoded, _, err := image.Decode(bytes.NewReader(rateAsset))
	if err != nil {
		log.Fatal(err)
	}
	rateImage = ebiten.NewImageFromImage(rateAssetDecoded)

	ex1AssetDecoded, _, err := image.Decode(bytes.NewReader(ex1Asset))
	if err != nil {
		log.Fatal(err)
	}
	ex1Image = ebiten.NewImageFromImage(ex1AssetDecoded)

	title1AssetDecoded, _, err := image.Decode(bytes.NewReader(title1Asset))
	if err != nil {
		log.Fatal(err)
	}
	title1Image = ebiten.NewImageFromImage(title1AssetDecoded)

	ex2AssetDecoded, _, err := image.Decode(bytes.NewReader(ex2Asset))
	if err != nil {
		log.Fatal(err)
	}
	ex2Image = ebiten.NewImageFromImage(ex2AssetDecoded)

	title2AssetDecoded, _, err := image.Decode(bytes.NewReader(title2Asset))
	if err != nil {
		log.Fatal(err)
	}
	title2Image = ebiten.NewImageFromImage(title2AssetDecoded)

	decodedAsset, _, err := image.Decode(bytes.NewReader(existPathMatrAsset))
	if err != nil {
		log.Fatal(err)
	}
	existPathMatrImage = ebiten.NewImageFromImage(decodedAsset)

	decodedAsset, _, err = image.Decode(bytes.NewReader(titleExistPathMatrAsset))
	if err != nil {
		log.Fatal(err)
	}
	titleExistPathMatrImage = ebiten.NewImageFromImage(decodedAsset)
}

// split the graphElementsImage
var (
	nodeSubimage                   = image.Rect(0, 0, spriteSide, spriteSide)
	nodeSelectedSubimage           = image.Rect(spriteSide, 0, 2*spriteSide, spriteSide)
	loopSubimage                   = image.Rect(2*spriteSide, 0, 3*spriteSide, spriteSide)
	loopSelectedSubimage           = image.Rect(3*spriteSide, 0, 4*spriteSide, spriteSide)
	edgeToSubimage                 = image.Rect(4*spriteSide, 0, 5*spriteSide, spriteSide)
	edgeToSelectedSubimage         = image.Rect(5*spriteSide, 0, 6*spriteSide, spriteSide)
	edgeSubimage                   = image.Rect(6*spriteSide, 0, 7*spriteSide, spriteSide)
	edgeSelectedSubimage           = image.Rect(7*spriteSide, 0, 8*spriteSide, spriteSide)
	matrixTopSubimage              = image.Rect(8*spriteSide, 0, 9*spriteSide, spriteSide)
	matrixLeftSubimage             = image.Rect(9*spriteSide, 0, 10*spriteSide, spriteSide)
	graphLayoutTopLeftSubimage     = image.Rect(0, 4*spriteSide, spriteSide, 5*spriteSide)
	graphLayoutTopRightSubimage    = image.Rect(spriteSide, 4*spriteSide, 2*spriteSide, 5*spriteSide)
	graphLayoutBottomRightSubimage = image.Rect(2*spriteSide, 4*spriteSide, 3*spriteSide, 5*spriteSide)
	graphLayoutBottomLeftSubimage  = image.Rect(3*spriteSide, 4*spriteSide, 4*spriteSide, 5*spriteSide)
	buttonLeftSubimage             = image.Rect(4*spriteSide, 4*spriteSide, 5*spriteSide, 5*spriteSide)
	buttonLeftSelectedSubimage     = image.Rect(5*spriteSide, 4*spriteSide, 6*spriteSide, 5*spriteSide)
	buttonCenterSubimage           = image.Rect(6*spriteSide, 4*spriteSide, 7*spriteSide, 5*spriteSide)
	buttonCenterSelectedSubimage   = image.Rect(7*spriteSide, 4*spriteSide, 8*spriteSide, 5*spriteSide)
	buttonRightSubimage            = image.Rect(8*spriteSide, 4*spriteSide, 9*spriteSide, 5*spriteSide)
	buttonRightSelectedSubimage    = image.Rect(9*spriteSide, 4*spriteSide, 10*spriteSide, 5*spriteSide)
	undoneQuestionSubimage         = image.Rect(0, 5*spriteSide, spriteSide, 6*spriteSide)
	doneQuestionSubimage           = image.Rect(spriteSide, 5*spriteSide, 2*spriteSide, 6*spriteSide)
)
