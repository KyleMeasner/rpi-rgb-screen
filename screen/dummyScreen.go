package screen

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
)

type DummyScreen struct {
	ctx   *gg.Context
	color color.RGBA
}

func NewDummyScreen() Screen {
	return &DummyScreen{
		ctx:   gg.NewContext(64, 32),
		color: color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
	}
}

func (s *DummyScreen) Render() image.Image {
	s.ctx.SetColor(color.Black)
	s.ctx.Clear()

	positions := []image.Point{image.Pt(0, 0), image.Pt(32, 16), image.Pt(0, 32), image.Pt(64, 0), image.Pt(64, 32)}
	for _, position := range positions {
		s.ctx.DrawCircle(float64(position.X), float64(position.Y), 5)
		s.ctx.SetColor(s.color)
		s.ctx.Fill()
	}

	return s.ctx.Image()
}
