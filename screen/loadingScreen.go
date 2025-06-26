package screen

import (
	"image"
	"image/color"
	"math/rand"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/transition"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

type LoadingScreen struct {
	State        ScreenState
	Ctx          *gg.Context
	KeyFrames    *animation.KeyFrames
	Fonts        *fonts.Fonts
	Color        color.RGBA
	SelectedFont font.Face
}

func NewLoadingScreen(fonts *fonts.Fonts) Screen {
	return &LoadingScreen{
		State:        StateNotDisplayed,
		Ctx:          gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:    animation.NewKeyFrames(3 * time.Second),
		Fonts:        fonts,
		Color:        color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
		SelectedFont: fonts.Size5x7,
	}
}

func (s *LoadingScreen) GetPreferredTransition() transition.Transition {
	return transition.NewNoOpTransition()
}

func (s *LoadingScreen) SetState(state ScreenState) {
	s.State = state
	if state == StateDisplayed {
		s.KeyFrames.Start()
	}
}

func (s *LoadingScreen) Render() (image.Image, bool) {
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	positions := []image.Point{image.Pt(0, 0), image.Pt(0, 32), image.Pt(64, 0), image.Pt(64, 32)}
	for _, position := range positions {
		s.Ctx.DrawCircle(float64(position.X), float64(position.Y), 5)
		s.Ctx.SetColor(s.Color)
		s.Ctx.Fill()
	}

	s.Ctx.SetFontFace(s.SelectedFont)
	s.Ctx.DrawStringAnchored("LOADING", 32, 16, 0.5, 0.5)

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *LoadingScreen) Refresh() chan bool {
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

		close(doneChan)
	}()

	return doneChan
}
