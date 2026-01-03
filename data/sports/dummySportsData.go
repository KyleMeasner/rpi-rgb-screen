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
