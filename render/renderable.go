package render

import (
	"image"
	"time"
)

type Renderable interface {
	Render(elapsed time.Duration) (image.Image, bool)
}
