package render

import (
	"image"
)

type Renderable interface {
	Render() (image.Image, bool)
}
