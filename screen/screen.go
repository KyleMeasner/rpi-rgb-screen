package screen

import (
	"image"
	"rpi-rgb-screen/transition"
)

type ScreenState int

const (
	StateNotDisplayed ScreenState = iota
	StateTransitionIn
	StateDisplayed
	StateTransitionOut
)

type Screen interface {
	Render() (image.Image, bool)
	Refresh() chan bool
	SetState(state ScreenState)
	GetPreferredTransition() transition.Transition
}
