package screen

import (
	"image"
	"image/color"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/sports"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/utils"
	"time"

	"github.com/fogleman/gg"
)

type SportsLeagueScreen struct {
	State      ScreenState
	Ctx        *gg.Context
	KeyFrames  *animation.KeyFrames
	Fonts      *fonts.Fonts
	SportsData sports.SportsData
	LeagueId   int
	League     *sports.League
	LeagueLogo image.Image
}

func NewSportsLeagueScreen(fonts *fonts.Fonts, sportsData sports.SportsData, leagueId int) Screen {
	return &SportsLeagueScreen{
		State:      StateNotDisplayed,
		Ctx:        gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:  animation.NewKeyFrames(5 * time.Second),
		Fonts:      fonts,
		SportsData: sportsData,
		LeagueId:   leagueId,
	}
}

func (s *SportsLeagueScreen) SetState(state ScreenState) {
	s.State = state
	if state == StateDisplayed {
		s.KeyFrames.Start()
	}
}

func (s *SportsLeagueScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		if s.League == nil {
			s.League = s.SportsData.GetLeague(s.LeagueId)
		}

		if s.LeagueLogo == nil {
			logo := s.SportsData.GetLogo(s.League.LogoUrl)
			if logo != nil {
				s.LeagueLogo = utils.ResizeImage(logo, 32)
			}
		}

		close(doneChan)
	}()

	return doneChan
}

func (s *SportsLeagueScreen) Render(elapsed time.Duration) (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	if s.LeagueLogo != nil {
		s.Ctx.DrawImageAnchored(s.LeagueLogo, 32, 16, 0.5, 0.5)
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}
