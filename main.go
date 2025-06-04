package main

import (
	"image"
	"image/color"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

func main() {
	config := &rgbmatrix.DefaultConfig
	config.Rows = 32
	config.Cols = 64
	config.Brightness = 100
	config.HardwareMapping = "adafruit-hat"
	config.ShowRefreshRate = true

	matrix, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		panic(err)
	}

	canvas := rgbmatrix.NewCanvas(matrix)
	defer canvas.Close()

	bounds := canvas.Bounds()
	pixels := []image.Point{}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			// fmt.Println("x", x, "y", y)
			// canvas.Set(x, y, color.RGBA{255, 0, 0, 255})

			pixels = append(pixels, image.Point{x, y})
			for _, point := range pixels {
				canvas.Set(point.X, point.Y, color.RGBA{255, 0, 0, 255})
			}

			canvas.Render()
		}
	}
}
