package sports

import (
	"image"
	"rpi-rgb-screen/config"
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
	Id              int
	LeagueId        int
	Name            string
	HomeTeamName    string
	AwayTeamName    string
	HomeTeamId      int
	AwayTeamId      int
	HomeTeamLogoUrl string
	AwayTeamLogoUrl string
	Time            time.Time
	HomeScore       int
	AwayScore       int
}

type SportsData interface {
	GetUpcomingEventsForLeague(leagueId int, onlyCacheFavoriteTeams bool) []*Event
	GetPastEventsForLeague(leagueId int, onlyCacheFavoriteTeams bool) []*Event
	GetLeague(leagueId int) *League
	GetTeam(teamName string) *Team
	GetLogo(logoUrl string) image.Image
	GetTeamShortName(teamId int) string
}

type SportsDataManager struct {
	TheSportsDbClient *TheSportsDbClient
	UpcomingEvents    *utils.ExpirableMap[int, map[int]*Event]
	PastEvents        *utils.ExpirableMap[int, map[int]*Event]
	Leagues           *utils.ExpirableMap[int, *League]
	Teams             *utils.ExpirableMap[string, *Team]
	Logos             *utils.ExpirableMap[string, image.Image]
}

func NewSportsData() SportsData {
	return &SportsDataManager{
		TheSportsDbClient: NewTheSportsDbClient(),
		UpcomingEvents:    utils.NewExpirableMap[int, map[int]*Event](time.Hour),
		PastEvents:        utils.NewExpirableMap[int, map[int]*Event](time.Hour),
		Leagues:           utils.NewExpirableMap[int, *League](time.Hour),
		Teams:             utils.NewExpirableMap[string, *Team](time.Hour * 24),
		Logos:             utils.NewExpirableMap[string, image.Image](time.Hour * 24),
	}
}

// Fetches and caches upcoming events for the given league. If onlyCacheFavoriteTeams is true, only events involving favorite teams are cached.
// All cached events are returned regardless of the onlyCacheFavoriteTeams flag.
func (s *SportsDataManager) GetUpcomingEventsForLeague(leagueId int, onlyCacheFavoriteTeams bool) []*Event {
	if s.UpcomingEvents.Get(leagueId) == nil {
		s.UpcomingEvents.Set(leagueId, map[int]*Event{})
	}

	teamIds := constants.LEAGUE_TEAMS[leagueId]
	if onlyCacheFavoriteTeams {
		teamIds = utils.GetFavoriteTeamsInLeague(config.Config.FavoriteTeams, leagueId)
	}

	events := *s.UpcomingEvents.Get(leagueId)
	if len(events) == 0 {
		for _, teamId := range teamIds {
			nextGame := s.TheSportsDbClient.GetNextGameForTeam(teamId)
			if nextGame != nil && time.Until(nextGame.Time) < time.Hour*24*7 { // Only include games within the next week
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

// Fetches and caches past events for the given league. If onlyCacheFavoriteTeams is true, only events involving favorite teams are cached.
// All cached events are returned regardless of the onlyCacheFavoriteTeams flag.
func (s *SportsDataManager) GetPastEventsForLeague(leagueId int, onlyCacheFavoriteTeams bool) []*Event {
	if s.PastEvents.Get(leagueId) == nil {
		s.PastEvents.Set(leagueId, map[int]*Event{})
	}

	teamIds := constants.LEAGUE_TEAMS[leagueId]
	if onlyCacheFavoriteTeams {
		teamIds = utils.GetFavoriteTeamsInLeague(config.Config.FavoriteTeams, leagueId)
	}

	events := *s.PastEvents.Get(leagueId)
	if len(events) == 0 {
		for _, teamId := range teamIds {
			lastGame := s.TheSportsDbClient.GetLastGameForTeam(teamId)
			if lastGame != nil && lastGame.LeagueId == leagueId && time.Since(lastGame.Time) < time.Hour*24*7 { // Only include games within the last week
				events[lastGame.Id] = lastGame
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

	if shortName, ok := constants.TEAM_SHORT_NAMES[team.Id]; ok {
		team.ShortName = shortName
	}

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

func (s *SportsDataManager) GetTeamShortName(teamId int) string {
	if shortName, ok := constants.TEAM_SHORT_NAMES[teamId]; ok {
		return shortName
	}

	return "???"
}
