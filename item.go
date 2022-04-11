package main

type (
	item struct {
		name       string
		done       bool
		charBuffer []byte

		sessionRequired  int
		sessionCompleted int
		sessionLength    minute

		// Render info
		rect         rectangle
		textPosition point
		checkRect    rectangle
	}
)
