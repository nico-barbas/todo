package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type optWindow struct {
	active            bool
	rect              rectangle
	addBtnRect        rectangle
	nameRect          rectangle
	nameInputSelected bool
	nameTextBox       textBox
	blinkTimer        int

	incrementRect  rectangle
	decrementRect  rectangle
	lengthRect     rectangle
	lengthValue    int
	lengthText     [3]rune
	lengthTextRect rectangle
	lengthCount    int

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

	o.nameRect = rectangle{
		x:      o.rect.x + btnPadding,
		y:      o.rect.y + 50,
		width:  o.rect.width - (btnPadding * 2),
		height: 30,
	}

	o.lengthRect = rectangle{
		x:      o.rect.x + btnPadding*4,
		y:      o.rect.y + 80 + btnPadding,
		width:  o.rect.width - (btnPadding * 8),
		height: 60,
	}

	advance := font.GlyphAdvance('>', textSize) + 3
	o.decrementRect = rectangle{
		x:      o.lengthRect.x,
		y:      o.lengthRect.y,
		width:  advance,
		height: o.lengthRect.height,
	}
	o.incrementRect = rectangle{
		x:      o.lengthRect.x + o.lengthRect.width - advance,
		y:      o.lengthRect.y,
		width:  advance,
		height: o.lengthRect.height,
	}
	o.lengthTextRect = rectangle{
		x:      o.decrementRect.x + o.decrementRect.width,
		y:      o.decrementRect.y,
		width:  (o.decrementRect.x + o.decrementRect.width) - (o.incrementRect.x),
		height: o.lengthRect.height,
	}

	o.addBtnRect = rectangle{
		x:      o.rect.x + btnPadding*6,
		y:      o.rect.y + o.rect.height - btnHeight - 10,
		width:  o.rect.width - (btnPadding * 12),
		height: btnHeight,
	}

	o.lengthValue = int(minSessionLength)
	o.formatLength()

	o.font = font
	o.rectOutline = outline
	o.outlineConstr = constraint{2, 2, 2, 2}

	o.nameTextBox.init(font, textSize)
}

func (o *optWindow) update(mPos point, mLeft bool) {
	if o.active {
		o.shouldHighlight = false
		inBounds := o.rect.boundCheck(mPos)
		if !inBounds && mLeft {
			o.active = false
			return
		}
		if inBounds {
			switch {
			case o.nameRect.boundCheck(mPos):
				if mLeft {
					o.nameInputSelected = true
				}

			case o.incrementRect.boundCheck(mPos):
				o.shouldHighlight = true
				o.highlightRect = o.incrementRect
				if mLeft {
					o.lengthValue += 1
					if o.lengthValue > 999 {
						o.lengthValue = 999
					}
					o.formatLength()
				}

			case o.decrementRect.boundCheck(mPos):
				o.shouldHighlight = true
				o.highlightRect = o.decrementRect
				if mLeft {
					o.lengthValue -= 1
					if o.lengthValue < 0 {
						o.lengthValue = 0
					}
					o.formatLength()
				}

			case o.addBtnRect.boundCheck(mPos):
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
							sessionRequired: 1,
							sessionLength:   minute(o.lengthValue),
						},
					)
					o.lengthValue = int(minSessionLength)
					o.nameTextBox.Clear()
					o.active = false
				}
			default:
				if mLeft {
					o.nameInputSelected = false
				}
			}
		}

		if o.nameInputSelected {
			var runes []rune
			runes = ebiten.AppendInputChars(runes[:0])

			for _, r := range runes {
				o.nameTextBox.AppendChar(r)
			}

			if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
				o.nameTextBox.DeleteChar()
			}
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

		textPos := point{o.rect.x + 5, o.rect.y + 5}
		drawText(dst, o.font, "Add new task", textPos, textSize, White)

		drawImageSlice(dst, o.nameRect, o.rectOutline, o.outlineConstr, White)
		if o.nameInputSelected {
			cursor := o.nameTextBox.cursor
			cursor.x += o.nameRect.x + 2
			cursor.y += o.nameRect.y + 5
			drawRect(dst, cursor, White)
			drawRect(dst, o.nameRect, Color{255, 255, 255, 120})

		}

		if o.nameTextBox.charCount > 0 {
			textPos = point{o.nameRect.x + 2, o.nameRect.y + 5}
			drawText(dst, o.font, string(o.nameTextBox.GetText()), textPos, textSize, White)
		} else {
			textPos = point{o.nameRect.x + 2, o.nameRect.y + 5}
			drawText(dst, o.font, "Name", textPos, textSize, Color{255, 255, 255, 120})
		}

		drawImageSlice(dst, o.lengthRect, o.rectOutline, o.outlineConstr, White)
		textPos = point{o.decrementRect.x + 2, o.decrementRect.y + o.font.Ascent(textSize)}
		drawText(dst, o.font, "<", textPos, textSize, Color{255, 255, 255, 120})
		textPos[0] = o.incrementRect.x
		drawText(dst, o.font, ">", textPos, textSize, Color{255, 255, 255, 120})

		lSize := o.font.MeasureText(string(o.lengthText[:o.lengthCount]), largeTextSize)[0]
		tSize := lSize + o.font.MeasureText("min", smallTextSize)[0]
		textPos = point{
			o.lengthRect.x + (o.lengthRect.width/2 - tSize/2),
			o.lengthRect.y + o.font.Ascent(largeTextSize)/2,
		}
		drawText(dst, o.font, string(o.lengthText[:o.lengthCount]), textPos, largeTextSize, White)
		textPos[0] += lSize + 4
		textPos[1] += (o.font.Ascent(largeTextSize) - o.font.Ascent(smallTextSize))
		drawText(dst, o.font, "min", textPos, smallTextSize, Color{255, 255, 255, 120})

		drawImageSlice(dst, o.addBtnRect, o.rectOutline, o.outlineConstr, White)
		textPos = point{
			o.addBtnRect.x + (o.addBtnRect.width/2 - o.font.MeasureText("Add", textSize)[0]/2),
			o.addBtnRect.y + o.font.Ascent(textSize)/2,
		}
		drawText(dst, o.font, "Add", textPos, textSize, White)
	}
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
