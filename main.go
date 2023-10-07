package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	winWidth  = 800
	winHeight = 600
	gridSize  = 20
	snakeSpeed = time.Millisecond * 100
)

type point struct {
	x, y int
}

type GameState int

const (
	Menu GameState = iota
	Playing
	GameOver
)

var (
	snake      []point
	direction  point
	food       point
	gameState  GameState
	score      int
	lastUpdate time.Time
	atlas      *text.Atlas
	restart    bool
	gameOver   bool
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Snake Game",
		Bounds: pixel.R(0, 0, winWidth, winHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Cargar una fuente TTF
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	// Crear un font.Face a partir de ttfFont
	face := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 12, // Tamaño de fuente deseado
		DPI:  72, // Resolución DPI deseada
	})

	// Crear un atlas de texto con el font.Face y la tabla ASCII
	atlas = text.NewAtlas(
		face,
		text.ASCII,
	)

	initializeGame()

	for !win.Closed() {
		switch gameState {
		case Menu:
			drawMenu(win)
			if win.Pressed(pixelgl.KeyEnter) {
				gameState = Playing
			} else if win.Pressed(pixelgl.KeyQ) {
				win.SetClosed(true)
			}
		case Playing:
			if restart {
				initializeGame()
				restart = false
			}
			if !gameOver {
				update(win)
			}
			win.Clear(colornames.Black)
			if !gameOver {
				draw(win)
			} else {
				gameState = GameOver
			}
		case GameOver:
			drawGameOver(win)
			if win.Pressed(pixelgl.KeyR) {
				gameState = Playing
				restart = true
			} else if win.Pressed(pixelgl.KeyQ) {
				win.SetClosed(true)
			}
		}

		win.Update()
	}
}

func initializeGame() {
	snake = []point{{5, 5}}
	direction = point{1, 0}
	generateFood()
	score = 0
	gameOver = false
	gameState = Menu
	lastUpdate = time.Now()
}

func update(win *pixelgl.Window) {
    currentTime := time.Now()
    elapsedTime := currentTime.Sub(lastUpdate)

    if elapsedTime >= snakeSpeed {
        head := snake[len(snake)-1]

        if win.Pressed(pixelgl.KeyLeft) {
            direction = point{-1, 0}
        }
        if win.Pressed(pixelgl.KeyRight) {
            direction = point{1, 0}
        }
        if win.Pressed(pixelgl.KeyUp) {
            direction = point{0, 1}
        }
        if win.Pressed(pixelgl.KeyDown) {
            direction = point{0, -1}
        }

        newHead := point{head.x + direction.x, head.y + direction.y}
        snake = append(snake, newHead)

        if newHead == food {
            score++
            generateFood()
        } else {
            snake = snake[1:]
        }

        if checkCollision(newHead) {
            gameOver = true
        }

        lastUpdate = currentTime
    }
}


func draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	// Dibujar serpiente
	for _, p := range snake {
		imd.Color = colornames.Green
		imd.Push(pixel.V(float64(p.x*gridSize), float64(p.y*gridSize)))
		imd.Push(pixel.V(float64((p.x+1)*gridSize), float64((p.y+1)*gridSize)))
		imd.Rectangle(0)
	}

	// Dibujar comida
	imd.Color = colornames.Red
	imd.Push(pixel.V(float64(food.x*gridSize), float64(food.y*gridSize)))
	imd.Push(pixel.V(float64((food.x+1)*gridSize), float64((food.y+1)*gridSize)))
	imd.Rectangle(0)

	// Dibujar puntaje
	drawText(win, pixel.V(10, winHeight-20), pixel.RGBA{1, 1, 1, 1}, "Score: "+strconv.Itoa(score))

	imd.Draw(win)
}

func drawText(win *pixelgl.Window, pos pixel.Vec, col pixel.RGBA, textStr string) {
	txt := text.New(pos, atlas)
	txt.Color = col
	txt.WriteString(textStr)
	txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))
}

func generateFood() {
	rand.Seed(time.Now().UnixNano())
	for {
		randX := rand.Intn(winWidth / gridSize)
		randY := rand.Intn(winHeight / gridSize)
		food = point{randX, randY}
		collision := false
		for _, p := range snake {
			if p == food {
				collision = true
				break
			}
		}
		if !collision {
			break
		}
	}
}

func checkCollision(head point) bool {
	if head.x < 0 || head.x >= winWidth/gridSize || head.y < 0 || head.y >= winHeight/gridSize {
		return true
	}
	for _, p := range snake[:len(snake)-1] {
		if p == head {
			return true
		}
	}
	return false
}

func drawMenu(win *pixelgl.Window) {
	drawText(win, pixel.V(winWidth/2-100, winHeight/2+20), pixel.RGBA{1, 1, 1, 1}, "Snake Game")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-20), pixel.RGBA{1, 1, 1, 1}, "Press Enter to Start")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-60), pixel.RGBA{1, 1, 1, 1}, "Press Q to Quit")
}

func drawGameOver(win *pixelgl.Window) {
	drawText(win, pixel.V(winWidth/2-100, winHeight/2+20), pixel.RGBA{1, 1, 1, 1}, "Game Over")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-20), pixel.RGBA{1, 1, 1, 1}, "Press R to Restart")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-60), pixel.RGBA{1, 1, 1, 1}, "Press Q to Quit")
}

func main() {
	pixelgl.Run(run)
}