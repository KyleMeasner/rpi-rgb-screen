package animation

import (
	"fmt"
	"image"
	"image/color"
	"slices"
	"time"
)

type KeyFrames struct {
	StartTime time.Time
	Points    map[string]*AnimatedPoint
	Colors    map[string]*AnimatedColor

	totalDuration time.Duration
}

type AnimatedPoint struct {
	StartValue  image.Point
	Transitions []AnimatedPointTransition
}

type AnimatedPointTransition struct {
	Offset   time.Duration
	Duration time.Duration
	EndValue image.Point
}

type AnimatedColor struct {
	StartValue  color.RGBA
	Transitions []AnimatedColorTransition
}

type AnimatedColorTransition struct {
	Offset   time.Duration
	Duration time.Duration
	EndValue color.RGBA
}

func NewKeyFrames() *KeyFrames {
	return &KeyFrames{
		Points: map[string]*AnimatedPoint{},
		Colors: map[string]*AnimatedColor{},
	}
}

func (k *KeyFrames) Start() {
	k.StartTime = time.Now()
}

func (k *KeyFrames) HasStarted() bool {
	return !k.StartTime.IsZero()
}

func (k *KeyFrames) RunTime() time.Duration {
	if !k.HasStarted() {
		return 0
	}
	return time.Since(k.StartTime)
}

func (k *KeyFrames) HasEnded() bool {
	if !k.HasStarted() {
		return false
	}
	return k.RunTime() > k.totalDuration
}

func (k *KeyFrames) AddPoint(key string, startValue image.Point) error {
	if _, ok := k.Points[key]; ok {
		return fmt.Errorf("KeyFrames already has a Point registered under key '%s'", key)
	}

	k.Points[key] = &AnimatedPoint{
		StartValue:  startValue,
		Transitions: []AnimatedPointTransition{},
	}

	return nil
}

func (k *KeyFrames) AddPointTransitions(key string, transitions ...AnimatedPointTransition) error {
	animatedPoint, ok := k.Points[key]
	if !ok {
		return fmt.Errorf("KeyFrames does not have a Point registered under key '%s'", key)
	}

	for _, transition := range transitions {
		transitionOverlaps := slices.ContainsFunc(animatedPoint.Transitions, func(existing AnimatedPointTransition) bool {
			return transition.Offset >= existing.Offset && transition.Offset < existing.Offset+existing.Duration
		})
		if transitionOverlaps {
			return fmt.Errorf(
				"unable to add Point transition to key '%s' due to overlap with existing transition. Offset of new transition: %d", key, transition.Offset)
		}

		animatedPoint.Transitions = append(animatedPoint.Transitions, transition)

		if k.totalDuration < transition.Offset+transition.Duration {
			k.totalDuration = transition.Offset + transition.Duration
		}
	}

	slices.SortFunc(animatedPoint.Transitions, func(a, b AnimatedPointTransition) int {
		return int(a.Offset - b.Offset)
	})
	return nil
}

func (k *KeyFrames) GetPoint(key string) image.Point {
	animatedPoint, ok := k.Points[key]
	if !ok {
		return image.Point{}
	}

	if !k.HasStarted() {
		return animatedPoint.StartValue
	}

	timeSinceStart := time.Since(k.StartTime)
	currValue := animatedPoint.StartValue
	for _, transition := range animatedPoint.Transitions {
		// All transitions past this point start after the current time
		if timeSinceStart < transition.Offset {
			return currValue
		}

		// We're in the middle of this transition
		if timeSinceStart >= transition.Offset && timeSinceStart < transition.Offset+transition.Duration {
			percentComplete := float64(timeSinceStart-transition.Offset) / float64(transition.Duration)
			return image.Point{
				X: computeValue(currValue.X, transition.EndValue.X, percentComplete),
				Y: computeValue(currValue.Y, transition.EndValue.Y, percentComplete),
			}
		}

		// This transition has already ended
		currValue = transition.EndValue
	}

	return currValue
}

func (k *KeyFrames) AddColor(key string, startValue color.RGBA) error {
	if _, ok := k.Colors[key]; ok {
		return fmt.Errorf("KeyFrames already has a Color registered under key '%s'", key)
	}

	k.Colors[key] = &AnimatedColor{
		StartValue:  startValue,
		Transitions: []AnimatedColorTransition{},
	}

	return nil
}

func (k *KeyFrames) AddColorTransitions(key string, transitions ...AnimatedColorTransition) error {
	animatedColor, ok := k.Colors[key]
	if !ok {
		return fmt.Errorf("KeyFrames does not have a Color registered under key '%s'", key)
	}

	for _, transition := range transitions {
		transitionOverlaps := slices.ContainsFunc(animatedColor.Transitions, func(existing AnimatedColorTransition) bool {
			return transition.Offset >= existing.Offset && transition.Offset < existing.Offset+existing.Duration
		})
		if transitionOverlaps {
			return fmt.Errorf(
				"unable to add Color transition to key '%s' due to overlap with existing transition. Offset of new transition: %d", key, transition.Offset)
		}

		animatedColor.Transitions = append(animatedColor.Transitions, transition)

		if k.totalDuration < transition.Offset+transition.Duration {
			k.totalDuration = transition.Offset + transition.Duration
		}
	}

	slices.SortFunc(animatedColor.Transitions, func(a, b AnimatedColorTransition) int {
		return int(a.Offset - b.Offset)
	})
	return nil
}

func (k *KeyFrames) GetColor(key string) color.RGBA {
	animatedColor, ok := k.Colors[key]
	if !ok {
		return color.RGBA{}
	}

	if !k.HasStarted() {
		return animatedColor.StartValue
	}

	timeSinceStart := time.Since(k.StartTime)
	currValue := animatedColor.StartValue
	for _, transition := range animatedColor.Transitions {
		// All transitions past this point start after the current time
		if timeSinceStart < transition.Offset {
			return currValue
		}

		// We're in the middle of this transition
		if timeSinceStart >= transition.Offset && timeSinceStart < transition.Offset+transition.Duration {
			percentComplete := float64(timeSinceStart-transition.Offset) / float64(transition.Duration)
			return color.RGBA{
				R: uint8(computeValue(int(currValue.R), int(transition.EndValue.R), percentComplete)),
				G: uint8(computeValue(int(currValue.G), int(transition.EndValue.G), percentComplete)),
				B: uint8(computeValue(int(currValue.B), int(transition.EndValue.B), percentComplete)),
				A: uint8(computeValue(int(currValue.A), int(transition.EndValue.A), percentComplete)),
			}
		}

		// This transition has already ended
		currValue = transition.EndValue
	}

	return currValue
}

func computeValue(start, end int, percentComplete float64) int {
	return start + int(float64(end-start)*percentComplete)
}
