package screen

import (
	"image"
	"time"
)

type Screen interface {
	Render(elapsed time.Duration) (image.Image, bool)
	Refresh() chan bool
	TransitionStart()
	TransitionEnd(isDisplayed bool)
}
