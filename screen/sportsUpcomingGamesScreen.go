package screen

import (
	"image"
	"image/color"
	"log"
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
	Event      sports.Event
	Logo1      image.Image
	Logo2      image.Image
	Team1      *sports.Team
	Team2      *sports.Team

	// Animation state
	KeyFrames *animation.KeyFrames
}

func NewSportsUpcomingGamesScreen(fonts *fonts.Fonts, sportsData sports.SportsData, event sports.Event) Screen {
	keyFrames := animation.NewKeyFrames()

	keyFrames.AddPoint("logoHome", image.Point{0, 0})
	keyFrames.AddPointTransitions("logoHome",
		animation.AnimatedPointTransition{Offset: 1000 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{-16, 0}},
		animation.AnimatedPointTransition{Offset: 4750 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{0, 0}},
		animation.AnimatedPointTransition{Offset: 5250 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{-16, 0}},
		animation.AnimatedPointTransition{Offset: 9500 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{0, 0}},
	)

	keyFrames.AddPoint("logoAway", image.Point{64, 0})
	keyFrames.AddPointTransitions("logoAway",
		animation.AnimatedPointTransition{Offset: 1000 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{80, 0}},
		animation.AnimatedPointTransition{Offset: 4750 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{64, 0}},
		animation.AnimatedPointTransition{Offset: 5250 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{80, 0}},
		animation.AnimatedPointTransition{Offset: 9500 * time.Millisecond, Duration: 500 * time.Millisecond, EndValue: image.Point{64, 0}},
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
		if s.Logo1 == nil {
			logo1 := s.SportsData.GetLogo(s.Event.HomeTeamName)
			if logo1 != nil {
				s.Logo1 = utils.ResizeImage(logo1, 32)
			}
		}
		if s.Logo2 == nil {
			logo2 := s.SportsData.GetLogo(s.Event.AwayTeamName)
			if logo2 != nil {
				s.Logo2 = utils.ResizeImage(logo2, 32)
			}
		}

		if s.Team1 == nil {
			s.Team1 = s.SportsData.GetTeam(s.Event.HomeTeamName)
		}
		if s.Team2 == nil {
			s.Team2 = s.SportsData.GetTeam(s.Event.AwayTeamName)
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
		s.renderTransitionDisplay()
	case StateTransitionOut:
		fallthrough
	case StateDisplayed:
		s.renderText()
		s.renderLogoAnimation()
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *SportsUpcomingGamesScreen) renderTransitionDisplay() {
	s.drawLogos(image.Point{0, 0}, image.Point{64, 0})
}

func (s *SportsUpcomingGamesScreen) renderLogoAnimation() {
	logo1Position := s.KeyFrames.GetPoint("logoHome")
	logo2Position := s.KeyFrames.GetPoint("logoAway")
	s.drawLogos(logo1Position, logo2Position)
}

func (s *SportsUpcomingGamesScreen) drawLogos(logo1Position, logo2Position image.Point) {
	// Draw the away team logo first so it is underneath the home logo
	if s.Logo2 != nil {
		resizedLogo2 := utils.ResizeImage(s.Logo2, 32)
		s.Ctx.DrawImageAnchored(resizedLogo2, logo2Position.X, logo2Position.Y, 1, 0)
	}
	if s.Logo1 != nil {
		resizedLogo1 := utils.ResizeImage(s.Logo1, 32)
		s.Ctx.DrawImageAnchored(resizedLogo1, logo1Position.X, logo1Position.Y, 0, 0)
	}
}

func (s *SportsUpcomingGamesScreen) renderText() {
	s.Ctx.SetFontFace(s.Fonts.Size5x7)

	teamNamesColor := s.KeyFrames.GetColor("teamNames")
	s.Ctx.SetColor(teamNamesColor)
	s.Ctx.DrawStringAnchored(s.Team1.ShortName, 32, 1, 0.5, 1)
	s.Ctx.DrawStringAnchored("VS", 32, 8, 0.5, 1)
	s.Ctx.DrawStringAnchored(s.Team2.ShortName, 32, 15, 0.5, 1)

	eventTime, err := time.Parse("2006-01-02T15:04:05", s.Event.Timestamp)
	if err != nil {
		log.Printf("Error reading event time from timestamp '%s'.", s.Event.Timestamp)
		return
	}
	eventTime = eventTime.Local()

	dateAndTimeColor := s.KeyFrames.GetColor("dateAndTime")
	s.Ctx.SetColor(dateAndTimeColor)
	s.Ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Mon")), 32, 1, 0.5, 1)
	s.Ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Jan 2")), 32, 8, 0.5, 1)
	s.Ctx.DrawStringAnchored(eventTime.Format("3:04"), 32, 15, 0.5, 1)
	s.Ctx.DrawStringAnchored(eventTime.Format("PM"), 32, 22, 0.5, 1)
}
