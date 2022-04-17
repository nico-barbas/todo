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

func drawSlider(dst *ebiten.Image, rect, incRect, decRect rectangle, t string) {
	drawImageSlice(dst, rect, rectOutline, rectConstraint, White)
	drawTextCenter(dst, textOptions{
		font: &defaultFont, text: "<", bounds: decRect,
		size: textSize, clr: Color{255, 255, 255, 120},
	})
	drawTextCenter(dst, textOptions{
		font: &defaultFont, text: ">", bounds: incRect,
		size: textSize, clr: Color{255, 255, 255, 120},
	})

	// lSize := defaultFont.MeasureText(t, largeTextSize)[0]
	// tSize := lSize + defaultFont.MeasureText("count", smallTextSize)[0]
	// textPos := point{
	// 	a.countRect.x + (a.countRect.width/2 - tSize/2),
	// 	a.countRect.y + a.font.Ascent(largeTextSize)/2,
	// }
	// drawText(dst, textOptions{
	// 	font: &defaultFont, text: txt, pos: textPos,
	// 	size: largeTextSize, clr: White,
	// })
	// textPos[0] += lSize + 4
	// textPos[1] += (defaultFont.Ascent(largeTextSize) - defaultFont.Ascent(smallTextSize))
	// drawText(dst, textOptions{
	// 	font: &defaultFont, text: "count", pos: textPos,
	// 	size: smallTextSize, clr: Color{255, 255, 255, 120},
	// })
	drawTextCenter(dst, textOptions{
		font: &defaultFont, text: t, bounds: rect,
		size: largeTextSize, clr: White,
	})
}
