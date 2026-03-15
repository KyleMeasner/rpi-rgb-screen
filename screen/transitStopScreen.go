package screen

import (
	"fmt"
	"image"
	"image/color"
	"rpi-rgb-screen/animation"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data/transit"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/transition"
	"rpi-rgb-screen/utils"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type TransitStopScreen struct {
	State       ScreenState
	Ctx         *gg.Context
	KeyFrames   *animation.KeyFrames
	Fonts       *fonts.Fonts
	TransitData transit.TransitData
	StopId      string
	Arrivals    []*transit.Arrival
	Routes      map[string]*transit.Route
}

func NewTransitStopScreen(fonts *fonts.Fonts, transitData transit.TransitData, stopId string) Screen {
	keyFrames := animation.NewKeyFrames(10 * time.Second)

	return &TransitStopScreen{
		State:       StateNotDisplayed,
		Ctx:         gg.NewContext(constants.SCREEN_WIDTH, constants.SCREEN_HEIGHT),
		KeyFrames:   keyFrames,
		Fonts:       fonts,
		TransitData: transitData,
		StopId:      stopId,
		Routes:      map[string]*transit.Route{},
	}
}

func (s *TransitStopScreen) GetPreferredTransition() transition.Transition {
	return transition.NewSlideAndZoomTransition()
}

func (s *TransitStopScreen) SetState(state ScreenState) {
	s.State = state

	switch state {
	case StateDisplayed:
		s.KeyFrames.Start()
	case StateTransitionIn:
		s.KeyFrames.Reset()
	}
}

func (s *TransitStopScreen) Refresh() chan bool {
	doneChan := make(chan bool)

	go func() {
		s.Arrivals = s.TransitData.GetArrivalsForStop(s.StopId)

		// Fetch route data for all arrivals
		s.Routes = make(map[string]*transit.Route)
		for _, arrival := range s.Arrivals {
			if _, exists := s.Routes[arrival.RouteId]; !exists {
				s.Routes[arrival.RouteId] = s.TransitData.GetRoute(arrival.RouteId)
			}
		}

		close(doneChan)
	}()

	return doneChan
}

func (s *TransitStopScreen) Render() (image.Image, bool) {
	// Clear image context
	s.Ctx.Identity()
	s.Ctx.SetColor(color.Black)
	s.Ctx.Clear()

	// Draw arrivals
	for i, arrival := range s.Arrivals {
		s.renderArrival(arrival, float64(i*16))
	}

	// Draw "no arrivals" message if empty
	if len(s.Arrivals) == 0 {
		s.Ctx.SetFontFace(s.Fonts.Size4x6)
		s.Ctx.SetColor(color.RGBA{200, 100, 100, 255})
		s.Ctx.DrawStringAnchored("No arrivals", 32, 16, 0.5, 0.5)
	}

	return s.Ctx.Image(), s.KeyFrames.HasEnded()
}

func (s *TransitStopScreen) renderArrival(arrival *transit.Arrival, offset float64) {

	// Draw route number
	s.Ctx.SetFontFace(s.Fonts.Size5x7)
	route := s.Routes[arrival.RouteId]
	s.Ctx.SetColor(route.Color)

	if route.Type == 0 { // Train
		s.Ctx.DrawCircle(5, offset+8, 5)
		s.Ctx.Fill()
		s.Ctx.SetColor(color.White)
	}

	routeNumber := strings.SplitN(arrival.RouteName, " ", 2)[0]
	s.Ctx.DrawStringAnchored(routeNumber, 6, offset+8, 0.5, 0.5)

	// Draw headsign (destination)
	s.Ctx.SetColor(color.White)
	s.Ctx.SetFontFace(s.Fonts.Size4x6)
	headsignText := arrival.Headsign
	if len(headsignText) > 12 {
		headsignText = headsignText[:12] + "…"
	}
	s.Ctx.DrawStringAnchored(headsignText, 12, offset, 0, 1)

	// Draw time until arrival
	s.Ctx.SetFontFace(s.Fonts.Size5x7)
	s.Ctx.SetColor(color.RGBA{100, 200, 255, 255}) // #64C8FF

	isPredictedTime := true
	arrivalTime := arrival.PredictedTime
	if arrivalTime.Unix() == 0 {
		isPredictedTime = false
		arrivalTime = arrival.ScheduledTime
	}

	timeUntil := time.Until(arrivalTime)
	timeText := "Now"
	if timeUntil >= time.Minute {
		minutes := int(timeUntil.Minutes())
		timeText = fmt.Sprintf("%dm", minutes)
	}
	s.Ctx.DrawStringAnchored(timeText, 64, offset+7, 1, 1)

	if isPredictedTime {
		icon, _ := utils.ReadImageFromFile("./resources/transitIcons/predicted.png")
		s.Ctx.DrawImageAnchored(icon, 62-(len(timeText)*5), int(offset+8), 1, 0)
	}
}
