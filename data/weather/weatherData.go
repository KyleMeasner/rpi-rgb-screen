package weather

import (
	"rpi-rgb-screen/utils"
	"time"
)

type CurrentWeather struct {
	Temperature              float64
	WeatherCode              int
	PrecipitationProbability int
	FeelsLike                float64
}

type HourlyWeather struct {
	PrecipitationProbability int
	UVIndex                  int
	WindSpeed                float64
	WindGust                 float64
	Temperature              float64
	FeelsLike                float64
}

type WeatherForecast struct {
	Date           time.Time
	TemperatureMin float64
	TemperatureMax float64
	WeatherCode    int
	SunriseTime    time.Time
	SunsetTime     time.Time
}

type WeatherData interface {
	GetCurrentWeather(location string) *CurrentWeather
	GetHourlyWeather(location string) []*HourlyWeather
	GetForecast(location string) []*WeatherForecast
}

type WeatherDataManager struct {
	TomorrowIoClient *TomorrowIoClient
	CurrentWeather   *utils.ExpirableMap[string, *CurrentWeather]
	HourlyWeather    *utils.ExpirableMap[string, []*HourlyWeather]
	Forecast         *utils.ExpirableMap[string, []*WeatherForecast]
}

func NewWeatherData() WeatherData {
	return &WeatherDataManager{
		TomorrowIoClient: NewTomorrowIoClient(),
		CurrentWeather:   utils.NewExpirableMap[string, *CurrentWeather](time.Minute * 15),
		HourlyWeather:    utils.NewExpirableMap[string, []*HourlyWeather](time.Minute * 15),
		Forecast:         utils.NewExpirableMap[string, []*WeatherForecast](time.Minute * 15),
	}
}

func (w *WeatherDataManager) GetCurrentWeather(location string) *CurrentWeather {
	cachedCurrentWeather := w.CurrentWeather.Get(location)
	if cachedCurrentWeather != nil {
		return *cachedCurrentWeather
	}

	currentWeather := w.TomorrowIoClient.GetCurrentWeather(location)
	if currentWeather == nil {
		return nil
	}

	w.CurrentWeather.Set(location, currentWeather)
	return currentWeather
}

func (w *WeatherDataManager) GetHourlyWeather(location string) []*HourlyWeather {
	cachedHourlyWeather := w.HourlyWeather.Get(location)
	if cachedHourlyWeather != nil {
		return *cachedHourlyWeather
	}

	hourlyWeather := w.TomorrowIoClient.GetHourlyWeather(location)
	if len(hourlyWeather) < 24 {
		return nil
	}

	w.HourlyWeather.Set(location, hourlyWeather)
	return hourlyWeather
}

func (w *WeatherDataManager) GetForecast(location string) []*WeatherForecast {
	cachedForecast := w.Forecast.Get(location)
	if cachedForecast != nil {
		return *cachedForecast
	}

	forecast := w.TomorrowIoClient.GetForecast(location)
	if forecast == nil {
		return nil
	}

	w.Forecast.Set(location, forecast)
	return forecast
}
