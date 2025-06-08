package sports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"net/url"
)

const baseUrl = "https://www.thesportsdb.com/api/v1/json/123"

type TheSportsDbClient interface {
	GetLogo(logoUrl string) image.Image
	GetTeam(teamName string) *Team
	GetUpcomingEventsForLeague(league int) []Event
}

type TheSportsDb struct{}

func NewTheSportsDbClient() TheSportsDbClient {
	return &TheSportsDb{}
}

func (t *TheSportsDb) GetLogo(logoUrl string) image.Image {
	badgeUrl := logoUrl + "/tiny"
	logoBytes, err := sendGetRequest(badgeUrl)
	if err != nil {
		log.Printf("Logo fetch failed. Error: %s", err)
		return nil
	}

	logoImage, _, err := image.Decode(bytes.NewReader(logoBytes))
	if err != nil {
		log.Printf("Logo fetch failed. Error decoding logo bytes: %s", err)
		return nil
	}

	return logoImage
}

type TeamSearchResponse struct {
	Teams []Team `json:"teams"`
}

func (t *TheSportsDb) GetTeam(teamName string) *Team {
	url := fmt.Sprintf("%s/searchteams.php?t=%s", baseUrl, url.QueryEscape(teamName))
	responseBody, err := sendGetRequest(url)
	if err != nil {
		log.Printf("Team fetch failed for team name %s. Error: %s", teamName, err)
		return nil
	}

	var teamSearchResponse TeamSearchResponse
	err = json.Unmarshal(responseBody, &teamSearchResponse)
	if err != nil {
		log.Printf("Team fetch failed for team name %s. Error: %s", teamName, err)
		return nil
	}
	if len(teamSearchResponse.Teams) < 1 {
		log.Printf("Team fetch failed for team name %s. No results found.", teamName)
		return nil
	}

	return &teamSearchResponse.Teams[0]
}

type EventSearchResponse struct {
	Events []Event `json:"events"`
}

func (t *TheSportsDb) GetUpcomingEventsForLeague(league int) []Event {
	url := fmt.Sprintf("%s/eventsseason.php?id=%d&s=2025", baseUrl, league)
	responseBody, err := sendGetRequest(url)
	if err != nil {
		log.Printf("Upcoming event fetch failed for league %d. Error: %s", league, err)
		return nil
	}

	var eventSearchResponse EventSearchResponse
	err = json.Unmarshal(responseBody, &eventSearchResponse)
	if err != nil {
		log.Printf("Upcoming event fetch failed for league %d. Error: %s", league, err)
		return nil
	}

	return eventSearchResponse.Events
}

func sendGetRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}
