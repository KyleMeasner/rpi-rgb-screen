package fonts

import (
	"os"

	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
)

const fontDirectory = "/home/kyle/repos/bitmap-fonts/bitmap/"

type Fonts struct {
	Bitocra     font.Face
	Lemon       font.Face
	Scientifica font.Face
}

func LoadFonts() *Fonts {
	return &Fonts{
		Bitocra:     loadFont("bitocra/bitocra7.bdf"),
		Lemon:       loadFont("phallus/lemon.bdf"),
		Scientifica: loadFont("scientifica/scientifica-11.bdf"),
	}
}

func loadFont(fontLocation string) font.Face {
	fontFile, err := os.ReadFile(fontDirectory + fontLocation)
	if err != nil {
		panic(err)
	}
	bdfFont, err := bdf.Parse(fontFile)
	if err != nil {
		panic(err)
	}
	return bdfFont.NewFace()
}
