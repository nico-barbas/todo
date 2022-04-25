package main

import (
	"time"

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
	todoTaskStarted
	todoTaskStopped
)

var todo *Todo

type (
	Todo struct {
		tasks   taskBuffer
		archive taskBuffer
		taskID  int

		selected *task

		font        Font
		rectOutline *ebiten.Image

		// List Window data
		list listWindow

		// Main Window data
		mainWindow mainWindow

		// Optional windows data
		addWindow addWindow

		// Archive window
		archiveWindow archiveWindow

		signals signalDispatcher
	}
)

func (t *Todo) Init() {
	todo = t
	loadTheme()

	// Caching all the rects possible
	// and init the subsytems
	t.tasks = newTaskBuffer()

	tnow := time.Now()
	t.taskID = tnow.Year() + int(tnow.Month()) + tnow.Day() + tnow.Hour() + tnow.Minute()
	t.signals.init()

	t.signals.addListener(todoTaskAdded, t)
	t.signals.addListener(todoTaskStarted, t)
	t.signals.addListener(todoTaskStopped, t)
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

	// Add archive window
	t.archiveWindow.init(&t.font, t.rectOutline)
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
		t.selected = t.tasks.getTask(selected)
	}

	t.mainWindow.update(mPos, mLeft, t.selected)

	// Advance all the timer and check for completed sessions
	for i := 0; i < t.tasks.count; i += 1 {
		task := t.tasks.getTask(i)
		task.update()
	}

	return nil
}

func (t *Todo) Draw(screen *ebiten.Image) {
	screen.Fill(darkBackground1)
	t.list.draw(screen, t.tasks.items[:t.tasks.count])

	t.mainWindow.draw(screen, t.selected)

	t.addWindow.draw(screen)
	t.archiveWindow.draw(screen)
}

func (t *Todo) Layout(outW, outH int) (int, int) {
	return windowWidth, windowHeight
}

func (t *Todo) addTask(_t task) {
	newTask := _t
	newTask.id = t.genID()
	newTask.init()
	t.tasks.addTask(newTask)
	t.list.addItem()
}

func (t *Todo) OnSignal(s Signal) {
	switch s.Kind {
	case todoTaskAdded:
		t.addTask(s.Value.(task))
	case todoTaskStarted:
		t.selected.startWork()
	case todoTaskStopped:
		t.selected.stopWork()
	case todoTaskRemoveAnimationDone:
		// This is always the currently selected one
		copied := t.tasks.copyTask(t.selected.id)
		at := t.tasks.removeTask(t.selected.id)
		if at > -1 {
			t.list.removeItem(at)
			t.selected = nil
			t.archive.addTask(copied)
		}
	}
}

func (t *Todo) genID() int {
	t.taskID += 1
	return t.taskID
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
