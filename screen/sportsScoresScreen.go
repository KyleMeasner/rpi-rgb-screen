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

const animationDuration = 3 * time.Second

type SportsScoresScreen struct {
	ctx              *gg.Context
	fonts            *fonts.Fonts
	sportsData       sports.SportsData
	event            sports.Event
	logo1            image.Image
	logo2            image.Image
	team1            *sports.Team
	team2            *sports.Team
	animationDone    bool
	animationPercent float64
	transitionDone   bool
}

func NewSportsScoresScreen(fonts *fonts.Fonts, sportsData sports.SportsData, event sports.Event) Screen {
	return &SportsScoresScreen{
		ctx:        gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		fonts:      fonts,
		sportsData: sportsData,
		event:      event,
	}
}

func (s *SportsScoresScreen) Refresh() {
	if s.logo1 == nil {
		s.logo1 = s.sportsData.GetLogo(s.event.HomeTeamName)
	}
	if s.logo2 == nil {
		s.logo2 = s.sportsData.GetLogo(s.event.AwayTeamName)
	}

	if s.team1 == nil {
		s.team1 = s.sportsData.GetTeam(s.event.HomeTeamName)
	}
	if s.team2 == nil {
		s.team2 = s.sportsData.GetTeam(s.event.AwayTeamName)
	}
}

func (s *SportsScoresScreen) Render(elapsed time.Duration) image.Image {
	s.ctx.Identity()
	s.ctx.SetColor(color.Black)
	s.ctx.Clear()

	if !s.animationDone {
		if s.transitionDone {
			return s.renderAnimation(elapsed)
		}
		return s.renderAnimation(0)
	}

	if s.logo1 != nil {
		resizedLogo1 := utils.ResizeImage(s.logo1, 32)
		s.ctx.DrawImageAnchored(resizedLogo1, -16, 0, 0, 0)
	}
	if s.logo2 != nil {
		resizedLogo2 := utils.ResizeImage(s.logo2, 32)
		s.ctx.DrawImageAnchored(resizedLogo2, 64+16, 0, 1, 0)
	}

	eventTime, err := time.Parse("2006-01-02T15:04:05", s.event.Timestamp)
	if err == nil {
		eventTime = eventTime.Local()

		s.ctx.SetFontFace(s.fonts.Size5x7)
		s.ctx.SetColor(color.White)
		s.ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Mon")), 32, 1, 0.5, 1)
		s.ctx.DrawStringAnchored(strings.ToUpper(eventTime.Format("Jan 2")), 32, 8, 0.5, 1)
		s.ctx.DrawStringAnchored(eventTime.Format("3:04"), 32, 15, 0.5, 1)
		s.ctx.DrawStringAnchored(eventTime.Format("PM"), 32, 22, 0.5, 1)
	} else {
		log.Printf("Error reading event time from timestamp '%s'.", s.event.Timestamp)
	}

	return s.ctx.Image()
}

func (s *SportsScoresScreen) renderAnimation(elapsed time.Duration) image.Image {
	s.animationPercent += float64(elapsed.Milliseconds()) / float64(animationDuration.Milliseconds())
	if s.animationPercent >= 1 {
		s.animationPercent = 1
		s.animationDone = true
	}

	x1 := 16 - int(32*s.animationPercent)
	x2 := 48 + int(32*s.animationPercent)

	// Draw the away team logo first so it is underneath the home logo
	if s.logo2 != nil {
		resizedLogo2 := utils.ResizeImage(s.logo2, 32)
		s.ctx.DrawImageAnchored(resizedLogo2, x2, 0, 1, 0)
	}
	if s.logo1 != nil {
		resizedLogo1 := utils.ResizeImage(s.logo1, 32)
		s.ctx.DrawImageAnchored(resizedLogo1, x1, 0, 0, 0)
	}

	return s.ctx.Image()
}

func (s *SportsScoresScreen) TransitionStart() {
	// Nothing to do here
}

func (s *SportsScoresScreen) TransitionEnd(isDisplayed bool) {
	s.transitionDone = true
}
