package transition

import (
	"image"
	"io"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/screen"
	"time"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
	"github.com/fogleman/gg"
)

type SlideIn struct {
	ctx       *gg.Context
	position  image.Point
	oldScreen screen.Screen
	newScreen screen.Screen
}

func NewSlideIn(oldScreen, newScreen screen.Screen) rgbmatrix.Animation {
	return &SlideIn{
		ctx:       gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		position:  image.Pt(constants.SCREEN_WIDTH, 0),
		oldScreen: oldScreen,
		newScreen: newScreen,
	}
}

func (s *SlideIn) Next() (image.Image, <-chan time.Time, error) {
	s.position = s.position.Sub(image.Pt(1, 0))

	if s.position.X < 0 {
		// Animation is done
		return &image.NRGBA{}, make(<-chan time.Time), io.EOF
	}

	renderedOldScreen := s.oldScreen.Render()
	renderedNewScreen := s.newScreen.Render()

	s.ctx.DrawImage(renderedOldScreen, s.position.X-64, s.position.Y)
	s.ctx.DrawImage(renderedNewScreen, s.position.X, s.position.Y)
	return s.ctx.Image(), time.After(time.Millisecond * 50), nil
}
