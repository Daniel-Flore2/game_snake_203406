package models

import (
  "math/rand"
  "time"
)

const GridSize = 20

const (
  GameRunning GameStateType = iota
)



type Obstacle struct {
  X, Y int
}

var Obstacles []Obstacle

func GenerateObstacles() {
    rand.Seed(time.Now().UnixNano())

    for {
        if GameState != GameRunning {
            return
        }

        o := Obstacle{
            X: rand.Intn(WinWidth/GridSize),
            Y: rand.Intn(WinHeight/GridSize),
        }

        collision := false
        for _, s := range Snake {
            if s.X == o.X && s.Y == o.Y {
                collision = true
                break
            }
        }

        if !collision {
            Obstacles = append(Obstacles, o)
        }

        time.Sleep(5 * time.Second)
    }
}

func CheckCollision() {
    for _, s := range Snake {
        for _, o := range Obstacles {
            if s.X == o.X && s.Y == o.Y {
                GameState = SnakeHitObstacle
                return
            }
        }
    }
}

func StartObstacleGenerator() {
  GameState = GameRunning
   go GenerateObstacles()
}

