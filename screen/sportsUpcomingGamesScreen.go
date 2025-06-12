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

// const animationDuration = 2 * time.Second
const sportsScreenDuration = 5 * time.Second

type SportsUpcomingGamesScreen struct {
	Ctx                 *gg.Context
	Fonts               *fonts.Fonts
	SportsData          sports.SportsData
	Event               sports.Event
	Logo1               image.Image
	Logo2               image.Image
	Team1               *sports.Team
	Team2               *sports.Team
	AnimationDone       bool
	AnimationPercent    float64
	TransitionDone      bool
	ScreenDisplayedTime time.Time
	LogoAnimation       *animation.Animation
	TextAnimation       *animation.Animation
}

func NewSportsUpcomingGamesScreen(fonts *fonts.Fonts, sportsData sports.SportsData, event sports.Event) Screen {
	return &SportsUpcomingGamesScreen{
		Ctx:                 gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		Fonts:               fonts,
		SportsData:          sportsData,
		Event:               event,
		ScreenDisplayedTime: time.Now(),
		LogoAnimation: animation.NewAnimation(2*time.Second, map[string]animation.AnimationValue{
			"x1": {Start: 16, End: -16},
			"x2": {Start: 48, End: 80},
		}),
		TextAnimation: animation.NewAnimation(500*time.Millisecond, map[string]animation.AnimationValue{
			"r": {Start: 0, End: 255},
			"g": {Start: 0, End: 255},
			"b": {Start: 0, End: 255},
		}),
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

		s.AnimationDone = false
		s.AnimationPercent = 0
		s.TransitionDone = false

		close(doneChan)
	}()

	return doneChan
}

func (s *SportsUpcomingGamesScreen) Render(elapsed time.Duration) (image.Image, bool) {
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	isLogoAnimationDone := s.TransitionDone && s.LogoAnimation.IsDone(time.Since(s.ScreenDisplayedTime))
	if !isLogoAnimationDone {
		return s.renderAnimation()
	}

	if s.Logo1 != nil {
		s.Ctx.DrawImageAnchored(s.Logo1, -16, 0, 0, 0)
	}
	if s.Logo2 != nil {
		s.Ctx.DrawImageAnchored(s.Logo2, 80, 0, 1, 0)
	}

	eventTime, err := time.Parse("2006-01-02T15:04:05", s.Event.Timestamp)
	if err == nil {
		eventTime = eventTime.Local()

		colors := s.TextAnimation.GetValues(time.Since(s.ScreenDisplayedTime.Add(s.LogoAnimation.Duration)))
		s.Ctx.SetColor(color.RGBA{uint8(colors["r"]), uint8(colors["g"]), uint8(colors["b"]), 255})

		s.Ctx.SetFontFace(s.Fonts.Size5x7)
		s.Ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Mon")), 32, 1, 0.5, 1)
		s.Ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Jan 2")), 32, 8, 0.5, 1)
		s.Ctx.DrawStringAnchored(eventTime.Format("3:04"), 32, 15, 0.5, 1)
		s.Ctx.DrawStringAnchored(eventTime.Format("PM"), 32, 22, 0.5, 1)
	} else {
		log.Printf("Error reading event time from timestamp '%s'.", s.Event.Timestamp)
	}

	return s.Ctx.Image(), time.Since(s.ScreenDisplayedTime) > sportsScreenDuration
}

func (s *SportsUpcomingGamesScreen) renderAnimation() (image.Image, bool) {
	var currentAnimationValues map[string]float64
	if !s.TransitionDone {
		currentAnimationValues = s.LogoAnimation.GetValues(0)
	} else {
		currentAnimationValues = s.LogoAnimation.GetValues(time.Since(s.ScreenDisplayedTime))
	}

	x1 := int(currentAnimationValues["x1"])
	x2 := int(currentAnimationValues["x2"])

	// Draw the away team logo first so it is underneath the home logo
	if s.Logo2 != nil {
		resizedLogo2 := utils.ResizeImage(s.Logo2, 32)
		s.Ctx.DrawImageAnchored(resizedLogo2, x2, 0, 1, 0)
	}
	if s.Logo1 != nil {
		resizedLogo1 := utils.ResizeImage(s.Logo1, 32)
		s.Ctx.DrawImageAnchored(resizedLogo1, x1, 0, 0, 0)
	}

	return s.Ctx.Image(), time.Since(s.ScreenDisplayedTime) > sportsScreenDuration
}

func (s *SportsUpcomingGamesScreen) TransitionStart() {
	// Nothing to do here
}

func (s *SportsUpcomingGamesScreen) TransitionEnd(isDisplayed bool) {
	if isDisplayed {
		s.TransitionDone = true
		s.ScreenDisplayedTime = time.Now()
	}
}
