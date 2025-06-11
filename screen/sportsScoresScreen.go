package screen

import (
	"image"
	"image/color"
	"log"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/sports"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/utils"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

const animationDuration = 2 * time.Second
const sportsScreenDuration = 5 * time.Second

type SportsScoresScreen struct {
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
}

func NewSportsScoresScreen(fonts *fonts.Fonts, sportsData sports.SportsData, event sports.Event) Screen {
	return &SportsScoresScreen{
		Ctx:                 gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		Fonts:               fonts,
		SportsData:          sportsData,
		Event:               event,
		ScreenDisplayedTime: time.Now(),
	}
}

func (s *SportsScoresScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		if s.Logo1 == nil {
			s.Logo1 = s.SportsData.GetLogo(s.Event.HomeTeamName)
		}
		if s.Logo2 == nil {
			s.Logo2 = s.SportsData.GetLogo(s.Event.AwayTeamName)
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

func (s *SportsScoresScreen) Render(elapsed time.Duration) (image.Image, bool) {
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	if !s.AnimationDone {
		if s.TransitionDone {
			return s.renderAnimation(elapsed)
		}
		return s.renderAnimation(0)
	}

	if s.Logo1 != nil {
		resizedLogo1 := utils.ResizeImage(s.Logo1, 32)
		s.Ctx.DrawImageAnchored(resizedLogo1, -16, 0, 0, 0)
	}
	if s.Logo2 != nil {
		resizedLogo2 := utils.ResizeImage(s.Logo2, 32)
		s.Ctx.DrawImageAnchored(resizedLogo2, 64+16, 0, 1, 0)
	}

	eventTime, err := time.Parse("2006-01-02T15:04:05", s.Event.Timestamp)
	if err == nil {
		eventTime = eventTime.Local()

		s.Ctx.SetFontFace(s.Fonts.Size5x7)
		s.Ctx.SetColor(color.White)
		s.Ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Mon")), 32, 1, 0.5, 1)
		s.Ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Jan 2")), 32, 8, 0.5, 1)
		s.Ctx.DrawStringAnchored(eventTime.Format("3:04"), 32, 15, 0.5, 1)
		s.Ctx.DrawStringAnchored(eventTime.Format("PM"), 32, 22, 0.5, 1)
	} else {
		log.Printf("Error reading event time from timestamp '%s'.", s.Event.Timestamp)
	}

	return s.Ctx.Image(), time.Now().Sub(s.ScreenDisplayedTime) > sportsScreenDuration
}

func (s *SportsScoresScreen) renderAnimation(elapsed time.Duration) (image.Image, bool) {
	s.AnimationPercent += float64(elapsed.Milliseconds()) / float64(animationDuration.Milliseconds())
	if s.AnimationPercent >= 1 {
		s.AnimationPercent = 1
		s.AnimationDone = true
	}

	x1 := 16 - int(32*s.AnimationPercent)
	x2 := 48 + int(32*s.AnimationPercent)

	// Draw the away team logo first so it is underneath the home logo
	if s.Logo2 != nil {
		resizedLogo2 := utils.ResizeImage(s.Logo2, 32)
		s.Ctx.DrawImageAnchored(resizedLogo2, x2, 0, 1, 0)
	}
	if s.Logo1 != nil {
		resizedLogo1 := utils.ResizeImage(s.Logo1, 32)
		s.Ctx.DrawImageAnchored(resizedLogo1, x1, 0, 0, 0)
	}

	return s.Ctx.Image(), time.Now().Sub(s.ScreenDisplayedTime) > sportsScreenDuration
}

func (s *SportsScoresScreen) TransitionStart() {
	// Nothing to do here
}

func (s *SportsScoresScreen) TransitionEnd(isDisplayed bool) {
	if isDisplayed {
		s.TransitionDone = true
		s.ScreenDisplayedTime = time.Now()
	}
}
