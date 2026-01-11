package sports

import (
	"image"
	"rpi-rgb-screen/constants"
)

type DummySportsDataManager struct {
	PassThrough SportsData
}

func NewDummySportsData() SportsData {
	return &DummySportsDataManager{
		PassThrough: NewSportsData(false), // Call through for unimplemented functions
	}
}

func (d *DummySportsDataManager) GetUpcomingEventsForLeague(leagueId int, onlyCacheFavoriteTeams bool) []*Event {
	return d.PassThrough.GetUpcomingEventsForLeague(leagueId, onlyCacheFavoriteTeams)
}

func (d *DummySportsDataManager) GetPastEventsForLeague(leagueId int, onlyCacheFavoriteTeams bool) []*Event {
	return d.PassThrough.GetPastEventsForLeague(leagueId, onlyCacheFavoriteTeams)
}

func (d *DummySportsDataManager) GetLeague(leagueId int) *League {
	return d.PassThrough.GetLeague(leagueId)
}

func (d *DummySportsDataManager) GetTeam(teamName string) *Team {
	return d.PassThrough.GetTeam(teamName)
}

func (d *DummySportsDataManager) GetLogo(logoUrl string) image.Image {
	return d.PassThrough.GetLogo(logoUrl)
}

func (d *DummySportsDataManager) GetTeamShortName(teamId int) string {
	if shortName, ok := constants.TEAM_SHORT_NAMES[teamId]; ok {
		return shortName
	}

	return "???"
}

