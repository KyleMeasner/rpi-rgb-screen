package weather

import "time"

type DummyWeatherDataManager struct{}

func (d *DummyWeatherDataManager) GetCurrentWeather(location string) *CurrentWeather {
	return &CurrentWeather{
		Temperature:              9.7,
		WeatherCode:              1000,
		PrecipitationProbability: 0,
		FeelsLike:                7.3,
	}
}

func (d *DummyWeatherDataManager) GetHourlyWeather(location string) []*HourlyWeather {
	return []*HourlyWeather{
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 9, WindGust: 17, Temperature: 10.2, FeelsLike: 10.2},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 7.7, WindGust: 16, Temperature: 9.8, FeelsLike: 9.8},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 6.9, WindGust: 13.9, Temperature: 10, FeelsLike: 10},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 4.1, WindGust: 10.8, Temperature: 9.9, FeelsLike: 9.9},
		{PrecipitationProbability: 20, UVIndex: 0, WindSpeed: 1, WindGust: 9.6, Temperature: 8.9, FeelsLike: 8.9},
		{PrecipitationProbability: 50, UVIndex: 0, WindSpeed: 3.7, WindGust: 12.8, Temperature: 9.1, FeelsLike: 9.1},
		{PrecipitationProbability: 60, UVIndex: 0, WindSpeed: 7.2, WindGust: 15.1, Temperature: 9.4, FeelsLike: 9.4},
		{PrecipitationProbability: 50, UVIndex: 0, WindSpeed: 8.1, WindGust: 14.9, Temperature: 9.5, FeelsLike: 9.5},
		{PrecipitationProbability: 10, UVIndex: 0, WindSpeed: 8.2, WindGust: 15.2, Temperature: 9.5, FeelsLike: 9.5},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 6.7, WindGust: 14.7, Temperature: 9.6, FeelsLike: 9.6},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 7.2, WindGust: 15.5, Temperature: 9.4, FeelsLike: 9.4},
		{PrecipitationProbability: 0, UVIndex: 1, WindSpeed: 7, WindGust: 13.8, Temperature: 9.7, FeelsLike: 9.7},
		{PrecipitationProbability: 0, UVIndex: 2, WindSpeed: 5.8, WindGust: 13, Temperature: 9.4, FeelsLike: 9.4},
		{PrecipitationProbability: 0, UVIndex: 3, WindSpeed: 6.4, WindGust: 13, Temperature: 10.4, FeelsLike: 10.4},
		{PrecipitationProbability: 0, UVIndex: 3, WindSpeed: 5.4, WindGust: 12.1, Temperature: 9.6, FeelsLike: 9.6},
		{PrecipitationProbability: 0, UVIndex: 3, WindSpeed: 5.5, WindGust: 11.4, Temperature: 9.5, FeelsLike: 9.5},
		{PrecipitationProbability: 0, UVIndex: 1, WindSpeed: 5.6, WindGust: 12.5, Temperature: 9, FeelsLike: 9},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 5.2, WindGust: 12.3, Temperature: 8.6, FeelsLike: 8.6},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 4.9, WindGust: 11.6, Temperature: 8.3, FeelsLike: 8.3},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 3.1, WindGust: 8.9, Temperature: 7.8, FeelsLike: 7.8},
		{PrecipitationProbability: 0, UVIndex: 0, WindSpeed: 2, WindGust: 5.6, Temperature: 7, FeelsLike: 7},
		{PrecipitationProbability: 20, UVIndex: 0, WindSpeed: 2.6, WindGust: 5.4, Temperature: 6.7, FeelsLike: 6.7},
		{PrecipitationProbability: 50, UVIndex: 0, WindSpeed: 3.1, WindGust: 6.3, Temperature: 7, FeelsLike: 7},
		{PrecipitationProbability: 100, UVIndex: 0, WindSpeed: 3.3, WindGust: 6.4, Temperature: 7.1, FeelsLike: 7.1},
	}
}

func (d *DummyWeatherDataManager) GetForecast(location string) []*WeatherForecast {
	forecast := []*WeatherForecast{
		{TemperatureMin: 6.7, TemperatureMax: 10.4, WeatherCode: 4205},
		{TemperatureMin: 8.1, TemperatureMax: 11.3, WeatherCode: 4205},
		{TemperatureMin: 8.5, TemperatureMax: 12.6, WeatherCode: 4200},
		{TemperatureMin: 8.1, TemperatureMax: 12.6, WeatherCode: 4210},
		{TemperatureMin: 9.1, TemperatureMax: 12.5, WeatherCode: 4200},
		{TemperatureMin: 12.4, TemperatureMax: 13.7, WeatherCode: 1001},
	}

	for i, f := range forecast {
		now := time.Now()
		f.Date = time.Date(now.Year(), now.Month(), now.Day()+i, 6, 0, 0, 0, time.Local)
		f.SunriseTime = time.Date(now.Year(), now.Month(), now.Day()+i, 7, 27+i, 0, 0, time.Local)
		f.SunsetTime = time.Date(now.Year(), now.Month(), now.Day()+i, 16, 34+i, 0, 0, time.Local)
	}

	return forecast
}
