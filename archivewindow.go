package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type (
	archiveWindow struct {
		active   bool
		dirty    bool
		canvas   *ebiten.Image
		position point

		rect      rectLayout
		titleRect rectLayout
		listRect  rectLayout

		font          *Font
		rectOutline   *ebiten.Image
		outlineConstr constraint
		bgFill        *ebiten.Image

		items []archiveItem
		count int
		cap   int
	}

	archiveItem struct {
		offset           point
		rect             rectLayout
		nameRect         rectangle
		finishedRect     rectangle
		archivedDateRect rectangle
	}
)

func (a *archiveWindow) init(font *Font, outline *ebiten.Image) {
	const archiveWindowPadding = 10
	const addWindowMargin = 50

	AddSignalListener(todoArchiveBtnPressed, a)
	a.items = make([]archiveItem, initialTaskCap)

	a.active = false
	a.rect = newRectLayout(rectangle{
		x:      0,
		y:      0,
		width:  windowWidth - 200,
		height: windowHeight - 100,
	})
	a.position = point{
		windowWidth/2 - a.rect.full.width/2,
		windowHeight/2 - a.rect.full.height/2,
	}
	a.rect.cut(rectCutUp, archiveWindowPadding, 0)
	a.rect.cut(rectCutDown, archiveWindowPadding, 0)

	a.titleRect = a.rect.cut(rectCutUp, textSize, 15)
	a.titleRect.cut(rectCutLeft, archiveWindowPadding, 0)
	a.titleRect.cut(rectCutRight, archiveWindowPadding, 0)

	a.rect.cut(rectCutUp, addWindowMargin, 0)

	a.listRect = a.rect.cut(rectCutUp, a.rect.remaining.height, 0)
	a.listRect.cut(rectCutLeft, addWindowMargin, 0)
	a.listRect.cut(rectCutRight, addWindowMargin, 0)

	a.dirty = true
	a.canvas = ebiten.NewImage(int(a.rect.full.width), int(a.rect.full.height))
	a.font = font
	a.rectOutline = outline
	a.outlineConstr = constraint{2, 2, 2, 2}
	a.bgFill = ebiten.NewImage(1, 1)
	a.bgFill.Fill(darkBackground3)
}

func (a *archiveWindow) update(mPos point, mLeft bool) {
	if a.active {
		relPos := mPos.sub(a.position)
		if !a.rect.full.boundCheck(relPos) && mLeft {
			a.active = false
			return
		}
	}
}

func (a *archiveWindow) draw(dst *ebiten.Image, archivedTasks []task) {
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
		// a.elements.highlight(dst)

		rect = a.listRect.remaining.addPoint(a.position)
		drawRect(dst, rect, darkBackground2)
		drawImageSlice(dst, rect, a.rectOutline, a.outlineConstr, darkSeparator)

		if a.dirty {
			a.redraw(archivedTasks)
		}
		drawImage(dst, a.canvas, a.position)
	}
}

func (a *archiveWindow) redraw(archivedTasks []task) {
	a.canvas.Clear()

	drawTextCenter(a.canvas, textOptions{
		font: a.font, text: "Archive", bounds: a.titleRect.remaining,
		size: largeTextSize, clr: White,
	})

	for i := 0; i < a.count; i += 1 {
		item := &a.items[i]
		task := &archivedTasks[i]

		drawTextCenter(a.canvas, textOptions{
			font: a.font, text: task.name, bounds: item.nameRect,
			size: textSize, clr: White,
		})
		drawRect(
			a.canvas,
			rectangle{
				item.rect.full.x,
				item.rect.full.y + item.rect.full.height,
				item.rect.full.width,
				1,
			},
			darkSeparator,
		)
	}
	a.dirty = false
}

func (a *archiveWindow) addItem() {
	rect := newRectLayout(rectangle{
		x:      a.listRect.remaining.x,
		y:      a.listRect.remaining.y + float64(a.count)*itemHeight,
		width:  a.listRect.remaining.width,
		height: itemHeight,
	})
	item := archiveItem{
		rect:     rect,
		nameRect: rect.cut(rectCutLeft, rect.remaining.width-100, 10).full,
	}
	if a.count > len(a.items) {
		newSlice := make([]archiveItem, a.cap*2)
		copy(newSlice[:], a.items[:])
		a.items = newSlice
	}
	a.items[a.count] = item
	a.count += 1

	a.dirty = true
}

func (a *archiveWindow) OnSignal(s Signal) {
	switch s.Kind {
	case todoArchiveBtnPressed:
		a.active = true
	}
}
