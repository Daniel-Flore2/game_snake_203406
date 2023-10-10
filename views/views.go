package views

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"juego/models" 
)

const (
	winWidth  = 800
	winHeight = 600
	gridSize  = 20
)

var (
	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
)

// drawText es una función auxiliar para dibujar texto en la ventana.
func drawText(win *pixelgl.Window, pos pixel.Vec, col pixel.RGBA, textStr string) {
	txt := text.New(pos, atlas)
	txt.Color = col
	txt.WriteString(textStr)
	txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))
}

// DrawMenu dibuja la pantalla de inicio del juego.
func DrawMenu(win *pixelgl.Window) {
	drawText(win, pixel.V(winWidth/2-100, winHeight/2+20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Snake Game")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Presiona Enter para iniciar")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-60), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Presiona Q para salir")
}

// Draw dibuja el estado actual del juego, incluyendo la serpiente y la comida.
func Draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	// Dibuja la serpiente.
	for _, p := range models.Snake {
		imd.Color = colornames.Green
		imd.Push(pixel.V(float64(p.X*gridSize), float64(p.Y*gridSize)))
		imd.Push(pixel.V(float64((p.X+1)*gridSize), float64((p.Y+1)*gridSize)))
		imd.Rectangle(0)
	}

	// Dibuja la comida.
	imd.Color = colornames.Red
	imd.Push(pixel.V(float64(models.Food.X*gridSize), float64(models.Food.Y*gridSize)))
	imd.Push(pixel.V(float64((models.Food.X+1)*gridSize), float64((models.Food.Y+1)*gridSize)))
	imd.Rectangle(0)

	// Dibuja el puntaje en la esquina superior izquierda.
	drawText(win, pixel.V(10, winHeight-20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Score: "+fmt.Sprintf("%d", models.Score))

	imd.Draw(win)
}

// DrawGameOver dibuja la pantalla de fin del juego cuando el jugador pierde.
func DrawGameOver(win *pixelgl.Window) {
	drawText(win, pixel.V(winWidth/2-100, winHeight/2+20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Game Over")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-20), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Presiona enter para reiniciar")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-60), pixel.RGBA{R: 1, G: 1, B: 1, A: 1}, "Presiona Q para Salir")
}
