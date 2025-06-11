package screen

import (
	"image"
	"image/color"
	"math/rand"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/fonts"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

const speed = 1.2 // Pixels per second
const textUpperBound = 5
const textLowerBound = 10

// A dummy screen used to test drawing screens to the display
type DummyScreen struct {
	ctx           *gg.Context
	color         color.RGBA
	fonts         *fonts.Fonts
	selectedFont  font.Face
	textYPosition float64
	textDirection int
}

func NewDummyScreen(fonts *fonts.Fonts) Screen {
	return &DummyScreen{
		ctx:           gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		color:         color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
		fonts:         fonts,
		selectedFont:  fonts.Size5x7,
		textYPosition: textUpperBound,
		textDirection: 1,
	}
}

func (s *DummyScreen) Render(elapsed time.Duration) image.Image {
	s.ctx.SetColor(color.Black)
	s.ctx.Clear()

	positions := []image.Point{image.Pt(0, 0), image.Pt(0, 32), image.Pt(64, 0), image.Pt(64, 32)}
	for _, position := range positions {
		s.ctx.DrawCircle(float64(position.X), float64(position.Y), 5)
		s.ctx.SetColor(s.color)
		s.ctx.Fill()
	}

	s.textYPosition += speed * elapsed.Seconds() * float64(s.textDirection)
	if s.textYPosition >= textLowerBound {
		s.textYPosition = textLowerBound
		s.textDirection = -1
	} else if s.textYPosition <= textUpperBound {
		s.textYPosition = textUpperBound
		s.textDirection = 1
	}

	s.ctx.SetFontFace(s.selectedFont)
	s.ctx.DrawStringAnchored("Dummy", 32, s.textYPosition, 0.5, 0.5)
	s.ctx.DrawStringAnchored("Screen", 32, s.textYPosition+13, 0.5, 0.5)

	return s.ctx.Image()
}

func (s *DummyScreen) Refresh() {
	// Change font
	switch rand.Intn(3) {
	case 0:
		s.selectedFont = s.fonts.Size5x7
	case 1:
		s.selectedFont = s.fonts.Size6x10
	case 2:
		s.selectedFont = s.fonts.Size8x13B
	}

	// Change color
	s.color = color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}
}

func (s *DummyScreen) TransitionStart() {

}

func (s *DummyScreen) TransitionEnd(isDisplayed bool) {

}
