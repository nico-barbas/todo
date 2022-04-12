package main

import "github.com/hajimehoshi/ebiten/v2"

type optWindow struct {
	active bool
	rect   rectangle

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

	o.font = font
	o.rectOutline = outline
	o.outlineConstr = constraint{2, 2, 2, 2}
}

func (o *optWindow) draw(dst *ebiten.Image) {
	if o.active {
		drawRect(dst, o.rect, Black)
		drawImageSlice(dst, o.rect, o.rectOutline, o.outlineConstr, White)

		textPos := point{o.rect.x + 5, o.rect.y + 5}
		drawText(dst, o.font, "Add new task", textPos, textSize, White)
	}
}

func (o *optWindow) OnSignal(s Signal) {
	switch s.Kind {
	case todoAddBtnPressed:
		o.active = true
	}
}
