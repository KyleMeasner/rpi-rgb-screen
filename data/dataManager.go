package data

import (
	"rpi-rgb-screen/data/sports"
	"rpi-rgb-screen/data/transit"
	"rpi-rgb-screen/data/weather"
)

type DataManager struct {
	SportsData  sports.SportsData
	WeatherData weather.WeatherData
	TransitData transit.TransitData
}

func NewDataManager(useDummyData bool) *DataManager {
	return &DataManager{
		SportsData:  sports.NewSportsData(useDummyData),
		WeatherData: weather.NewWeatherData(useDummyData),
		TransitData: transit.NewTransitData(useDummyData),
	}
}
