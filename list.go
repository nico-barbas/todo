package main

import (
	"todo/anim"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	itemPadding = 4
	itemHeight  = textSize + itemPadding*2
)

const (
	listItemAddAnimation = iota
	listItemRemoveAnimation
	listItemMoveAnimation
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

		items       []listItem
		count       int
		cap         int
		removeIndex int
	}

	listItem struct {
		rect         rectangle
		textPosition point
		checkRect    rectangle
		animations   [3]anim.Animation
	}
)

func (l *listWindow) init(font *Font, outline *ebiten.Image) {
	AddSignalListener(todoTaskRemoved, l)
	l.rect = rectangle{0, 0, 200, windowHeight}
	l.addBtnRect = rectangle{
		x:      l.rect.x + btnPadding,
		y:      l.rect.y + l.rect.height - textSize - (10 * 2) - 6,
		width:  l.rect.width - (btnPadding * 2),
		height: textSize + (10 * 2),
	}
	l.items = make([]listItem, initialTaskCap)
	l.cap = initialTaskCap

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

	for i := 0; i < l.count; i += 1 {
		item := &l.items[i]
		for i := range item.animations {
			item.animations[i].Update()
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

		drawText(dst, textOptions{
			font: l.font, text: task.name, pos: item.textPosition,
			size: textSize, clr: White,
		})
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

	drawTextBtn(dst, l.addBtnRect, "NewTask", textSize)
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
	l.items[l.count] = i
	l.count += 1

	func(li *listItem) {
		li.animations = [3]anim.Animation{
			anim.NewAnimation("add", l),
			anim.NewAnimation("remove", l),
			anim.NewAnimation("move", l),
		}
		li.animations[listItemAddAnimation].AddProperty("rectx", &li.rect.x, -200, false)
		li.animations[listItemAddAnimation].AddKey("rectx", anim.AnimationKey{
			Easing:   anim.EaseOutCubic,
			Duration: anim.SecondsToTicks(0.2),
			Change:   li.rect.width,
		})
		li.animations[listItemAddAnimation].AddProperty("textx", &li.textPosition[0], -200, false)
		li.animations[listItemAddAnimation].AddKey("textx", anim.AnimationKey{
			Easing:   anim.EaseOutCubic,
			Duration: anim.SecondsToTicks(0.2),
			Change:   li.rect.width,
		})
		li.animations[listItemAddAnimation].AddProperty("checkrectx", &li.checkRect.x, li.checkRect.x-200, false)
		li.animations[listItemAddAnimation].AddKey("checkrectx", anim.AnimationKey{
			Easing:   anim.EaseOutCubic,
			Duration: anim.SecondsToTicks(0.2),
			Change:   li.rect.width,
		})

		// Remove animation
		li.animations[listItemRemoveAnimation].AddProperty("rectx", &li.rect.x, li.rect.x, false)
		li.animations[listItemRemoveAnimation].AddKey("rectx", anim.AnimationKey{
			Easing:   anim.EaseOutCubic,
			Duration: anim.SecondsToTicks(0.2),
			Change:   -li.rect.width,
		})
		li.animations[listItemRemoveAnimation].AddProperty("textx", &li.textPosition[0], li.textPosition[0], false)
		li.animations[listItemRemoveAnimation].AddKey("textx", anim.AnimationKey{
			Easing:   anim.EaseOutCubic,
			Duration: anim.SecondsToTicks(0.2),
			Change:   -li.rect.width,
		})
		li.animations[listItemRemoveAnimation].AddProperty("checkrectx", &li.checkRect.x, li.checkRect.x, false)
		li.animations[listItemRemoveAnimation].AddKey("checkrectx", anim.AnimationKey{
			Easing:   anim.EaseOutCubic,
			Duration: anim.SecondsToTicks(0.2),
			Change:   -li.rect.width,
		})

		li.animations[listItemAddAnimation].Play()
	}(&l.items[l.count-1])
}

func (l *listWindow) removeItem(at int) {
	copy(l.items[l.removeIndex:], l.items[l.removeIndex+1:])
	l.count -= 1
	l.orderItems()
}

func (l *listWindow) orderItems() {
	for i := 0; i < l.count; i += 1 {
		rect := rectangle{0, float64(i * itemHeight), 200, itemHeight}
		textPos := point{rect.x, rect.y + itemPadding}
		item := &l.items[i]
		item.rect = rect
		item.textPosition = textPos
		item.checkRect = rectangle{(rect.x + rect.width) - itemHeight, textPos[1], textSize, textSize}
		item.animations[listItemAddAnimation].SetPropertyRef("rectx", &item.rect.x)
		item.animations[listItemAddAnimation].SetPropertyRef("textx", &item.textPosition[0])
		item.animations[listItemAddAnimation].SetPropertyRef("checkrectx", &item.checkRect.x)

		item.animations[listItemRemoveAnimation].SetPropertyRef("rectx", &item.rect.x)
		item.animations[listItemRemoveAnimation].SetPropertyRef("textx", &item.textPosition[0])
		item.animations[listItemRemoveAnimation].SetPropertyRef("checkrectx", &item.checkRect.x)

	}
}

func (l *listWindow) OnSignal(s Signal) {
	switch s.Kind {
	case todoTaskRemoved:
		l.selected.animations[listItemRemoveAnimation].Play()
	}
}

func (l *listWindow) OnAnimationEnd(name string) {
	switch name {
	case "add":
	case "remove":
		FireSignal(todoTaskRemoveAnimationDone, SignalNoArgs)
	}
}
