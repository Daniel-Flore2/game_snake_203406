package main

import (
    "log"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
    "juego/models"
    "juego/views"
)

func main() {
    pixelgl.Run(run)
}

func run() {
    cfg := pixelgl.WindowConfig{
        Title:  "Snake Game",
        Bounds: pixel.R(0, 0, float64(models.WinWidth), float64(models.WinHeight)),
        VSync:  true,
    }
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        log.Fatal(err)
    }

    models.InitializeGame()

    go models.HandleInput(win)
    go models.GameLogic()

    for !win.Closed() {
        win.Clear(colornames.Black)
        models.mu.Lock()
        switch models.gameState {
        case models.Menu:
            views.DrawMenu(win)
        case models.Playing:
            if models.gameState != models.GameOver {
                views.Draw(win)
            } else {
                models.gameState = models.GameOver
            }
        case models.GameOver:
            views.DrawGameOver(win)
        }
        models.mu.Unlock()
        win.Update()
    }
}