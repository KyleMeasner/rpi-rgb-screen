package sports

import (
	"image"
	"rpi-rgb-screen/constants"
	"slices"
	"time"
)

const LEAGE_CFL = 4405

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
	GetUpcomingEvents() []*Event
	GetTeam(teamName string) *Team
	GetLogo(teamName string) image.Image
}

type SportsDataManager struct {
	TheSportsDbClient *TheSportsDbClient
	Events            map[int]*Event
	Leagues           map[int]*League
	Teams             map[string]*Team
	Logos             map[string]image.Image
}

func NewSportsData() SportsData {
	return &SportsDataManager{
		TheSportsDbClient: NewTheSportsDbClient(),
		Events:            map[int]*Event{},
		Teams:             map[string]*Team{},
		Logos:             map[string]image.Image{},
	}
}

func (s *SportsDataManager) GetUpcomingEvents() []*Event {
	if len(s.Events) == 0 {
		for _, leagueId := range constants.LEAGUES {
			for _, teamId := range constants.LEAGUE_TEAMS[leagueId] {
				nextGame := s.TheSportsDbClient.GetNextGameForTeam(teamId)
				if nextGame != nil {
					s.Events[nextGame.Id] = nextGame
				}
			}
		}
	}

	eventsSlice := []*Event{}
	for _, event := range s.Events {
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

func (s *SportsDataManager) GetLogo(teamName string) image.Image {
	if logo, ok := s.Logos[teamName]; ok {
		return logo
	}

	team := s.GetTeam(teamName)
	if team == nil {
		return nil
	}

	logo := s.TheSportsDbClient.GetLogo(team.LogoUrl)
	if logo == nil {
		return nil
	}

	s.Logos[teamName] = logo
	return logo
}
