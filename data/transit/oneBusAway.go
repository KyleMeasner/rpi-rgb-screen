package transit

import (
	"context"
	"image/color"
	"log"
	"rpi-rgb-screen/config"
	"rpi-rgb-screen/utils"
	"time"

	onebusaway "github.com/OneBusAway/go-sdk"
	"github.com/OneBusAway/go-sdk/option"
)

// Sound transit: 40
// 1 Line: 40_100479
// 2 Line: 40_2LINE
// Capitol Hill: 40_99603 (N), 40_99610 (S)

// Metro: 1
// Bus 8: 1_100275
// 12th & John: 1_29277 (E), 1_29301 (W)

type OneBusAwayClient struct {
	OneBusAway *onebusaway.Client
}

func NewOneBusAwayClient() *OneBusAwayClient {
	return &OneBusAwayClient{
		OneBusAway: onebusaway.NewClient(
			option.WithAPIKey(config.Config.OneBusAwayApiKey),
			option.WithBaseURL("https://api.pugetsound.onebusaway.org/"),
		),
	}
}

func (o *OneBusAwayClient) GetStop(stopId string) *Stop {
	response, err := o.OneBusAway.Stop.Get(context.Background(), stopId)
	if err != nil {
		log.Printf("Failed to get stop data for stop ID '%s'. Error: %s", stopId, err)
		return nil
	}

	stop := response.Data.Entry
	return &Stop{
		Id:        stop.ID,
		Name:      stop.Name,
		Direction: stop.Direction,
	}
}

func (o *OneBusAwayClient) GetRoute(routeId string) *Route {
	response, err := o.OneBusAway.Route.Get(context.Background(), routeId)
	if err != nil {
		log.Printf("Failed to get route data for route ID '%s'. Error: %s", routeId, err)
		return nil
	}

	route := response.Data.Entry

	var routeColor color.Color
	if route.Color != "" {
		routeColor, err = utils.ParseColorFromHex(route.Color)
		if err != nil {
			log.Printf("Failed to parse route color from color '%s' for route ID '%s'. Error: %s", route.Color, routeId, err)
			return nil
		}
	}

	return &Route{
		Id:    route.ID,
		Name:  route.ShortName,
		Color: routeColor,
		Type:  int(route.Type),
	}
}

func (o *OneBusAwayClient) GetArrivalsForStop(stopId string) []*Arrival {
	response, err := o.OneBusAway.ArrivalAndDeparture.List(context.Background(), stopId, onebusaway.ArrivalAndDepartureListParams{
		MinutesBefore: onebusaway.Int(1),
		MinutesAfter:  onebusaway.Int(35),
	})
	if err != nil {
		log.Printf("Failed to get arrival and departure data for stop '%s'. Error: %s", stopId, err)
		return []*Arrival{}
	}

	var arrivals []*Arrival
	for _, arrival := range response.Data.Entry.ArrivalsAndDepartures {
		arrival := &Arrival{
			RouteName:     arrival.RouteShortName,
			RouteId:       arrival.RouteID,
			Headsign:      arrival.TripHeadsign,
			PredictedTime: time.Unix(arrival.PredictedArrivalTime/1000, 0),
			ScheduledTime: time.Unix(arrival.ScheduledArrivalTime/1000, 0),
		}
		arrivals = append(arrivals, arrival)
	}
	return arrivals
}
