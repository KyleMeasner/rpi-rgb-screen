package screen

import (
	"image"
	"time"
)

type Screen interface {
	Render(elapsed time.Duration) image.Image
	Refresh()
	TransitionStart()
	TransitionEnd(isDisplayed bool)
}
