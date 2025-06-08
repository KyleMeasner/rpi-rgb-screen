package data

import "rpi-rgb-screen/data/sports"

type DataManager struct {
	SportsData sports.SportsData
}

func NewDataManager() *DataManager {
	return &DataManager{
		SportsData: sports.NewSportsData(),
	}
}
