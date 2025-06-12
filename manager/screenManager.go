package manager

import (
	"image"
	"image/draw"
	"rpi-rgb-screen/data"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/screen"
	"rpi-rgb-screen/transition"
	"time"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

type ScreenManager struct {
	Screens     []screen.Screen
	Canvas      *rgbmatrix.Canvas
	DataManager *data.DataManager
}

func NewScreenManager(fonts *fonts.Fonts, canvas *rgbmatrix.Canvas, dataManager *data.DataManager) *ScreenManager {
	screens := []screen.Screen{
		screen.NewDummyScreen(fonts),
	}

	events := dataManager.SportsData.GetUpcomingEvents()
	for _, event := range events {
		screens = append(screens, screen.NewSportsUpcomingGamesScreen(fonts, dataManager.SportsData, event))
	}

	return &ScreenManager{
		Screens:     screens,
		Canvas:      canvas,
		DataManager: dataManager,
	}
}

func (s *ScreenManager) Run() {
	// Prep the first screen before we start the loop
	<-s.Screens[0].Refresh()
	s.Screens[0].TransitionEnd(true)

	i := 0
	for {
		currScreen := s.Screens[i]

		i = (i + 1) % len(s.Screens)
		nextScreen := s.Screens[i]
		nextScreenReadyChan := nextScreen.Refresh()

		s.DisplayScreen(currScreen, nextScreenReadyChan)

		currScreen.TransitionStart()
		nextScreen.TransitionStart()

		s.DisplayTransition(transition.NewSlideInTransition(currScreen, nextScreen))

		currScreen.TransitionEnd(false)
		nextScreen.TransitionEnd(true)
	}
}

func (s *ScreenManager) DisplayScreen(screen screen.Screen, nextScreenReadyChan <-chan bool) {
	nextScreenReady := false
	currScreenDone := false
	go func() {
		<-nextScreenReadyChan
		nextScreenReady = true
	}()

	lastRenderTime := time.Now()
	var renderedScreen image.Image
	for !currScreenDone || !nextScreenReady {
		renderTime := time.Now()
		timeDiff := renderTime.Sub(lastRenderTime)

		renderedScreen, currScreenDone = screen.Render(timeDiff)
		draw.Draw(s.Canvas, s.Canvas.Bounds(), renderedScreen, image.Point{}, draw.Over)
		s.Canvas.Render()

		lastRenderTime = renderTime
	}
}

func (s *ScreenManager) DisplayTransition(transition transition.Transition) {
	var renderedTransition image.Image
	transitionDone := false
	lastRenderTime := time.Now()
	for !transitionDone {
		renderTime := time.Now()
		timeDiff := renderTime.Sub(lastRenderTime)

		renderedTransition, transitionDone = transition.Render(timeDiff)
		draw.Draw(s.Canvas, s.Canvas.Bounds(), renderedTransition, image.Point{}, draw.Over)
		s.Canvas.Render()

		lastRenderTime = renderTime
	}
}
