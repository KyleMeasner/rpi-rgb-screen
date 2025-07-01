package screen

import (
	"fmt"
	"image"
	"image/color"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/config"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/weather"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/transition"
	"rpi-rgb-screen/utils"
	"time"

	"github.com/fogleman/gg"
)

type WeatherCurrentScreen struct {
	State          ScreenState
	Ctx            *gg.Context
	KeyFrames      *animation.KeyFrames
	Fonts          *fonts.Fonts
	WeatherData    weather.WeatherData
	CurrentWeather *weather.CurrentWeather
}

func NewWeatherCurrentScreen(fonts *fonts.Fonts, weatherData weather.WeatherData) Screen {
	keyFrames := animation.NewKeyFrames(10 * time.Second)

	keyFrames.AddNumber("C", 255)
	keyFrames.AddNumberTransitions("C", animation.AnimatedNumberTransition{Offset: 4500 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 0})

	keyFrames.AddNumber("F", 0)
	keyFrames.AddNumberTransitions("F", animation.AnimatedNumberTransition{Offset: 4500 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 255})

	return &WeatherCurrentScreen{
		State:       StateNotDisplayed,
		Ctx:         gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:   keyFrames,
		Fonts:       fonts,
		WeatherData: weatherData,
	}
}

func (s *WeatherCurrentScreen) GetPreferredTransition() transition.Transition {
	return transition.NewSlideInTransition()
}

func (s *WeatherCurrentScreen) SetState(state ScreenState) {
	s.State = state

	switch state {
	case StateDisplayed:
		s.KeyFrames.Start()
	case StateTransitionIn:
		s.KeyFrames.Reset()
	}
}

func (s *WeatherCurrentScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		s.CurrentWeather = s.WeatherData.GetCurrentWeather(config.Config.Location)
		close(doneChan)
	}()

	return doneChan
}

func (s *WeatherCurrentScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	s.Ctx.SetFontFace(s.Fonts.Size5x7)

	celciusOpacity := uint8(s.KeyFrames.GetNumber("C"))
	fahrenheitOpacity := uint8(s.KeyFrames.GetNumber("F"))

	if celciusOpacity > 0 {
		s.Ctx.SetColor(color.RGBA{255, 255, 255, celciusOpacity})
		tempString := fmt.Sprintf("%.f°C", s.CurrentWeather.Temperature)
		s.Ctx.DrawStringAnchored(tempString, 32, 2, 0.5, 1)
	}
	if fahrenheitOpacity > 0 {
		s.Ctx.SetColor(color.RGBA{255, 255, 255, fahrenheitOpacity})
		tempString := fmt.Sprintf("%.f°F", utils.GetFahrenheit(s.CurrentWeather.Temperature))
		s.Ctx.DrawStringAnchored(tempString, 32, 2, 0.5, 1)
	}

	filePath := fmt.Sprintf("./resources/weatherIcons/%d.png", s.CurrentWeather.WeatherCode)
	icon, err := utils.ReadImageFromFile(filePath)
	if err == nil {
		s.Ctx.DrawImageAnchored(icon, 32, 20, 0.5, 0.5)
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}
