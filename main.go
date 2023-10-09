package main

import (
    "log"
    "sync"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "juego/models"
    "juego/views"
)

var (
    mu sync.Mutex
)

func handleInput(win *pixelgl.Window) {
    for !win.Closed() {
        if win.JustPressed(pixelgl.KeyEnter) {
            models.InitializeGame()
            models.GameStateValue = models.Playing
        }
        if win.JustPressed(pixelgl.KeyQ) {
            win.SetClosed(true)
        }

        if models.GameStateValue == models.Playing {
            if win.Pressed(pixelgl.KeyLeft) {
                models.SetDirectionValue(models.LeftDirection)
            } else if win.Pressed(pixelgl.KeyRight) {
                models.SetDirectionValue(models.RightDirection)
            } else if win.Pressed(pixelgl.KeyUp) {
                models.SetDirectionValue(models.UpDirection)
            } else if win.Pressed(pixelgl.KeyDown) {
                models.SetDirectionValue(models.DownDirection)
            }
        }
    }
}

func run() {
	models.PlayBackgroundMusic()
    winWidth := 800
    winHeight := 600

    cfg := pixelgl.WindowConfig{
        Title:  "Snake Game",
        Bounds: pixel.R(0, 0, float64(winWidth), float64(winHeight)),
        VSync:  true,
    }
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Inicializar el juego
	
    models.InitializeGame()

    // Iniciar la goroutine para manejar la entrada del usuario
    go handleInput(win)

    for !win.Closed() {
        win.Clear(pixel.RGB(0, 0, 0)) // Cambia el color de fondo si lo deseas

        mu.Lock()

        switch models.GameStateValue {
        case models.Menu:
            views.DrawMenu(win)
        case models.Playing:
            if !models.GameOverValue {
                models.Update() // Actualizar el juego
                views.Draw(win)
            } else {
                models.GameStateValue = models.GameOver
            }
        case models.GameOver:
            views.DrawGameOver(win)
        }

        mu.Unlock()
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}
