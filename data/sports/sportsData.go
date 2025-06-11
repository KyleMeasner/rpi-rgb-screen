package sports

import (
	"image"
)

const LEAGE_CFL = 4405

type Team struct {
	Id        string `json:"idTeam"`
	Name      string `json:"strTeam"`
	ShortName string `json:"strTeamShort"`
	BadgeUrl  string `json:"strBadge"`
}

type Event struct {
	Id           string `json:"idEvent"`
	Name         string `json:"strEvent"`
	HomeTeamName string `json:"strHomeTeam"`
	AwayTeamName string `json:"strAwayTeam"`
	Timestamp    string `json:"strTimestamp"`
}

type SportsData interface {
	GetUpcomingEvents() []Event
	GetTeam(teamName string) *Team
	GetLogo(teamName string) image.Image
}

type SportsDataManager struct {
	TheSportsDb TheSportsDbClient
	Events      []Event
	Teams       map[string]*Team
	Logos       map[string]image.Image
}

func NewSportsData() SportsData {
	theSportsDb := NewTheSportsDbClient()
	cflEvents := theSportsDb.GetUpcomingEventsForLeague(LEAGE_CFL)

	return &SportsDataManager{
		TheSportsDb: theSportsDb,
		Events:      cflEvents,
		Teams:       map[string]*Team{},
		Logos:       map[string]image.Image{},
	}
}

func (s *SportsDataManager) GetUpcomingEvents() []Event {
	if len(s.Events) < 3 {
		return s.Events
	}
	return s.Events[0:3]
}

func (s *SportsDataManager) GetTeam(teamName string) *Team {
	if team, ok := s.Teams[teamName]; ok {
		return team
	}

	team := s.TheSportsDb.GetTeam(teamName)
	if team == nil {
		return nil
	}

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

	logo := s.TheSportsDb.GetLogo(team.BadgeUrl)
	if logo == nil {
		return nil
	}

	s.Logos[teamName] = logo
	return logo
}
