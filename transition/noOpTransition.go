package transition

import (
	"image"
	"rpi-rgb-screen/render"
)

type NoOpTransition struct {
	NewScreen render.Renderable
}

func NewNoOpTransition() Transition {
	return &NoOpTransition{}
}

func (n *NoOpTransition) SetScreens(oldScreen render.Renderable, newScreen render.Renderable) {
	n.NewScreen = newScreen
}

func (n *NoOpTransition) Start() {

}

func (n *NoOpTransition) Render() (image.Image, bool) {
	renderedScreen, _ := n.NewScreen.Render()
	return renderedScreen, true
}
