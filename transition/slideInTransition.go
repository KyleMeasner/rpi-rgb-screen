package transition

import (
	"image"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/screen"
	"time"

	"github.com/fogleman/gg"
)

const animationDuration = 3 * time.Second

// This transition slides the new screen in from the right-hand side of the display
type SlideInTransition struct {
	Ctx              *gg.Context
	Position         image.Point
	OldScreen        screen.Screen
	NewScreen        screen.Screen
	AnimationPercent float64
}

func NewSlideInTransition(oldScreen, newScreen screen.Screen) Transition {
	return &SlideInTransition{
		Ctx:       gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		Position:  image.Pt(constants.SCREEN_WIDTH, 0),
		OldScreen: oldScreen,
		NewScreen: newScreen,
	}
}

func (s *SlideInTransition) Render(elapsed time.Duration) (image.Image, bool) {
	s.AnimationPercent += float64(elapsed.Milliseconds()) / float64(animationDuration.Milliseconds())
	if s.AnimationPercent >= 1 {
		s.AnimationPercent = 1
	}

	renderedOldScreen, _ := s.OldScreen.Render(elapsed)
	renderedNewScreen, _ := s.NewScreen.Render(elapsed)

	offset := int(constants.SCREEN_WIDTH * s.AnimationPercent)

	s.Ctx.DrawImage(renderedOldScreen, -offset, 0)
	s.Ctx.DrawImage(renderedNewScreen, constants.SCREEN_WIDTH-offset, 0)
	return s.Ctx.Image(), s.AnimationPercent == 1
}
