package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Get work done")

	todo := new(Todo)
	todo.Init()

	todo.addItem("Clean house")
	todo.addItem("Clean desk")
	if err := ebiten.RunGame(todo); err != nil {
		e := err.(exitStatus)
		if e.kind != exitNoError {
			panic(e)
		}
	}
}
