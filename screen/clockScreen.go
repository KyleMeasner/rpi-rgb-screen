package screen

import (
	"image"
	"image/color"
	"math"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/transition"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type ClockScreen struct {
	State     ScreenState
	Ctx       *gg.Context
	KeyFrames *animation.KeyFrames
	Fonts     *fonts.Fonts
}

func NewClockScreen(fonts *fonts.Fonts) Screen {
	return &ClockScreen{
		State:     StateNotDisplayed,
		Ctx:       gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames: animation.NewKeyFrames(30 * time.Second),
		Fonts:     fonts,
	}
}

func (s *ClockScreen) GetPreferredTransition() transition.Transition {
	return transition.NewSlideInTransition()
}

func (s *ClockScreen) SetState(state ScreenState) {
	s.State = state
	if state == StateDisplayed {
		s.KeyFrames.Start()
	}
}

func (s *ClockScreen) Refresh() chan bool {
	doneChan := make(chan bool)
	close(doneChan)
	return doneChan
}

func (s *ClockScreen) Render() (image.Image, bool) {
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	currentTime := time.Now()

	// Clock on the left
	s.Ctx.SetColor(color.White)
	for i := range 12 {
		x, y := positionOnCircle(float64(i) / 12)
		s.Ctx.SetPixel(int(math.Round(15+14*x)), int(math.Round(15+14*y)))
	}

	hours, minutes, seconds := currentTime.Clock()
	if hours >= 12 {
		hours -= 12
	}

	s.Ctx.SetColor(color.RGBA{0, 255, 0, 255})
	hoursPercentage := (float64(hours)*60 + float64(minutes)) / 720
	hoursX, hoursY := positionOnCircle(hoursPercentage)
	s.drawLine(15, 15, hoursX, hoursY, 8)

	s.Ctx.SetColor(color.RGBA{0, 0, 255, 255})
	minutesX, minutesY := positionOnCircle(float64(minutes) / 60)
	s.drawLine(15, 15, minutesX, minutesY, 11)

	s.Ctx.SetColor(color.RGBA{255, 0, 0, 255})
	secondsX, secondsY := positionOnCircle(float64(seconds) / 60)
	s.drawLine(15, 15, secondsX, secondsY, 13)

	// Date and time on the right
	s.Ctx.SetColor(color.White)
	s.Ctx.SetFontFace(s.Fonts.Size5x7)

	timeString := currentTime.Format("03:04")
	s.Ctx.DrawStringAnchored(timeString, 64, 0, 1, 1)
	amPmString := currentTime.Format("PM")
	s.Ctx.DrawStringAnchored(amPmString, 64, 7, 1, 1)

	dayOfWeekString := strings.ToUpper(currentTime.Format("Mon"))
	s.Ctx.DrawStringAnchored(dayOfWeekString, 64, 24, 1, 0)
	monthString := strings.ToUpper(currentTime.Format("Jan02"))
	s.Ctx.DrawStringAnchored(monthString, 64, 31, 1, 0)

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *ClockScreen) drawLine(xStart, yStart, xOffset, yOffset float64, length int) {
	for i := range length {
		x := int(math.Round(xStart + xOffset*float64(i)))
		y := int(math.Round(yStart + yOffset*float64(i)))
		s.Ctx.SetPixel(x, y)
	}
}

func positionOnCircle(percentage float64) (float64, float64) {
	radians := math.Pi - (2 * math.Pi * percentage)
	return math.Sin(radians), math.Cos(radians)
}
