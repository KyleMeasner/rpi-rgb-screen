package weather

import "time"

type CurrentWeather struct {
	Temperature float64
	WeatherCode int
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
	GetForecast(location string) []*WeatherForecast
}

type WeatherDataManager struct {
	TomorrowIoClient *TomorrowIoClient
	CurrentWeather   map[string]*CurrentWeather
	Forecast         map[string][]*WeatherForecast
}

func NewWeatherData() WeatherData {
	return &WeatherDataManager{
		TomorrowIoClient: NewTomorrowIoClient(),
		CurrentWeather:   map[string]*CurrentWeather{},
		Forecast:         map[string][]*WeatherForecast{},
	}
}

func (w *WeatherDataManager) GetCurrentWeather(location string) *CurrentWeather {
	if currentWeather, ok := w.CurrentWeather[location]; ok {
		return currentWeather
	}

	currentWeather := w.TomorrowIoClient.GetCurrentWeather(location)
	if currentWeather == nil {
		return nil
	}

	w.CurrentWeather[location] = currentWeather
	return currentWeather
}

func (w *WeatherDataManager) GetForecast(location string) []*WeatherForecast {
	if forecast, ok := w.Forecast[location]; ok {
		return forecast
	}

	forecast := w.TomorrowIoClient.GetForecast(location)
	if forecast == nil {
		return nil
	}

	w.Forecast[location] = forecast
	return forecast
}
