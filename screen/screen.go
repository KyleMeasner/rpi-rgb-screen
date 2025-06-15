package screen

import (
	"image"
	"time"
)

type ScreenState int

const (
	StateNotDisplayed ScreenState = iota
	StateTransitionIn
	StateDisplayed
	StateTransitionOut
)

type Screen interface {
	Render(elapsed time.Duration) (image.Image, bool)
	Refresh() chan bool
	SetState(state ScreenState)
}
