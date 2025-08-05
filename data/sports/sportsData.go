package sports

import (
	"image"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/utils"
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
	Events            *utils.ExpirableMap[int, map[int]*Event]
	Leagues           *utils.ExpirableMap[int, *League]
	Teams             *utils.ExpirableMap[string, *Team]
	Logos             *utils.ExpirableMap[string, image.Image]
}

func NewSportsData() SportsData {
	return &SportsDataManager{
		TheSportsDbClient: NewTheSportsDbClient(),
		Events:            utils.NewExpirableMap[int, map[int]*Event](time.Hour),
		Leagues:           utils.NewExpirableMap[int, *League](time.Hour),
		Teams:             utils.NewExpirableMap[string, *Team](time.Hour * 24),
		Logos:             utils.NewExpirableMap[string, image.Image](time.Hour * 24),
	}
}

func (s *SportsDataManager) GetUpcomingEventsForLeague(leagueId int) []*Event {
	if s.Events.Get(leagueId) == nil {
		s.Events.Set(leagueId, map[int]*Event{})
	}

	events := *s.Events.Get(leagueId)
	if len(events) == 0 {
		for _, teamId := range constants.LEAGUE_TEAMS[leagueId] {
			nextGame := s.TheSportsDbClient.GetNextGameForTeam(teamId)
			if nextGame != nil {
				events[nextGame.Id] = nextGame
			}
		}
	}

	eventsSlice := []*Event{}
	for _, event := range events {
		eventsSlice = append(eventsSlice, event)
	}

	slices.SortFunc(eventsSlice, func(eventA, eventB *Event) int {
		return int(eventA.Time.Sub(eventB.Time))
	})

	return eventsSlice
}

func (s *SportsDataManager) GetLeague(leagueId int) *League {
	cachedLeague := s.Leagues.Get(leagueId)
	if cachedLeague != nil {
		return *cachedLeague
	}

	league := s.TheSportsDbClient.GetLeague(leagueId)
	if league == nil {
		return nil
	}

	s.Leagues.Set(leagueId, league)
	return league
}

func (s *SportsDataManager) GetTeam(teamName string) *Team {
	cachedTeam := s.Teams.Get(teamName)
	if cachedTeam != nil {
		return *cachedTeam
	}

	team := s.TheSportsDbClient.GetTeam(teamName)
	if team == nil {
		return nil
	}

	team.ShortName = constants.TEAM_SHORT_NAMES[team.Id]

	s.Teams.Set(teamName, team)
	return team
}

func (s *SportsDataManager) GetLogo(logoUrl string) image.Image {
	cachedLogo := s.Logos.Get(logoUrl)
	if cachedLogo != nil {
		return *cachedLogo
	}

	logo := s.TheSportsDbClient.GetLogo(logoUrl)
	if logo == nil {
		return nil
	}

	s.Logos.Set(logoUrl, logo)
	return logo
}
