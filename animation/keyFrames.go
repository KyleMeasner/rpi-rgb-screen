package animation

import (
	"image"
	"image/color"
	"time"
)

type KeyFrames struct {
	Animations map[string]*Animation
	StartTime  time.Time
}

func NewKeyFrames(animations map[string]*Animation) *KeyFrames {
	return &KeyFrames{
		Animations: animations,
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

func (k *KeyFrames) GetActiveAnimations() map[string]*Animation {
	if !k.HasStarted() {
		return map[string]*Animation{}
	}

	active := map[string]*Animation{}
	for key, animation := range k.Animations {
		animationStart := k.StartTime.Add(animation.StartAt)
		animationEnd := animationStart.Add(animation.Duration)
		if animationStart.Before(time.Now()) && animationEnd.After(time.Now()) {
			active[key] = animation
		}
	}
	return active
}

func (k *KeyFrames) GetNumber(animationName string, key string) int {
	animation, ok := k.Animations[animationName]
	if !ok {
		return 0
	}

	timeSinceAnimationStarted := time.Since(k.StartTime.Add(animation.StartAt))
	return animation.GetNumber(key, timeSinceAnimationStarted)
}

func (k *KeyFrames) GetColor(animationName string, key string) color.Color {
	animation, ok := k.Animations[animationName]
	if !ok {
		return color.RGBA{}
	}

	timeSinceAnimationStarted := time.Since(k.StartTime.Add(animation.StartAt))
	return animation.GetColor(key, timeSinceAnimationStarted)
}

func (k *KeyFrames) GetPoint(animationName string, key string) image.Point {
	animation, ok := k.Animations[animationName]
	if !ok {
		return image.Point{}
	}

	timeSinceAnimationStarted := time.Since(k.StartTime.Add(animation.StartAt))
	return animation.GetPoint(key, timeSinceAnimationStarted)
}
