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

func run() {
    models.PlayBackgroundMusic() // Reproduce la música de fondo
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

    models.InitializeGame()

    // Iniciar la goroutine para manejar la entrada del usuario
    go models.HandleInput(win)

    // Iniciar la goroutine para gestionar el puntaje
    go func() {
        for {
            select {
            case scoreIncrement := <-models.ScoreChan:
                models.Score += scoreIncrement // Actualiza el puntaje del juego
            }
        }
    }()

    for !win.Closed() {
        win.Clear(pixel.RGB(0, 0, 0)) // Limpia la ventana con un fondo negro

        mu.Lock()

        switch models.GameStateValue {
        case models.Menu:
            views.DrawMenu(win) // Dibuja la pantalla de menú
        case models.Playing:
            if !models.GameOverValue {
                models.Update() // Actualiza el juego
                views.Draw(win) // Dibuja el estado actual del juego
            } else {
                models.GameStateValue = models.GameOver
            }
        case models.GameOver:
            views.DrawGameOver(win) // Dibuja la pantalla de fin del juego
        }

        mu.Unlock()
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}
