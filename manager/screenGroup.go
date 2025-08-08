package manager

import (
	"rpi-rgb-screen/screen"
)

type ScreenGroup struct {
	Screens     []screen.Screen
	ScreenIndex int
	InitFunc    func() []screen.Screen
}

func NewScreenGroup(initFunc func() []screen.Screen) *ScreenGroup {
	return &ScreenGroup{
		Screens:  []screen.Screen{},
		InitFunc: initFunc,
	}
}

func (s *ScreenGroup) GetNextScreen() screen.Screen {
	if s.ScreenIndex >= len(s.Screens) {
		return nil
	}

	screen := s.Screens[s.ScreenIndex]
	s.ScreenIndex++
	return screen
}

func (s *ScreenGroup) Initialize() {
	newScreens := s.InitFunc()

	s.Screens = newScreens
	s.ScreenIndex = 0
}
