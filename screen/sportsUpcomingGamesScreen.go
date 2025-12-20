package screen

import (
	"fmt"
	"image"
	"image/color"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/sports"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/transition"
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
	SportsData   sports.SportsData
	Event        *sports.Event
	LogoHome     image.Image
	LogoAway     image.Image
	HomeTeamName string
	AwayTeamName string

	// Animation state
	KeyFrames *animation.KeyFrames
}

func NewSportsUpcomingGamesScreen(fonts *fonts.Fonts, sportsData sports.SportsData, event *sports.Event) Screen {
	keyFrames := animation.NewKeyFrames(10 * time.Second)

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

func (s *SportsUpcomingGamesScreen) GetPreferredTransition() transition.Transition {
	return transition.NewSlideAndZoomTransition()
}

func (s *SportsUpcomingGamesScreen) SetState(state ScreenState) {
	s.State = state

	switch state {
	case StateDisplayed:
		s.KeyFrames.Start()
	case StateTransitionIn:
		s.KeyFrames.Reset()
	}
}

func (s *SportsUpcomingGamesScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		s.HomeTeamName = s.SportsData.GetTeamShortName(s.Event.HomeTeamId)
		s.AwayTeamName = s.SportsData.GetTeamShortName(s.Event.AwayTeamId)

		if s.LogoHome == nil {
			var logoHome image.Image
			if s.Event.LeagueId == constants.LEAGUE_PWHL {
				filePath := fmt.Sprintf("./resources/pwhlLogos/%d.png", s.Event.HomeTeamId)
				logoHome, _ = utils.ReadImageFromFile(filePath)
			} else {
				logoHome = s.SportsData.GetLogo(s.Event.HomeTeamLogoUrl)
			}
			if logoHome != nil {
				s.LogoHome = utils.ResizeImageSquare(logoHome, 32)
			}
		}
		if s.LogoAway == nil {
			var logoAway image.Image
			if s.Event.LeagueId == constants.LEAGUE_PWHL {
				filePath := fmt.Sprintf("./resources/pwhlLogos/%d.png", s.Event.AwayTeamId)
				logoAway, _ = utils.ReadImageFromFile(filePath)
			} else {
				logoAway = s.SportsData.GetLogo(s.Event.AwayTeamLogoUrl)
			}
			if logoAway != nil {
				s.LogoAway = utils.ResizeImageSquare(logoAway, 32)
			}
		}

		close(doneChan)
	}()

	return doneChan
}

func (s *SportsUpcomingGamesScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	switch s.State {
	case StateTransitionIn:
		s.renderLogos()
	case StateTransitionOut:
		fallthrough
	case StateDisplayed:
		s.renderText()
		s.renderLogos()
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *SportsUpcomingGamesScreen) renderLogos() {
	// Draw the away team logo first so it is underneath the home logo if they overlap
	if s.LogoAway != nil {
		logoAwayPosition := s.KeyFrames.GetPoint("logoAway")
		s.Ctx.DrawImageAnchored(s.LogoAway, logoAwayPosition.X, logoAwayPosition.Y, 0, 0)
	}
	if s.LogoHome != nil {
		logoHomePosition := s.KeyFrames.GetPoint("logoHome")
		s.Ctx.DrawImageAnchored(s.LogoHome, logoHomePosition.X, logoHomePosition.Y, 0, 0)
	}
}

func (s *SportsUpcomingGamesScreen) renderText() {
	teamNamesColor := s.KeyFrames.GetColor("teamNames")
	if teamNamesColor.A > 0 {
		s.Ctx.SetFontFace(s.Fonts.Size8x13B)
		s.Ctx.SetColor(teamNamesColor)
		s.Ctx.DrawStringAnchored(s.AwayTeamName, 32, -3, 0.5, 1)
		s.Ctx.DrawStringAnchored(s.HomeTeamName, 32, 31, 0.5, 0)
		s.Ctx.SetFontFace(s.Fonts.Size6x10)
		s.Ctx.DrawStringAnchored("@", 32, 14, 0.5, 0.5)
	}

	dateAndTimeColor := s.KeyFrames.GetColor("dateAndTime")
	if dateAndTimeColor.A > 0 {
		s.Ctx.SetFontFace(s.Fonts.Size5x7)
		s.Ctx.SetColor(dateAndTimeColor)
		s.Ctx.DrawStringAnchored(strings.ToUpper(s.Event.Time.Format("Mon")), 32, 1, 0.5, 1)
		s.Ctx.DrawStringAnchored(strings.ToUpper(s.Event.Time.Format("Jan 2")), 32, 8, 0.5, 1)
		if s.Event.IsTBD {
			s.Ctx.DrawStringAnchored("TBD", 32, 18, 0.5, 1)
		} else {
			s.Ctx.DrawStringAnchored(s.Event.Time.Format("3:04"), 32, 15, 0.5, 1)
			s.Ctx.DrawStringAnchored(s.Event.Time.Format("PM"), 32, 22, 0.5, 1)
		}
	}
}
