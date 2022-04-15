package anim

import "math"

const (
	animationInitialSize = 10
)

type (
	Animation struct {
		name       string
		Playing    bool
		callback   AnimationCallback
		timer      int
		properties []AnimationProperty
	}

	AnimationProperty struct {
		name              string
		startValue        float64
		resetToStartValue bool
		property          *float64
		keys              []AnimationKey
		keyStartValue     float64
		keyIndex          int
	}

	AnimationKey struct {
		Easing    EaseKind
		StartTime int
		Duration  int
		Change    float64
	}

	AnimationCallback interface {
		OnAnimationEnd(name string)
	}

	EaseKind int
)

const (
	EaseLinear EaseKind = iota
	EaseInCubic
	EaseOutCubic
)

func ease(e EaseKind, start, change float64, time, duration int) (result float64) {
	scale := float64(time) / float64(duration)
	if scale > 1 {
		scale = 1
	}
	switch e {
	case EaseLinear:
		result = start + (change * scale)

	case EaseInCubic:
		result = start + (change * math.Pow(scale, 3))

	case EaseOutCubic:
		result = start + (change * (math.Pow(scale-1, 3) + 1))
	}
	return
}

func SecondsToTicks(t float64) int {
	return int(t * 60)
}

func (a *AnimationProperty) advance(timer int) (finished bool) {
	length := len(a.keys)
	if a.keyIndex >= length {
		return true
	}
	key := &a.keys[a.keyIndex]
	keyTime := timer - key.StartTime
	*a.property = ease(
		key.Easing,
		a.keyStartValue, key.Change,
		keyTime, key.Duration,
	)
	if keyTime > key.Duration {
		a.keyIndex += 1
		a.keyStartValue = *a.property
		if a.keyIndex >= length {
			finished = true
		}
	}
	return
}

func (a *AnimationProperty) reset() {
	a.keyIndex = 0
	a.keyStartValue = 0
	if a.resetToStartValue {
		*a.property = a.startValue
	}
}

func NewAnimation(name string, callback AnimationCallback) Animation {
	return Animation{
		name:       name,
		callback:   callback,
		properties: make([]AnimationProperty, 0, animationInitialSize),
	}
}

// could return an error
func (a *Animation) AddProperty(name string, propertyRef *float64, startValue float64, reset bool) {
	var exist bool
	for index := range a.properties {
		if a.properties[index].name == name {
			exist = true
			break
		}
	}

	if !exist {
		a.properties = append(a.properties, AnimationProperty{
			name:              name,
			startValue:        startValue,
			resetToStartValue: reset,
			keyStartValue:     startValue,
			property:          propertyRef,
			keys:              make([]AnimationKey, 0, animationInitialSize),
		})
	}
}

func (a *Animation) SetPropertyRef(name string, ref *float64) {
	for index := range a.properties {
		property := &a.properties[index]
		if property.name == name {
			property.property = ref
			break
		}
	}
}

func (a *Animation) AddKey(name string, key AnimationKey) {
	for index := range a.properties {
		property := &a.properties[index]
		if property.name == name {
			length := len(property.keys)
			if length > 1 {
				key.StartTime = property.keys[length-1].StartTime + property.keys[length-1].Duration
			}
			property.keys = append(property.keys, key)
			break
		}
	}
}

func (a *Animation) Play() {
	a.Playing = true
}

func (a *Animation) Reset() {
	for index := range a.properties {
		a.properties[index].reset()
	}
}

func (a *Animation) Update() {
	if !a.Playing {
		return
	}
	a.timer += 1

	finished := false
	for index := range a.properties {
		keyDone := a.properties[index].advance(a.timer)
		if !finished {
			finished = keyDone
		}
	}

	if finished {
		a.timer = 0
		for index := range a.properties {
			a.properties[index].reset()
		}
		a.Playing = false
		a.callback.OnAnimationEnd(a.name)
	}
}
