package manager

import (
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/screen"
	"rpi-rgb-screen/transition"
	"time"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

type ScreenManager struct {
	Screens []screen.Screen
	ToolKit *rgbmatrix.ToolKit
}

func NewScreenManager(fonts *fonts.Fonts, toolKit *rgbmatrix.ToolKit) *ScreenManager {
	screens := []screen.Screen{
		screen.NewDummyScreen(fonts),
		screen.NewDummyScreen(fonts),
	}

	return &ScreenManager{
		Screens: screens,
		ToolKit: toolKit,
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
