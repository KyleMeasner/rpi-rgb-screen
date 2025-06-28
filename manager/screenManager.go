package manager

import (
	"image"
	"image/draw"
	"log"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/screen"
	"rpi-rgb-screen/transition"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

type ScreenManager struct {
	Screens          [][]screen.Screen
	ScreenGroupIndex int
	ScreenIndex      int
	Canvas           *rgbmatrix.Canvas
	DataManager      *data.DataManager
	Fonts            *fonts.Fonts
}

func NewScreenManager(fonts *fonts.Fonts, canvas *rgbmatrix.Canvas, dataManager *data.DataManager) *ScreenManager {
	return &ScreenManager{
		Screens:     [][]screen.Screen{},
		Canvas:      canvas,
		DataManager: dataManager,
		Fonts:       fonts,
	}
}

func (s *ScreenManager) Initialize() {
	log.Println("Initializing Screen Manager")
	s.Screens = [][]screen.Screen{
		{screen.NewLoadingScreen(s.Fonts)},
	}
	s.initializeWeatherScreens()
	go s.initializeSportsLeagues()
}

func (s *ScreenManager) initializeSportsLeagues() {
	for _, leagueId := range constants.LEAGUES {
		events := s.DataManager.SportsData.GetUpcomingEventsForLeague(leagueId)
		if len(events) == 0 {
			return
		}

		screenGroup := []screen.Screen{}
		screenGroup = append(screenGroup, screen.NewSportsLeagueScreen(s.Fonts, s.DataManager.SportsData, leagueId))
		for _, event := range events {
			screenGroup = append(screenGroup, screen.NewSportsUpcomingGamesScreen(s.Fonts, s.DataManager.SportsData, event))
		}
		s.Screens = append(s.Screens, screenGroup)
	}
}

func (s *ScreenManager) initializeWeatherScreens() {
	currentWeatherScreen := screen.NewWeatherCurrentScreen(s.Fonts, s.DataManager.WeatherData)
	forecastScreen := screen.NewWeatherForecastScreen(s.Fonts, s.DataManager.WeatherData)
	s.Screens = append(s.Screens, []screen.Screen{currentWeatherScreen, forecastScreen})
}

func (s *ScreenManager) Run() {
	s.ScreenGroupIndex = 0
	s.ScreenIndex = -1

	// Prep the first screen before we start the loop
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
	screenGroup := s.Screens[s.ScreenGroupIndex]
	if s.ScreenIndex < len(screenGroup)-1 {
		s.ScreenIndex++
	} else {
		s.ScreenGroupIndex = (s.ScreenGroupIndex + 1) % len(s.Screens)
		s.ScreenIndex = 0
		screenGroup = s.Screens[s.ScreenGroupIndex]
	}

	return screenGroup[s.ScreenIndex]
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
