package utils

import (
	"image"
	"math"

	"github.com/fogleman/gg"
)

// Resizes the image to be contained in a square of newWidthOrHeight.
func ResizeImage(image image.Image, newWidthOrHeight int) image.Image {
	scaleFactor := float64(newWidthOrHeight) / math.Max(float64(image.Bounds().Dx()), float64(image.Bounds().Dy()))
	center := newWidthOrHeight / 2

	resizeCtx := gg.NewContext(newWidthOrHeight, newWidthOrHeight)
	resizeCtx.ScaleAbout(scaleFactor, scaleFactor, float64(center), float64(center))
	resizeCtx.DrawImageAnchored(image, center, center, 0.5, 0.5)
	return resizeCtx.Image()
}
