package screen

import (
	"image"
	"image/color"
	"math/rand"
	"rpi-rgb-screen/fonts"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

// A dummy screen used to test drawing screens to the display
type DummyScreen struct {
	ctx          *gg.Context
	color        color.RGBA
	fonts        *fonts.Fonts
	selectedFont font.Face
}

func NewDummyScreen(fonts *fonts.Fonts) Screen {
	return &DummyScreen{
		ctx:          gg.NewContext(64, 32),
		color:        color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
		fonts:        fonts,
		selectedFont: fonts.Scientifica,
	}
}

func (s *DummyScreen) Render() image.Image {
	s.ctx.SetColor(color.Black)
	s.ctx.Clear()

	positions := []image.Point{image.Pt(0, 0), image.Pt(0, 32), image.Pt(64, 0), image.Pt(64, 32)}
	for _, position := range positions {
		s.ctx.DrawCircle(float64(position.X), float64(position.Y), 5)
		s.ctx.SetColor(s.color)
		s.ctx.Fill()
	}

	s.ctx.SetFontFace(s.selectedFont)
	s.ctx.DrawStringAnchored("Dummy Screen", 32, 14, 0.5, 0.5)

	return s.ctx.Image()
}

func (s *DummyScreen) Refresh() {
	// Change font
	switch rand.Intn(3) {
	case 0:
		s.selectedFont = s.fonts.Bitocra
	case 1:
		s.selectedFont = s.fonts.Lemon
	case 2:
		s.selectedFont = s.fonts.Scientifica
	}

	// Change color
	s.color = color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}
}
