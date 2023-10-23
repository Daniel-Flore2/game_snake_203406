package models

import (
    "math/rand"
    "time"
)

const GridSize = 20

type Obstacle struct {
    X, Y int
}


var Obstacles []Point

// Agrega un nuevo tipo de estado para el generador de obstáculos
const (
    GameRunning GameStateType = iota
    GeneratingObstacles
)

func GenerateObstacles(obstacleCh chan<- bool) {
    rand.Seed(time.Now().UnixNano())

    for GameState == GeneratingObstacles {
        x := rand.Intn(WinWidth / GridSize)
        y := rand.Intn(WinHeight / GridSize)
        obstacle := Point{x, y}

        collision := false
        for _, s := range Snake {
            if obstacle == s {
                collision = true
                GameState = GameOver
                break
            }
        }

        if !collision {
            Obstacles = append(Obstacles, obstacle)
            obstacleCh <- true
        }

        // Espera 10 segundos antes de generar otro obstáculo
        time.Sleep(10 * time.Second)

        if GameState == GameOver {
            break
        }

        Obstacles = nil
    }

    // Informa al canal que se ha terminado la generación de obstáculos
    obstacleCh <- true
}

func StartObstacleGenerator(obstacleCh chan<- bool) {
    GameState = GeneratingObstacles
    go GenerateObstacles(obstacleCh)
}

func CheckCollisionWithObstacles() bool {
    head := Snake[len(Snake)-1]

    for _, o := range Obstacles {
        if o.X == head.X && o.Y == head.Y {
            return true
        }
    }

    return false
}