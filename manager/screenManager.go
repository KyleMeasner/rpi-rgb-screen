package manager

import (
	"image"
	"image/draw"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/screen"
	"rpi-rgb-screen/transition"

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

	for _, screenGroup := range s.ScreenGroups {
		screenGroup.Initialize()
	}

	for _, leagueId := range constants.LEAGUES {
		initFunction := func() []screen.Screen {
			return s.initializeSportsLeague(leagueId)
		}
		screenGroup := NewScreenGroup(initFunction)
		go screenGroup.Initialize()
		s.ScreenGroups = append(s.ScreenGroups, screenGroup)
	}
}

func (s *ScreenManager) initializeClockScreen() []screen.Screen {
	return []screen.Screen{screen.NewClockScreen(s.Fonts)}
}

func (s *ScreenManager) initializeSportsLeague(leagueId int) []screen.Screen {
	events := s.DataManager.SportsData.GetUpcomingEventsForLeague(leagueId)
	if len(events) == 0 {
		return []screen.Screen{}
	}

	screens := []screen.Screen{}
	screens = append(screens, screen.NewSportsLeagueScreen(s.Fonts, s.DataManager.SportsData, leagueId))
	for _, event := range events {
		screens = append(screens, screen.NewSportsUpcomingGamesScreen(s.Fonts, s.DataManager.SportsData, event))
	}
	return screens
}

func (s *ScreenManager) initializeWeatherScreens() []screen.Screen {
	return []screen.Screen{
		screen.NewWeatherCurrentScreen(s.Fonts, s.DataManager.WeatherData),
		screen.NewWeatherForecastScreen(s.Fonts, s.DataManager.WeatherData),
	}
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
