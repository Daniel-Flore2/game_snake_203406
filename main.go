package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
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
	winWidth   = 800
	winHeight  = 600
	gridSize   = 20
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
	mu         sync.Mutex
	inputChan  chan pixelgl.Button
)

func initializeGame() {
	snake = []point{{5, 5}}
	direction = point{1, 0}
	generateFood()
	score = 0
	gameOver = false
	gameState = Menu
	lastUpdate = time.Now()
}

func update() {
	currentTime := time.Now()
	elapsedTime := currentTime.Sub(lastUpdate)

	if elapsedTime >= snakeSpeed {
		head := snake[len(snake)-1]

		var newHead point

		switch direction {
		case point{-1, 0}:
			newHead = point{head.x - 1, head.y}
		case point{1, 0}:
			newHead = point{head.x + 1, head.y}
		case point{0, 1}:
			newHead = point{head.x, head.y + 1}
		case point{0, -1}:
			newHead = point{head.x, head.y - 1}
		}

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

func handleInput(win *pixelgl.Window) {
	for !win.Closed() {
		if gameState == Menu && win.Pressed(pixelgl.KeyEnter) {
			gameState = Playing
		}

		if gameState == Playing {
			if win.Pressed(pixelgl.KeyLeft) && direction != (point{1, 0}) {
				direction = point{-1, 0}
			}
			if win.Pressed(pixelgl.KeyRight) && direction != (point{-1, 0}) {
				direction = point{1, 0}
			}
			if win.Pressed(pixelgl.KeyUp) && direction != (point{0, -1}) {
				direction = point{0, 1}
			}
			if win.Pressed(pixelgl.KeyDown) && direction != (point{0, 1}) {
				direction = point{0, -1}
			}
		}

		if gameState == GameOver && win.Pressed(pixelgl.KeyR) {
			initializeGame()
			gameState = Playing
		}

		if win.Pressed(pixelgl.KeyQ) {
			win.SetClosed(true)
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

func draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)

	for _, p := range snake {
		imd.Color = colornames.Green
		imd.Push(pixel.V(float64(p.x*gridSize), float64(p.y*gridSize)))
		imd.Push(pixel.V(float64((p.x+1)*gridSize), float64((p.y+1)*gridSize)))
		imd.Rectangle(0)
	}

	imd.Color = colornames.Red
	imd.Push(pixel.V(float64(food.x*gridSize), float64(food.y*gridSize)))
	imd.Push(pixel.V(float64((food.x+1)*gridSize), float64((food.y+1)*gridSize)))
	imd.Rectangle(0)

	drawText(win, pixel.V(10, winHeight-20), pixel.RGBA{1, 1, 1, 1}, "Score: "+fmt.Sprintf("%d", score))

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
		randX := rand.Intn(winWidth/gridSize)
		randY := rand.Intn(winHeight/gridSize)
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

func drawGameOver(win *pixelgl.Window) {
	drawText(win, pixel.V(winWidth/2-100, winHeight/2+20), pixel.RGBA{1, 1, 1, 1}, "Game Over")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-20), pixel.RGBA{1, 1, 1, 1}, "Press R to Restart")
	drawText(win, pixel.V(winWidth/2-100, winHeight/2-60), pixel.RGBA{1, 1, 1, 1}, "Press Q to Quit")
}

func gameLogic() {
	for {
		mu.Lock()
		if gameState == Playing && !gameOver {
			update()
		}
		mu.Unlock()
		time.Sleep(snakeSpeed)
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Snake Game",
		Bounds: pixel.R(0, 0, winWidth, winHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Cargar una fuente TTF
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	// Crear un font.Face a partir de ttfFont
	face := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 12,
		DPI:  72,
	})
	atlas = text.NewAtlas(
		face,
		text.ASCII,
	)

	initializeGame()

	go handleInput(win)
	go gameLogic()

	for !win.Closed() {
		win.Clear(colornames.Black)
		mu.Lock()
		switch gameState {
		case Menu:
			drawMenu(win)
		case Playing:
			if !gameOver {
				draw(win)
			} else {
				gameState = GameOver
			}
		case GameOver:
			drawGameOver(win)
		}
		mu.Unlock()
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
