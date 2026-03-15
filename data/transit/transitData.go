package transit

import (
	"image/color"
	"rpi-rgb-screen/utils"
	"time"
)

type Arrival struct {
	RouteName     string
	RouteId       string
	Headsign      string
	PredictedTime time.Time
	ScheduledTime time.Time
}

type Stop struct {
	Id        string
	Name      string
	Direction string
}

type Route struct {
	Id    string
	Name  string
	Color color.Color
	Type  int
}

type TransitData interface {
	GetStop(stopId string) *Stop
	GetRoute(routeId string) *Route
	GetArrivalsForStop(stopId string) []*Arrival
}

type TransitDataManager struct {
	OneBusAwayClient *OneBusAwayClient
	Stops            *utils.ExpirableMap[string, *Stop]
	Routes           *utils.ExpirableMap[string, *Route]
}

func NewTransitData(useDummyData bool) TransitData {
	return &TransitDataManager{
		OneBusAwayClient: NewOneBusAwayClient(),
		Stops:            utils.NewExpirableMap[string, *Stop](time.Hour * 24),
		Routes:           utils.NewExpirableMap[string, *Route](time.Hour * 24),
	}
}

// GetArrivalsForStop implements TransitData.
func (t *TransitDataManager) GetArrivalsForStop(stopId string) []*Arrival {
	return t.OneBusAwayClient.GetArrivalsForStop(stopId)
}

// GetRoute implements TransitData.
func (t *TransitDataManager) GetRoute(routeId string) *Route {
	cachedRoute := t.Routes.Get(routeId)
	if cachedRoute != nil {
		return *cachedRoute
	}

	route := t.OneBusAwayClient.GetRoute(routeId)
	if route == nil {
		return nil
	}

	t.Routes.Set(routeId, route)
	return route
}

// GetStop implements TransitData.
func (t *TransitDataManager) GetStop(stopId string) *Stop {
	cachedStop := t.Stops.Get(stopId)
	if cachedStop != nil {
		return *cachedStop
	}

	stop := t.OneBusAwayClient.GetStop(stopId)
	if stop == nil {
		return nil
	}

	t.Stops.Set(stopId, stop)
	return stop
}
