package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	itemPadding = 4
	itemHeight  = textSize + itemPadding*2
)

type (
	listWindow struct {
		rect       rectangle
		activeRect rectangle
		addBtnRect rectangle

		font          *Font
		rectOutline   *ebiten.Image
		outlineConstr constraint

		hovered  *listItem
		selected *listItem

		shouldHighlight bool
		highlightRect   rectangle

		items []listItem
		count int
		cap   int
	}

	listItem struct {
		rect         rectangle
		textPosition point
		checkRect    rectangle
	}
)

func (l *listWindow) init(font *Font, outline *ebiten.Image) {
	l.rect = rectangle{0, 0, 200, windowHeight}
	l.addBtnRect = rectangle{
		x:      l.rect.x + btnPadding,
		y:      l.rect.y + l.rect.height - textSize - (10 * 2) - 6,
		width:  l.rect.width - (btnPadding * 2),
		height: textSize + (10 * 2),
	}

	l.font = font
	l.rectOutline = outline
	l.outlineConstr = constraint{2, 2, 2, 2}
}

func (l *listWindow) update(mPos point, mLeft bool) (selected int) {
	selected = -1
	l.shouldHighlight = false

	if l.rect.boundCheck(mPos) {
		index := int(mPos[1] / itemHeight)
		if index < l.count {
			l.hovered = &l.items[index]
			l.shouldHighlight = true
			if l.hovered.checkRect.boundCheck(mPos) {
				l.highlightRect = l.hovered.checkRect
			} else {
				l.highlightRect = l.hovered.rect
			}
			if mLeft {
				selected = index
				l.selected = l.hovered
			}
		} else if l.addBtnRect.boundCheck(mPos) {
			l.shouldHighlight = true
			l.highlightRect = l.addBtnRect
			if mLeft {
				FireSignal(todoAddBtnPressed, SignalNoArgs)
			}
		} else if mLeft {
			l.selected = nil
		}
	}

	return
}

func (l *listWindow) draw(dst *ebiten.Image, tasks []task) {
	if l.shouldHighlight {
		ebitenutil.DrawRect(
			dst,
			l.highlightRect.x,
			l.highlightRect.y,
			l.highlightRect.width,
			l.highlightRect.height,
			Color{255, 255, 255, 120},
		)
	}

	for i := 0; i < l.count; i += 1 {
		item := &l.items[i]
		task := &tasks[i]

		drawText(dst, l.font, task.name, item.textPosition, textSize, White)
		ebitenutil.DrawLine(
			dst,
			item.rect.x,
			item.rect.y+item.rect.height,
			item.rect.x+item.rect.width,
			item.rect.y+item.rect.height,
			White,
		)

		drawImageSlice(dst, item.checkRect, l.rectOutline, l.outlineConstr, White)
		if task.done {
			rect := rectangle{item.checkRect.x + 3, item.checkRect.y + 3, item.checkRect.width - 6, item.checkRect.height - 6}
			drawRect(dst, rect, White)
		}
	}

	drawImageSlice(dst, l.addBtnRect, l.rectOutline, l.outlineConstr, White)
	textpos := point{
		l.addBtnRect.x + (l.addBtnRect.width/2 - l.font.MeasureText("New Task", textSize)[0]/2),
		l.addBtnRect.y + l.font.Ascent(textSize)/2,
	}
	drawText(dst, l.font, "New Task", textpos, textSize, White)
}

func (l *listWindow) addItem() {
	rect := rectangle{0, float64(l.count) * itemHeight, 200, itemHeight}
	textPos := point{rect.x, rect.y + itemPadding}
	i := listItem{
		rect:         rect,
		textPosition: textPos,
		checkRect:    rectangle{(rect.x + rect.width) - itemHeight, textPos[1], textSize, textSize},
	}
	if l.count > len(l.items) {
		newSlice := make([]listItem, l.cap*2)
		copy(newSlice[:], l.items[:])
		l.items = newSlice
	}
	l.items = append(l.items, i)
	l.count += 1
}
