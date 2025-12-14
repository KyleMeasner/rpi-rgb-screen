package screen

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/weather"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/transition"
	"rpi-rgb-screen/utils"
	"time"

	"github.com/fogleman/gg"
)

type WeatherForecastDayScreen struct {
	State           ScreenState
	Ctx             *gg.Context
	KeyFrames       *animation.KeyFrames
	Fonts           *fonts.Fonts
	WeatherForecast *weather.WeatherForecast
	WeatherIcon     image.Image
}

func NewWeatherForecastDayScreen(fonts *fonts.Fonts, weatherForecast *weather.WeatherForecast) Screen {
	keyFrames := animation.NewKeyFrames(11 * time.Second)

	keyFrames.AddNumber("offsetTensHigh", 0)
	keyFrames.AddNumberTransitions("offsetTensHigh",
		animation.AnimatedNumberTransition{Offset: 4000 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: -2},
		animation.AnimatedNumberTransition{Offset: 4200 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 14},
		animation.AnimatedNumberTransition{Offset: 5200 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: 11},
	)

	keyFrames.AddNumber("offsetOnesHigh", 0)
	keyFrames.AddNumberTransitions("offsetOnesHigh",
		animation.AnimatedNumberTransition{Offset: 4200 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: -2},
		animation.AnimatedNumberTransition{Offset: 4400 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 14},
		animation.AnimatedNumberTransition{Offset: 5400 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: 11},
	)

	keyFrames.AddNumber("offsetUnits", 0)
	keyFrames.AddNumberTransitions("offsetUnits",
		animation.AnimatedNumberTransition{Offset: 4400 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: -2},
		animation.AnimatedNumberTransition{Offset: 4600 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 14},
		animation.AnimatedNumberTransition{Offset: 5600 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: 11},
	)

	keyFrames.AddNumber("offsetTensLow", 0)
	keyFrames.AddNumberTransitions("offsetTensLow",
		animation.AnimatedNumberTransition{Offset: 4600 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: -2},
		animation.AnimatedNumberTransition{Offset: 4800 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 14},
		animation.AnimatedNumberTransition{Offset: 5800 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: 11},
	)

	keyFrames.AddNumber("offsetOnesLow", 0)
	keyFrames.AddNumberTransitions("offsetOnesLow",
		animation.AnimatedNumberTransition{Offset: 4800 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: -2},
		animation.AnimatedNumberTransition{Offset: 5000 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 14},
		animation.AnimatedNumberTransition{Offset: 6000 * time.Millisecond, Duration: 200 * time.Millisecond, EndValue: 11},
	)

	return &WeatherForecastDayScreen{
		State:           StateNotDisplayed,
		Ctx:             gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:       keyFrames,
		Fonts:           fonts,
		WeatherForecast: weatherForecast,
	}
}

func (s *WeatherForecastDayScreen) GetPreferredTransition() transition.Transition {
	return transition.NewRotateTransition()
}

func (s *WeatherForecastDayScreen) SetState(state ScreenState) {
	s.State = state

	switch state {
	case StateDisplayed:
		s.KeyFrames.Start()
	case StateTransitionIn:
		s.KeyFrames.Reset()
	}
}

func (s *WeatherForecastDayScreen) Refresh() chan bool {
	doneChan := make(chan bool)
	go func() {
		filePath := fmt.Sprintf("./resources/weatherIcons/%d.png", s.WeatherForecast.WeatherCode)
		icon, err := utils.ReadImageFromFile(filePath)
		if err == nil {
			s.WeatherIcon = utils.ResizeImageSquare(icon, 14)
		} else {
			s.WeatherIcon = nil
			log.Printf("Failed to load weather icon for weather code %d: %v", s.WeatherForecast.WeatherCode, err)
		}

		close(doneChan)
	}()
	return doneChan
}