func (d *DummySportsDataManager) GetLeagueStandings(leagueId int) []*Conference {
	switch leagueId {
	case constants.LEAGUE_NHL:
		return []*Conference{
			{
				Name: "Western Conference",
				Standings: []*Standing{
					{TeamId: 134855, Rank: 1, GamesPlayed: 36, Points: 61, RegulationWins: 25, RegulationLosses: 2, OvertimeWins: 2, OvertimeLosses: 7, Ties: 0},
					{TeamId: 134856, Rank: 2, GamesPlayed: 38, Points: 56, RegulationWins: 21, RegulationLosses: 7, OvertimeWins: 4, OvertimeLosses: 6, Ties: 0},
					{TeamId: 134857, Rank: 3, GamesPlayed: 38, Points: 50, RegulationWins: 16, RegulationLosses: 10, OvertimeWins: 6, OvertimeLosses: 6, Ties: 0},
					{TeamId: 135913, Rank: 4, GamesPlayed: 35, Points: 44, RegulationWins: 13, RegulationLosses: 8, OvertimeWins: 4, OvertimeLosses: 10, Ties: 0},
					{TeamId: 134846, Rank: 5, GamesPlayed: 37, Points: 44, RegulationWins: 13, RegulationLosses: 14, OvertimeWins: 8, OvertimeLosses: 2, Ties: 0},
					{TeamId: 134849, Rank: 6, GamesPlayed: 38, Points: 44, RegulationWins: 14, RegulationLosses: 13, OvertimeWins: 5, OvertimeLosses: 6, Ties: 0},
					{TeamId: 134852, Rank: 7, GamesPlayed: 36, Points: 39, RegulationWins: 10, RegulationLosses: 12, OvertimeWins: 5, OvertimeLosses: 9, Ties: 0},
					{TeamId: 148494, Rank: 8, GamesPlayed: 39, Points: 39, RegulationWins: 13, RegulationLosses: 18, OvertimeWins: 5, OvertimeLosses: 3, Ties: 0},
					{TeamId: 134853, Rank: 9, GamesPlayed: 37, Points: 37, RegulationWins: 9, RegulationLosses: 17, OvertimeWins: 8, OvertimeLosses: 3, Ties: 0},
					{TeamId: 140082, Rank: 10, GamesPlayed: 35, Points: 36, RegulationWins: 10, RegulationLosses: 14, OvertimeWins: 5, OvertimeLosses: 6, Ties: 0},
					{TeamId: 134858, Rank: 11, GamesPlayed: 36, Points: 36, RegulationWins: 11, RegulationLosses: 16, OvertimeWins: 5, OvertimeLosses: 4, Ties: 0},
					{TeamId: 134859, Rank: 12, GamesPlayed: 38, Points: 36, RegulationWins: 14, RegulationLosses: 16, OvertimeWins: 0, OvertimeLosses: 8, Ties: 0},
					{TeamId: 134848, Rank: 13, GamesPlayed: 37, Points: 34, RegulationWins: 12, RegulationLosses: 18, OvertimeWins: 3, OvertimeLosses: 4, Ties: 0},
					{TeamId: 134851, Rank: 14, GamesPlayed: 35, Points: 33, RegulationWins: 13, RegulationLosses: 17, OvertimeWins: 2, OvertimeLosses: 3, Ties: 0},
					{TeamId: 134850, Rank: 15, GamesPlayed: 36, Points: 33, RegulationWins: 10, RegulationLosses: 18, OvertimeWins: 5, OvertimeLosses: 3, Ties: 0},
					{TeamId: 134854, Rank: 16, GamesPlayed: 36, Points: 32, RegulationWins: 12, RegulationLosses: 17, OvertimeWins: 1, OvertimeLosses: 6, Ties: 0},
				},
			},
			{
				Name: "Eastern Conference",
				Standings: []*Standing{
					{TeamId: 134838, Rank: 1, GamesPlayed: 36, Points: 47, RegulationWins: 14, RegulationLosses: 11, OvertimeWins: 8, OvertimeLosses: 3, Ties: 0},
					{TeamId: 134832, Rank: 2, GamesPlayed: 38, Points: 47, RegulationWins: 16, RegulationLosses: 13, OvertimeWins: 6, OvertimeLosses: 3, Ties: 0},
					{TeamId: 134843, Rank: 3, GamesPlayed: 36, Points: 45, RegulationWins: 12, RegulationLosses: 10, OvertimeWins: 7, OvertimeLosses: 7, Ties: 0},
					{TeamId: 134834, Rank: 4, GamesPlayed: 37, Points: 45, RegulationWins: 13, RegulationLosses: 12, OvertimeWins: 7, OvertimeLosses: 5, Ties: 0},
					{TeamId: 134841, Rank: 5, GamesPlayed: 37, Points: 44, RegulationWins: 14, RegulationLosses: 13, OvertimeWins: 6, OvertimeLosses: 4, Ties: 0},
					{TeamId: 134836, Rank: 6, GamesPlayed: 36, Points: 43, RegulationWins: 17, RegulationLosses: 13, OvertimeWins: 3, OvertimeLosses: 3, Ties: 0},
					{TeamId: 134845, Rank: 7, GamesPlayed: 37, Points: 43, RegulationWins: 18, RegulationLosses: 13, OvertimeWins: 1, OvertimeLosses: 5, Ties: 0},
					{TeamId: 134833, Rank: 8, GamesPlayed: 36, Points: 42, RegulationWins: 17, RegulationLosses: 14, OvertimeWins: 3, OvertimeLosses: 2, Ties: 0},
					{TeamId: 134842, Rank: 9, GamesPlayed: 39, Points: 42, RegulationWins: 12, RegulationLosses: 16, OvertimeWins: 7, OvertimeLosses: 4, Ties: 0},
					{TeamId: 134835, Rank: 10, GamesPlayed: 36, Points: 41, RegulationWins: 13, RegulationLosses: 13, OvertimeWins: 5, OvertimeLosses: 5, Ties: 0},
					{TeamId: 134840, Rank: 11, GamesPlayed: 37, Points: 41, RegulationWins: 13, RegulationLosses: 16, OvertimeWins: 7, OvertimeLosses: 1, Ties: 0},
					{TeamId: 134830, Rank: 12, GamesPlayed: 38, Points: 41, RegulationWins: 14, RegulationLosses: 17, OvertimeWins: 6, OvertimeLosses: 1, Ties: 0},
					{TeamId: 134831, Rank: 13, GamesPlayed: 36, Points: 40, RegulationWins: 12, RegulationLosses: 14, OvertimeWins: 6, OvertimeLosses: 4, Ties: 0},
					{TeamId: 134844, Rank: 14, GamesPlayed: 36, Points: 39, RegulationWins: 13, RegulationLosses: 12, OvertimeWins: 2, OvertimeLosses: 9, Ties: 0},
					{TeamId: 134837, Rank: 15, GamesPlayed: 36, Points: 37, RegulationWins: 12, RegulationLosses: 15, OvertimeWins: 4, OvertimeLosses: 5, Ties: 0},
					{TeamId: 134839, Rank: 16, GamesPlayed: 36, Points: 36, RegulationWins: 8, RegulationLosses: 15, OvertimeWins: 7, OvertimeLosses: 6, Ties: 0},
				},
			},
		}
	case constants.LEAGUE_NFL:
		return []*Conference{
			{
				Name: "American Football",
				Standings: []*Standing{
					{TeamId: 134920, Rank: 1, GamesPlayed: 15, Points: 0, RegulationWins: 12, RegulationLosses: 3, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134930, Rank: 2, GamesPlayed: 15, Points: 0, RegulationWins: 11, RegulationLosses: 3, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134928, Rank: 3, GamesPlayed: 15, Points: 0, RegulationWins: 9, RegulationLosses: 4, OvertimeWins: 2, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134918, Rank: 4, GamesPlayed: 15, Points: 0, RegulationWins: 11, RegulationLosses: 4, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 135908, Rank: 5, GamesPlayed: 15, Points: 0, RegulationWins: 10, RegulationLosses: 4, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134926, Rank: 6, GamesPlayed: 15, Points: 0, RegulationWins: 10, RegulationLosses: 5, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134925, Rank: 7, GamesPlayed: 15, Points: 0, RegulationWins: 9, RegulationLosses: 6, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134927, Rank: 8, GamesPlayed: 15, Points: 0, RegulationWins: 7, RegulationLosses: 6, OvertimeWins: 1, OvertimeLosses: 1, Ties: 0},
					{TeamId: 134922, Rank: 9, GamesPlayed: 15, Points: 0, RegulationWins: 7, RegulationLosses: 8, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134931, Rank: 10, GamesPlayed: 15, Points: 0, RegulationWins: 5, RegulationLosses: 9, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134923, Rank: 12, GamesPlayed: 15, Points: 0, RegulationWins: 5, RegulationLosses: 10, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134924, Rank: 13, GamesPlayed: 15, Points: 0, RegulationWins: 3, RegulationLosses: 12, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134921, Rank: 14, GamesPlayed: 15, Points: 0, RegulationWins: 3, RegulationLosses: 12, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134929, Rank: 15, GamesPlayed: 15, Points: 0, RegulationWins: 3, RegulationLosses: 12, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134932, Rank: 16, GamesPlayed: 15, Points: 0, RegulationWins: 2, RegulationLosses: 12, OvertimeWins: 0, OvertimeLosses: 1, Ties: 0},
				},
			},
			{
				Name: "National Football",
				Standings: []*Standing{
					{TeamId: 134949, Rank: 1, GamesPlayed: 15, Points: 0, RegulationWins: 11, RegulationLosses: 3, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 135907, Rank: 2, GamesPlayed: 15, Points: 0, RegulationWins: 11, RegulationLosses: 2, OvertimeWins: 0, OvertimeLosses: 2, Ties: 0},
					{TeamId: 134948, Rank: 3, GamesPlayed: 15, Points: 0, RegulationWins: 10, RegulationLosses: 4, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134938, Rank: 4, GamesPlayed: 15, Points: 0, RegulationWins: 10, RegulationLosses: 4, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134936, Rank: 5, GamesPlayed: 15, Points: 0, RegulationWins: 10, RegulationLosses: 4, OvertimeWins: 0, OvertimeLosses: 1, Ties: 0},
					{TeamId: 134940, Rank: 6, GamesPlayed: 15, Points: 0, RegulationWins: 9, RegulationLosses: 4, OvertimeWins: 0, OvertimeLosses: 1, Ties: 0},
					{TeamId: 134939, Rank: 7, GamesPlayed: 15, Points: 0, RegulationWins: 7, RegulationLosses: 7, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134943, Rank: 8, GamesPlayed: 15, Points: 0, RegulationWins: 7, RegulationLosses: 7, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134941, Rank: 9, GamesPlayed: 15, Points: 0, RegulationWins: 7, RegulationLosses: 8, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134945, Rank: 10, GamesPlayed: 15, Points: 0, RegulationWins: 7, RegulationLosses: 8, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134934, Rank: 11, GamesPlayed: 15, Points: 0, RegulationWins: 5, RegulationLosses: 8, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134942, Rank: 12, GamesPlayed: 15, Points: 0, RegulationWins: 6, RegulationLosses: 7, OvertimeWins: 0, OvertimeLosses: 2, Ties: 0},
					{TeamId: 134944, Rank: 13, GamesPlayed: 15, Points: 0, RegulationWins: 5, RegulationLosses: 10, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 134937, Rank: 14, GamesPlayed: 15, Points: 0, RegulationWins: 4, RegulationLosses: 9, OvertimeWins: 0, OvertimeLosses: 2, Ties: 0},
					{TeamId: 134946, Rank: 15, GamesPlayed: 15, Points: 0, RegulationWins: 3, RegulationLosses: 11, OvertimeWins: 0, OvertimeLosses: 1, Ties: 0},
					{TeamId: 134935, Rank: 16, GamesPlayed: 15, Points: 0, RegulationWins: 2, RegulationLosses: 11, OvertimeWins: 0, OvertimeLosses: 2, Ties: 0},
				},
			},
		}
	case constants.LEAGUE_PWHL:
		return []*Conference{
			{
				Name: "League Standings",
				Standings: []*Standing{
					{TeamId: 1, Rank: 1, GamesPlayed: 7, Points: 18, RegulationWins: 6, RegulationLosses: 1, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
					{TeamId: 3, Rank: 2, GamesPlayed: 6, Points: 11, RegulationWins: 3, RegulationLosses: 2, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 6, Rank: 3, GamesPlayed: 7, Points: 11, RegulationWins: 3, RegulationLosses: 2, OvertimeWins: 0, OvertimeLosses: 2, Ties: 0},
					{TeamId: 2, Rank: 4, GamesPlayed: 6, Points: 10, RegulationWins: 3, RegulationLosses: 2, OvertimeWins: 0, OvertimeLosses: 1, Ties: 0},
					{TeamId: 8, Rank: 5, GamesPlayed: 6, Points: 10, RegulationWins: 3, RegulationLosses: 2, OvertimeWins: 0, OvertimeLosses: 1, Ties: 0},
					{TeamId: 9, Rank: 6, GamesPlayed: 7, Points: 8, RegulationWins: 2, RegulationLosses: 4, OvertimeWins: 1, OvertimeLosses: 0, Ties: 0},
					{TeamId: 5, Rank: 7, GamesPlayed: 8, Points: 7, RegulationWins: 1, RegulationLosses: 5, OvertimeWins: 2, OvertimeLosses: 0, Ties: 0},
					{TeamId: 4, Rank: 8, GamesPlayed: 7, Points: 6, RegulationWins: 2, RegulationLosses: 5, OvertimeWins: 0, OvertimeLosses: 0, Ties: 0},
				},
			},
		}
	default:
		return []*Conference{}
	}
}
