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
	"time"

	"github.com/fogleman/gg"
)

type WeatherForecastScreen struct {
	State           ScreenState
	Ctx             *gg.Context
	KeyFrames       *animation.KeyFrames
	Fonts           *fonts.Fonts
	WeatherData     weather.WeatherData
	WeatherForecast []*weather.WeatherForecast
}

func NewWeatherForecastScreen(fonts *fonts.Fonts, weatherData weather.WeatherData) Screen {
	return &WeatherForecastScreen{
		State:       StateNotDisplayed,
		Ctx:         gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:   animation.NewKeyFrames(2 * time.Second),
		Fonts:       fonts,
		WeatherData: weatherData,
	}
}

func (s *WeatherForecastScreen) GetPreferredTransition() transition.Transition {
	return transition.NewSlideInTransition()
}

func (s *WeatherForecastScreen) SetState(state ScreenState) {
	s.State = state
	if state == StateDisplayed {
		s.KeyFrames.Start()
	}
}

func (s *WeatherForecastScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		s.WeatherForecast = s.WeatherData.GetForecast(config.Config.Location)
		close(doneChan)
	}()

	return doneChan
}

func (s *WeatherForecastScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	s.Ctx.SetFontFace(s.Fonts.Size5x7)
	s.Ctx.SetColor(color.White)

	tomorrowsWeather := s.WeatherForecast[1]

	tempString := fmt.Sprintf("%.f°C - %.f°C", tomorrowsWeather.TemperatureMin, tomorrowsWeather.TemperatureMax)
	s.Ctx.DrawStringAnchored(tempString, 32, 2, 0.5, 1)

	weatherCodeString := fmt.Sprintf("%d", tomorrowsWeather.WeatherCode)
	s.Ctx.DrawStringAnchored(weatherCodeString, 32, 15, 0.5, 1)

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}
