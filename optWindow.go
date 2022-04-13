package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type optWindow struct {
	active            bool
	dirty             bool
	canvas            *ebiten.Image
	rect              rectangle
	addBtnRect        rectangle
	nameRect          rectangle
	nameInputSelected bool
	nameTextBox       textBox
	blinkTimer        int

	incrementLengthRect rectangle
	decrementLengthRect rectangle
	lengthRect          rectangle
	lengthValue         int
	lengthText          [3]rune
	lengthTextRect      rectangle
	lengthCount         int

	incrementCountRect rectangle
	decrementCountRect rectangle
	countRect          rectangle
	countValue         int
	countText          [3]rune
	countTextRect      rectangle
	countCount         int

	shouldHighlight bool
	highlightRect   rectangle

	// resources
	font          *Font
	rectOutline   *ebiten.Image
	outlineConstr constraint
}

func (o *optWindow) init(font *Font, outline *ebiten.Image) {
	AddSignalListener(todoAddBtnPressed, o)

	o.rect = rectangle{
		x:      windowWidth/2 - 150,
		y:      windowHeight/2 - 150,
		width:  300,
		height: 300,
	}
	o.dirty = true
	o.canvas = ebiten.NewImage(int(o.rect.width), int(o.rect.height))

	o.nameRect = rectangle{
		x:      btnPadding,
		y:      50,
		width:  o.rect.width - (btnPadding * 2),
		height: 30,
	}

	// Increment widgets for setting the task info (length and count)
	advance := font.GlyphAdvance('>', textSize) + 3
	o.lengthRect = rectangle{
		x:      btnPadding * 4,
		y:      80 + btnPadding,
		width:  o.rect.width - (btnPadding * 8),
		height: 60,
	}
	o.decrementLengthRect = rectangle{
		x:      o.lengthRect.x,
		y:      o.lengthRect.y,
		width:  advance,
		height: o.lengthRect.height,
	}
	o.incrementLengthRect = rectangle{
		x:      o.lengthRect.x + o.lengthRect.width - advance,
		y:      o.lengthRect.y,
		width:  advance,
		height: o.lengthRect.height,
	}
	o.lengthTextRect = rectangle{
		x:      o.decrementLengthRect.x + o.decrementLengthRect.width,
		y:      o.decrementLengthRect.y,
		width:  (o.decrementLengthRect.x + o.decrementLengthRect.width) - (o.incrementLengthRect.x),
		height: o.lengthRect.height,
	}

	o.countRect = rectangle{
		x:      btnPadding * 4,
		y:      o.lengthRect.y + o.lengthRect.height + btnPadding,
		width:  o.rect.width - (btnPadding * 8),
		height: 60,
	}
	o.decrementCountRect = rectangle{
		x:      o.countRect.x,
		y:      o.countRect.y,
		width:  advance,
		height: o.countRect.height,
	}
	o.incrementCountRect = rectangle{
		x:      o.countRect.x + o.countRect.width - advance,
		y:      o.countRect.y,
		width:  advance,
		height: o.countRect.height,
	}
	o.countTextRect = rectangle{
		x:      o.decrementCountRect.x + o.decrementCountRect.width,
		y:      o.decrementCountRect.y,
		width:  (o.decrementCountRect.x + o.decrementCountRect.width) - (o.incrementCountRect.x),
		height: o.countRect.height,
	}

	// Bottom button
	o.addBtnRect = rectangle{
		x:      btnPadding * 6,
		y:      o.rect.height - btnHeight - 10,
		width:  o.rect.width - (btnPadding * 12),
		height: btnHeight,
	}

	o.lengthValue = int(minSessionLength)
	o.formatLength()
	o.countValue = minSessionCount
	o.formatCount()

	o.font = font
	o.rectOutline = outline
	o.outlineConstr = constraint{2, 2, 2, 2}

	o.nameTextBox.init(font, textSize)
}

