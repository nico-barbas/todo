package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	addInputBoxID rectID = iota
	addDecCountID
	addIncCountID
	addDecWorkID
	addIncWorkID
	addIncRestID
	addDecRestID
	addAddBtnID
)

type addWindow struct {
	active   bool
	dirty    bool
	canvas   *ebiten.Image
	position point

	rect              rectLayout
	titleRect         rectLayout
	inputBoxRect      rectLayout
	addBtnRect        rectLayout
	nameInputSelected bool
	nameInput         textBox
	blinkTimer        int

	countRect    rectLayout
	incCountRect rectangle
	decCountRect rectangle
	countValue   int
	countBuf     [3]rune
	countText    string

	workLengthRect    rectLayout
	incWorkLengthRect rectangle
	decWorkLengthRect rectangle
	workLengthValue   int
	workLengthBuf     [3]rune
	workLengthText    string

	restLengthRect    rectLayout
	incRestLengthRect rectangle
	decRestLengthRect rectangle
	restLengthValue   int
	restLengthBuf     [3]rune
	restLengthText    string

	elements rectArray

	// resources
	font          *Font
	rectOutline   *ebiten.Image
	outlineConstr constraint
	bgFill        *ebiten.Image
}

func (a *addWindow) init(font *Font, outline *ebiten.Image) {
	const addWindowPadding = 10
	const addWindowNoPadding = 0
	const addWindowMargin = 20
	AddSignalListener(todoAddBtnPressed, a)
	a.elements.init(a, 10)

	a.rect = newRectLayout(rectangle{
		x:      0,
		y:      0,
		width:  300,
		height: 300,
	})
	a.position = point{
		windowWidth/2 - 150,
		windowHeight/2 - 150,
	}
	a.rect.cut(rectCutUp, addWindowPadding, 0)
	a.rect.cut(rectCutDown, addWindowPadding, 0)

	a.titleRect = a.rect.cut(rectCutUp, textSize, 15)
	a.titleRect.cut(rectCutLeft, addWindowPadding, 0)
	a.titleRect.cut(rectCutRight, addWindowPadding, 0)

	a.inputBoxRect = a.rect.cut(rectCutUp, 30, addWindowPadding)
	a.inputBoxRect.cut(rectCutLeft, addWindowMargin, 0)
	a.inputBoxRect.cut(rectCutRight, addWindowMargin, 0)

	tWidth := (a.rect.full.width-(addWindowMargin*2))/2 - addWindowPadding/2
	advance := font.GlyphAdvance('>', textSize) + 3
	a.countRect = a.rect.cut(rectCutUp, 60, addWindowPadding)
	{
		toCut := (a.rect.full.width - tWidth) / 2
		a.countRect.cut(rectCutLeft, toCut, 0)
		a.countRect.cut(rectCutRight, toCut, 0)
		a.countRect = newRectLayout(a.countRect.remaining)
		a.decCountRect = a.countRect.cut(rectCutLeft, advance, 0).remaining
		a.incCountRect = a.countRect.cut(rectCutRight, advance, 0).remaining
	}
	{
		timeSelectRect := a.rect.cut(rectCutUp, 60, addWindowPadding)
		timeSelectRect.cut(rectCutLeft, addWindowMargin, 0)
		timeSelectRect.cut(rectCutRight, addWindowMargin, 0)

		a.workLengthRect = timeSelectRect.cut(rectCutLeft, tWidth, 0)
		a.decWorkLengthRect = a.workLengthRect.cut(rectCutLeft, advance, 0).remaining
		a.incWorkLengthRect = a.workLengthRect.cut(rectCutRight, advance, 0).remaining
		a.restLengthRect = timeSelectRect.cut(rectCutRight, tWidth, 0)
		a.decRestLengthRect = a.restLengthRect.cut(rectCutLeft, advance, 0).remaining
		a.incRestLengthRect = a.restLengthRect.cut(rectCutRight, advance, 0).remaining
	}

	a.addBtnRect = a.rect.cut(rectCutDown, btnHeight, 0)
	a.addBtnRect.cut(rectCutLeft, 75, 0)
	a.addBtnRect.cut(rectCutRight, 75, 0)

	a.elements.setOffset(a.position)
	a.elements.add(a.inputBoxRect.remaining, addInputBoxID)
	a.elements.add(a.incCountRect, addIncCountID)
	a.elements.add(a.decCountRect, addDecCountID)
	a.elements.add(a.incWorkLengthRect, addIncWorkID)
	a.elements.add(a.decWorkLengthRect, addDecWorkID)
	a.elements.add(a.incRestLengthRect, addIncRestID)
	a.elements.add(a.decRestLengthRect, addDecRestID)
	a.elements.add(a.addBtnRect.remaining, addAddBtnID)

	a.workLengthValue = int(minSessionLength)
	a.restLengthValue = int(minSessionLength)
	a.formatLength()
	a.countValue = minSessionCount
	a.formatCount()

	a.dirty = true
	a.canvas = ebiten.NewImage(int(a.rect.full.width), int(a.rect.full.height))
	a.font = font
	a.rectOutline = outline
	a.outlineConstr = constraint{2, 2, 2, 2}
	a.bgFill = ebiten.NewImage(1, 1)
	a.bgFill.Fill(darkBackground3)

	a.nameInput.init(font, textSize)
}

