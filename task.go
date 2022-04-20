package main

import "fmt"

const (
	minSessionLength minute = 1
	minSessionCount  int    = 1
)

type (
	task struct {
		name       string
		id         int
		done       bool
		charBuffer []byte

		sessionRequired  int
		sessionCompleted int
		sessionLength    minute
		restLength       minute

		timer timer
	}
)

func (t *task) init() {
	t.timer.setDuration(t.sessionLength, 0)
}

func (t *task) completeSession() {
	t.sessionCompleted += 1
	if t.sessionCompleted == t.sessionRequired {
		t.done = true
	}
	t.timer.setDuration(t.sessionLength, 0)
}

func (t *task) startWork() {
	t.timer.setDuration(t.sessionLength, 0)
	t.timer.start()
}

func (t task) ToString() string {
	return "task"
}

func (t task) progress() float64 {
	sProg := float64(t.timer.sec) / 60
	prog := (float64(t.timer.min) + sProg)
	fmt.Println(prog)
	return prog/float64(t.sessionLength) - 1
}
