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

type WeatherForecastScreen struct {
	State           ScreenState
	Ctx             *gg.Context
	KeyFrames       *animation.KeyFrames
	Fonts           *fonts.Fonts
	WeatherData     weather.WeatherData
	WeatherForecast []*weather.WeatherForecast
}

func NewWeatherForecastScreen(fonts *fonts.Fonts, weatherData weather.WeatherData) Screen {
	keyFrames := animation.NewKeyFrames(30 * time.Second)

	keyFrames.AddNumber("C", 255)
	keyFrames.AddNumberTransitions("C", animation.AnimatedNumberTransition{Offset: 14500 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 0})

	keyFrames.AddNumber("F", 0)
	keyFrames.AddNumberTransitions("F", animation.AnimatedNumberTransition{Offset: 14500 * time.Millisecond, Duration: 1000 * time.Millisecond, EndValue: 255})

	return &WeatherForecastScreen{
		State:       StateNotDisplayed,
		Ctx:         gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:   keyFrames,
		Fonts:       fonts,
		WeatherData: weatherData,
	}
}

func (s *WeatherForecastScreen) GetPreferredTransition() transition.Transition {
	return transition.NewSlideAndZoomTransition()
}

func (s *WeatherForecastScreen) SetState(state ScreenState) {
	s.State = state

	switch state {
	case StateDisplayed:
		s.KeyFrames.Start()
	case StateTransitionIn:
		s.KeyFrames.Reset()
	}
}

func (s *WeatherForecastScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		forecast := s.WeatherData.GetForecast(config.Config.Location)
		if len(forecast) >= 4 {
			s.WeatherForecast = forecast[1:4]
		} else if len(forecast) > 0 {
			s.WeatherForecast = forecast[1:]
		} else {
			s.WeatherForecast = []*weather.WeatherForecast{}
		}
		close(doneChan)
	}()

	return doneChan
}

func (s *WeatherForecastScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	s.Ctx.SetFontFace(s.Fonts.Size4x6)

	for i, forecast := range s.WeatherForecast {
		x := 11 + 21*i

		dayString := forecast.Date.Format("Mon")
		s.Ctx.SetColor(color.White)
		s.Ctx.DrawStringAnchored(dayString, float64(x), 1, 0.5, 1)

		filePath := fmt.Sprintf("./resources/weatherIcons/%d.png", forecast.WeatherCode)
		icon, err := utils.ReadImageFromFile(filePath)
		if err == nil {
			s.Ctx.DrawImageAnchored(icon, x, 8, 0.5, 0)
		}

		celciusOpacity := uint8(s.KeyFrames.GetNumber("C"))
		fahrenheitOpacity := uint8(s.KeyFrames.GetNumber("F"))

		// High temp
		r, g, b, _ := colorForTemp(forecast.TemperatureMax).RGBA()
		if celciusOpacity > 0 {
			highTempCelsius := fmt.Sprintf("%.f", forecast.TemperatureMax)
			s.Ctx.SetColor(color.RGBA{uint8(r), uint8(g), uint8(b), celciusOpacity})
			s.Ctx.DrawStringAnchored(highTempCelsius, float64(x-4), 18, 0.5, 1)
		}
		if fahrenheitOpacity > 0 {
			highTempFahrenheit := fmt.Sprintf("%.f", utils.GetFahrenheit(forecast.TemperatureMax))
			s.Ctx.SetColor(color.RGBA{uint8(r), uint8(g), uint8(b), fahrenheitOpacity})
			s.Ctx.DrawStringAnchored(highTempFahrenheit, float64(x-4), 18, 0.5, 1)
		}

		// Low temp
		r, g, b, _ = colorForTemp(forecast.TemperatureMin).RGBA()
		if celciusOpacity > 0 {
			lowTempCelsius := fmt.Sprintf("%.f", forecast.TemperatureMin)
			s.Ctx.SetColor(color.RGBA{uint8(r), uint8(g), uint8(b), celciusOpacity})
			s.Ctx.DrawStringAnchored(lowTempCelsius, float64(x+5), 18, 0.5, 1)
		}
		if fahrenheitOpacity > 0 {
			lowTempFahrenheit := fmt.Sprintf("%.f", utils.GetFahrenheit(forecast.TemperatureMin))
			s.Ctx.SetColor(color.RGBA{uint8(r), uint8(g), uint8(b), fahrenheitOpacity})
			s.Ctx.DrawStringAnchored(lowTempFahrenheit, float64(x+5), 18, 0.5, 1)
		}

		if celciusOpacity > 0 {
			s.Ctx.SetColor(color.RGBA{255, 255, 255, celciusOpacity})
			s.Ctx.DrawStringAnchored("°C", 64, 32, 1, 0)
		}
		if fahrenheitOpacity > 0 {
			s.Ctx.SetColor(color.RGBA{255, 255, 255, fahrenheitOpacity})
			s.Ctx.DrawStringAnchored("°F", 64, 32, 1, 0)
		}
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func colorForTemp(temperature float64) color.Color {
	if temperature > 35 {
		temperature = 35
	} else if temperature < -5 {
		temperature = -5
	}

	percentComplete := (temperature + 5) / 20 // Range of -5 to 15
	start := color.RGBA{0, 0, 0xFF, 0xFF}     // #0000FF
	end := color.RGBA{0xBB, 0xBB, 0xFF, 0xFF} // #BBBBFF

	if temperature >= 15 {
		percentComplete = (temperature - 15) / 20  // Range of 15 to 35
		start = color.RGBA{0xFF, 0xBB, 0xBB, 0xFF} // #FFBBBB
		end = color.RGBA{0xFF, 0, 0, 0xFF}         // #FF0000
	}

	return color.RGBA{
		R: uint8(utils.ComputeValue(int(start.R), int(end.R), percentComplete)),
		G: uint8(utils.ComputeValue(int(start.G), int(end.G), percentComplete)),
		B: uint8(utils.ComputeValue(int(start.B), int(end.B), percentComplete)),
		A: uint8(utils.ComputeValue(int(start.A), int(end.A), percentComplete)),
	}
}
