package transition

import (
	"image"
	"image/color"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/render"
	"time"

	"github.com/fogleman/gg"
)

// This transition rotates in the new screen from 90 degrees
type RotateTransition struct {
	Ctx       *gg.Context
	OldScreen render.Renderable
	NewScreen render.Renderable
	KeyFrames *animation.KeyFrames
}

func NewRotateTransition() Transition {
	keyFrames := animation.NewKeyFrames(1500 * time.Millisecond)

	keyFrames.AddNumber("scale", 100)
	keyFrames.AddNumberTransitions("scale",
		animation.AnimatedNumberTransition{Offset: 0, Duration: 100 * time.Millisecond, EndValue: 90},
		animation.AnimatedNumberTransition{Offset: 1400 * time.Millisecond, Duration: 100 * time.Millisecond, EndValue: 100},
	)

	keyFrames.AddNumber("angle", 0)
	keyFrames.AddNumberTransitions("angle",
		animation.AnimatedNumberTransition{Offset: 100 * time.Millisecond, Duration: 50 * time.Millisecond, EndValue: 2},
		animation.AnimatedNumberTransition{Offset: 150 * time.Millisecond, Duration: 1200 * time.Millisecond, EndValue: -62},
		animation.AnimatedNumberTransition{Offset: 1350 * time.Millisecond, Duration: 50 * time.Millisecond, EndValue: -60},
	)

	keyFrames.Start()

	return &RotateTransition{
		Ctx:       gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames: keyFrames,
	}
}

func (s *RotateTransition) SetScreens(oldScreen render.Renderable, newScreen render.Renderable) {
	s.OldScreen = oldScreen
	s.NewScreen = newScreen
}

func (s *RotateTransition) Start() {
	s.KeyFrames.Start()
}

func (s *RotateTransition) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	scale := float64(s.KeyFrames.GetNumber("scale")) / 100
	angle := float64(s.KeyFrames.GetNumber("angle"))

	// Render old screen
	renderedOldScreen, _ := s.OldScreen.Render()
	s.Ctx.Push()
	s.Ctx.Translate(2*constants.SCREEN_WIDTH, -constants.SCREEN_MIDDLE_Y)
	s.Ctx.Rotate(gg.Radians(angle))
	s.Ctx.Translate(-2*constants.SCREEN_WIDTH, constants.SCREEN_MIDDLE_Y)
	s.Ctx.ScaleAbout(scale, scale, constants.SCREEN_MIDDLE_X, constants.SCREEN_MIDDLE_Y)
	s.Ctx.DrawImage(renderedOldScreen, 0, 0)
	s.Ctx.Pop()

	// Render new screen
	renderedNewScreen, _ := s.NewScreen.Render()
	s.Ctx.Push()
	s.Ctx.Translate(2*constants.SCREEN_WIDTH, -constants.SCREEN_MIDDLE_Y)
	s.Ctx.Rotate(gg.Radians(angle + 60))
	s.Ctx.Translate(-2*constants.SCREEN_WIDTH, constants.SCREEN_MIDDLE_Y)
	s.Ctx.ScaleAbout(scale, scale, constants.SCREEN_MIDDLE_X, constants.SCREEN_MIDDLE_Y)
	s.Ctx.DrawImage(renderedNewScreen, 0, 0)
	s.Ctx.Pop()

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}
