package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	windowWidth  = 800
	windowHeight = 600
	initialCap   = 20
	textSize     = 20
	itemPadding  = 4
	itemHeight   = textSize + itemPadding*2

	maxSessionPerLine   = 10
	sessionCheckPadding = 4
	sessionCheckHeight  = 600 / (maxSessionPerLine + sessionCheckPadding*2)
	sessionCheckSize    = sessionCheckHeight - (sessionCheckPadding * 2)
)

type (
	Todo struct {
		items    []item
		count    int
		cap      int
		listRect rectangle
		mainRect rectangle

		hovered           *item
		selected          *item
		listHighlightRect rectangle

		font        Font
		rectOutline *ebiten.Image

		sessionRects     [maxSessionPerLine]rectangle
		sessionTimer     timer
		sessionTimerText string
	}
)

func (t *Todo) Init() {
	t.items = make([]item, 0, initialCap)
	t.cap = initialCap
	t.listRect = rectangle{0, 0, 200, windowHeight}
	t.mainRect = rectangle{200, 0, 600, windowWidth}
	t.font = NewFont("assets/FiraSans-Regular.ttf", 72, []int{14, textSize})
	t.rectOutline, _, _ = ebitenutil.NewImageFromFile("assets/uiRectOutline.png")
}

func (t *Todo) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return exitStatus{kind: exitNoError}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fmt.Println("Stop")
	}

	mx, my := ebiten.CursorPosition()
	mPos := point{float64(mx), float64(my)}
	t.hovered = nil
	if t.listRect.boundCheck(mPos) {
		index := int(mPos[1] / itemHeight)
		if index < t.count {
			t.hovered = &t.items[index]
			if t.hovered.checkRect.boundCheck(mPos) {
				t.listHighlightRect = t.hovered.checkRect
			} else {
				t.listHighlightRect = t.hovered.rect
			}
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if t.hovered != nil {
			t.selected = t.hovered
			t.sessionTimer.setDuration(t.selected.sessionLength, 0)
			t.sessionTimerText = t.sessionTimer.toString()
		} else {
			t.selected = nil
		}
	}
	return nil
}

func (t *Todo) Draw(screen *ebiten.Image) {
	// Draw list
	if t.hovered != nil {
		ebitenutil.DrawRect(
			screen,
			t.listHighlightRect.x,
			t.listHighlightRect.y,
			t.listHighlightRect.width,
			t.listHighlightRect.height,
			Color{255, 255, 255, 120},
		)
	}

	for i := 0; i < t.count; i += 1 {
		item := &t.items[i]
		t.drawItem(screen, item)
	}

	// Separator
	ebitenutil.DrawLine(screen,
		t.mainRect.x,
		t.mainRect.y,
		t.mainRect.x,
		t.mainRect.y+t.mainRect.height,
		color.White,
	)

	// Draw selected item infos
	if t.selected != nil {
		checkWidth := float64(t.selected.sessionRequired * sessionCheckHeight)
		startPos := t.mainRect.x + (t.mainRect.width/2 - checkWidth/2)

		tSize := t.font.MeasureText("- Sessions -", textSize)
		tPos := point{t.mainRect.x + (t.mainRect.width/2 - tSize[0]/2), 70}
		drawText(screen, &t.font, "- Sessions -", tPos, textSize, White)
		for i := 0; i < t.selected.sessionRequired; i += 1 {
			checkBoxRect := rectangle{startPos + float64(i*sessionCheckHeight), 100, sessionCheckSize, sessionCheckSize}
			drawImageSlice(screen, checkBoxRect, t.rectOutline, constraint{2, 2, 2, 2}, White)
		}

		tSize = t.font.MeasureText(t.sessionTimerText, textSize)
		tPos = point{t.mainRect.x + (t.mainRect.width/2 - tSize[0]/2), 130}
		drawText(screen, &t.font, t.sessionTimerText, tPos, textSize, White)

		startRect := rectangle{t.mainRect.x + (t.mainRect.width/2 - 125/2), 160, 125, 38}
		tSize = t.font.MeasureText("Start Timer", textSize)
		tPos = point{startRect.x + (startRect.width/2 - tSize[0]/2), startRect.y + (startRect.height/2 - t.font.Ascent(textSize)/2 - 2)}
		drawText(screen, &t.font, "Start Timer", tPos, textSize, White)
		drawImageSlice(screen, startRect, t.rectOutline, constraint{2, 2, 2, 2}, White)
	}

}

func (t *Todo) Layout(outW, outH int) (int, int) {
	return windowWidth, windowHeight
}

func (t *Todo) addItem(name string) {
	rect := rectangle{0, float64(t.count) * itemHeight, 200, itemHeight}
	textPos := point{rect.x, rect.y + itemPadding}
	i := item{
		name:            name,
		done:            false,
		rect:            rect,
		textPosition:    textPos,
		checkRect:       rectangle{(rect.x + rect.width) - itemHeight, textPos[1], textSize, textSize},
		sessionRequired: 5,
		sessionLength:   25,
	}
	if t.count > len(t.items) {
		newSlice := make([]item, t.cap*2)
		copy(newSlice[:], t.items[:])
		t.items = newSlice
	}
	t.items = append(t.items, i)
	t.count += 1
}

func (t *Todo) drawItem(dst *ebiten.Image, item *item) {
	drawText(dst, &t.font, item.name, item.textPosition, textSize, White)
	ebitenutil.DrawLine(
		dst,
		item.rect.x,
		item.rect.y+item.rect.height,
		item.rect.x+item.rect.width,
		item.rect.y+item.rect.height,
		White,
	)

	drawImageSlice(dst, item.checkRect, t.rectOutline, constraint{2, 2, 2, 2}, White)
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Get work done")

	todo := new(Todo)
	todo.Init()

	todo.addItem("Clean house")
	todo.addItem("Clean desk")
	if err := ebiten.RunGame(todo); err != nil {
		e := err.(exitStatus)
		if e.kind != exitNoError {
			panic(e)
		}
	}
}
