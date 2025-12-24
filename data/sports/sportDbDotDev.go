package sports

import (
	"fmt"
	"log"
	"rpi-rgb-screen/config"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/utils"
)

type SportDbDotDevClient struct {
	BaseUrl string
}

type LeagueResponse struct {
	Seasons []struct {
		Season string `json:"season"`
	} `json:"seasons"`
}

type LeagueStandingsResponse struct {
	RoundType string                  `json:"roundType"`
	Teams     []TeamStandingsResponse `json:"teams"`
}

type TeamStandingsResponse struct {
	TeamId                         string `json:"teamId"`
	LossesAfterPenalties           string `json:"lossesAfterPenalties"`
	LossesInOvertimeAfterPenalties string `json:"lossesInOvertimeAfterPenalties"`
	LossesOvertime                 string `json:"lossesOvertime"`
	LossesRegular                  string `json:"lossesRegular"`
	Wins                           string `json:"wins"`
	WinsAfterPenalties             string `json:"winsAfterPenalties"`
	WinsInOvertimeAfterPenalties   string `json:"winsInOvertimeAfterPenalties"`
	WinsOvertime                   string `json:"winsOvertime"`
	WinsRegular                    string `json:"winsRegular"`
	Draws                          string `json:"draws"`
	Matches                        string `json:"matches"`
	Points                         string `json:"points"`
	Rank                           string `json:"rank"`
}

func NewSportDbDotDevClient() *SportDbDotDevClient {
	return &SportDbDotDevClient{
		BaseUrl: "https://api.sportdb.dev/api/flashscore/",
	}
}

func (s *SportDbDotDevClient) GetLeagueStandings(leagueId int, season string) []*Conference {
	url := fmt.Sprintf("%s%s/%s/standings", s.BaseUrl, constants.LEAGUE_STANDINGS_SLUGS[leagueId], season)
	headers := map[string]string{"X-API-Key": config.Config.SportDbDotDevApiKey}
	var leagueStandingsResponse []LeagueStandingsResponse

	// PWHL and NWSL are special cases where the API returns a single standings for the entire league
	if leagueId == constants.LEAGUE_NWSL || leagueId == constants.LEAGUE_PWHL {
		var teamStandingsResponse []TeamStandingsResponse
		err := utils.GetAndUnmarshal(url, headers, &teamStandingsResponse, nil)
		if err != nil {
			log.Printf("GetLeagueStandings failed. Error: %s", err)
			return []*Conference{}
		}
		leagueStandingsResponse = append(leagueStandingsResponse, LeagueStandingsResponse{
			RoundType: "League Standings",
			Teams:     teamStandingsResponse,
		})
	} else {
		err := utils.GetAndUnmarshal(url, headers, &leagueStandingsResponse, nil)
		if err != nil {
			log.Printf("GetLeagueStandings failed. Error: %s", err)
			return []*Conference{}
		}
	}

	if len(leagueStandingsResponse) == 0 {
		return []*Conference{}
	}
	if len(leagueStandingsResponse) > 2 {
		leagueStandingsResponse = leagueStandingsResponse[:2] // Only first two are conferences, rest are divisions
	}

	conferences := []*Conference{}
	for _, conference := range leagueStandingsResponse {
		standings := []*Standing{}
		for _, team := range conference.Teams {
			standing := &Standing{
				TeamId:           constants.TEAM_STANDINGS_IDS[team.TeamId],
				Rank:             utils.StringToInt(team.Rank),
				GamesPlayed:      utils.StringToInt(team.Matches),
				Points:           utils.StringToInt(team.Points),
				RegulationWins:   utils.StringToInt(team.WinsRegular),
				OvertimeWins:     utils.StringToInt(team.WinsOvertime) + utils.StringToInt(team.WinsAfterPenalties),
				RegulationLosses: utils.StringToInt(team.LossesRegular),
				OvertimeLosses:   utils.StringToInt(team.LossesOvertime) + utils.StringToInt(team.LossesAfterPenalties),
				Ties:             utils.StringToInt(team.Draws),
			}
			standings = append(standings, standing)
		}
		conferences = append(conferences, &Conference{
			Name:      conference.RoundType,
			Standings: standings,
		})
	}
	return conferences
}

func (s *SportDbDotDevClient) GetLatestSeason(leagueId int) string {
	url := fmt.Sprintf("%s%s", s.BaseUrl, constants.LEAGUE_STANDINGS_SLUGS[leagueId])
	headers := map[string]string{"X-API-Key": config.Config.SportDbDotDevApiKey}
	var leagueResponse LeagueResponse
	err := utils.GetAndUnmarshal(url, headers, &leagueResponse, nil)
	if err != nil {
		log.Printf("GetLatestSeason failed. Error: %s", err)
		return ""
	}
	if len(leagueResponse.Seasons) > 0 {
		return leagueResponse.Seasons[0].Season
	}
	return ""
}
