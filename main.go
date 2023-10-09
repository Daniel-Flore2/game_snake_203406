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

// handleInput maneja las entradas del usuario, como iniciar el juego y salir.
func handleInput(win *pixelgl.Window) {
    for !win.Closed() {
        if win.JustPressed(pixelgl.KeyEnter) {
            models.InitializeGame() // Inicializa el juego
            models.GameStateValue = models.Playing // Cambia el estado del juego a "Jugando"
        }
        if win.JustPressed(pixelgl.KeyQ) {
            win.SetClosed(true) // Cierra la ventana si se presiona "Q"
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

    // Inicializar el juego
    models.InitializeGame()

    // Iniciar la goroutine para manejar la entrada del usuario
    go handleInput(win)

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
    pixelgl.Run(run) // Ejecuta la función "run" utilizando el motor de PixelGL
}
