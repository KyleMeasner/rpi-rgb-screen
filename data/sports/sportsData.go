package sports

import (
	"image"
	"rpi-rgb-screen/constants"
	"slices"
	"time"
)

type League struct {
	Id            int
	Name          string
	CurrentSeason string
	LogoUrl       string
}

type Team struct {
	Id        int
	Name      string
	ShortName string
	LogoUrl   string
}

type Event struct {
	Id           int
	Name         string
	HomeTeamName string
	AwayTeamName string
	Time         time.Time
}

type SportsData interface {
	GetUpcomingEventsForLeague(leagueId int) []*Event
	GetLeague(leagueId int) *League
	GetTeam(teamName string) *Team
	GetLogo(teamName string) image.Image
}

type SportsDataManager struct {
	TheSportsDbClient *TheSportsDbClient
	Events            map[int]map[int]*Event
	Leagues           map[int]*League
	Teams             map[string]*Team
	Logos             map[string]image.Image
}

func NewSportsData() SportsData {
	return &SportsDataManager{
		TheSportsDbClient: NewTheSportsDbClient(),
		Events:            map[int]map[int]*Event{},
		Leagues:           map[int]*League{},
		Teams:             map[string]*Team{},
		Logos:             map[string]image.Image{},
	}
}

func (s *SportsDataManager) GetUpcomingEventsForLeague(leagueId int) []*Event {
	_, ok := s.Events[leagueId]
	if !ok {
		s.Events[leagueId] = map[int]*Event{}
	}

	if len(s.Events[leagueId]) == 0 {
		for _, teamId := range constants.LEAGUE_TEAMS[leagueId] {
			nextGame := s.TheSportsDbClient.GetNextGameForTeam(teamId)
			if nextGame != nil {
				s.Events[leagueId][nextGame.Id] = nextGame
			}
		}
	}

	eventsSlice := []*Event{}
	for _, event := range s.Events[leagueId] {
		eventsSlice = append(eventsSlice, event)
	}

	slices.SortFunc(eventsSlice, func(eventA, eventB *Event) int {
		return int(eventA.Time.Sub(eventB.Time))
	})

	return eventsSlice
}

func (s *SportsDataManager) GetLeague(leagueId int) *League {
	if league, ok := s.Leagues[leagueId]; ok {
		return league
	}

	league := s.TheSportsDbClient.GetLeague(leagueId)
	if league == nil {
		return nil
	}

	s.Leagues[leagueId] = league
	return league
}

func (s *SportsDataManager) GetTeam(teamName string) *Team {
	if team, ok := s.Teams[teamName]; ok {
		return team
	}

	team := s.TheSportsDbClient.GetTeam(teamName)
	if team == nil {
		return nil
	}

	team.ShortName = constants.TEAM_SHORT_NAMES[team.Id]

	s.Teams[teamName] = team
	return team
}

func (s *SportsDataManager) GetLogo(logoUrl string) image.Image {
	if logo, ok := s.Logos[logoUrl]; ok {
		return logo
	}

	logo := s.TheSportsDbClient.GetLogo(logoUrl)
	if logo == nil {
		return nil
	}

	s.Logos[logoUrl] = logo
	return logo
}
