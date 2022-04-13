package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	blinkTime = 30
)

type optWindow struct {
	active            bool
	rect              rectangle
	addBtnRect        rectangle
	nameRect          rectangle
	nameInputSelected bool
	nameTextBox       textBox
	blinkTimer        int

	incrementRect rectangle
	decrementRect rectangle
	lengthRect    rectangle

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

	o.addBtnRect = rectangle{
		x:      o.rect.x + btnPadding*6,
		y:      o.rect.y + o.rect.height - btnHeight - 10,
		width:  o.rect.width - (btnPadding * 12),
		height: btnHeight,
	}

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
				o.shouldHighlight = true
				o.highlightRect = o.nameRect
				if mLeft {
					o.nameInputSelected = true
				}
			case o.addBtnRect.boundCheck(mPos):
				o.shouldHighlight = true
				o.highlightRect = o.addBtnRect
				if mLeft {
					FireSignal(
						todoTaskAdded,
						task{
							name:            string(o.nameTextBox.GetText()),
							sessionRequired: 1,
							sessionLength:   1,
						},
					)
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
			ebitenutil.DrawRect(
				dst,
				o.highlightRect.x,
				o.highlightRect.y,
				o.highlightRect.width,
				o.highlightRect.height,
				Color{255, 255, 255, 120},
			)
		}

		textPos := point{o.rect.x + 5, o.rect.y + 5}
		drawText(dst, o.font, "Add new task", textPos, textSize, White)

		drawImageSlice(dst, o.nameRect, o.rectOutline, o.outlineConstr, White)
		if o.nameInputSelected {
			cursor := o.nameTextBox.cursor
			cursor.x += o.nameRect.x + 2
			cursor.y += o.nameRect.y + 5
			drawRect(dst, cursor, White)
			if o.nameTextBox.charCount > 0 {
				textPos = point{o.nameRect.x + 2, o.nameRect.y + 5}
				drawText(dst, o.font, string(o.nameTextBox.GetText()), textPos, textSize, White)
			}
		} else {
			textPos = point{o.nameRect.x + 2, o.nameRect.y + 5}
			drawText(dst, o.font, "Name", textPos, textSize, Color{255, 255, 255, 120})
		}

		drawImageSlice(dst, o.lengthRect, o.rectOutline, o.outlineConstr, White)
		textPos = point{o.lengthRect.x + 2, o.lengthRect.y + o.font.Ascent(textSize)}
		drawText(dst, o.font, "<", textPos, textSize, Color{255, 255, 255, 120})
		textPos[0] = o.lengthRect.x + o.lengthRect.width - (o.font.GlyphAdvance('>', textSize) + 2)
		drawText(dst, o.font, ">", textPos, textSize, Color{255, 255, 255, 120})

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
