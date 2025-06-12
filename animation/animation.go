package animation

import (
	"time"
)

type Animation struct {
	Duration time.Duration
	Values   map[string]AnimationValue
}

type AnimationValue struct {
	Start float64
	End   float64
}

func NewAnimation(duration time.Duration, values map[string]AnimationValue) *Animation {
	return &Animation{
		Duration: duration,
		Values:   values,
	}
}

func (a *Animation) GetValue(name string, timeSinceStart time.Duration) float64 {
	percentComplete := min(float64(timeSinceStart)/float64(a.Duration), 1)

	if value, ok := a.Values[name]; ok {
		return value.Start + ((value.End - value.Start) * percentComplete)
	}
	return 0
}

func (a *Animation) GetValues(timeSinceStart time.Duration) map[string]float64 {
	percentComplete := min(float64(timeSinceStart.Milliseconds())/float64(a.Duration.Milliseconds()), 1)
	currentValues := map[string]float64{}
	for key, value := range a.Values {
		currentValues[key] = value.Start + ((value.End - value.Start) * percentComplete)
	}
	return currentValues
}

func (a *Animation) IsDone(timeSinceStart time.Duration) bool {
	return timeSinceStart >= a.Duration
}
