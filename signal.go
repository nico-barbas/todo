package main

import "fmt"

const listenerBufferCap = 10

const SignalNoArgs SignalInt = -1

type signalDispatcher struct {
	signals map[SignalKind][]SignalListener
}

func (s *signalDispatcher) init() {
	s.signals = make(map[SignalKind][]SignalListener)
}

func (s *signalDispatcher) addListener(k SignalKind, l SignalListener) {
	if _, exist := s.signals[k]; !exist {
		s.signals[k] = make([]SignalListener, 0, listenerBufferCap)
	}
	s.signals[k] = append(s.signals[k], l)
}

func (s *signalDispatcher) dispatch(k SignalKind, v SignalValue) {
	if _, exist := s.signals[k]; !exist {
		// Silently returning, probably not the best
		return
	}
	for _, l := range s.signals[k] {
		l.OnSignal(Signal{
			Value: v,
			Kind:  k,
		})
	}
}

type Signal struct {
	Value SignalValue
	Kind  SignalKind
}

type SignalListener interface {
	OnSignal(s Signal)
}

// The actual enum def is left to the user
type SignalKind uint8

type SignalValue interface {
	ToString() string
}

// Some general purpose type wrappers
type (
	SignalFloat  float64
	SignalInt    int
	SignalString string
	SignalArray  []SignalValue
)

func (f SignalFloat) ToString() string  { return fmt.Sprint(f) }
func (i SignalInt) ToString() string    { return fmt.Sprint(i) }
func (s SignalString) ToString() string { return fmt.Sprint(s) }
func (a SignalArray) ToString() string  { return fmt.Sprint(a) }
