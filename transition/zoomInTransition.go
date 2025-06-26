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
	keyFrames := animation.NewKeyFrames(1000 * time.Millisecond)

	keyFrames.AddNumber("scale", 10)
	keyFrames.AddNumberTransitions("scale", animation.AnimatedNumberTransition{Offset: 333, Duration: 667 * time.Millisecond, EndValue: 100})

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

	renderedNewScreen, _ := z.NewScreen.Render()

	scale := float64(z.KeyFrames.GetNumber("scale")) / 100
	z.Ctx.ScaleAbout(scale, scale, 32, 16)
	z.Ctx.DrawImageAnchored(renderedNewScreen, 32, 16, 0.5, 0.5)
	return z.Ctx.Image(), z.KeyFrames.HasEnded()
}
