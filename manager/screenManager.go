package manager

import (
	"image"
	"image/draw"
	"log"
	"rpi-rgb-screen/config"
	"rpi-rgb-screen/data"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/screen"
	"rpi-rgb-screen/transition"
	"rpi-rgb-screen/utils"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

type ScreenManager struct {
	ScreenGroups     []*ScreenGroup
	ScreenGroupIndex int
	Canvas           *rgbmatrix.Canvas
	DataManager      *data.DataManager
	Fonts            *fonts.Fonts
}

func NewScreenManager(fonts *fonts.Fonts, canvas *rgbmatrix.Canvas, dataManager *data.DataManager) *ScreenManager {
	return &ScreenManager{
		ScreenGroups: []*ScreenGroup{},
		Canvas:       canvas,
		DataManager:  dataManager,
		Fonts:        fonts,
	}
}

func (s *ScreenManager) Initialize() {
	s.ScreenGroups = []*ScreenGroup{
		NewScreenGroup(s.initializeClockScreen),
		NewScreenGroup(s.initializeWeatherScreens),
	}

	numQuickLoadingScreens := len(s.ScreenGroups)

	// Initialize the fast loading screen groups
	for _, screenGroup := range s.ScreenGroups {
		screenGroup.Initialize()
	}

	favoriteLeagues := utils.GetFavoriteLeagues(config.Config.FavoriteTeams)

	// Upcoming games screens
	for _, leagueId := range favoriteLeagues {
		initFunction := func() []screen.Screen {
			return s.initializeSportsLeagueUpcomingGames(leagueId)
		}
		screenGroup := NewScreenGroup(initFunction)
		s.ScreenGroups = append(s.ScreenGroups, screenGroup)
	}

	// Past games screens
	for _, leagueId := range favoriteLeagues {
		initFunction := func() []screen.Screen {
			return s.initializeSportsLeaguePastGames(leagueId)
		}
		screenGroup := NewScreenGroup(initFunction)
		s.ScreenGroups = append(s.ScreenGroups, screenGroup)
	}

	// Standings screens
	for _, leagueId := range favoriteLeagues {
		initFunction := func() []screen.Screen {
			return s.initializeSportsLeagueStandings(leagueId)
		}
		screenGroup := NewScreenGroup(initFunction)
		s.ScreenGroups = append(s.ScreenGroups, screenGroup)
	}

	// Initialize the slow loading screen groups in the background
	go func() {
		for _, screenGroup := range s.ScreenGroups[numQuickLoadingScreens:] {
			screenGroup.Initialize()
		}
	}()
}

func (s *ScreenManager) initializeClockScreen() []screen.Screen {
	log.Printf("Initializing clock screen")
	return []screen.Screen{screen.NewClockScreen(s.Fonts)}
}

func (s *ScreenManager) initializeSportsLeagueUpcomingGames(leagueId int) []screen.Screen {
	log.Printf("Initializing sports league %d upcoming games screen", leagueId)
	events := s.DataManager.SportsData.GetUpcomingEventsForLeague(leagueId, true)
	if len(events) == 0 {
		return []screen.Screen{}
	}

	favoriteTeamsMap := map[int]bool{}
	for _, teamId := range config.Config.FavoriteTeams {
		favoriteTeamsMap[teamId] = true
	}

	screens := []screen.Screen{}
	for _, event := range events {
		if !favoriteTeamsMap[event.HomeTeamId] && !favoriteTeamsMap[event.AwayTeamId] {
			continue
		}
		screens = append(screens, screen.NewSportsUpcomingGamesScreen(s.Fonts, s.DataManager.SportsData, event))
	}

	if len(screens) > 0 {
		return append([]screen.Screen{screen.NewSportsLeagueScreen(s.Fonts, s.DataManager.SportsData, leagueId, "UPCOMING GAMES")}, screens...)
	}
	return []screen.Screen{}
}

func (s *ScreenManager) initializeSportsLeaguePastGames(leagueId int) []screen.Screen {
	log.Printf("Initializing sports league %d past games screen", leagueId)
	events := s.DataManager.SportsData.GetPastEventsForLeague(leagueId, true)
	if len(events) == 0 {
		return []screen.Screen{}
	}

	favoriteTeamsMap := map[int]bool{}
	for _, teamId := range config.Config.FavoriteTeams {
		favoriteTeamsMap[teamId] = true
	}

	screens := []screen.Screen{}
	for _, event := range events {
		if !favoriteTeamsMap[event.HomeTeamId] && !favoriteTeamsMap[event.AwayTeamId] {
			continue
		}
		screens = append(screens, screen.NewSportsScoresScreen(s.Fonts, s.DataManager.SportsData, event))
	}

	if len(screens) > 0 {
		return append([]screen.Screen{screen.NewSportsLeagueScreen(s.Fonts, s.DataManager.SportsData, leagueId, "LATEST SCORES")}, screens...)
	}
	return []screen.Screen{}
}

func (s *ScreenManager) initializeSportsLeagueStandings(leagueId int) []screen.Screen {
	log.Printf("Initializing sports league %d standings screen", leagueId)
	standings := s.DataManager.SportsData.GetLeagueStandings(leagueId)
	if len(standings) == 0 {
		return []screen.Screen{}
	}

	screens := []screen.Screen{}
	for _, conference := range standings {
		screens = append(screens, screen.NewSportsStandingsScreen(s.Fonts, s.DataManager.SportsData, leagueId, conference))
	}

	if len(screens) > 0 {
		return append([]screen.Screen{screen.NewSportsLeagueScreen(s.Fonts, s.DataManager.SportsData, leagueId, "STANDINGS")}, screens...)
	}
	return []screen.Screen{}
}

func (s *ScreenManager) initializeWeatherScreens() []screen.Screen {
	log.Printf("Initializing weather screens")

	screens := []screen.Screen{
		screen.NewWeatherCurrentScreen(s.Fonts, s.DataManager.WeatherData),
	}

	forecast := s.DataManager.WeatherData.GetForecast(config.Config.Location)
	if len(forecast) > 0 {
		for _, dayForecast := range forecast[1:] {
			screens = append(screens, screen.NewWeatherForecastDayScreen(s.Fonts, dayForecast))
		}
	}

	return screens
}

func (s *ScreenManager) Run() {
	s.ScreenGroupIndex = 0

	// Prep the first screen before we start the loop
	s.ScreenGroups[s.ScreenGroupIndex].Initialize()
	currScreen := s.GetNextScreen()
	<-currScreen.Refresh()

	for {
		nextScreen := s.GetNextScreen()
		nextScreenReadyChan := nextScreen.Refresh()

		currScreen.SetState(screen.StateDisplayed)
		s.DisplayScreen(currScreen, nextScreenReadyChan)

		currScreen.SetState(screen.StateTransitionOut)
		nextScreen.SetState(screen.StateTransitionIn)

		transition := nextScreen.GetPreferredTransition()
		transition.SetScreens(currScreen, nextScreen)
		s.DisplayTransition(transition)

		currScreen.SetState(screen.StateNotDisplayed)
		currScreen = nextScreen
	}
}

func (s *ScreenManager) GetNextScreen() screen.Screen {
	currScreenGroup := s.ScreenGroups[s.ScreenGroupIndex]
	nextScreen := currScreenGroup.GetNextScreen()
	if nextScreen != nil {
		return nextScreen
	}

	for {
		// 1. Reset the current screen group for next time
		go currScreenGroup.Initialize()

		// 2. Get the first screen from the next screen group
		s.ScreenGroupIndex = (s.ScreenGroupIndex + 1) % len(s.ScreenGroups)
		nextScreenGroup := s.ScreenGroups[s.ScreenGroupIndex]
		screen := nextScreenGroup.GetNextScreen()
		if screen != nil {
			return screen
		}
	}
}

func (s *ScreenManager) DisplayScreen(screen screen.Screen, nextScreenReadyChan <-chan bool) {
	nextScreenReady := false
	currScreenDone := false
	go func() {
		<-nextScreenReadyChan
		nextScreenReady = true
	}()

	var renderedScreen image.Image
	for !currScreenDone || !nextScreenReady {
		renderedScreen, currScreenDone = screen.Render()
		draw.Draw(s.Canvas, s.Canvas.Bounds(), renderedScreen, image.Point{}, draw.Over)
		s.Canvas.Render()
	}
}

func (s *ScreenManager) DisplayTransition(transition transition.Transition) {
	transition.Start()

	var renderedTransition image.Image
	transitionDone := false
	for !transitionDone {
		renderedTransition, transitionDone = transition.Render()
		draw.Draw(s.Canvas, s.Canvas.Bounds(), renderedTransition, image.Point{}, draw.Over)
		s.Canvas.Render()
	}
}
