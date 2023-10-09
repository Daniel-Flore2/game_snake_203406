package models

import (
    "math/rand"
    "sync"
    "time"

    "github.com/faiface/beep"
    "github.com/faiface/beep/mp3"
    "github.com/faiface/beep/speaker"
    "log"
    "os"
)

const (
    GridSize   = 20
    SnakeSpeed = time.Millisecond * 100
    WinWidth   = 800 // Define el ancho de la ventana aquí
    WinHeight  = 600 // Define la altura de la ventana aquí
)

type Point struct {
    X, Y int
}

type GameState int

const (
    Menu GameState = iota
    Playing
    GameOver
    
)

var (
    Snake      []Point
    Direction  Point
    Food       Point
    GameStateValue GameState // Cambia GameState a GameStateValue
    Score      int
    LastUpdate time.Time
    Restart    bool
    GameOverValue bool // Cambia GameOver a GameOverValue
    Mu         sync.Mutex
)

var (
    LeftDirection  = Point{-1, 0}
    RightDirection = Point{1, 0}
    UpDirection    = Point{0, 1}
    DownDirection  = Point{0, -1}
)

// Función para establecer la dirección
func SetDirectionValue(direction Point) {
    Direction = direction
}

func PlayBackgroundMusic() {
    f, err := os.Open("assets/background_music.mp3")
    if err != nil {
        log.Fatal(err)
        return
    }
    defer f.Close()

    streamer, format, err := mp3.Decode(f)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer streamer.Close()

    speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

    ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}

    speaker.Play(ctrl)
}


func InitializeGame() {
    Snake = []Point{{5, 5}}
    Direction = Point{1, 0}
    
    GenerateFood()
    Score = 0
    GameOverValue = false // Cambia GameOver a GameOverValue

    GameStateValue = Menu // Cambia GameState a GameStateValue
    LastUpdate = time.Now()

    Snake = []Point{{5, 5}}
    Direction = Point{1, 0}
    Score = 0
    GameOverValue = false

    GameStateValue = Menu
    LastUpdate = time.Now()

    // Inicia la goroutine de generación de comida
    StartFoodGenerator()

    // Genera comida inicialmente
    GenerateFoodAsync()
}

func Update() {
    currentTime := time.Now()
    elapsedTime := currentTime.Sub(LastUpdate)

    if elapsedTime >= SnakeSpeed {
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
            GenerateFood()
        } else {
            Snake = Snake[1:]
        }

        if CheckCollision(newHead) {
            GameOverValue = true // Cambia GameOver a GameOverValue
        }

        LastUpdate = currentTime
    }
}

func GenerateFood() {
    rand.Seed(time.Now().UnixNano())
    randX := rand.Intn(WinWidth/GridSize)
    randY := rand.Intn(WinHeight/GridSize)
    Food = Point{randX, randY}
}


var foodGeneratorChan = make(chan bool)

func StartFoodGenerator() {
    go foodGenerator()
}

func foodGenerator() {
    rand.Seed(time.Now().UnixNano())
    <-foodGeneratorChan // Espera a que se solicite la generación de comida
    randX := rand.Intn(WinWidth/GridSize)
    randY := rand.Intn(WinHeight/GridSize)
    Food = Point{randX, randY}
}

func GenerateFoodAsync() {
    foodGeneratorChan <- true // Solicita la generación de comida de forma asíncrona
}

func CheckCollision(head Point) bool {
    if head.X < 0 || head.X >= WinWidth/GridSize || head.Y < 0 || head.Y >= WinHeight/GridSize {
        return true
    }
    for _, p := range Snake[:len(Snake)-1] {
        if p == head {
            return true
        }
    }
    return false
}
