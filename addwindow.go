package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	nameTextBox       textBox
	blinkTimer        int

	countRect    rectLayout
	incCountRect rectangle
	decCountRect rectangle
	countValue   int
	countText    [3]rune
	countCount   int
	// countTextRect      rectangle

	workLengthRect    rectLayout
	incWorkLengthRect rectangle
	decWorkLengthRect rectangle
	workLengthValue   int
	workLengthText    [3]rune
	workLengthCount   int
	// lengthTextRect      rectangle

	restLengthRect    rectLayout
	incRestLengthRect rectangle
	decRestLengthRect rectangle
	restLengthValue   int
	restLengthText    [3]rune
	restLengthCount   int
	// lengthTextRect      rectangle

	// shouldHighlight bool
	// highlightRect   rectangle

	// resources
	font          *Font
	rectOutline   *ebiten.Image
	outlineConstr constraint
	bgFill        *ebiten.Image
}

func (a *addWindow) init(font *Font, outline *ebiten.Image) {
	const addWindowPadding = 10
	const addWindowNoPadding = 0
	AddSignalListener(todoAddBtnPressed, a)

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

	a.titleRect = a.rect.cut(rectCutUp, textSize, addWindowPadding)
	a.titleRect.cut(rectCutLeft, addWindowPadding, 0)
	a.titleRect.cut(rectCutRight, addWindowPadding, 0)

	a.inputBoxRect = a.rect.cut(rectCutUp, 30, addWindowPadding)
	a.countRect = a.rect.cut(rectCutUp, 60, addWindowPadding)

	timeSelectRect := a.rect.cut(rectCutUp, 60, addWindowPadding)
	timeSelectRect.cut(rectCutLeft, addWindowPadding, 0)
	timeSelectRect.cut(rectCutRight, addWindowPadding, 0)
	tWidth := timeSelectRect.remaining.width/2 - btnPadding
	a.workLengthRect = timeSelectRect.cut(rectCutLeft, tWidth, 0)
	a.restLengthRect = timeSelectRect.cut(rectCutRight, tWidth, 0)

	a.addBtnRect = a.rect.cut(rectCutDown, btnHeight, 0)

	// a.rect = rectangle{
	// 	x:      windowWidth/2 - 150,
	// 	y:      windowHeight/2 - 150,
	// 	width:  300,
	// 	height: 300,
	// }

	// a.nameRect = rectangle{
	// 	x:      btnPadding,
	// 	y:      50,
	// 	width:  a.rect.width - (btnPadding * 2),
	// 	height: 30,
	// }

	// // Increment widgets for setting the task info (length and count)
	// advance := font.GlyphAdvance('>', textSize) + 3
	// a.lengthRect = rectangle{
	// 	x:      btnPadding * 4,
	// 	y:      80 + btnPadding,
	// 	width:  a.rect.width - (btnPadding * 8),
	// 	height: 60,
	// }
	// a.decrementLengthRect = rectangle{
	// 	x:      a.lengthRect.x,
	// 	y:      a.lengthRect.y,
	// 	width:  advance,
	// 	height: a.lengthRect.height,
	// }
	// a.incrementLengthRect = rectangle{
	// 	x:      a.lengthRect.x + a.lengthRect.width - advance,
	// 	y:      a.lengthRect.y,
	// 	width:  advance,
	// 	height: a.lengthRect.height,
	// }
	// a.lengthTextRect = rectangle{
	// 	x:      a.decrementLengthRect.x + a.decrementLengthRect.width,
	// 	y:      a.decrementLengthRect.y,
	// 	width:  (a.decrementLengthRect.x + a.decrementLengthRect.width) - (a.incrementLengthRect.x),
	// 	height: a.lengthRect.height,
	// }

	// a.countRect = rectangle{
	// 	x:      btnPadding * 4,
	// 	y:      a.lengthRect.y + a.lengthRect.height + btnPadding,
	// 	width:  a.rect.width - (btnPadding * 8),
	// 	height: 60,
	// }
	// a.decrementCountRect = rectangle{
	// 	x:      a.countRect.x,
	// 	y:      a.countRect.y,
	// 	width:  advance,
	// 	height: a.countRect.height,
	// }
	// a.incrementCountRect = rectangle{
	// 	x:      a.countRect.x + a.countRect.width - advance,
	// 	y:      a.countRect.y,
	// 	width:  advance,
	// 	height: a.countRect.height,
	// }
	// a.countTextRect = rectangle{
	// 	x:      a.decrementCountRect.x + a.decrementCountRect.width,
	// 	y:      a.decrementCountRect.y,
	// 	width:  (a.decrementCountRect.x + a.decrementCountRect.width) - (a.incrementCountRect.x),
	// 	height: a.countRect.height,
	// }

	// Bottom button
	// a.addBtnRect = rectangle{
	// 	x:      btnPadding * 6,
	// 	y:      a.rect.height - btnHeight - 10,
	// 	width:  a.rect.width - (btnPadding * 12),
	// 	height: btnHeight,
	// }

	// a.lengthValue = int(minSessionLength)
	// a.formatLength()
	// a.countValue = minSessionCount
	// a.formatCount()

	a.dirty = true
	a.canvas = ebiten.NewImage(int(a.rect.full.width), int(a.rect.full.height))
	a.font = font
	a.rectOutline = outline
	a.outlineConstr = constraint{2, 2, 2, 2}
	a.bgFill = ebiten.NewImage(1, 1)
	a.bgFill.Fill(darkBackground3)

	a.nameTextBox.init(font, textSize)
}

