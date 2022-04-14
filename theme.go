package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	darkBackground1 = Color{48, 41, 51, 255}
	darkBackground2 = Color{24, 20, 26, 255}
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
