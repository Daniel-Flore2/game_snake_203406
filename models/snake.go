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
    GameState  GameStateType
    Score      int
    lastUpdate time.Time
    Mu         sync.Mutex
    gameOver   bool
    inputChan  chan pixelgl.Button
)

var WinWidth = 800
var WinHeight = 600

func InitializeGame() {
    Snake = []Point{{5, 5}}
    Direction = Point{1, 0}
    generateFood()
    Score = 0
    gameOver = false
    GameState = Menu
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
            GameState = GameOver
        }

        lastUpdate = currentTime
    }
}

func HandleInput(win *pixelgl.Window) {
    for !win.Closed() {
        if GameState == Menu && win.Pressed(pixelgl.KeyEnter) {
            GameState = Playing
        }

        if GameState == Playing {
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

        if GameState == GameOver && win.Pressed(pixelgl.KeyR) {
            InitializeGame()
            GameState = Playing
        }

        if win.Pressed(pixelgl.KeyQ) {
            win.SetClosed(true)
        }
    }
}

func GameLogic() {
    for {
        Mu.Lock()
        if gameOver {
            Mu.Unlock()
            time.Sleep(time.Second)
            continue
        }

        if GameState == Playing {
            update()
        }
        Mu.Unlock()
        time.Sleep(snakeSpeed)
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