func (a *addWindow) update(mPos point, mLeft bool) {
	// if a.active {
	// 	a.shouldHighlight = false
	// 	relPos := point{mPos[0] - a.rect.x, mPos[1] - a.rect.y}
	// 	inBounds := a.rect.boundCheck(mPos)
	// 	if !inBounds && mLeft {
	// 		a.active = false
	// 		return
	// 	}
	// 	if inBounds {
	// 		switch {
	// 		case a.nameRect.boundCheck(relPos):
	// 			if mLeft {
	// 				a.nameInputSelected = true
	// 				a.dirty = true
	// 			}

	// 		case a.incrementLengthRect.boundCheck(relPos):
	// 			a.shouldHighlight = true
	// 			a.highlightRect = a.incrementLengthRect
	// 			if mLeft {
	// 				a.lengthValue += 1
	// 				if a.lengthValue > 999 {
	// 					a.lengthValue = 999
	// 				}
	// 				a.formatLength()
	// 				a.dirty = true
	// 			}

	// 		case a.decrementLengthRect.boundCheck(relPos):
	// 			a.shouldHighlight = true
	// 			a.highlightRect = a.decrementLengthRect
	// 			if mLeft {
	// 				a.lengthValue -= 1
	// 				if a.lengthValue < 0 {
	// 					a.lengthValue = 0
	// 				}
	// 				a.formatLength()
	// 				a.dirty = true
	// 			}

	// 		case a.incrementCountRect.boundCheck(relPos):
	// 			a.shouldHighlight = true
	// 			a.highlightRect = a.incrementCountRect
	// 			if mLeft {
	// 				a.countValue += 1
	// 				if a.countValue > 999 {
	// 					a.countValue = 999
	// 				}
	// 				a.formatCount()
	// 				a.dirty = true
	// 			}

	// 		case a.decrementCountRect.boundCheck(relPos):
	// 			a.shouldHighlight = true
	// 			a.highlightRect = a.decrementCountRect
	// 			if mLeft {
	// 				a.countValue -= 1
	// 				if a.countValue < 0 {
	// 					a.countValue = 0
	// 				}
	// 				a.formatCount()
	// 				a.dirty = true
	// 			}

	// 		case a.addBtnRect.boundCheck(relPos):
	// 			a.shouldHighlight = true
	// 			a.highlightRect = a.addBtnRect
	// 			if mLeft {
	// 				var name string
	// 				if a.nameTextBox.charCount == 0 {
	// 					name = "Unnamed Task"
	// 				} else {
	// 					name = string(a.nameTextBox.GetText())
	// 				}
	// 				FireSignal(
	// 					todoTaskAdded,
	// 					task{
	// 						name:            name,
	// 						sessionRequired: a.countValue,
	// 						sessionLength:   minute(a.lengthValue),
	// 					},
	// 				)
	// 				a.nameInputSelected = false
	// 				a.lengthValue = int(minSessionLength)
	// 				a.formatLength()
	// 				a.countValue = minSessionCount
	// 				a.formatCount()
	// 				a.nameTextBox.Clear()
	// 				a.dirty = true
	// 				a.active = false
	// 			}
	// 		default:
	// 			if mLeft {
	// 				a.nameInputSelected = false
	// 				a.dirty = true
	// 			}
	// 		}
	// 	}

	// 	if a.nameInputSelected {
	// 		previousCharCout := a.nameTextBox.charCount
	// 		var runes []rune
	// 		runes = ebiten.AppendInputChars(runes[:0])

	// 		for _, r := range runes {
	// 			a.nameTextBox.AppendChar(r)
	// 		}

	// 		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
	// 			a.nameTextBox.DeleteChar()
	// 		}

	// 		currentCharCount := a.nameTextBox.charCount
	// 		if previousCharCout != currentCharCount {
	// 			a.dirty = true
	// 		}
	// 	}
	// 	if a.shouldHighlight {
	// 		a.highlightRect.x += a.rect.x
	// 		a.highlightRect.y += a.rect.y
	// 	}
	// }
}

