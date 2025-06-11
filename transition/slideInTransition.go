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
	ctx              *gg.Context
	position         image.Point
	oldScreen        screen.Screen
	newScreen        screen.Screen
	animationPercent float64
}

func NewSlideInTransition(oldScreen, newScreen screen.Screen) Transition {
	return &SlideInTransition{
		ctx:       gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		position:  image.Pt(constants.SCREEN_WIDTH, 0),
		oldScreen: oldScreen,
		newScreen: newScreen,
	}
}

func (s *SlideInTransition) Render(elapsed time.Duration) image.Image {
	s.animationPercent += float64(elapsed.Milliseconds()) / float64(animationDuration.Milliseconds())
	if s.animationPercent >= 1 {
		s.animationPercent = 1
	}

	renderedOldScreen := s.oldScreen.Render(elapsed)
	renderedNewScreen := s.newScreen.Render(elapsed)

	offset := int(constants.SCREEN_WIDTH * s.animationPercent)

	s.ctx.DrawImage(renderedOldScreen, -offset, 0)
	s.ctx.DrawImage(renderedNewScreen, constants.SCREEN_WIDTH-offset, 0)
	return s.ctx.Image()
}
