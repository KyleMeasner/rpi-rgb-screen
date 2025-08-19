package sports

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"net/url"
	"rpi-rgb-screen/utils"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

// Rate limits:
// 30 requests per minute
// https://www.thesportsdb.com/documentation#rate_limit

const baseUrl = "https://www.thesportsdb.com/api/v1/json/123"

type TheSportsDbClient struct {
	RateLimiter *rate.Limiter
}

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

type NextEventSearchResponse struct {
	Events []struct {
		Id           string `json:"idEvent"`
		Name         string `json:"strEvent"`
		HomeTeamName string `json:"strHomeTeam"`
		AwayTeamName string `json:"strAwayTeam"`
		Timestamp    string `json:"strTimestamp"`
	} `json:"events"`
}

type LastEventSearchResponse struct {
	Events []struct {
		Id           string `json:"idEvent"`
		Name         string `json:"strEvent"`
		HomeTeamName string `json:"strHomeTeam"`
		AwayTeamName string `json:"strAwayTeam"`
		Timestamp    string `json:"strTimestamp"`
		HomeScore    string `json:"intHomeScore"`
		AwayScore    string `json:"intAwayScore"`
	} `json:"results"`
}

func NewTheSportsDbClient() *TheSportsDbClient {
	return &TheSportsDbClient{
		RateLimiter: rate.NewLimiter(0.5, 1), // 30 requests per minute
	}
}

func (t *TheSportsDbClient) GetLogo(logoUrl string) image.Image {
	badgeUrl := logoUrl + "/tiny"
	logoBytes, err := utils.SendGetRequest(badgeUrl, nil) // Not a rate limited API
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
	err := utils.GetAndUnmarshal(url, &leagueSearchResponse, t.RateLimiter)
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
		LogoUrl:       rawLeague.BadgeUrl,
	}
}

func (t *TheSportsDbClient) GetTeam(teamName string) *Team {
	var teamSearchResponse TeamSearchResponse
	url := fmt.Sprintf("%s/searchteams.php?t=%s", baseUrl, url.QueryEscape(teamName))
	err := utils.GetAndUnmarshal(url, &teamSearchResponse, t.RateLimiter)
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
		LogoUrl:   rawTeam.BadgeUrl,
	}
}

func (t *TheSportsDbClient) GetNextGameForTeam(teamId int) *Event {
	var nextEventSearchResponse NextEventSearchResponse
	url := fmt.Sprintf("%s/eventsnext.php?id=%d", baseUrl, teamId)
	err := utils.GetAndUnmarshal(url, &nextEventSearchResponse, t.RateLimiter)
	if err != nil {
		log.Printf("GetNextGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}
	if len(nextEventSearchResponse.Events) < 1 {
		log.Printf("GetNextGameForTeam failed. Team ID %d. No results found.", teamId)
		return nil
	}

	rawEvent := nextEventSearchResponse.Events[0]
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

func (t *TheSportsDbClient) GetLastGameForTeam(teamId int) *Event {
	var lastEventSearchResponse LastEventSearchResponse
	url := fmt.Sprintf("%s/eventslast.php?id=%d", baseUrl, teamId)
	err := utils.GetAndUnmarshal(url, &lastEventSearchResponse, t.RateLimiter)
	if err != nil {
		log.Printf("GetLastGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}
	if len(lastEventSearchResponse.Events) < 1 {
		log.Printf("GetLastGameForTeam failed. Team ID %d. No results found.", teamId)
		return nil
	}

	rawEvent := lastEventSearchResponse.Events[0]
	eventId, err := strconv.Atoi(rawEvent.Id)
	if err != nil {
		log.Printf("GetLastGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}
	eventTime, err := time.Parse("2006-01-02T15:04:05", rawEvent.Timestamp)
	if err != nil {
		log.Printf("GetLastGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}
	homeScore, err := strconv.Atoi(rawEvent.HomeScore)
	if err != nil {
		log.Printf("GetLastGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}
	awayScore, err := strconv.Atoi(rawEvent.AwayScore)
	if err != nil {
		log.Printf("GetLastGameForTeam failed. Team ID %d. Error: %s", teamId, err)
		return nil
	}

	return &Event{
		Id:           eventId,
		Name:         rawEvent.Name,
		HomeTeamName: rawEvent.HomeTeamName,
		AwayTeamName: rawEvent.AwayTeamName,
		Time:         eventTime.Local(),
		HomeScore:    homeScore,
		AwayScore:    awayScore,
	}
}
