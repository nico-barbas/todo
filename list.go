package main

import (
	"fmt"
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
	listItemHoverAnimation
)

type (
	listWindow struct {
		rect       rectLayout
		listRect   rectLayout
		addBtnRect rectLayout

		font          *Font
		rectOutline   *ebiten.Image
		outlineConstr constraint

		previousHovered *listItem
		hovered         *listItem
		selected        *listItem

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
		animations   [4]anim.Animation
	}
)

func (l *listWindow) init(font *Font, outline *ebiten.Image) {
	AddSignalListener(todoTaskRemoved, l)
	l.rect = newRectLayout(rectangle{0, 0, 200, windowHeight})
	l.addBtnRect = l.rect.cut(rectCutDown, textSize*2, 0)
	l.listRect = l.rect.cut(rectCutUp, l.rect.remaining.height, 0)
	l.items = make([]listItem, initialTaskCap)
	l.cap = initialTaskCap

	l.font = font
	l.rectOutline = outline
	l.outlineConstr = constraint{2, 2, 2, 2}
}

func (l *listWindow) update(mPos point, mLeft bool) (selected int) {
	selected = -1
	l.shouldHighlight = false
	l.previousHovered = l.hovered
	l.hovered = nil

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		fmt.Println("")
	}

	if l.rect.full.boundCheck(mPos) {
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
		} else if l.addBtnRect.remaining.boundCheck(mPos) {
			l.shouldHighlight = true
			l.highlightRect = l.addBtnRect.remaining
			if mLeft {
				FireSignal(todoAddBtnPressed, SignalNoArgs)
			}
		} else if mLeft {
			l.selected = nil
		}
	}

	if l.hovered != l.previousHovered {
		if l.hovered != nil {
			l.hovered.animations[listItemHoverAnimation].Play()
		}
	}

	for i := 0; i < l.count; i += 1 {
		item := &l.items[i]
		for i := range item.animations {
			item.animations[i].Update()
		}
		if item != l.hovered {
			if !item.animations[listItemAddAnimation].Playing || !item.animations[listItemRemoveAnimation].Playing {
				item.animations[listItemHoverAnimation].Reset()
				item.textPosition[0] = item.rect.x + itemPadding
			}
		}
	}
	return
}

func (l *listWindow) draw(dst *ebiten.Image, tasks []task) {
	if l.shouldHighlight {
		drawRect(dst, l.highlightRect, WhiteA125)
	}

	for i := 0; i < l.count; i += 1 {
		item := &l.items[i]
		task := &tasks[i]

		// Draw progress in case it is running
		if task.timer.running {
			progress := task.progress()
			drawRect(dst, rectangle{
				item.rect.x, item.rect.y,
				item.rect.width * progress, item.rect.height,
			}, WhiteA125)
		}

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
			darkSeparator,
		)

		drawImageSlice(dst, item.checkRect, l.rectOutline, l.outlineConstr, White)
		if task.done {
			rect := rectangle{item.checkRect.x + 3, item.checkRect.y + 3, item.checkRect.width - 6, item.checkRect.height - 6}
			drawRect(dst, rect, White)
		}
	}

	// drawTextBtn(dst, l.addBtnRect, "NewTask", textSize)
	drawRect(dst, rectangle{l.addBtnRect.remaining.x, l.addBtnRect.remaining.y, l.addBtnRect.remaining.width, 1}, darkSeparator)
	drawTextCenter(dst, textOptions{
		font: l.font, text: "New Task", bounds: l.addBtnRect.remaining,
		size: textSize, clr: White,
	})
}

func (l *listWindow) addItem() {
	rect := rectangle{0, float64(l.count) * itemHeight, 200, itemHeight}
	textPos := point{rect.x + itemPadding, rect.y + itemPadding}
	i := listItem{
		rect:         rect,
		textPosition: textPos,
		checkRect: rectangle{
			(rect.x + rect.width) - itemHeight,
			textPos[1] + itemPadding/2,
			textSize - itemPadding,
			textSize - itemPadding,
		},
	}
	if l.count > len(l.items) {
		newSlice := make([]listItem, l.cap*2)
		copy(newSlice[:], l.items[:])
		l.items = newSlice
	}
	l.items[l.count] = i
	l.count += 1

	func(li *listItem) {
		li.animations = [4]anim.Animation{
			anim.NewAnimation("add", l),
			anim.NewAnimation("remove", l),
			anim.NewAnimation("hoverStart", l),
			anim.NewAnimation("hoverEnd", l),
		}
		li.animations[listItemAddAnimation].AddProperty("rectx", &li.rect.x, li.rect.x-200, false)
		li.animations[listItemAddAnimation].AddKey("rectx", anim.AnimationKey{
			Easing:   anim.EaseOutCubic,
			Duration: anim.SecondsToTicks(0.2),
			Change:   li.rect.width,
		})
		li.animations[listItemAddAnimation].AddProperty("textx", &li.textPosition[0], li.textPosition[0]-200, false)
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

		// Remove animation
		li.animations[listItemHoverAnimation].AddProperty("textx", &li.textPosition[0], li.textPosition[0], false)
		li.animations[listItemHoverAnimation].AddKey("textx", anim.AnimationKey{
			Easing:   anim.EaseInCubic,
			Duration: anim.SecondsToTicks(0.1),
			Change:   20,
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
		textPos := point{rect.x + itemPadding, rect.y + itemPadding}
		item := &l.items[i]
		item.rect = rect
		item.textPosition = textPos
		item.checkRect = rectangle{
			(rect.x + rect.width) - itemHeight,
			textPos[1] + itemPadding/2,
			textSize - itemPadding,
			textSize - itemPadding,
		}
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
