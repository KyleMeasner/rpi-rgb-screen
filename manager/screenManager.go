package manager

import (
	"rpi-rgb-screen/data"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/screen"
	"rpi-rgb-screen/transition"
	"time"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

type ScreenManager struct {
	Screens     []screen.Screen
	ToolKit     *rgbmatrix.ToolKit
	DataManager *data.DataManager
}

func NewScreenManager(fonts *fonts.Fonts, toolKit *rgbmatrix.ToolKit, dataManager *data.DataManager) *ScreenManager {
	screens := []screen.Screen{
		screen.NewDummyScreen(fonts),
	}

	events := dataManager.SportsData.GetUpcomingEvents()
	for _, event := range events {
		screens = append(screens, screen.NewSportsScoresScreen(fonts, dataManager.SportsData, event))
	}

	return &ScreenManager{
		Screens:     screens,
		ToolKit:     toolKit,
		DataManager: dataManager,
	}
}

func (s *ScreenManager) Run() {
	i := 0
	currScreen := s.Screens[i]
	s.ToolKit.PlayImage(currScreen.Render(), 0)

	for {
		time.Sleep(5 * time.Second)

		i = (i + 1) % len(s.Screens)
		nextScreen := s.Screens[i]
		nextScreen.Refresh()

		err := s.ToolKit.PlayAnimation(transition.NewSlideIn(currScreen, nextScreen))
		if err != nil {
			panic(err)
		}

		currScreen = nextScreen
	}
}
