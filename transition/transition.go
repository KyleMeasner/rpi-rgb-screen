package transition

import (
	"image"
	"rpi-rgb-screen/render"
)

type Transition interface {
	SetScreens(oldScreen, newScreen render.Renderable)
	Start()
	Render() (image.Image, bool)
}
