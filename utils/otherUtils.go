package utils

import (
	"fmt"
	"image/color"
)

func ComputeValue(start, end int, percentComplete float64) int {
	return start + int(float64(end-start)*percentComplete)
}

func GetFahrenheit(celsius float64) float64 {
	return (celsius * 9 / 5) + 32
}

// Parses a hex color string (e.g. "#RRGGBB") and returns a color.Color object
func ParseColorFromHex(hex string) (color.Color, error) {
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return nil, err
	}
	return color.RGBA{R: r, G: g, B: b, A: 0xFF}, nil
}
