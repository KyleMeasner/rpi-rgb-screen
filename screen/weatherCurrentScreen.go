package screen

import (
	"fmt"
	"image"
	"image/color"
	"math"
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
	HourlyWeather  []*weather.HourlyWeather
}

func NewWeatherCurrentScreen(fonts *fonts.Fonts, weatherData weather.WeatherData) Screen {
	keyFrames := animation.NewKeyFrames(30 * time.Second)

	keyFrames.AddNumber("C", 255)
	keyFrames.AddNumberTransitions("C", animation.AnimatedNumberTransition{Offset: 14500 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 0})

	keyFrames.AddNumber("F", 0)
	keyFrames.AddNumberTransitions("F", animation.AnimatedNumberTransition{Offset: 14500 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 255})

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
		s.HourlyWeather = s.WeatherData.GetHourlyWeather(config.Config.Location)
		close(doneChan)
	}()

	return doneChan
}

func (s *WeatherCurrentScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	celciusOpacity := uint8(s.KeyFrames.GetNumber("C"))
	fahrenheitOpacity := uint8(s.KeyFrames.GetNumber("F"))

	if celciusOpacity > 0 {
		s.Ctx.SetColor(color.RGBA{255, 255, 255, celciusOpacity})

		s.Ctx.SetFontFace(s.Fonts.Size8x13B)
		tempString := fmt.Sprintf("%.f°C", s.CurrentWeather.Temperature)
		s.Ctx.DrawStringAnchored(tempString, 33, -2, 1, 1)

		s.Ctx.SetFontFace(s.Fonts.Size4x6)
		feelsLikeString := fmt.Sprintf("Feels %.f°C", s.CurrentWeather.FeelsLike)
		s.Ctx.DrawStringAnchored(feelsLikeString, 1, 13, 0, 1)
	}
	if fahrenheitOpacity > 0 {
		s.Ctx.SetColor(color.RGBA{255, 255, 255, fahrenheitOpacity})

		s.Ctx.SetFontFace(s.Fonts.Size8x13B)
		tempString := fmt.Sprintf("%.f°F", utils.GetFahrenheit(s.CurrentWeather.Temperature))
		s.Ctx.DrawStringAnchored(tempString, 33, -2, 1, 1)

		s.Ctx.SetFontFace(s.Fonts.Size4x6)
		feelsLikeString := fmt.Sprintf("Feels %.f°F", utils.GetFahrenheit(s.CurrentWeather.FeelsLike))
		s.Ctx.DrawStringAnchored(feelsLikeString, 1, 13, 0, 1)
	}

	filePath := fmt.Sprintf("./resources/weatherIcons/%d.png", s.CurrentWeather.WeatherCode)
	icon, err := utils.ReadImageFromFile(filePath)
	if err == nil {
		s.Ctx.DrawImageAnchored(icon, 57, 2, 1, 0)
	}

	s.renderPrecipitationGraph(celciusOpacity)
	s.renderUVIndexGraph(fahrenheitOpacity)

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *WeatherCurrentScreen) renderPrecipitationGraph(opacity uint8) {
	if opacity == 0 {
		return
	}

	s.Ctx.SetFontFace(s.Fonts.Size4x6)
	s.Ctx.SetColor(color.RGBA{255, 255, 255, opacity})
	s.Ctx.DrawStringAnchored("POP", 1, 29, 0, 0)

	s.Ctx.SetColor(color.RGBA{0x22, 0x22, 0x22, opacity})
	for x := range 48 {
		s.Ctx.SetPixel(14+x, 31)
	}

	currentHour := time.Now().Hour()

	for x, hourlyWeather := range s.HourlyWeather {
		if x < currentHour {
			s.Ctx.SetColor(color.RGBA{0x69, 0x6E, 0x71, opacity})
		} else if x > currentHour {
			s.Ctx.SetColor(color.RGBA{0, 0, 255, opacity})
		} else {
			// Show highlight
			s.Ctx.SetColor(color.RGBA{255, 0, 0, opacity})
			s.Ctx.SetPixel(14+x*2, 31)
			s.Ctx.SetPixel(15+x*2, 31)

			s.Ctx.SetColor(color.RGBA{0x97, 0xE7, 0xF5, opacity})
		}
		height := int(math.Round(hourlyWeather.PrecipitationProbability * 10))
		for y := range height {
			s.Ctx.SetPixel(14+x*2, 30-y)
			s.Ctx.SetPixel(15+x*2, 30-y)
		}
	}
}

func (s *WeatherCurrentScreen) renderUVIndexGraph(opacity uint8) {
	if opacity == 0 {
		return
	}

	s.Ctx.SetFontFace(s.Fonts.Size4x6)
	s.Ctx.SetColor(color.RGBA{255, 255, 255, opacity})
	s.Ctx.DrawStringAnchored("UV", 3, 29, 0, 0)

	s.Ctx.SetColor(color.RGBA{0x22, 0x22, 0x22, opacity})
	for x := range 48 {
		s.Ctx.SetPixel(14+x, 31)
	}

	currentHour := time.Now().Hour()

	for x, hourlyWeather := range s.HourlyWeather {
		if x < currentHour {
			s.Ctx.SetColor(color.RGBA{0x69, 0x6E, 0x71, opacity})
		} else if x > currentHour {
			s.Ctx.SetColor(color.RGBA{0xFC, 0xB4, 0x04, opacity})
		} else {
			// Show highlight
			s.Ctx.SetColor(color.RGBA{255, 0, 0, opacity})
			s.Ctx.SetPixel(14+x*2, 31)
			s.Ctx.SetPixel(15+x*2, 31)

			s.Ctx.SetColor(color.RGBA{0xFF, 0xF8, 0x04, opacity})
		}
		height := int(math.Round(hourlyWeather.UVIndex))
		if height > 10 {
			height = 10
		}
		for y := range height {
			s.Ctx.SetPixel(14+x*2, 30-y)
			s.Ctx.SetPixel(15+x*2, 30-y)
		}
	}
}
