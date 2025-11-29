package screen

import (
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

type SportsLeagueScreen struct {
	State      ScreenState
	Ctx        *gg.Context
	KeyFrames  *animation.KeyFrames
	Fonts      *fonts.Fonts
	SportsData sports.SportsData
	LeagueId   int
	League     *sports.League
	LeagueLogo image.Image
	Message    string
}

func NewSportsLeagueScreen(fonts *fonts.Fonts, sportsData sports.SportsData, leagueId int, message string) Screen {
	keyFrames := animation.NewKeyFrames(10 * time.Second)

	messageWidth := 8 * len(message) // 8 pixels per character
	keyFrames.AddPoint("message", image.Point{constants.SCREEN_WIDTH, constants.SCREEN_MIDDLE_Y})
	keyFrames.AddPointTransitions("message",
		animation.AnimatedPointTransition{Offset: 2000 * time.Millisecond, Duration: 6000 * time.Millisecond, EndValue: image.Point{-messageWidth, constants.SCREEN_MIDDLE_Y}},
	)

	return &SportsLeagueScreen{
		State:      StateNotDisplayed,
		Ctx:        gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:  keyFrames,
		Fonts:      fonts,
		SportsData: sportsData,
		LeagueId:   leagueId,
		Message:    message,
	}
}

func (s *SportsLeagueScreen) GetPreferredTransition() transition.Transition {
	return transition.NewZoomInTransition()
}

func (s *SportsLeagueScreen) SetState(state ScreenState) {
	s.State = state

	switch state {
	case StateDisplayed:
		s.KeyFrames.Start()
	case StateTransitionIn:
		s.KeyFrames.Reset()
	}
}

func (s *SportsLeagueScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		if s.League == nil {
			s.League = s.SportsData.GetLeague(s.LeagueId)
		}

		if s.LeagueLogo == nil && s.League != nil {
			logo := s.SportsData.GetLogo(s.League.LogoUrl)
			if logo != nil {
				s.LeagueLogo = utils.ResizeImageSquare(logo, 32)
			}
		}

		close(doneChan)
	}()

	return doneChan
}

func (s *SportsLeagueScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	if s.LeagueLogo != nil {
		s.Ctx.DrawImageAnchored(s.LeagueLogo, 32, 16, 0.5, 0.5)
	}

	if len(s.Message) > 0 {
		messagePosition := s.KeyFrames.GetPoint("message")
		messageWidth := 8 * len(s.Message)

		s.Ctx.SetColor(color.RGBA{0, 0, 0, 192})
		s.Ctx.DrawRectangle(float64(messagePosition.X)-1, float64(messagePosition.Y)-5, float64(messageWidth)+1, 12)
		s.Ctx.Fill()

		s.Ctx.SetFontFace(s.Fonts.Size8x13B)
		s.Ctx.SetColor(color.White)
		s.Ctx.DrawStringAnchored(s.Message, float64(messagePosition.X), float64(messagePosition.Y), 0, 0.5)
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}