func (a *addWindow) update(mPos point, mLeft bool) {
	if a.active {
		relPos := mPos.sub(a.position)
		if !a.rect.full.boundCheck(relPos) && mLeft {
			a.active = false
			FireSignal(todoAddWindowClosed, SignalNoArgs)
			return
		}
		a.elements.update(mPos, mLeft)

		if a.nameInputSelected {
			previousCharCout := a.nameInput.charCount
			var runes []rune
			runes = ebiten.AppendInputChars(runes[:0])

			for _, r := range runes {
				a.nameInput.AppendChar(r)
			}

			if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
				a.nameInput.DeleteChar()
			}

			currentCharCount := a.nameInput.charCount
			if previousCharCout != currentCharCount {
				a.dirty = true
			}
		}
	}
}

func (a *addWindow) draw(dst *ebiten.Image) {
	if a.active {
		fillOpt := ebiten.DrawImageOptions{}
		fillOpt.CompositeMode = ebiten.CompositeModeSourceOver
		fillOpt.GeoM.Scale(800, 600)
		fillOpt.ColorM.Scale(1.0, 1.0, 1.0, 0.5)
		dst.DrawImage(a.bgFill, &fillOpt)

		// Draw the background
		rect := a.rect.full.addPoint(a.position)
		drawRect(dst, rect, darkBackground1)
		drawImageSlice(dst, rect, a.rectOutline, a.outlineConstr, White)
		a.elements.highlight(dst)

		if a.dirty {
			a.redraw()
			a.dirty = false
		}
		drawImage(dst, a.canvas, a.position)
	}
}

