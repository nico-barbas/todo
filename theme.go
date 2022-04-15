package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	darkBackground1 = Color{32, 32, 32, 255}
	darkBackground2 = Color{25, 25, 25, 255}
	darkBackground3 = Color{12, 12, 12, 255}
	darkSeparator   = Color{100, 92, 87, 255}
)

var (
	rectOutline    *ebiten.Image
	rectConstraint = constraint{2, 2, 2, 2}
	defaultFont    Font
)

func loadTheme() {
	rectOutline, _, _ = ebitenutil.NewImageFromFile("assets/uiRectOutline.png")
	defaultFont = NewFont("assets/FiraSans-Regular.ttf", 72, []int{smallTextSize, textSize, largeTextSize})
}
