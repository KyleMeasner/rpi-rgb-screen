package screen

import (
	"image"
	"image/color"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/sports"
	"rpi-rgb-screen/fonts"

	"github.com/fogleman/gg"
)

type SportsScoresScreen struct {
	ctx        *gg.Context
	fonts      *fonts.Fonts
	sportsData sports.SportsData
	event      sports.Event
	logo1      image.Image
	logo2      image.Image
	team1      *sports.Team
	team2      *sports.Team
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

func (s *SportsScoresScreen) Render() image.Image {
	s.ctx.Identity()
	s.ctx.SetColor(color.Black)
	s.ctx.Clear()

	if s.logo1 != nil {
		s.ctx.DrawImageAnchored(s.logo1, 0, 32, 0, 1)
	}
	if s.logo2 != nil {
		s.ctx.DrawImageAnchored(s.logo2, 64, 32, 1, 1)
	}

	if s.team1 != nil {
		s.ctx.SetFontFace(s.fonts.Lemon)
		s.ctx.SetColor(color.White)
		s.ctx.DrawStringAnchored(s.team1.ShortName, 0, 0, 0, 1)
	}
	if s.team2 != nil {
		s.ctx.SetFontFace(s.fonts.Lemon)
		s.ctx.SetColor(color.White)
		s.ctx.DrawStringAnchored(s.team2.ShortName, 64, 0, 1, 1)
	}

	return s.ctx.Image()
}
