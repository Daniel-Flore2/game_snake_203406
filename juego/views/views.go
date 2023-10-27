package views

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"juego/models"
)

// Constants
const (
	GridSize = 20
)

// Global variables
var (
	atlas *text.Atlas
)

// Initialization
func init() {
	// Load a basic font
	fontFace := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	atlas = fontFace
}

// Draw menu
func DrawMenu(win *pixelgl.Window) {
	// Draw menu text
	drawText(win, pixel.V(float64(models.WinWidth/2-100), float64(models.WinHeight/2+20)), colornames.White, "Snake Game")
	drawText(win, pixel.V(float64(models.WinWidth/2-100), float64(models.WinHeight/2-20)), colornames.White, "Press Enter to Start")
	drawText(win, pixel.V(float64(models.WinWidth/2-100), float64(models.WinHeight/2-60)), colornames.White, "Press Q to Quit")
}

// Draw game over screen
func DrawGameOver(win *pixelgl.Window) {
	// Draw game over text
	drawText(win, pixel.V(float64(models.WinWidth/2-100), float64(models.WinHeight/2+20)), colornames.White, "Game Over")
	drawText(win, pixel.V(float64(models.WinWidth/2-100), float64(models.WinHeight/2-20)), colornames.White, "Presiona R para iniciar")
	drawText(win, pixel.V(float64(models.WinWidth/2-100), float64(models.WinHeight/2-60)), colornames.White, "Presiona Q para  quitar")
}

// Draw game board
func Draw(win *pixelgl.Window, score int) {
    imd := imdraw.New(nil)

    // Draw the snake
    for _, s := range models.Snake {
        imd.Color = colornames.Green
        imd.Push(pixel.V(float64(s.X*GridSize), float64(s.Y*GridSize)))
        imd.Push(pixel.V(float64((s.X+1)*GridSize), float64((s.Y+1)*GridSize)))
        imd.Rectangle(0)
    }

    // Draw the obstacles
    for _, o := range models.Obstacles {
        imd.Color = colornames.Blue
        imd.Push(pixel.V(float64(o.X*GridSize), float64(o.Y*GridSize)))
        imd.Push(pixel.V(float64((o.X+1)*GridSize), float64((o.Y+1)*GridSize)))
        imd.Rectangle(0)
    }

    // Draw the food
    imd.Color = colornames.Red
    imd.Push(pixel.V(float64(models.Food.X*GridSize), float64(models.Food.Y*GridSize)))
    imd.Push(pixel.V(float64((models.Food.X+1)*GridSize), float64((models.Food.Y+1)*GridSize)))
    imd.Rectangle(0)

    // Draw the score
    drawText(win, pixel.V(10, float64(models.WinHeight-20)), colornames.White, "Score: "+fmt.Sprint(score))

    imd.Draw(win)
}

// Draw text
func drawText(win *pixelgl.Window, pos pixel.Vec, col color.Color, txt string) {
    basicTxt := text.New(pos, atlas)
    basicTxt.Color = col
    basicTxt.WriteString(txt)
    basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
}
