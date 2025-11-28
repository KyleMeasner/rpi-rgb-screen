package transition

import (
	"image"
	"image/color"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/render"
	"rpi-rgb-screen/utils"
	"time"

	"github.com/fogleman/gg"
)

// This transition slides the new screen in from the right-hand side of the display
type SlideAndZoomTransition struct {
	Ctx       *gg.Context
	OldScreen render.Renderable
	NewScreen render.Renderable
	KeyFrames *animation.KeyFrames
}

func NewSlideAndZoomTransition() Transition {
	keyFrames := animation.NewKeyFrames(1500 * time.Millisecond)

	keyFrames.AddNumber("scale", 100)
	keyFrames.AddNumberTransitions("scale",
		animation.AnimatedNumberTransition{Offset: 0, Duration: 150 * time.Millisecond, EndValue: 80},
		animation.AnimatedNumberTransition{Offset: 1350 * time.Millisecond, Duration: 150 * time.Millisecond, EndValue: 100},
	)

	keyFrames.AddPoint("oldScreen", image.Point{constants.SCREEN_MIDDLE_X, constants.SCREEN_MIDDLE_Y})
	keyFrames.AddPointTransitions("oldScreen",
		animation.AnimatedPointTransition{Offset: 150 * time.Millisecond, Duration: 1200 * time.Millisecond, EndValue: image.Point{-constants.SCREEN_MIDDLE_X, constants.SCREEN_MIDDLE_Y}},
	)

	keyFrames.AddPoint("newScreen", image.Point{constants.SCREEN_WIDTH + constants.SCREEN_MIDDLE_X, constants.SCREEN_MIDDLE_Y})
	keyFrames.AddPointTransitions("newScreen",
		animation.AnimatedPointTransition{Offset: 150 * time.Millisecond, Duration: 1200 * time.Millisecond, EndValue: image.Point{constants.SCREEN_MIDDLE_X, constants.SCREEN_MIDDLE_Y}},
	)

	keyFrames.Start()

	return &SlideAndZoomTransition{
		Ctx:       gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames: keyFrames,
	}
}

func (s *SlideAndZoomTransition) SetScreens(oldScreen render.Renderable, newScreen render.Renderable) {
	s.OldScreen = oldScreen
	s.NewScreen = newScreen
}

func (s *SlideAndZoomTransition) Start() {
	s.KeyFrames.Start()
}

func (s *SlideAndZoomTransition) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	renderedOldScreen, _ := s.OldScreen.Render()
	renderedNewScreen, _ := s.NewScreen.Render()

	scale := float64(s.KeyFrames.GetNumber("scale")) / 100
	if scale < 1 {
		renderedOldScreen = utils.ResizeImage(renderedOldScreen, int(float64(constants.SCREEN_WIDTH)*scale), int(float64(constants.SCREEN_HEIGHT)*scale))
		renderedNewScreen = utils.ResizeImage(renderedNewScreen, int(float64(constants.SCREEN_WIDTH)*scale), int(float64(constants.SCREEN_HEIGHT)*scale))
	}

	oldScreenPos := s.KeyFrames.GetPoint("oldScreen")
	newScreenPos := s.KeyFrames.GetPoint("newScreen")

	s.Ctx.DrawImageAnchored(renderedOldScreen, oldScreenPos.X, oldScreenPos.Y, 0.5, 0.5)
	s.Ctx.DrawImageAnchored(renderedNewScreen, newScreenPos.X, newScreenPos.Y, 0.5, 0.5)
	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}