func (o *optWindow) update(mPos point, mLeft bool) {
	if o.active {
		o.shouldHighlight = false
		relPos := point{mPos[0] - o.rect.x, mPos[1] - o.rect.y}
		inBounds := o.rect.boundCheck(mPos)
		if !inBounds && mLeft {
			o.active = false
			return
		}
		if inBounds {
			switch {
			case o.nameRect.boundCheck(relPos):
				if mLeft {
					o.nameInputSelected = true
					o.dirty = true
				}

			case o.incrementLengthRect.boundCheck(relPos):
				o.shouldHighlight = true
				o.highlightRect = o.incrementLengthRect
				if mLeft {
					o.lengthValue += 1
					if o.lengthValue > 999 {
						o.lengthValue = 999
					}
					o.formatLength()
					o.dirty = true
				}

			case o.decrementLengthRect.boundCheck(relPos):
				o.shouldHighlight = true
				o.highlightRect = o.decrementLengthRect
				if mLeft {
					o.lengthValue -= 1
					if o.lengthValue < 0 {
						o.lengthValue = 0
					}
					o.formatLength()
					o.dirty = true
				}

			case o.incrementCountRect.boundCheck(relPos):
				o.shouldHighlight = true
				o.highlightRect = o.incrementCountRect
				if mLeft {
					o.countValue += 1
					if o.countValue > 999 {
						o.countValue = 999
					}
					o.formatCount()
					o.dirty = true
				}

			case o.decrementCountRect.boundCheck(relPos):
				o.shouldHighlight = true
				o.highlightRect = o.decrementCountRect
				if mLeft {
					o.countValue -= 1
					if o.countValue < 0 {
						o.countValue = 0
					}
					o.formatCount()
					o.dirty = true
				}

			case o.addBtnRect.boundCheck(relPos):
				o.shouldHighlight = true
				o.highlightRect = o.addBtnRect
				if mLeft {
					var name string
					if o.nameTextBox.charCount == 0 {
						name = "Unnamed Task"
					} else {
						name = string(o.nameTextBox.GetText())
					}
					FireSignal(
						todoTaskAdded,
						task{
							name:            name,
							sessionRequired: o.countValue,
							sessionLength:   minute(o.lengthValue),
						},
					)
					o.nameInputSelected = false
					o.lengthValue = int(minSessionLength)
					o.formatLength()
					o.countValue = minSessionCount
					o.formatCount()
					o.nameTextBox.Clear()
					o.dirty = true
					o.active = false
				}
			default:
				if mLeft {
					o.nameInputSelected = false
					o.dirty = true
				}
			}
		}

		if o.nameInputSelected {
			previousCharCout := o.nameTextBox.charCount
			var runes []rune
			runes = ebiten.AppendInputChars(runes[:0])

			for _, r := range runes {
				o.nameTextBox.AppendChar(r)
			}

			if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
				o.nameTextBox.DeleteChar()
			}

			currentCharCount := o.nameTextBox.charCount
			if previousCharCout != currentCharCount {
				o.dirty = true
			}
		}
		if o.shouldHighlight {
			o.highlightRect.x += o.rect.x
			o.highlightRect.y += o.rect.y
		}
	}
}

func (o *optWindow) draw(dst *ebiten.Image) {
	if o.active {
		// Draw the background
		drawRect(dst, o.rect, Black)
		drawImageSlice(dst, o.rect, o.rectOutline, o.outlineConstr, White)

		// Draw the content of the window
		if o.shouldHighlight {
			drawRect(dst, o.highlightRect, Color{255, 255, 255, 120})
		}

		if o.dirty {
			o.redraw()
			o.dirty = false
		}
		drawImage(dst, o.canvas, point{o.rect.x, o.rect.y})
	}
}

