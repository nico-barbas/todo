package main

const (
	minSessionLength minute = 1
)

type (
	task struct {
		name       string
		done       bool
		charBuffer []byte

		sessionRequired  int
		sessionCompleted int
		sessionLength    minute

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
