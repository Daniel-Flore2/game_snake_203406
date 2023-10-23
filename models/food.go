package models

import (
    "math/rand"
    "time"
)

var Food Point

func generateFood() {
    rand.Seed(time.Now().UnixNano())
    for {
        randX := rand.Intn(WinWidth/gridSize)
        randY := rand.Intn(WinHeight/gridSize)
        Food = Point{randX, randY}
        collision := false
        for _, p := range Snake {
            if p == Food {
                collision = true
                break
            }
        }
        if !collision {
            break
        }
    }
}
