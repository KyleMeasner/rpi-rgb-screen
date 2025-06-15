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
	State          ScreenState
	ScreenDuration time.Duration
	Ctx            *gg.Context
	Fonts          *fonts.Fonts

	// Data
	SportsData sports.SportsData
	Event      sports.Event
	Logo1      image.Image
	Logo2      image.Image
	Team1      *sports.Team
	Team2      *sports.Team

	// Animation state
	LogoAnimation *animation.Animation
	TextAnimation *animation.Animation
	KeyFrames     *animation.KeyFrames
}

func NewSportsUpcomingGamesScreen(fonts *fonts.Fonts, sportsData sports.SportsData, event sports.Event) Screen {
	logoAnimation := animation.NewAnimation(1*time.Second, 1*time.Second)
	logoAnimation.Points = map[string]animation.AnimationPoint{
		"logo1": {Start: image.Point{0, 0}, End: image.Point{-16, 0}},
		"logo2": {Start: image.Point{64, 0}, End: image.Point{80, 0}},
	}

	teamNamesAnimation := animation.NewAnimation(1750*time.Millisecond, 500*time.Millisecond)
	teamNamesAnimation.Colors = map[string]animation.AnimationColor{
		"teamNames": {Start: color.RGBA{255, 255, 255, 0}, End: color.RGBA{255, 255, 255, 255}},
	}

	dateAndTimeAnimation := animation.NewAnimation(5250*time.Millisecond, 500*time.Millisecond)
	dateAndTimeAnimation.Colors = map[string]animation.AnimationColor{
		"dateAndTime": {Start: color.RGBA{255, 255, 255, 0}, End: color.RGBA{255, 255, 255, 255}},
	}

	keyFrames := animation.NewKeyFrames(map[string]*animation.Animation{
		"logo":        logoAnimation,
		"teamNames":   teamNamesAnimation,
		"dateAndTime": dateAndTimeAnimation,
	})

	return &SportsUpcomingGamesScreen{
		State:          StateNotDisplayed,
		ScreenDuration: 10 * time.Second,
		Ctx:            gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		Fonts:          fonts,
		SportsData:     sportsData,
		Event:          event,
		KeyFrames:      keyFrames,
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
	isScreenDone := time.Since(s.KeyFrames.StartTime) > s.ScreenDuration

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
		s.renderLogoAnimation()
		s.renderText()
	}

	return s.Ctx.Image(), isScreenDone
}

func (s *SportsUpcomingGamesScreen) renderTransitionDisplay() {
	s.drawLogos(image.Point{0, 0}, image.Point{64, 0})
}

func (s *SportsUpcomingGamesScreen) renderLogoAnimation() {
	logo1Position := s.KeyFrames.GetPoint("logo", "logo1")
	logo2Position := s.KeyFrames.GetPoint("logo", "logo2")
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
	teamNamesColor := s.KeyFrames.GetColor("teamNames", "teamNames")
	dateAndTimeColor := s.KeyFrames.GetColor("dateAndTime", "dateAndTime")

	if _, _, _, alpha := dateAndTimeColor.RGBA(); alpha == 0 {
		s.Ctx.SetColor(teamNamesColor)
		s.Ctx.SetFontFace(s.Fonts.Size5x7)

		s.Ctx.DrawStringAnchored(s.Team1.ShortName, 32, 1, 0.5, 1)
		s.Ctx.DrawStringAnchored("VS", 32, 8, 0.5, 1)
		s.Ctx.DrawStringAnchored(s.Team2.ShortName, 32, 15, 0.5, 1)
		return
	}

	eventTime, err := time.Parse("2006-01-02T15:04:05", s.Event.Timestamp)
	if err != nil {
		log.Printf("Error reading event time from timestamp '%s'.", s.Event.Timestamp)
		return
	}
	eventTime = eventTime.Local()

	s.Ctx.SetColor(dateAndTimeColor)
	s.Ctx.SetFontFace(s.Fonts.Size5x7)

	s.Ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Mon")), 32, 1, 0.5, 1)
	s.Ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Jan 2")), 32, 8, 0.5, 1)
	s.Ctx.DrawStringAnchored(eventTime.Format("3:04"), 32, 15, 0.5, 1)
	s.Ctx.DrawStringAnchored(eventTime.Format("PM"), 32, 22, 0.5, 1)
}
