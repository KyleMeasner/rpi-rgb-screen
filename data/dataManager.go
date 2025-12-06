package data

import (
	"rpi-rgb-screen/data/sports"
	"rpi-rgb-screen/data/weather"
)

type DataManager struct {
	SportsData  sports.SportsData
	WeatherData weather.WeatherData
}

func NewDataManager(useDummyData bool) *DataManager {
	return &DataManager{
		SportsData:  sports.NewSportsData(),
		WeatherData: weather.NewWeatherData(useDummyData),
	}
}
