package transition

import (
	"image"
	"time"
)

type Transition interface {
	Render(elapsed time.Duration) image.Image
}
