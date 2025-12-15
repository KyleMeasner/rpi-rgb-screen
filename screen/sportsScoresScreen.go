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
	"time"

	"github.com/fogleman/gg"
)

type SportsScoresScreen struct {
	// State
	State     ScreenState
	Ctx       *gg.Context
	Fonts     *fonts.Fonts
	KeyFrames *animation.KeyFrames
	// Data
	SportsData   sports.SportsData
	Event        *sports.Event
	LogoHome     image.Image
	LogoAway     image.Image
	HomeTeamName string
	AwayTeamName string
}

func NewSportsScoresScreen(fonts *fonts.Fonts, sportsData sports.SportsData, event *sports.Event) Screen {
	keyFrames := animation.NewKeyFrames(10 * time.Second)

	keyFrames.AddPoint("logoAway", image.Point{0, 0})
	keyFrames.AddPointTransitions("logoAway",
		animation.AnimatedPointTransition{Offset: 1000 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{-16, 0}},
		animation.AnimatedPointTransition{Offset: 9500 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{0, 0}},
	)

	keyFrames.AddPoint("logoHome", image.Point{32, 0})
	keyFrames.AddPointTransitions("logoHome",
		animation.AnimatedPointTransition{Offset: 1000 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{48, 0}},
		animation.AnimatedPointTransition{Offset: 9500 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{32, 0}},
	)

	keyFrames.AddColor("scores", color.RGBA{255, 255, 255, 0})
	keyFrames.AddColorTransitions("scores",
		animation.AnimatedColorTransition{Offset: 1000 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: color.RGBA{255, 255, 255, 255}},
		animation.AnimatedColorTransition{Offset: 9250 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: color.RGBA{255, 255, 255, 0}},
	)

	keyFrames.AddColor("loserFade", color.RGBA{0, 0, 0, 0})
	keyFrames.AddColorTransitions("loserFade",
		animation.AnimatedColorTransition{Offset: 1500 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: color.RGBA{0, 0, 0, 192}},
		animation.AnimatedColorTransition{Offset: 9250 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: color.RGBA{0, 0, 0, 0}},
	)

	return &SportsScoresScreen{
		State:      StateNotDisplayed,
		Ctx:        gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		Fonts:      fonts,
		KeyFrames:  keyFrames,
		SportsData: sportsData,
		Event:      event,
	}
}

func (s *SportsScoresScreen) GetPreferredTransition() transition.Transition {
	return transition.NewSlideAndZoomTransition()
}

func (s *SportsScoresScreen) SetState(state ScreenState) {
	s.State = state

	switch state {
	case StateDisplayed:
		s.KeyFrames.Start()
	case StateTransitionIn:
		s.KeyFrames.Reset()
	}
}

func (s *SportsScoresScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		s.HomeTeamName = s.SportsData.GetTeamShortName(s.Event.HomeTeamId)
		s.AwayTeamName = s.SportsData.GetTeamShortName(s.Event.AwayTeamId)

		if s.LogoHome == nil {
			logoHome := s.SportsData.GetLogo(s.Event.HomeTeamLogoUrl)
			if logoHome != nil {
				s.LogoHome = utils.ResizeImageSquare(logoHome, 32)
			}
		}
		if s.LogoAway == nil {
			logoAway := s.SportsData.GetLogo(s.Event.AwayTeamLogoUrl)
			if logoAway != nil {
				s.LogoAway = utils.ResizeImageSquare(logoAway, 32)
			}
		}

		close(doneChan)
	}()

	return doneChan
}

func (s *SportsScoresScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	s.Ctx.SetFontFace(s.Fonts.Size8x13B)
	s.Ctx.SetColor(color.White)

	switch s.State {
	case StateTransitionIn:
		s.renderLogos()
	case StateTransitionOut:
		fallthrough
	case StateDisplayed:
		s.renderText()
		s.renderLogos()
		s.renderLoserFade()
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *SportsScoresScreen) renderLogos() {
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

func (s *SportsScoresScreen) renderText() {
	scoreColor := s.KeyFrames.GetColor("scores")
	if scoreColor.A > 0 {
		s.Ctx.SetFontFace(s.Fonts.Size6x10)
		s.Ctx.SetColor(scoreColor)

		s.Ctx.DrawStringAnchored(fmt.Sprintf("%s", s.AwayTeamName), 17, -3, 0, 1)
		s.Ctx.DrawStringAnchored(fmt.Sprintf("%d", s.Event.AwayScore), 17, 5, 0, 1)

		s.Ctx.DrawStringAnchored(fmt.Sprintf("%s", s.HomeTeamName), 47, 32, 1, 0)
		s.Ctx.DrawStringAnchored(fmt.Sprintf("%d", s.Event.HomeScore), 47, 24, 1, 0)
	}
}

func (s *SportsScoresScreen) renderLoserFade() {
	loserFadeColor := s.KeyFrames.GetColor("loserFade")
	s.Ctx.SetColor(loserFadeColor)
	if s.Event.AwayScore < s.Event.HomeScore {
		s.Ctx.DrawRectangle(0, 0, 16, 32)
		s.Ctx.DrawRectangle(17, 0, 31, 16)
		s.Ctx.Fill()
	} else if s.Event.HomeScore < s.Event.AwayScore {
		s.Ctx.DrawRectangle(48, 0, 16, 32)
		s.Ctx.DrawRectangle(17, 17, 31, 16)
		s.Ctx.Fill()
	}
}
