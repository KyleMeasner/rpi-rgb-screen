package transition

import (
	"image"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/render"
	"time"

	"github.com/fogleman/gg"
)

// This transition slides the new screen in from the right-hand side of the display
type SlideInTransition struct {
	Ctx       *gg.Context
	OldScreen render.Renderable
	NewScreen render.Renderable
	KeyFrames *animation.KeyFrames
}

func NewSlideInTransition() Transition {
	keyFrames := animation.NewKeyFrames(1500 * time.Millisecond)

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
		KeyFrames: keyFrames,
	}
}

func (s *SlideInTransition) SetScreens(oldScreen render.Renderable, newScreen render.Renderable) {
	s.OldScreen = oldScreen
	s.NewScreen = newScreen
}

func (s *SlideInTransition) Start() {
	s.KeyFrames.Start()
}

func (s *SlideInTransition) Render() (image.Image, bool) {
	renderedOldScreen, _ := s.OldScreen.Render()
	renderedNewScreen, _ := s.NewScreen.Render()

	oldScreenPos := s.KeyFrames.GetPoint("oldScreen")
	newScreenPos := s.KeyFrames.GetPoint("newScreen")

	s.Ctx.DrawImage(renderedOldScreen, oldScreenPos.X, oldScreenPos.Y)
	s.Ctx.DrawImage(renderedNewScreen, newScreenPos.X, newScreenPos.Y)
	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}
