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

const speed = 5 // Pixels per second
const textUpperBound = 5
const textLowerBound = 10
const screenDuration = 3 * time.Second

// A dummy screen used to test drawing screens to the display
type DummyScreen struct {
	Ctx                 *gg.Context
	Color               color.RGBA
	Fonts               *fonts.Fonts
	SelectedFont        font.Face
	TextYPosition       float64
	TextDirection       int
	ScreenDisplayedTime time.Time
}

func NewDummyScreen(fonts *fonts.Fonts) Screen {
	return &DummyScreen{
		Ctx:                 gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		Color:               color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
		Fonts:               fonts,
		SelectedFont:        fonts.Size5x7,
		TextYPosition:       textUpperBound,
		TextDirection:       1,
		ScreenDisplayedTime: time.Now(),
	}
}

func (s *DummyScreen) Render(elapsed time.Duration) (image.Image, bool) {
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	positions := []image.Point{image.Pt(0, 0), image.Pt(0, 32), image.Pt(64, 0), image.Pt(64, 32)}
	for _, position := range positions {
		s.Ctx.DrawCircle(float64(position.X), float64(position.Y), 5)
		s.Ctx.SetColor(s.Color)
		s.Ctx.Fill()
	}

	s.TextYPosition += speed * elapsed.Seconds() * float64(s.TextDirection)
	if s.TextYPosition >= textLowerBound {
		s.TextYPosition = textLowerBound
		s.TextDirection = -1
	} else if s.TextYPosition <= textUpperBound {
		s.TextYPosition = textUpperBound
		s.TextDirection = 1
	}

	s.Ctx.SetFontFace(s.SelectedFont)
	s.Ctx.DrawStringAnchored("DUMMY", 32, s.TextYPosition, 0.5, 0.5)
	s.Ctx.DrawStringAnchored("SCREEN", 32, s.TextYPosition+13, 0.5, 0.5)

	return s.Ctx.Image(), time.Now().Sub(s.ScreenDisplayedTime) > screenDuration
}

func (s *DummyScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		// Change font
		switch rand.Intn(3) {
		case 0:
			s.SelectedFont = s.Fonts.Size5x7
		case 1:
			s.SelectedFont = s.Fonts.Size6x10
		case 2:
			s.SelectedFont = s.Fonts.Size8x13B
		}

		// Change color
		s.Color = color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}

		s.TextYPosition = textUpperBound
		s.TextDirection = 1
		close(doneChan)
	}()

	return doneChan
}

func (s *DummyScreen) TransitionStart() {

}

func (s *DummyScreen) TransitionEnd(isDisplayed bool) {
	if isDisplayed {
		s.ScreenDisplayedTime = time.Now()
	}
}
