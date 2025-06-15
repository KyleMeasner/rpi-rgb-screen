package animation

import (
	"image"
	"image/color"
	"time"
)

type Animation struct {
	Duration time.Duration
	StartAt  time.Duration // Offset from start of KeyFrames for this animation to start
	Numbers  map[string]AnimationNumber
	Colors   map[string]AnimationColor
	Points   map[string]AnimationPoint
}

type AnimationNumber struct {
	Start int
	End   int
}

type AnimationColor struct {
	Start color.RGBA
	End   color.RGBA
}

type AnimationPoint struct {
	Start image.Point
	End   image.Point
}

func NewAnimation(startAt time.Duration, duration time.Duration) *Animation {
	return &Animation{
		Duration: duration,
		StartAt:  startAt,
		Numbers:  map[string]AnimationNumber{},
		Colors:   map[string]AnimationColor{},
		Points:   map[string]AnimationPoint{},
	}
}

func (a *Animation) IsDone(timeSinceStart time.Duration) bool {
	return timeSinceStart >= a.Duration
}

func (a *Animation) GetNumber(key string, timeSinceStart time.Duration) int {
	percent := percentComplete(timeSinceStart, a.Duration)
	if value, ok := a.Numbers[key]; ok {
		return computeValue(value.Start, value.End, percent)
	}
	return 0
}

func (a *Animation) GetColor(key string, timeSinceStart time.Duration) color.Color {
	percent := percentComplete(timeSinceStart, a.Duration)
	value, ok := a.Colors[key]
	if !ok {
		return color.RGBA{}
	}

	return color.RGBA{
		R: uint8(computeValue(int(value.Start.R), int(value.End.R), percent)),
		G: uint8(computeValue(int(value.Start.G), int(value.End.G), percent)),
		B: uint8(computeValue(int(value.Start.B), int(value.End.B), percent)),
		A: uint8(computeValue(int(value.Start.A), int(value.End.A), percent)),
	}
}

func (a *Animation) GetPoint(key string, timeSinceStart time.Duration) image.Point {
	percent := percentComplete(timeSinceStart, a.Duration)
	value, ok := a.Points[key]
	if !ok {
		return image.Point{}
	}

	return image.Point{
		X: computeValue(value.Start.X, value.End.X, percent),
		Y: computeValue(value.Start.Y, value.End.Y, percent),
	}
}

func computeValue(start, end int, percentComplete float64) int {
	return start + int(float64(end-start)*percentComplete)
}

func percentComplete(timeSinceStart, duration time.Duration) float64 {
	if timeSinceStart < 0 {
		return 0
	}
	if timeSinceStart > duration {
		return 1
	}
	return float64(timeSinceStart) / float64(duration)
}
