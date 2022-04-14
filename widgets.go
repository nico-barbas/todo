package main

import "github.com/hajimehoshi/ebiten/v2"

func drawTextBtn(dst *ebiten.Image, rect rectangle, text string, size float64) {
	drawImageSlice(dst, rect, rectOutline, rectConstraint, White)
	drawTextCenter(dst, textOptions{
		font: &defaultFont, text: text, bounds: rect,
		size: size, clr: White,
	})
}

func drawIcontBtn(dst *ebiten.Image, rect rectangle, icon *ebiten.Image) {
	drawImageSlice(dst, rect, rectOutline, rectConstraint, White)
	drawImageCentered(dst, icon, rect, 1, White)
}
