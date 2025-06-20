package transition

import (
	"image"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/screen"
	"time"

	"github.com/fogleman/gg"
)

// This transition slides the new screen in from the right-hand side of the display
type SlideInTransition struct {
	Ctx       *gg.Context
	Position  image.Point
	OldScreen screen.Screen
	NewScreen screen.Screen
	KeyFrames *animation.KeyFrames
}

func NewSlideInTransition(oldScreen, newScreen screen.Screen) Transition {
	keyFrames := animation.NewKeyFrames()

	keyFrames.AddPoint("oldScreen", image.Point{0, 0})
	keyFrames.AddPointTransitions("oldScreen",
		animation.AnimatedPointTransition{Offset: 0, Duration: 1500 * time.Millisecond, EndValue: image.Point{-constants.SCREEN_WIDTH, 0}},
	)

	keyFrames.AddPoint("newScreen", image.Point{constants.SCREEN_WIDTH, 0})
	keyFrames.AddPointTransitions("newScreen",
		animation.AnimatedPointTransition{Offset: 0, Duration: 1500 * time.Millisecond, EndValue: image.Point{0, 0}},
	)

	keyFrames.Start()

	return &SlideInTransition{
		Ctx:       gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		Position:  image.Pt(constants.SCREEN_WIDTH, 0),
		OldScreen: oldScreen,
		NewScreen: newScreen,
		KeyFrames: keyFrames,
	}
}

func (s *SlideInTransition) Render(elapsed time.Duration) (image.Image, bool) {
	renderedOldScreen, _ := s.OldScreen.Render(elapsed)
	renderedNewScreen, _ := s.NewScreen.Render(elapsed)

	oldScreenPos := s.KeyFrames.GetPoint("oldScreen")
	newScreenPos := s.KeyFrames.GetPoint("newScreen")

	s.Ctx.DrawImage(renderedOldScreen, oldScreenPos.X, oldScreenPos.Y)
	s.Ctx.DrawImage(renderedNewScreen, newScreenPos.X, newScreenPos.Y)
	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}
