package views

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"juego/models" // Asegúrate de que este es el camino correcto a tu paquete models
)

const (
	winWidth  = 800
	winHeight = 600
	gridSize  = 20
)

var (
	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
)

func drawText(win *pixelgl.Window, pos pixel.Vec, col pixel.RGBA, textStr string) {
	txt := text.New(pos, atlas)
	txt.Color = col
	txt.WriteString(textStr)
	txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))
}

func DrawMenu(win *pixelgl.Window) {
	drawText(win, pixel.V(winWidth/2-100, winHeight/2+20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Snake Game")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Pressiona Enter para iniciar")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-60), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Pressiona Q para salir")
}

func Draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	for _, p := range models.Snake {
		imd.Color = colornames.Green
		imd.Push(pixel.V(float64(p.X*gridSize), float64(p.Y*gridSize)))
		imd.Push(pixel.V(float64((p.X+1)*gridSize), float64((p.Y+1)*gridSize)))
		imd.Rectangle(0)
	}

	imd.Color = colornames.Red
	imd.Push(pixel.V(float64(models.Food.X*gridSize), float64(models.Food.Y*gridSize)))
	imd.Push(pixel.V(float64((models.Food.X+1)*gridSize), float64((models.Food.Y+1)*gridSize)))
	imd.Rectangle(0)

	drawText(win, pixel.V(10, winHeight-20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Score: "+fmt.Sprintf("%d", models.Score))

	imd.Draw(win)
}

func DrawGameOver(win *pixelgl.Window) {
	drawText(win, pixel.V(winWidth/2-100, winHeight/2+20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Game Over")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Presiona R para reinniiciar")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-60), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Pressiona Q para Salir")
}
