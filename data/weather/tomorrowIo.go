package weather

import (
	"fmt"
	"log"
	"net/url"
	"rpi-rgb-screen/config"
	"rpi-rgb-screen/utils"
	"time"

	"golang.org/x/time/rate"
)

// Rate limits:
// 500 requests per day
// 25 requests per hour
// 3 requests per second
// https://support.tomorrow.io/hc/en-us/articles/20273728362644-Free-API-Plan-Rate-Limits

const baseUrl = "https://api.tomorrow.io/v4"

type TomorrowIoClient struct {
	RateLimiter *rate.Limiter
}

type RealtimeWeatherResponse struct {
	Data struct {
		Values struct {
			Temperature              float64 `json:"temperature"`
			WeatherCode              int     `json:"weatherCode"`
			PrecipitationProbability int     `json:"precipitationProbability"`
			TemperatureApparent      float64 `json:"temperatureApparent"`
		} `json:"values"`
	} `json:"data"`
}

type WeatherForecastResponse struct {
	Timelines struct {
		Daily []struct {
			Time   string          `json:"time"`
			Values WeatherForecast `json:"values"`
		} `json:"daily"`
	} `json:"timelines"`
}

type WeatherTimelinesRequest struct {
	Location       string   `json:"location"`
	Fields         []string `json:"fields"`
	Units          string   `json:"units"`
	Timesteps      []string `json:"timesteps"`
	StartTime      string   `json:"startTime"`
	EndTime        string   `json:"endTime"`
	DailyStartHour int      `json:"dailyStartHour"`
}

type WeatherTimelinesResponse struct {
	Data struct {
		Timelines []struct {
			Intervals []struct {
				StartTime string `json:"startTime"`
				Values    struct {
					TemperatureMin float64 `json:"temperatureMin"`
					TemperatureMax float64 `json:"temperatureMax"`
					WeatherCode    int     `json:"weatherCodeFullDay"`
					SunriseTime    string  `json:"sunriseTime"`
					SunsetTime     string  `json:"sunsetTime"`
				} `json:"values"`
			} `json:"intervals"`
		} `json:"timelines"`
	} `json:"data"`
}

func NewTomorrowIoClient() *TomorrowIoClient {
	return &TomorrowIoClient{
		RateLimiter: rate.NewLimiter(3, 1), // 3 requests per second max
	}
}

func (t *TomorrowIoClient) GetCurrentWeather(location string) *CurrentWeather {
	url := fmt.Sprintf("%s/weather/realtime?apikey=%s&units=metric&location=%s", baseUrl, config.Config.TomorrowIoApiKey, url.QueryEscape(location))

	var realtimeWeatherResponse RealtimeWeatherResponse
	err := utils.GetAndUnmarshal(url, &realtimeWeatherResponse, t.RateLimiter)
	if err != nil {
		log.Printf("Failed to get current weather for location '%s'. Error: %s", location, err)
		return nil
	}

	return &CurrentWeather{
		Temperature: realtimeWeatherResponse.Data.Values.Temperature,
		WeatherCode: realtimeWeatherResponse.Data.Values.WeatherCode,
	}
}

func (t *TomorrowIoClient) GetForecast(location string) []*WeatherForecast {
	url := fmt.Sprintf("%s/timelines?apikey=%s", baseUrl, config.Config.TomorrowIoApiKey)

	requestPayload := WeatherTimelinesRequest{
		Location:       location,
		Fields:         []string{"temperatureMin", "temperatureMax", "weatherCodeFullDay", "sunriseTime", "sunsetTime"},
		Units:          "metric",
		Timesteps:      []string{"1d"},
		StartTime:      "now",
		EndTime:        "nowPlus5d",
		DailyStartHour: 6,
	}

	var weatherTimelinesResponse WeatherTimelinesResponse
	err := utils.PostAndUnmarshal(url, "application/json", requestPayload, &weatherTimelinesResponse, t.RateLimiter)
	if err != nil {
		log.Printf("Failed to get weather timeline for location '%s'. Error: %s", location, err)
		return nil
	}
	if len(weatherTimelinesResponse.Data.Timelines) == 0 {
		log.Printf("Failed to get weather timeline for location '%s'. No results found.", location)
		return nil
	}

	weatherForecast := []*WeatherForecast{}
	for _, interval := range weatherTimelinesResponse.Data.Timelines[0].Intervals {
		date, _ := time.Parse("2006-01-02T15:04:05Z", interval.StartTime)
		sunrise, _ := time.Parse("2006-01-02T15:04:05Z", interval.Values.SunriseTime)
		sunset, _ := time.Parse("2006-01-02T15:04:05Z", interval.Values.SunsetTime)

		weatherForecast = append(weatherForecast, &WeatherForecast{
			Date:           date.Local(),
			TemperatureMin: interval.Values.TemperatureMin,
			TemperatureMax: interval.Values.TemperatureMax,
			WeatherCode:    interval.Values.WeatherCode,
			SunriseTime:    sunrise.Local(),
			SunsetTime:     sunset.Local(),
		})
	}

	return weatherForecast
}
