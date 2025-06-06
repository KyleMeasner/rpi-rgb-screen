package screen

import (
	"image"
)

type Screen interface {
	Render() image.Image
}
