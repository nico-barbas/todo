package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	windowWidth    = 800
	windowHeight   = 600
	initialTaskCap = 20
	smallTextSize  = 14
	textSize       = 20
	largeTextSize  = 30
	btnPadding     = 15
	btnHeight      = textSize + (10 * 2)
)

const (
	todoAddBtnPressed SignalKind = iota
	todoTaskAdded
	todoTaskRemoved
	todoTaskRemoveAnimationDone
)

var todo *Todo

type (
	Todo struct {
		tasks []task
		count int
		cap   int

		selected *task

		font        Font
		rectOutline *ebiten.Image

		// List Window data
		list listWindow

		// Main Window data
		mainWindow mainWindow

		// Optional windows data
		addWindow addWindow

		signals signalDispatcher
	}
)

func (t *Todo) Init() {
	todo = t
	loadTheme()

	// Caching all the rects possible
	// and init the subsytems
	t.tasks = make([]task, initialTaskCap)
	t.cap = initialTaskCap
	t.signals.init()

	t.signals.addListener(todoTaskAdded, t)
	t.signals.addListener(todoTaskRemoveAnimationDone, t)

	// Resources
	t.font = NewFont("assets/FiraSans-Regular.ttf", 72, []int{smallTextSize, textSize, largeTextSize})
	t.rectOutline, _, _ = ebitenutil.NewImageFromFile("assets/uiRectOutline.png")

	// List window init
	t.list.init(&t.font, t.rectOutline)

	// Main window init
	t.mainWindow.init(&t.font, t.rectOutline)

	// Add window
	t.addWindow.init(&t.font, t.rectOutline)
}

func (t *Todo) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return exitStatus{kind: exitNoError}
	}

	mx, my := ebiten.CursorPosition()
	mPos := point{float64(mx), float64(my)}
	mLeft := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	t.addWindow.update(mPos, mLeft)

	selected := t.list.update(mPos, mLeft)
	if selected >= 0 {
		t.selected = &t.tasks[selected]
	}

	startWork := t.mainWindow.update(mPos, mLeft, t.selected != nil)
	if startWork {
		t.selected.startWork()
	}

	// Advance all the timer and check for completed sessions
	for i := 0; i < t.count; i += 1 {
		task := &t.tasks[i]
		if task.timer.running {
			if finished := task.timer.advance(); finished {
				task.completeSession()
				task.timer.setDuration(task.sessionLength, 0)
			}
		}
	}

	return nil
}

func (t *Todo) Draw(screen *ebiten.Image) {
	screen.Fill(darkBackground1)
	t.list.draw(screen, t.tasks[:t.count])

	t.mainWindow.draw(screen, t.selected)

	t.addWindow.draw(screen)
}

func (t *Todo) Layout(outW, outH int) (int, int) {
	return windowWidth, windowHeight
}

func (t *Todo) addTask(_t task) {
	newTask := _t
	newTask.init()
	if t.count > len(t.tasks) {
		newSlice := make([]task, t.cap*2)
		copy(newSlice[:], t.tasks[:])
		t.tasks = newSlice
	}
	t.tasks[t.count] = newTask
	t.count += 1
	t.list.addItem()
}

func (t *Todo) OnSignal(s Signal) {
	switch s.Kind {
	case todoTaskAdded:
		t.addTask(s.Value.(task))
	case todoTaskRemoveAnimationDone:
		// This is always the currently selected one
		for i := 0; i < t.count; i += 1 {
			task := &t.tasks[i]
			if task.name == t.selected.name {
				copy(t.tasks[i:], t.tasks[i+1:])
				t.count -= 1
				t.list.removeItem(i)
				t.selected = nil
				break
			}
		}
	}
}

// Static wrapper over the signal dispatcher
//
func AddSignalListener(k SignalKind, l SignalListener) {
	todo.signals.addListener(k, l)
}

// Static wrapper over the signal dispatcher
//
func FireSignal(k SignalKind, v SignalValue) {
	todo.signals.dispatch(k, v)
}

func isInputHandled(mPos point) bool {
	if todo.addWindow.active && todo.addWindow.rect.remaining.boundCheck(mPos) {
		return true
	}
	return false
}
