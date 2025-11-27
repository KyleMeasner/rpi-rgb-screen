package utils

import (
	"bytes"
	"image"
	"math"
	"os"

	"github.com/fogleman/gg"
)

// Resizes the image to be contained in a square of newWidthOrHeight.
func ResizeImage(image image.Image, newWidth, newHeight int) image.Image {
	scaleFactorX := float64(newWidth) / float64(image.Bounds().Dx())
	scaleFactorY := float64(newHeight) / float64(image.Bounds().Dy())

	centerX := newWidth / 2
	centerY := newHeight / 2

	resizeCtx := gg.NewContext(newWidth, newHeight)
	resizeCtx.ScaleAbout(scaleFactorX, scaleFactorY, float64(centerX), float64(centerY))
	resizeCtx.DrawImageAnchored(image, centerX, centerY, 0.5, 0.5)
	return resizeCtx.Image()
}

// Resizes the image to be contained in a square of newWidthOrHeight.
func ResizeImageSquare(image image.Image, newWidthOrHeight int) image.Image {
	scaleFactor := float64(newWidthOrHeight) / math.Max(float64(image.Bounds().Dx()), float64(image.Bounds().Dy()))
	center := newWidthOrHeight / 2

	resizeCtx := gg.NewContext(newWidthOrHeight, newWidthOrHeight)
	resizeCtx.ScaleAbout(scaleFactor, scaleFactor, float64(center), float64(center))
	resizeCtx.DrawImageAnchored(image, center, center, 0.5, 0.5)
	return resizeCtx.Image()
}

func ReadImageFromFile(filePath string) (image.Image, error) {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	image, _, err := image.Decode(bytes.NewReader(fileContents))
	if err != nil {
		return nil, err
	}

	return image, nil
}
