package fonts

import (
	"os"

	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
)

const fontsFolder = "./resources/fonts/"

type Fonts struct {
	Size4x6   font.Face
	Size5x7   font.Face
	Size6x10  font.Face
	Size8x13B font.Face
}

func LoadFonts() *Fonts {
	return &Fonts{
		Size4x6:   loadFont("4x6.bdf"),
		Size5x7:   loadFont("5x7.bdf"),
		Size6x10:  loadFont("6x10.bdf"),
		Size8x13B: loadFont("8x13B.bdf"),
	}
}

func loadFont(fileName string) font.Face {
	fontFile, err := os.ReadFile(fontsFolder + fileName)
	if err != nil {
		panic(err)
	}
	bdfFont, err := bdf.Parse(fontFile)
	if err != nil {
		panic(err)
	}
	return bdfFont.NewFace()
}
