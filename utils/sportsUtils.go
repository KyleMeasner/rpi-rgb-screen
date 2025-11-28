package utils

import "rpi-rgb-screen/constants"

func GetFavoriteLeagues(favoriteTeams []int) []int {
	favoriteLeagues := map[int]bool{}
	for leagueId, teams := range constants.LEAGUE_TEAMS {
		for _, teamId := range teams {
			for _, favoriteTeamId := range favoriteTeams {
				if teamId == favoriteTeamId {
					favoriteLeagues[leagueId] = true
				}
			}
		}
	}

	leagues := []int{}
	for leagueId := range favoriteLeagues {
		leagues = append(leagues, leagueId)
	}
	return leagues
}
