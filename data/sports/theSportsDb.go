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
	"strconv"
	"time"
)

const baseUrl = "https://www.thesportsdb.com/api/v1/json/123"

type TheSportsDbClient struct{}

type LeagueSearchResponse struct {
	Leagues []struct {
		Id            string `json:"idLeague"`
		Name          string `json:"strLeague"`
		CurrentSeason string `json:"strCurrentSeason"`
		BadgeUrl      string `json:"strBadge"`
	} `json:"leagues"`
}

type TeamSearchResponse struct {
	Teams []struct {
		Id        string `json:"idTeam"`
		Name      string `json:"strTeam"`
		ShortName string `json:"strTeamShort"`
		BadgeUrl  string `json:"strBadge"`
	} `json:"teams"`
}

type EventSearchResponse struct {
	Events []struct {
		Id           string `json:"idEvent"`
		Name         string `json:"strEvent"`
		HomeTeamName string `json:"strHomeTeam"`
		AwayTeamName string `json:"strAwayTeam"`
		Timestamp    string `json:"strTimestamp"`
	} `json:"events"`
}

func NewTheSportsDbClient() *TheSportsDbClient {
	return &TheSportsDbClient{}
}

func (t *TheSportsDbClient) GetLogo(logoUrl string) image.Image {
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

func (t *TheSportsDbClient) GetLeague(leagueId int) *League {
	var leagueSearchResponse LeagueSearchResponse
	url := fmt.Sprintf("%s/lookupleague.php?id=%d", baseUrl, leagueId)
	err := getAndUnmarshal(url, &leagueSearchResponse)
	if err != nil {
		log.Printf("League fetch failed for league ID %d. Error: %s", leagueId, err)
		return nil
	}

	if len(leagueSearchResponse.Leagues) < 1 {
		log.Printf("League fetch failed for league ID %d. Error: %s", leagueId, err)
		return nil
	}

	rawLeague := leagueSearchResponse.Leagues[0]
	return &League{
		Id:            leagueId,
		Name:          rawLeague.Name,
		CurrentSeason: rawLeague.CurrentSeason,
		BadgeUrl:      rawLeague.BadgeUrl,
	}
}

func (t *TheSportsDbClient) GetTeam(teamName string) *Team {
	var teamSearchResponse TeamSearchResponse
	url := fmt.Sprintf("%s/searchteams.php?t=%s", baseUrl, url.QueryEscape(teamName))
	err := getAndUnmarshal(url, &teamSearchResponse)
	if err != nil {
		log.Printf("Team fetch failed for team name %s. Error: %s", teamName, err)
		return nil
	}

	if len(teamSearchResponse.Teams) < 1 {
		log.Printf("Team fetch failed for team name %s. No results found.", teamName)
		return nil
	}

	rawTeam := teamSearchResponse.Teams[0]
	teamId, err := strconv.Atoi(rawTeam.Id)
	if err != nil {
		log.Printf("Team fetch failed for team name %s. Error: %s", teamName, err)
		return nil
	}

	return &Team{
		Id:        teamId,
		Name:      rawTeam.Name,
		ShortName: rawTeam.ShortName,
		BadgeUrl:  rawTeam.BadgeUrl,
	}
}

func (t *TheSportsDbClient) GetNextGameForTeam(teamId int) *Event {
	var eventSearchResponse EventSearchResponse
	url := fmt.Sprintf("%s/eventsnext.php?id=%d", baseUrl, teamId)
	err := getAndUnmarshal(url, &eventSearchResponse)
	if err != nil {
		log.Printf("GetNextGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}
	if len(eventSearchResponse.Events) < 1 {
		log.Printf("GetNextGameForTeam failed. Team ID %d. No results found.", teamId)
		return nil
	}

	rawEvent := eventSearchResponse.Events[0]
	eventId, err := strconv.Atoi(rawEvent.Id)
	if err != nil {
		log.Printf("GetNextGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}
	eventTime, err := time.Parse("2006-01-02T15:04:05", rawEvent.Timestamp)
	if err != nil {
		log.Printf("GetNextGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}

	return &Event{
		Id:           eventId,
		Name:         rawEvent.Name,
		HomeTeamName: rawEvent.HomeTeamName,
		AwayTeamName: rawEvent.AwayTeamName,
		Time:         eventTime.Local(),
	}
}

func sendGetRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}

func getAndUnmarshal(url string, responseObject any) error {
	responseBody, err := sendGetRequest(url)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, responseObject)
}
