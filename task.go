package main

import "math"

const (
	minSessionLength minute = 1
	minSessionCount  int    = 1
)

const (
	taskStateIdle taskState = iota
	taskStateWork
	taskStateRest
)

type (
	task struct {
		name  string
		id    int
		done  bool
		state taskState

		sessionRequired  int
		sessionCompleted int
		sessionLength    minute
		restLength       minute

		timer timer
	}

	taskState int
)

func (t *task) init() {
	t.timer.setDuration(t.sessionLength, 0)
}

func (t *task) update() {
	if t.timer.running {
		if finished := t.timer.advance(); finished {
			switch t.state {
			case taskStateWork:
				t.state = taskStateRest
				t.timer.setDuration(t.restLength, 0)
			case taskStateRest:
				t.state = taskStateIdle
				t.completeSession()
				t.timer.setDuration(t.sessionLength, 0)
			default:
				// invalid state
			}
		}
	}
}

func (t *task) completeSession() {
	t.sessionCompleted += 1
	if t.sessionCompleted == t.sessionRequired {
		t.done = true
	}
	t.timer.setDuration(t.sessionLength, 0)
}

func (t *task) resetDuration() {
	t.timer.setDuration(t.sessionLength, 0)
}

func (t *task) startWork() {
	t.timer.start()
}

func (t *task) stopWork() {
	t.timer.stop()
}

func (t task) ToString() string {
	return "task"
}

func (t task) progress() float64 {
	sProg := float64(t.timer.sec) / 60
	prog := math.Abs((float64(t.timer.min) + sProg) - float64(t.sessionLength))
	return prog / float64(t.sessionLength)
}