func (s *WeatherForecastDayScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	s.Ctx.SetColor(color.White)
	s.Ctx.SetFontFace(s.Fonts.Size8x13B)
	dayString := s.WeatherForecast.Date.Format("Mon")
	s.Ctx.DrawStringAnchored(dayString, 2, -1, 0, 1)

	if s.WeatherIcon != nil {
		s.Ctx.DrawImageAnchored(s.WeatherIcon, 13, 22, 0.5, 0.5)
	}

	s.Ctx.SetColor(color.RGBA{0xAB, 0xAB, 0xAB, 255}) // #ABABAB
	s.Ctx.DrawRectangle(33, 4, 28, 11)
	s.Ctx.DrawRectangle(33, 17, 18, 11)
	s.Ctx.Fill()

	s.Ctx.SetColor(color.RGBA{0xE3, 0xE3, 0xE3, 255}) // #E3E3E3
	s.Ctx.DrawRectangle(33, 5, 28, 9)
	s.Ctx.DrawRectangle(33, 18, 18, 9)
	s.Ctx.Fill()

	s.renderTemperature(s.WeatherForecast.TemperatureMax, image.Point{33, 4}, color.RGBA{0x57, 0x0D, 0x0D, 255}, s.KeyFrames.GetNumber("offsetTensHigh"), s.KeyFrames.GetNumber("offsetOnesHigh"), true)
	s.renderTemperature(s.WeatherForecast.TemperatureMin, image.Point{33, 17}, color.RGBA{0x04, 0x1D, 0x56, 255}, s.KeyFrames.GetNumber("offsetTensLow"), s.KeyFrames.GetNumber("offsetOnesLow"), false)

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *WeatherForecastDayScreen) renderTemperature(temperature float64, position image.Point, textColor color.Color, offsetTens, offsetOnes int, withDegrees bool) {
	width := 18
	if withDegrees {
		width = 28
	}
	s.Ctx.DrawRectangle(float64(position.X), float64(position.Y), float64(width), 11)
	s.Ctx.Clip()

	s.Ctx.SetFontFace(s.Fonts.Size6x10)
	s.Ctx.SetColor(textColor)

	tensDigitString := fmt.Sprintf("%d", int(temperature)/10)
	if tensDigitString == "0" {
		tensDigitString = ""
	}
	s.renderDigit(tensDigitString, image.Point{X: position.X + 2, Y: position.Y}, offsetTens)

	onesDigitString := fmt.Sprintf("%d", int(temperature)%10)
	s.renderDigit(onesDigitString, image.Point{X: position.X + 11, Y: position.Y}, offsetOnes)

	fahrenheit := utils.GetFahrenheit(temperature)
	tensDigitString = fmt.Sprintf("%d", int(fahrenheit)/10)
	if tensDigitString == "0" {
		tensDigitString = ""
	}
	s.renderDigit(tensDigitString, image.Point{X: position.X + 2, Y: position.Y - 11}, offsetTens)

	onesDigitString = fmt.Sprintf("%d", int(fahrenheit)%10)
	s.renderDigit(onesDigitString, image.Point{X: position.X + 11, Y: position.Y - 11}, offsetOnes)

	if withDegrees {
		offsetUnits := s.KeyFrames.GetNumber("offsetUnits")
		s.Ctx.SetFontFace(s.Fonts.Size5x7)
		s.renderUnits("°C", image.Point{X: position.X + 18, Y: position.Y}, offsetUnits)
		s.renderUnits("°F", image.Point{X: position.X + 18, Y: position.Y - 11}, offsetUnits)
	}

	s.Ctx.ResetClip()
}

func (s *WeatherForecastDayScreen) renderDigit(digit string, position image.Point, offset int) {
	finalPosition := image.Point{X: position.X, Y: position.Y + offset}
	s.Ctx.DrawStringAnchored(digit, float64(finalPosition.X)+3, float64(finalPosition.Y)+4, 0.5, 0.5)
}

func (s *WeatherForecastDayScreen) renderUnits(units string, position image.Point, offset int) {
	finalPosition := image.Point{X: position.X, Y: position.Y + offset}
	s.Ctx.DrawStringAnchored(units, float64(finalPosition.X)+5, float64(finalPosition.Y)+5, 0.5, 0.5)
}
