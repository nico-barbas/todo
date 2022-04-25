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
	taskStatePaused
	taskStateCount
)

var (
	taskTransitionTable = [taskStateCount * taskStateCount]func(t *task){
		// idle -> idle
		nil,
		// work -> idle
		nil,
		// rest -> idle
		func(t *task) {
			t.completeSession()
		},
		// paused -> idle
		nil,
		//
		//
		// idle -> work
		func(t *task) {
			// This means that the task was not previously paused
			t.timer.setDuration(t.sessionLength, 0)
			t.timer.start()
		},
		// work -> work
		nil,
		// rest -> work
		nil,
		// paused -> work
		func(t *task) {
			t.timer.start()
		},
		//
		//
		// idle -> rest
		func(t *task) {
			t.timer.start()
		},
		// work -> rest
		func(t *task) {
			t.timer.setDuration(t.restLength, 0)
			t.timer.start()
		},
		// rest -> rest
		nil,
		// paused -> rest
		func(t *task) {
			t.timer.start()
		},
		//
		//
		// idle -> paused
		nil,
		// work -> paused
		func(t *task) {
			t.timer.stop()
		},
		// rest -> paused
		func(t *task) {
			t.timer.stop()
		},
		// paused -> paused
		nil,
	}
)

type (
	task struct {
		name          string
		id            int
		done          bool
		state         taskState
		previousState taskState
		transitions   [taskStateCount * taskStateCount]func(t *task)

		sessionRequired  int
		sessionCompleted int
		sessionLength    minute
		restLength       minute

		timer    timer
		workText [5]rune
		restText [5]rune
	}

	taskState int

	taskBuffer struct {
		items []task
		count int
		cap   int
	}
)

func (t *task) init() {
	t.transitions = taskTransitionTable
	t.timer.setDuration(t.sessionLength, 0)
}

func (t *task) update() {
	if t.timer.running {
		if finished := t.timer.advance(); finished {
			switch t.state {
			case taskStateWork:
				t.changeState(taskStateRest)
			case taskStateRest:
				t.changeState(taskStateIdle)
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

func (t *task) startWork() {
	switch t.previousState {
	case taskStateRest:
		t.changeState(taskStateRest)
	default:
		t.changeState(taskStateWork)
	}
}

func (t *task) stopWork() {
	t.changeState(taskStatePaused)
}

func (t task) ToString() string {
	return "task"
}

func (t task) progress() float64 {
	sProg := float64(t.timer.sec) / 60
	prog := math.Abs((float64(t.timer.min) + sProg) - float64(t.sessionLength))
	return prog / float64(t.sessionLength)
}

func (t task) isInProgress() bool {
	return t.state == taskStateWork || t.state == taskStateRest || t.state == taskStatePaused
}

func (t task) isWorkInProgress() bool {
	return t.state == taskStateWork || (t.state == taskStatePaused && t.previousState == taskStateWork)
}

func (t task) isRestInProgress() bool {
	return t.state == taskStateRest || (t.state == taskStatePaused && t.previousState == taskStateRest)
}

func (t *task) changeState(new taskState) {
	t.transitions[new*taskStateCount+t.state](t)
	t.previousState = t.state
	t.state = new
}

func (t *task) getWorkTime() string {
	switch t.state {
	case taskStateWork:
		numberToString(int(t.timer.min), t.workText[:])
		numberToString(int(t.timer.sec), t.workText[3:])
	case taskStatePaused:
		if t.previousState == taskStateWork {
			numberToString(int(t.timer.min), t.workText[:])
			numberToString(int(t.timer.sec), t.workText[3:])
		} else {
			numberToString(int(t.sessionLength), t.workText[:])
			numberToString(0, t.workText[3:])
		}
	default:
		numberToString(int(t.sessionLength), t.workText[:])
		numberToString(0, t.workText[3:])
	}
	t.workText[2] = ':'
	return string(t.workText[:])
}

func (t *task) getRestTime() string {
	switch t.state {
	case taskStateRest:
		numberToString(int(t.timer.min), t.restText[:])
		numberToString(int(t.timer.sec), t.restText[3:])
	case taskStatePaused:
		if t.previousState == taskStateRest {
			numberToString(int(t.timer.min), t.restText[:])
			numberToString(int(t.timer.sec), t.restText[3:])
		} else {
			numberToString(int(t.sessionLength), t.restText[:])
			numberToString(0, t.restText[3:])
		}
	default:
		numberToString(int(t.sessionLength), t.restText[:])
		numberToString(0, t.restText[3:])
	}
	t.restText[2] = ':'
	return string(t.restText[:])
}

func newTaskBuffer() taskBuffer {
	return taskBuffer{
		items: make([]task, initialTaskCap),
		cap:   initialTaskCap,
	}
}

func (t *taskBuffer) addTask(newTask task) {
	if t.count > len(t.items) {
		newSlice := make([]task, t.cap*2)
		copy(newSlice[:], t.items[:])
		t.items = newSlice
	}
	t.items[t.count] = newTask
	t.count += 1
}

func (t *taskBuffer) removeTask(id int) (at int) {
	for i := 0; i < t.count; i += 1 {
		task := &t.items[i]
		if task.id == id {
			copy(t.items[i:], t.items[i+1:])
			t.count -= 1
			return i
		}
	}
	return -1
}

func (t *taskBuffer) getTask(at int) *task {
	return &t.items[at]
}

func (t *taskBuffer) copyTask(id int) task {
	for i := 0; i < t.count; i += 1 {
		task := &t.items[i]
		if task.id == id {
			return t.items[i]
		}
	}
	return task{}
}
