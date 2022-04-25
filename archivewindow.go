package main

import "github.com/hajimehoshi/ebiten/v2"

type archiveWindow struct {
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
}

func (a *archiveWindow) init(font *Font, outline *ebiten.Image) {
	a.active = true
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

	a.dirty = true
	a.canvas = ebiten.NewImage(int(a.rect.full.width), int(a.rect.full.height))
	a.font = font
	a.rectOutline = outline
	a.outlineConstr = constraint{2, 2, 2, 2}
	a.bgFill = ebiten.NewImage(1, 1)
	a.bgFill.Fill(darkBackground3)
}

func (a *archiveWindow) draw(dst *ebiten.Image) {
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
	}
}