func (a *addWindow) draw(dst *ebiten.Image) {
	if a.active {
		fillOpt := ebiten.DrawImageOptions{}
		fillOpt.CompositeMode = ebiten.CompositeModeSourceOver
		fillOpt.GeoM.Scale(800, 600)
		fillOpt.ColorM.Scale(1.0, 1.0, 1.0, 0.5)
		dst.DrawImage(a.bgFill, &fillOpt)

		// Draw the background
		drawRect(dst, a.rect.full, darkBackground1)
		drawImageSlice(dst, a.rect.full, a.rectOutline, a.outlineConstr, White)

		// Draw the content of the window
		// if a.shouldHighlight {
		// 	drawRect(dst, a.highlightRect, Color{255, 255, 255, 120})
		// }

		if a.dirty {
			a.redraw()
			a.dirty = false
		}
		drawImage(dst, a.canvas, point{a.rect.full.x, a.rect.full.y})
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
		cursor := a.nameTextBox.cursor
		cursor.x += a.inputBoxRect.remaining.x + 2
		cursor.y += a.inputBoxRect.remaining.y + 5
		drawRect(a.canvas, cursor, White)
		drawRect(a.canvas, a.inputBoxRect.remaining, Color{255, 255, 255, 120})
	}
	if a.nameTextBox.charCount > 0 {
		drawText(a.canvas, textOptions{
			font: a.font, text: string(a.nameTextBox.GetText()),
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
		a.workLengthRect.remaining,
		a.incWorkLengthRect, a.decWorkLengthRect,
		string(a.workLengthText[:]),
	)
	drawSlider(
		a.canvas,
		a.restLengthRect.remaining,
		a.incRestLengthRect, a.decRestLengthRect,
		string(a.restLengthText[:]),
	)
	drawSlider(
		a.canvas,
		a.countRect.remaining,
		a.incCountRect, a.decCountRect,
		string(a.countText[:]),
	)
	// a.drawLengthWidget(a.canvas)
	// a.drawCountWidget(a.canvas)

	// Bottom button to validate the task addition
	// drawImageSlice(a.canvas, a.addBtnRect, a.rectOutline, a.outlineConstr, White)
	// drawTextCenter(a.canvas, textOptions{
	// 	font: a.font, text: "Add", bounds: a.addBtnRect,
	// 	size: textSize, clr: White,
	// })
	drawTextBtn(a.canvas, a.addBtnRect.remaining, "Add", textSize)
}

// func (a *addWindow) drawLengthWidget(dst *ebiten.Image) {
// 	// Given the fact that we want to use 2 different font size, we have to
// 	// calculate the position by hand and cannot rely on drawTextCenter()
// 	//
// 	// Session length setting widget
// 	drawImageSlice(dst, a.lengthRect, a.rectOutline, a.outlineConstr, White)
// 	drawTextCenter(dst, textOptions{
// 		font: a.font, text: "<", bounds: a.decrementLengthRect,
// 		size: textSize, clr: Color{255, 255, 255, 120},
// 	})
// 	drawTextCenter(dst, textOptions{
// 		font: a.font, text: ">", bounds: a.incrementLengthRect,
// 		size: textSize, clr: Color{255, 255, 255, 120},
// 	})
// 	lSize := a.font.MeasureText(string(a.lengthText[:a.lengthCount]), largeTextSize)[0]
// 	tSize := lSize + a.font.MeasureText("min", smallTextSize)[0]
// 	textPos := point{
// 		a.lengthRect.x + (a.lengthRect.width/2 - tSize/2),
// 		a.lengthRect.y + a.font.Ascent(largeTextSize)/2,
// 	}
// 	drawText(dst, textOptions{
// 		font: a.font, text: string(a.lengthText[:a.lengthCount]), pos: textPos,
// 		size: largeTextSize, clr: White,
// 	})
// 	textPos[0] += lSize + 4
// 	textPos[1] += (a.font.Ascent(largeTextSize) - a.font.Ascent(smallTextSize))
// 	drawText(dst, textOptions{
// 		font: a.font, text: "min", pos: textPos,
// 		size: smallTextSize, clr: Color{255, 255, 255, 120},
// 	})
// }

// func (a *addWindow) drawCountWidget(dst *ebiten.Image) {
// 	//
// 	// Count widget
// 	drawImageSlice(dst, a.countRect, a.rectOutline, a.outlineConstr, White)
// 	drawTextCenter(dst, textOptions{
// 		font: a.font, text: "<", bounds: a.decrementCountRect,
// 		size: textSize, clr: Color{255, 255, 255, 120},
// 	})
// 	drawTextCenter(dst, textOptions{
// 		font: a.font, text: ">", bounds: a.incrementCountRect,
// 		size: textSize, clr: Color{255, 255, 255, 120},
// 	})
// 	txt := string(a.countText[:a.countCount])
// 	lSize := a.font.MeasureText(txt, largeTextSize)[0]
// 	tSize := lSize + a.font.MeasureText("count", smallTextSize)[0]
// 	textPos := point{
// 		a.countRect.x + (a.countRect.width/2 - tSize/2),
// 		a.countRect.y + a.font.Ascent(largeTextSize)/2,
// 	}
// 	drawText(dst, textOptions{
// 		font: a.font, text: txt, pos: textPos,
// 		size: largeTextSize, clr: White,
// 	})
// 	textPos[0] += lSize + 4
// 	textPos[1] += (a.font.Ascent(largeTextSize) - a.font.Ascent(smallTextSize))
// 	drawText(dst, textOptions{
// 		font: a.font, text: "count", pos: textPos,
// 		size: smallTextSize, clr: Color{255, 255, 255, 120},
// 	})
// }

func (a *addWindow) OnSignal(s Signal) {
	switch s.Kind {
	case todoAddBtnPressed:
		a.active = true
	}
}

func (a *addWindow) isInputHandled(mPos point) bool {
	r := a.rect.full
	r.x += a.position[0]
	r.y += a.position[1]
	return r.boundCheck(mPos)
}

// func (a *addWindow) formatLength() {
// 	a.lengthCount = numberToString(a.lengthValue, a.lengthText[:])
// }

// func (a *addWindow) formatCount() {
// 	a.countCount = numberToString(a.countValue, a.countText[:])
// }
