package screen

import (
	"image"
	"image/color"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/sports"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/transition"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)

type SportsStandingsScreen struct {
	State      ScreenState
	Ctx        *gg.Context
	KeyFrames  *animation.KeyFrames
	Fonts      *fonts.Fonts
	SportsData sports.SportsData
	Conference *sports.Conference
	LeagueId   int
}

func NewSportsStandingsScreen(fonts *fonts.Fonts, sportsData sports.SportsData, leagueId int, conference *sports.Conference) Screen {
	keyFrames := animation.NewKeyFrames(time.Duration(len(conference.Standings)) * 3 * time.Second)

	keyFrames.AddNumber("stat", 0)
	for i := range len(conference.Standings) {
		if i == 0 {
			continue
		}

		keyFrames.AddNumberTransitions("stat", animation.AnimatedNumberTransition{
			Offset:   3 * time.Duration(i) * time.Second,
			Duration: 0,
			EndValue: i,
		})
	}

	return &SportsStandingsScreen{
		State:      StateNotDisplayed,
		Ctx:        gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:  keyFrames,
		Fonts:      fonts,
		SportsData: sportsData,
		Conference: conference,
		LeagueId:   leagueId,
	}
}

func (s *SportsStandingsScreen) GetPreferredTransition() transition.Transition {
	return transition.NewSlideAndZoomTransition()
}

func (s *SportsStandingsScreen) SetState(state ScreenState) {
	s.State = state

	switch state {
	case StateDisplayed:
		s.KeyFrames.Start()
	case StateTransitionIn:
		s.KeyFrames.Reset()
	}
}

func (s *SportsStandingsScreen) Refresh() chan bool {
	doneChan := make(chan bool)
	go func() {
		close(doneChan)
	}()
	return doneChan
}

