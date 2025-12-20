package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	TomorrowIoApiKey string `json:"tomorrowIoApiKey"`
	HockeyTechApiKey string `json:"hockeyTechApiKey"`
	Location         string `json:"location"`
	FavoriteTeams    []int  `json:"favoriteTeams"`
	UseDummyData     bool   `json:"useDummyData"`
}

var Config *Configuration = &Configuration{}

func LoadConfig() error {
	configFile, err := os.ReadFile("./config.json")
	if err != nil {
		log.Printf("Failed to load config. Error: %s", err)
		return err
	}

	err = json.Unmarshal(configFile, Config)
	if err != nil {
		log.Printf("Failed to load config. Error: %s", err)
		return err
	}

	return nil
}
