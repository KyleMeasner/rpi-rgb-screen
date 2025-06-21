package screen

import (
	"image"
	"image/color"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/sports"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/utils"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type SportsUpcomingGamesScreen struct {
	State ScreenState
	Ctx   *gg.Context
	Fonts *fonts.Fonts

	// Data
	SportsData sports.SportsData
	Event      *sports.Event
	LogoHome   image.Image
	LogoAway   image.Image
	TeamHome   *sports.Team
	TeamAway   *sports.Team

	// Animation state
	KeyFrames *animation.KeyFrames
}

func NewSportsUpcomingGamesScreen(fonts *fonts.Fonts, sportsData sports.SportsData, event *sports.Event) Screen {
	keyFrames := animation.NewKeyFrames()

	keyFrames.AddPoint("logoAway", image.Point{0, 0})
	keyFrames.AddPointTransitions("logoAway",
		animation.AnimatedPointTransition{Offset: 1000 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{-16, 0}},
		animation.AnimatedPointTransition{Offset: 4750 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{0, 0}},
		animation.AnimatedPointTransition{Offset: 5250 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{-16, 0}},
		animation.AnimatedPointTransition{Offset: 9500 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{0, 0}},
	)

	keyFrames.AddPoint("logoHome", image.Point{32, 0})
	keyFrames.AddPointTransitions("logoHome",
		animation.AnimatedPointTransition{Offset: 1000 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{48, 0}},
		animation.AnimatedPointTransition{Offset: 4750 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{32, 0}},
		animation.AnimatedPointTransition{Offset: 5250 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{48, 0}},
		animation.AnimatedPointTransition{Offset: 9500 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{32, 0}},
	)

	keyFrames.AddColor("teamNames", color.RGBA{255, 255, 255, 0})
	keyFrames.AddColorTransitions("teamNames",
		animation.AnimatedColorTransition{Offset: 1000 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: color.RGBA{255, 255, 255, 255}},
		animation.AnimatedColorTransition{Offset: 4750 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: color.RGBA{255, 255, 255, 0}},
	)

	keyFrames.AddColor("dateAndTime", color.RGBA{255, 255, 255, 0})
	keyFrames.AddColorTransitions("dateAndTime",
		animation.AnimatedColorTransition{Offset: 5250 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: color.RGBA{255, 255, 255, 255}},
		animation.AnimatedColorTransition{Offset: 9250 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: color.RGBA{255, 255, 255, 0}},
	)

	return &SportsUpcomingGamesScreen{
		State:      StateNotDisplayed,
		Ctx:        gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		Fonts:      fonts,
		SportsData: sportsData,
		Event:      event,
		KeyFrames:  keyFrames,
	}
}

func (s *SportsUpcomingGamesScreen) SetState(state ScreenState) {
	s.State = state
	if state == StateDisplayed {
		s.KeyFrames.Start()
	}
}

func (s *SportsUpcomingGamesScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		if s.LogoHome == nil {
			logo1 := s.SportsData.GetLogo(s.Event.HomeTeamName)
			if logo1 != nil {
				s.LogoHome = utils.ResizeImage(logo1, 32)
			}
		}
		if s.LogoAway == nil {
			logo2 := s.SportsData.GetLogo(s.Event.AwayTeamName)
			if logo2 != nil {
				s.LogoAway = utils.ResizeImage(logo2, 32)
			}
		}

		if s.TeamHome == nil {
			s.TeamHome = s.SportsData.GetTeam(s.Event.HomeTeamName)
		}
		if s.TeamAway == nil {
			s.TeamAway = s.SportsData.GetTeam(s.Event.AwayTeamName)
		}

		close(doneChan)
	}()

	return doneChan
}

func (s *SportsUpcomingGamesScreen) Render(elapsed time.Duration) (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	switch s.State {
	case StateTransitionIn:
		s.renderLogoAnimation()
	case StateTransitionOut:
		fallthrough
	case StateDisplayed:
		s.renderText()
		s.renderLogoAnimation()
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *SportsUpcomingGamesScreen) renderLogoAnimation() {
	logoHomePosition := s.KeyFrames.GetPoint("logoHome")
	logoAwayPosition := s.KeyFrames.GetPoint("logoAway")
	s.drawLogos(logoHomePosition, logoAwayPosition)
}

func (s *SportsUpcomingGamesScreen) drawLogos(logoHomePosition, logoAwayPosition image.Point) {
	// Draw the away team logo first so it is underneath the home logo if they overlap
	if s.LogoAway != nil {
		resizedLogo := utils.ResizeImage(s.LogoAway, 32)
		s.Ctx.DrawImageAnchored(resizedLogo, logoAwayPosition.X, logoAwayPosition.Y, 0, 0)
	}
	if s.LogoHome != nil {
		resizedLogo := utils.ResizeImage(s.LogoHome, 32)
		s.Ctx.DrawImageAnchored(resizedLogo, logoHomePosition.X, logoHomePosition.Y, 0, 0)
	}
}

func (s *SportsUpcomingGamesScreen) renderText() {
	teamNamesColor := s.KeyFrames.GetColor("teamNames")
	s.Ctx.SetFontFace(s.Fonts.Size8x13B)
	s.Ctx.SetColor(teamNamesColor)
	s.Ctx.DrawStringAnchored(s.TeamAway.ShortName, 32, -3, 0.5, 1)
	s.Ctx.DrawStringAnchored(s.TeamHome.ShortName, 32, 31, 0.5, 0)
	s.Ctx.SetFontFace(s.Fonts.Size6x10)
	s.Ctx.DrawStringAnchored("@", 32, 14, 0.5, 0.5)

	dateAndTimeColor := s.KeyFrames.GetColor("dateAndTime")
	s.Ctx.SetFontFace(s.Fonts.Size5x7)
	s.Ctx.SetColor(dateAndTimeColor)
	s.Ctx.DrawStringAnchored(strings.ToUpper(s.Event.Time.Format("Mon")), 32, 1, 0.5, 1)
	s.Ctx.DrawStringAnchored(strings.ToUpper(s.Event.Time.Format("Jan 2")), 32, 8, 0.5, 1)
	s.Ctx.DrawStringAnchored(s.Event.Time.Format("3:04"), 32, 15, 0.5, 1)
	s.Ctx.DrawStringAnchored(s.Event.Time.Format("PM"), 32, 22, 0.5, 1)
}