func (a *addWindow) redraw() {
	a.canvas.Clear()
	// Title
	drawText(a.canvas, textOptions{
		font: a.font, text: "Add new Task", pos: point{a.titleRect.remaining.x, a.titleRect.remaining.y},
		size: textSize, clr: White,
	})

	// Name input box
	drawImageSlice(a.canvas, a.inputBoxRect.remaining, a.rectOutline, a.outlineConstr, White)
	if a.nameInputSelected {
		cursor := a.nameInput.cursor
		cursor.x += a.inputBoxRect.remaining.x + 2
		cursor.y += a.inputBoxRect.remaining.y + 5
		drawRect(a.canvas, cursor, White)
		drawRect(a.canvas, a.inputBoxRect.remaining, Color{255, 255, 255, 120})
	}
	if a.nameInput.charCount > 0 {
		drawText(a.canvas, textOptions{
			font: a.font, text: string(a.nameInput.GetText()),
			pos:  point{a.inputBoxRect.remaining.x + 2, a.inputBoxRect.remaining.y + 5},
			size: textSize, clr: White,
		})
	} else {
		drawText(a.canvas, textOptions{
			font: a.font, text: "Name",
			pos:  point{a.inputBoxRect.remaining.x + 2, a.inputBoxRect.remaining.y + 5},
			size: textSize, clr: Color{255, 255, 255, 120},
		})
	}

	drawSlider(
		a.canvas,
		a.workLengthRect.full,
		a.incWorkLengthRect, a.decWorkLengthRect,
		a.workLengthText,
	)
	drawSlider(
		a.canvas,
		a.restLengthRect.full,
		a.incRestLengthRect, a.decRestLengthRect,
		a.restLengthText,
	)
	drawSlider(
		a.canvas,
		a.countRect.full,
		a.incCountRect, a.decCountRect,
		a.countText,
	)

	drawTextBtn(a.canvas, a.addBtnRect.remaining, "Add", textSize)
}

func (a *addWindow) OnSignal(s Signal) {
	switch s.Kind {
	case todoAddBtnPressed:
		a.active = true
	}
}

func (a *addWindow) onClick(userID rectID) {
	switch userID {
	case addInputBoxID:
		a.nameInputSelected = true
		a.dirty = true

	case addDecCountID:
		a.countValue -= 1
		if a.countValue < 1 {
			a.countValue = 1
		}
		a.formatCount()
		a.dirty = true

	case addIncCountID:
		a.countValue += 1
		if a.countValue > 999 {
			a.countValue = 999
		}
		a.formatCount()
		a.dirty = true

	case addDecWorkID:
		a.workLengthValue -= 1
		if a.workLengthValue < 1 {
			a.workLengthValue = 1
		}
		a.formatLength()
		a.dirty = true

	case addIncWorkID:
		a.workLengthValue += 1
		if a.workLengthValue > 999 {
			a.workLengthValue = 999
		}
		a.formatLength()
		a.dirty = true

	case addDecRestID:
		a.restLengthValue -= 1
		if a.restLengthValue < 1 {
			a.restLengthValue = 1
		}
		a.formatLength()
		a.dirty = true

	case addIncRestID:
		a.restLengthValue += 1
		if a.restLengthValue > 999 {
			a.restLengthValue = 999
		}
		a.formatLength()
		a.dirty = true

	case addAddBtnID:
		var name string
		if a.nameInput.charCount == 0 {
			name = "Unnamed Task"
		} else {
			name = string(a.nameInput.GetText())
		}
		FireSignal(
			todoTaskAdded,
			task{
				name:            name,
				sessionRequired: a.countValue,
				sessionLength:   minute(a.workLengthValue),
				restLength:      minute(a.restLengthValue),
			},
		)
		a.nameInputSelected = false
		a.workLengthValue = int(minSessionLength)
		a.restLengthValue = int(minSessionLength)
		a.formatLength()
		a.countValue = minSessionCount
		a.formatCount()
		a.nameInput.Clear()
		a.dirty = true
		a.active = false
	}
}

func (a *addWindow) isInputHandled(mPos point) bool {
	r := a.rect.full
	r.x += a.position[0]
	r.y += a.position[1]
	return r.boundCheck(mPos)
}

func (a *addWindow) formatLength() {
	workCount := numberToString(a.workLengthValue, a.workLengthBuf[:])
	a.workLengthText = string(a.workLengthBuf[:workCount])
	restCount := numberToString(a.restLengthValue, a.restLengthBuf[:])
	a.restLengthText = string(a.restLengthBuf[:restCount])
}

func (a *addWindow) formatCount() {
	count := numberToString(a.countValue, a.countBuf[:])
	a.countText = string(a.countBuf[:count])
}
