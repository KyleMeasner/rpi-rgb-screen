package sports

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"rpi-rgb-screen/config"
	"rpi-rgb-screen/utils"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type HockeyTechClient struct {
	RateLimiter *rate.Limiter
}

type HTLeagueScheduleResponse struct {
	SiteKit struct {
		Schedule []struct {
			ID                string `json:"id"`
			GameDate          string `json:"GameDateISO8601"`
			HomeTeam          string `json:"home_team"`
			VisitingTeam      string `json:"visiting_team"`
			HomeGoalCount     string `json:"home_goal_count"`     // "-" for games that have not been played yet
			VisitingGoalCount string `json:"visiting_goal_count"` // "-" for games that have not been played yet
			GameStatus        string `json:"game_status"`         // "6:00 pm EST", "TBD", "Final OT", "Final"
		} `json:"Schedule"`
	} `json:"SiteKit"`
}

func NewHockeyTechClient() *HockeyTechClient {
	client := &HockeyTechClient{
		RateLimiter: rate.NewLimiter(0.16, 1), // 10 requests per minute (don't know what the actual limit is)
	}
	client.GetLeagueSchedule()
	return client
}

func (h *HockeyTechClient) GetLogo(teamId int) image.Image {
	filePath := fmt.Sprintf("./resources/pwhlLogos/%d.png", teamId)
	logo, err := utils.ReadImageFromFile(filePath)
	if err != nil {
		log.Printf("Failed to read logo from file for teamId %d. Error: %s", teamId, err)
		return nil
	}
	return logo
}

func (h *HockeyTechClient) GetLeagueSchedule() ([]*Event, []*Event) {
	url := fmt.Sprintf("https://lscluster.hockeytech.com/feed/?feed=modulekit&view=schedule&key=%s&client_code=pwhl", config.Config.HockeyTechApiKey)
	responseBody, err := utils.SendGetRequest(url, nil, h.RateLimiter)
	if err != nil {
		log.Printf("GetLeagueSchedule failed. Error: %s", err)
		return nil, nil
	}

	var scheduleResponse HTLeagueScheduleResponse
	err = json.Unmarshal(responseBody, &scheduleResponse)
	if err != nil {
		log.Printf("GetLeagueSchedule JSON unmarshal failed. Error: %s", err)
		return nil, nil
	}

	if len(scheduleResponse.SiteKit.Schedule) == 0 {
		return []*Event{}, []*Event{}
	}

	pastEvents := []*Event{}
	upcomingEvents := []*Event{}
	for _, rawEvent := range scheduleResponse.SiteKit.Schedule {
		id, err := strconv.Atoi(rawEvent.ID)
		if err != nil {
			log.Printf("GetLeagueSchedule failed to parse game ID %s. Error: %s", rawEvent.ID, err)
			continue
		}

		homeTeamId, err := strconv.Atoi(rawEvent.HomeTeam)
		if err != nil {
			log.Printf("GetLeagueSchedule failed to parse home team ID %s. Error: %s", rawEvent.HomeTeam, err)
			continue
		}
		awayTeamId, err := strconv.Atoi(rawEvent.VisitingTeam)
		if err != nil {
			log.Printf("GetLeagueSchedule failed to parse away team ID %s. Error: %s", rawEvent.VisitingTeam, err)
			continue
		}

		gameTime, err := time.Parse("2006-01-02T15:04:05-07:00", rawEvent.GameDate)
		if err != nil {
			log.Printf("GetLeagueSchedule failed to parse date %s. Error: %s", rawEvent.GameDate, err)
			continue
		}

		event := &Event{
			Id:         id,
			LeagueId:   999,
			HomeTeamId: homeTeamId,
			AwayTeamId: awayTeamId,
			Time:       gameTime.Local(),
			IsTBD:      rawEvent.GameStatus == "TBD",
		}

		if rawEvent.GameStatus == "Final" || rawEvent.GameStatus == "Final OT" {
			homeGoals, err := strconv.Atoi(rawEvent.HomeGoalCount)
			if err != nil {
				log.Printf("GetLeagueSchedule failed to parse home goals %s. Error: %s", rawEvent.HomeGoalCount, err)
				continue
			}
			awayGoals, err := strconv.Atoi(rawEvent.VisitingGoalCount)
			if err != nil {
				log.Printf("GetLeagueSchedule failed to parse away goals %s. Error: %s", rawEvent.VisitingGoalCount, err)
				continue
			}
			event.HomeScore = homeGoals
			event.AwayScore = awayGoals

			pastEvents = append(pastEvents, event)
		} else {
			upcomingEvents = append(upcomingEvents, event)
		}
	}

	return pastEvents, upcomingEvents
}
