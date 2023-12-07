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

	score := &models.GameScore{}
	models.InitializeGame(score)

	obstacleCh := make(chan bool)

	go models.HandleInput(win)
	go models.UpdateSnake()

	// time.Sleep(10 * time.Second)
	go models.StartObstacleGenerator(obstacleCh)

	for !win.Closed() {
		select {
		case <-obstacleCh:
			models.Mu.Lock()
			if models.GameState == models.Playing {
				// Realiza la verificación de colisión con obstáculos
				collision := models.CheckCollisionWithObstacles()
				if collision {
					models.GameState = models.GameOver
				}
			}
			models.Mu.Unlock()
		default:
		}

		win.Clear(colornames.Black)
		models.Mu.Lock()
		switch models.GameState {
		case models.Menu:
			views.DrawMenu(win)
		case models.Playing:
			if models.GameState != models.GameOver {
				scoreValue := score.Get()
				// Pasar puntaje actualizado a Draw
				views.Draw(win, scoreValue)
			}
		case models.GameOver:
			views.DrawGameOver(win)
		}
		models.Mu.Unlock()
		win.Update()
	}
}