func (o *optWindow) redraw() {
	o.canvas.Clear()
	// Title
	drawText(o.canvas, textOptions{
		font: o.font, text: "Add new Task", pos: point{5, 5},
		size: textSize, clr: White,
	})

	// Name input box
	drawImageSlice(o.canvas, o.nameRect, o.rectOutline, o.outlineConstr, White)
	if o.nameInputSelected {
		cursor := o.nameTextBox.cursor
		cursor.x += o.nameRect.x + 2
		cursor.y += o.nameRect.y + 5
		drawRect(o.canvas, cursor, White)
		drawRect(o.canvas, o.nameRect, Color{255, 255, 255, 120})

	}
	if o.nameTextBox.charCount > 0 {
		drawText(o.canvas, textOptions{
			font: o.font, text: string(o.nameTextBox.GetText()), pos: point{o.nameRect.x + 2, o.nameRect.y + 5},
			size: textSize, clr: White,
		})
	} else {
		drawText(o.canvas, textOptions{
			font: o.font, text: "Name", pos: point{o.nameRect.x + 2, o.nameRect.y + 5},
			size: textSize, clr: Color{255, 255, 255, 120},
		})
	}

	o.drawLengthWidget(o.canvas)
	o.drawCountWidget(o.canvas)

	// Bottom button to validate the task addition
	drawImageSlice(o.canvas, o.addBtnRect, o.rectOutline, o.outlineConstr, White)
	drawTextCenter(o.canvas, textOptions{
		font: o.font, text: "Add", bounds: o.addBtnRect,
		size: textSize, clr: White,
	})
}

func (o *optWindow) drawLengthWidget(dst *ebiten.Image) {
	// Given the fact that we want to use 2 different font size, we have to
	// calculate the position by hand and cannot rely on drawTextCenter()
	//
	// Session length setting widget
	drawImageSlice(dst, o.lengthRect, o.rectOutline, o.outlineConstr, White)
	drawTextCenter(dst, textOptions{
		font: o.font, text: "<", bounds: o.decrementLengthRect,
		size: textSize, clr: Color{255, 255, 255, 120},
	})
	drawTextCenter(dst, textOptions{
		font: o.font, text: ">", bounds: o.incrementLengthRect,
		size: textSize, clr: Color{255, 255, 255, 120},
	})
	lSize := o.font.MeasureText(string(o.lengthText[:o.lengthCount]), largeTextSize)[0]
	tSize := lSize + o.font.MeasureText("min", smallTextSize)[0]
	textPos := point{
		o.lengthRect.x + (o.lengthRect.width/2 - tSize/2),
		o.lengthRect.y + o.font.Ascent(largeTextSize)/2,
	}
	drawText(dst, textOptions{
		font: o.font, text: string(o.lengthText[:o.lengthCount]), pos: textPos,
		size: largeTextSize, clr: White,
	})
	textPos[0] += lSize + 4
	textPos[1] += (o.font.Ascent(largeTextSize) - o.font.Ascent(smallTextSize))
	drawText(dst, textOptions{
		font: o.font, text: "min", pos: textPos,
		size: smallTextSize, clr: Color{255, 255, 255, 120},
	})
}

func (o *optWindow) drawCountWidget(dst *ebiten.Image) {
	//
	// Count widget
	drawImageSlice(dst, o.countRect, o.rectOutline, o.outlineConstr, White)
	drawTextCenter(dst, textOptions{
		font: o.font, text: "<", bounds: o.decrementCountRect,
		size: textSize, clr: Color{255, 255, 255, 120},
	})
	drawTextCenter(dst, textOptions{
		font: o.font, text: ">", bounds: o.incrementCountRect,
		size: textSize, clr: Color{255, 255, 255, 120},
	})
	txt := string(o.countText[:o.countCount])
	lSize := o.font.MeasureText(txt, largeTextSize)[0]
	tSize := lSize + o.font.MeasureText("count", smallTextSize)[0]
	textPos := point{
		o.countRect.x + (o.countRect.width/2 - tSize/2),
		o.countRect.y + o.font.Ascent(largeTextSize)/2,
	}
	drawText(dst, textOptions{
		font: o.font, text: txt, pos: textPos,
		size: largeTextSize, clr: White,
	})
	textPos[0] += lSize + 4
	textPos[1] += (o.font.Ascent(largeTextSize) - o.font.Ascent(smallTextSize))
	drawText(dst, textOptions{
		font: o.font, text: "count", pos: textPos,
		size: smallTextSize, clr: Color{255, 255, 255, 120},
	})
}

func (o *optWindow) OnSignal(s Signal) {
	switch s.Kind {
	case todoAddBtnPressed:
		o.active = true
	}
}

func (o *optWindow) formatLength() {
	o.lengthCount = numberToString(o.lengthValue, o.lengthText[:])
}

func (o *optWindow) formatCount() {
	o.countCount = numberToString(o.countValue, o.countText[:])
}