func (s *SportsStandingsScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	s.Ctx.SetFontFace(s.Fonts.Size4x6)
	s.Ctx.SetColor(color.White)
	conferenceAbbreviation := getConferenceAbbreviation(s.Conference.Name)
	if conferenceAbbreviation == "" {
		conferenceAbbreviation = constants.LEAGUE_NAMES[s.LeagueId]
	}
	s.Ctx.DrawStringAnchored(conferenceAbbreviation, constants.SCREEN_WIDTH, 0, 1, 1)

	s.Ctx.SetFontFace(s.Fonts.Size4x6)
	s.Ctx.SetColor(color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}) // #808080
	positions := getStatPositions(s.LeagueId)
	headings := getHeadings(s.LeagueId)
	for i, heading := range headings {
		s.Ctx.DrawStringAnchored(heading, positions[i], 17, 0.5, 0.5)
	}

	s.Ctx.SetColor(color.RGBA{R: 0xA6, G: 0xA6, B: 0xA6, A: 0xFF}) // #A6A6A6
	for i := range constants.SCREEN_WIDTH {
		s.Ctx.SetPixel(i, 13)
		s.Ctx.SetPixel(i, 21)
	}

	currentStanding := s.KeyFrames.GetNumber("stat")

	standing := s.Conference.Standings[currentStanding]
	s.Ctx.SetFontFace(s.Fonts.Size8x13B)
	s.Ctx.SetColor(color.White)
	s.Ctx.DrawStringAnchored(strconv.Itoa(standing.Rank), 9, -2, 0.5, 1)

	s.Ctx.SetFontFace(s.Fonts.Size6x10)
	teamName := s.SportsData.GetTeamShortName(standing.TeamId)
	s.Ctx.DrawStringAnchored(teamName, 18, 1, 0, 1)

	s.Ctx.SetFontFace(s.Fonts.Size5x7)
	toDisplay := getStats(s.LeagueId, standing, standing.RegulationWins+standing.OvertimeWins)
	for i, stat := range toDisplay {
		s.Ctx.DrawStringAnchored(strconv.Itoa(stat), positions[i], 27, 0.5, 0.5)
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func getConferenceAbbreviation(conferenceName string) string {
	switch conferenceName {
	// MLB
	case "American League":
		return "AL"
	case "National League":
		return "NL"
	// NHL, MLB, WNBA
	case "Western Conference":
		return "WEST"
	case "Eastern Conference":
		return "EAST"
	// NFL
	case "American Football":
		return "AFC"
	case "National Football":
		return "NFC"
	// CFL
	case "East Division":
		return "EAST"
	case "West Division":
		return "WEST"
	// PWHL, NWSL
	default:
		return ""
	}
}

func getStatPositions(leagueId int) []float64 {
	switch leagueId {
	case constants.LEAGUE_NHL:
		fallthrough
	case constants.LEAGUE_CFL:
		return []float64{6, 20, 34, 47, 59}
	case constants.LEAGUE_NFL:
		fallthrough
	case constants.LEAGUE_MLS:
		fallthrough
	case constants.LEAGUE_NWSL:
		fallthrough
	case constants.LEAGUE_WNBA:
		fallthrough
	case constants.LEAGUE_MLB:
		return []float64{9, 26, 43, 59}
	case constants.LEAGUE_PWHL:
		return []float64{6, 19, 32, 45, 58, 71}
	default:
		return nil
	}
}

func getHeadings(leagueId int) []string {
	switch leagueId {
	case constants.LEAGUE_NHL:
		return []string{"GP", "PTS", "W", "L", "OT"}
	case constants.LEAGUE_NFL:
		return []string{"GP", "W", "L", "T"}
	case constants.LEAGUE_CFL:
		return []string{"GP", "PTS", "W", "L", "T"}
	case constants.LEAGUE_MLS:
		return []string{"GP", "W", "D", "L"}
	case constants.LEAGUE_NWSL:
		return []string{"GP", "W", "D", "L"}
	case constants.LEAGUE_WNBA:
		return []string{"GP", "W", "L", "GB"}
	case constants.LEAGUE_MLB:
		return []string{"GP", "W", "L", "GB"}
	case constants.LEAGUE_PWHL:
		return []string{"GP", "PTS", "W", "OTW", "OTL", "L"}
	default:
		return nil
	}
}

func getStats(leagueId int, standing *sports.Standing, firstPlaceWins int) []int {
	switch leagueId {
	case constants.LEAGUE_NHL:
		return []int{standing.GamesPlayed, standing.Points, standing.RegulationWins + standing.OvertimeWins, standing.RegulationLosses, standing.OvertimeLosses}
	case constants.LEAGUE_NFL:
		return []int{standing.GamesPlayed, standing.RegulationWins + standing.OvertimeWins, standing.RegulationLosses + standing.OvertimeLosses, standing.Ties}
	case constants.LEAGUE_CFL:
		return []int{standing.GamesPlayed, standing.Points, standing.RegulationWins + standing.OvertimeWins, standing.RegulationLosses + standing.OvertimeLosses, standing.Ties}
	case constants.LEAGUE_MLS:
		return []int{standing.GamesPlayed, standing.RegulationWins + standing.OvertimeWins, standing.Ties, standing.RegulationLosses + standing.OvertimeLosses}
	case constants.LEAGUE_NWSL:
		return []int{standing.GamesPlayed, standing.RegulationWins + standing.OvertimeWins, standing.Ties, standing.RegulationLosses + standing.OvertimeLosses}
	case constants.LEAGUE_WNBA:
		return []int{standing.GamesPlayed, standing.RegulationWins + standing.OvertimeWins, standing.RegulationLosses + standing.OvertimeLosses, firstPlaceWins - (standing.RegulationWins + standing.OvertimeWins)}
	case constants.LEAGUE_MLB:
		return []int{standing.GamesPlayed, standing.RegulationWins + standing.OvertimeWins, standing.RegulationLosses + standing.OvertimeLosses, firstPlaceWins - (standing.RegulationWins + standing.OvertimeWins)}
	case constants.LEAGUE_PWHL:
		return []int{standing.GamesPlayed, standing.Points, standing.RegulationWins + standing.OvertimeWins, standing.OvertimeWins, standing.OvertimeLosses, standing.RegulationLosses}
	default:
		return nil
	}
}
