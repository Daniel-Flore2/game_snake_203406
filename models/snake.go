package models

import (
    "sync"
    "time"
    "github.com/faiface/pixel/pixelgl"
)

type Point struct {
    X, Y int
}

type GameStateType int

const (
    Menu GameStateType = iota
    Playing
    GameOver
)

const (
    gridSize   = 20
    snakeSpeed = time.Millisecond * 100
)

var (
    Snake      []Point
    Direction  Point
    gameState  GameStateType
    Score      int
    lastUpdate time.Time
    
    gameOver   bool
    inputChan  chan pixelgl.Button
)

var mu sync.Mutex


var WinWidth = 800
var WinHeight = 600

func InitializeGame() {
    Snake = []Point{{5, 5}}
    Direction = Point{1, 0}
    generateFood()
    Score = 0
    gameOver = false
    gameState = Menu
    lastUpdate = time.Now()
}

func update() {
    currentTime := time.Now()
    elapsedTime := currentTime.Sub(lastUpdate)

    if elapsedTime >= snakeSpeed {
        head := Snake[len(Snake)-1]

        var newHead Point

        switch Direction {
        case Point{-1, 0}:
            newHead = Point{head.X - 1, head.Y}
        case Point{1, 0}:
            newHead = Point{head.X + 1, head.Y}
        case Point{0, 1}:
            newHead = Point{head.X, head.Y + 1}
        case Point{0, -1}:
            newHead = Point{head.X, head.Y - 1}
        }

        Snake = append(Snake, newHead)

        if newHead == Food {
            Score++
            generateFood()
        } else {
            Snake = Snake[1:]
        }

        if checkCollision(newHead) {
            gameOver = true
        }

        lastUpdate = currentTime
    }
}

func HandleInput(win *pixelgl.Window) {
    for !win.Closed() {
        if gameState == Menu && win.Pressed(pixelgl.KeyEnter) {
            gameState = Playing
        }

        if gameState == Playing {
            if win.Pressed(pixelgl.KeyLeft) && Direction != (Point{1, 0}) {
                Direction = Point{-1, 0}
            }
            if win.Pressed(pixelgl.KeyRight) && Direction != (Point{-1, 0}) {
                Direction = Point{1, 0}
            }
            if win.Pressed(pixelgl.KeyUp) && Direction != (Point{0, -1}) {
                Direction = Point{0, 1}
            }
            if win.Pressed(pixelgl.KeyDown) && Direction != (Point{0, 1}) {
                Direction = Point{0, -1}
            }
        }

        if gameState == GameOver && win.Pressed(pixelgl.KeyR) {
            InitializeGame()
            gameState = Playing
        }

        if win.Pressed(pixelgl.KeyQ) {
            win.SetClosed(true)
        }
    }
}

func checkCollision(head Point) bool {
    if head.X < 0 || head.X >= WinWidth/gridSize || head.Y < 0 || head.Y >= WinHeight/gridSize {
        return true
    }
    for _, p := range Snake[:len(Snake)-1] {
        if p == head {
            return true
        }
    }
    return false
}

