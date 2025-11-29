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

type ZoomInTransition struct {
	Ctx       *gg.Context
	KeyFrames *animation.KeyFrames
	OldScreen render.Renderable
	NewScreen render.Renderable
}

func NewZoomInTransition() Transition {
	keyFrames := animation.NewKeyFrames(1500 * time.Millisecond)

	keyFrames.AddNumber("scaleOldScreen", 100)
	keyFrames.AddNumberTransitions("scaleOldScreen", animation.AnimatedNumberTransition{Offset: 0, Duration: 700 * time.Millisecond, EndValue: 0})

	keyFrames.AddNumber("scaleNewScreen", 0)
	keyFrames.AddNumberTransitions("scaleNewScreen", animation.AnimatedNumberTransition{Offset: 800 * time.Millisecond, Duration: 700 * time.Millisecond, EndValue: 100})

	return &ZoomInTransition{
		Ctx:       gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames: keyFrames,
	}
}

func (z *ZoomInTransition) SetScreens(oldScreen render.Renderable, newScreen render.Renderable) {
	z.OldScreen = oldScreen
	z.NewScreen = newScreen
}

func (z *ZoomInTransition) Start() {
	z.KeyFrames.Start()
}

func (z *ZoomInTransition) Render() (image.Image, bool) {
	// Clear image context
	z.Ctx.Identity()
	z.Ctx.SetColor(color.Black)
	z.Ctx.Clear()

	scaleOldScreen := float64(z.KeyFrames.GetNumber("scaleOldScreen")) / 100
	if scaleOldScreen > 0 {
		renderedOldScreen, _ := z.OldScreen.Render()
		z.Ctx.ScaleAbout(scaleOldScreen, scaleOldScreen, 32, 16)
		z.Ctx.DrawImageAnchored(renderedOldScreen, 32, 16, 0.5, 0.5)
		return z.Ctx.Image(), z.KeyFrames.HasEnded()
	}

	scaleNewScreen := float64(z.KeyFrames.GetNumber("scaleNewScreen")) / 100
	renderedNewScreen, _ := z.NewScreen.Render()
	z.Ctx.ScaleAbout(scaleNewScreen, scaleNewScreen, 32, 16)
	z.Ctx.DrawImageAnchored(renderedNewScreen, 32, 16, 0.5, 0.5)
	return z.Ctx.Image(), z.KeyFrames.HasEnded()
}
