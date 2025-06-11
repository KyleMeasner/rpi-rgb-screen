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

const screenDuration = 10 * time.Second
const transitionDuration = 3 * time.Second

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
		screens = append(screens, screen.NewSportsScoresScreen(fonts, dataManager.SportsData, event))
	}

	return &ScreenManager{
		Screens:     screens,
		Canvas:      canvas,
		DataManager: dataManager,
	}
}

func (s *ScreenManager) Run() {
	i := 0
	for {
		currScreen := s.Screens[i]
		currScreen.Refresh()
		s.DisplayScreen(currScreen)

		i = (i + 1) % len(s.Screens)
		nextScreen := s.Screens[i]
		nextScreen.Refresh()

		currScreen.TransitionStart()
		nextScreen.TransitionStart()

		s.DisplayTransition(transition.NewSlideInTransition(currScreen, nextScreen))

		currScreen.TransitionEnd(false)
		nextScreen.TransitionEnd(true)
	}
}

func (s *ScreenManager) DisplayScreen(screen screen.Screen) {
	lastRenderTime := time.Now()

	for start := time.Now(); time.Since(start) < screenDuration; {
		renderTime := time.Now()
		timeDiff := renderTime.Sub(lastRenderTime)

		renderedScreen := screen.Render(timeDiff)
		draw.Draw(s.Canvas, s.Canvas.Bounds(), renderedScreen, image.Point{}, draw.Over)
		s.Canvas.Render()

		lastRenderTime = renderTime
	}
}

func (s *ScreenManager) DisplayTransition(transition transition.Transition) {
	lastRenderTime := time.Now()
	for start := time.Now(); time.Since(start) < transitionDuration; {
		renderTime := time.Now()
		timeDiff := renderTime.Sub(lastRenderTime)

		renderedTransition := transition.Render(timeDiff)
		draw.Draw(s.Canvas, s.Canvas.Bounds(), renderedTransition, image.Point{}, draw.Over)
		s.Canvas.Render()

		lastRenderTime = renderTime
	}
}
