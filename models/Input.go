
package models

import "github.com/faiface/pixel/pixelgl"

var inputChan chan pixelgl.Button

func HandleInput(win *pixelgl.Window) {
	for !win.Closed() {
		if GameState == Menu && win.Pressed(pixelgl.KeyEnter) {
			GameState = Playing
		}

		if GameState == Playing {
			if win.Pressed(pixelgl.KeyLeft) && Direction != (Point{1, 0}) {
				Direction = Point{-1, 0}
			}
			if win.Pressed(pixelgl.KeyRight) && Direction != (Point{-1, 0}) {
				Direction = Point{1, 0}
			}
			if win.Pressed(pixelgl.KeyUp) && Direction != (Point{0, -1}) {
				Direction = Point{0, 1}
			}
			if win.Pressed(pixelgl.KeyDown) && Direction != (Point{0, 1}) {
				Direction = Point{0, -1}
			}
		}

		if GameState == GameOver && win.Pressed(pixelgl.KeyR) {
			InitializeGame(score)
			GameState = Playing
		}

		if win.Pressed(pixelgl.KeyQ) {
			win.SetClosed(true)
		}
	}
}
