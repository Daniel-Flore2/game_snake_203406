package models

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel/pixelgl"
)

// velocidad de la serpiente y dimensiones de la ventana.
const (
    GridSize   = 20
    SnakeSpeed = time.Millisecond * 100
    WinWidth   = 800 
    WinHeight  = 600 
)

// Estructura Point representa una posición en el juego con coordenadas X e Y.
type Point struct {
    X, Y int
}

// GameState es un tipo enumerado que representa el estado del juego.
type GameState int

const (
    Menu GameState = iota
    Playing
    GameOver
)

// Variables globales para el estado del juego, serpiente, dirección, comida, puntaje, etc.
var (
    Snake      []Point
    Direction  Point
    Food       Point
    GameStateValue GameState // Almacena el estado actual del juego
    Score      int
    LastUpdate time.Time
    Restart    bool
    GameOverValue bool // Indica si el juego ha terminado
    Mu         sync.Mutex
)

// Direcciones posibles para la serpiente.
var (
    LeftDirection  = Point{-1, 0}
    RightDirection = Point{1, 0}
    UpDirection    = Point{0, 1}
    DownDirection  = Point{0, -1}
)

// Canal para comunicar el puntaje entre goroutines.
var ScoreChan = make(chan int)

var Paused bool

// Función para establecer la dirección de la serpiente.
func SetDirectionValue(direction Point) {
    Direction = direction
}


func HandleInput(win *pixelgl.Window) {
    for !win.Closed() {
        if win.JustPressed(pixelgl.KeyEnter) {
            InitializeGame()
            GameStateValue = Playing
        }
        if win.JustPressed(pixelgl.KeyQ) {
            win.SetClosed(true)
        }
        
        if win.JustPressed(pixelgl.KeyP) {
            // Cambia el estado de pausa
            Paused = !Paused
        }

        if GameStateValue == Playing && !Paused {
            if win.Pressed(pixelgl.KeyLeft) {
                SetDirectionValue(LeftDirection)
            } else if win.Pressed(pixelgl.KeyRight) {
                SetDirectionValue(RightDirection)
            } else if win.Pressed(pixelgl.KeyUp) {
                SetDirectionValue(UpDirection)
            } else if win.Pressed(pixelgl.KeyDown) {
                SetDirectionValue(DownDirection)
            }
        }
    }
}

// PlayBackgroundMusic carga y reproduce música de fondo.
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

    // Inicializa el reproductor de audio.
    speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

    // Controla la reproducción de la música en bucle.
    ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}

    // Reproduce la música.
    speaker.Play(ctrl)
}

// InitializeGame inicializa el estado del juego al comienzo.
func InitializeGame() {
    Snake = []Point{{5, 5}}
    Direction = Point{1, 0}
    
    GenerateFood()
    Score = 0
    GameOverValue = false // Establece el estado del juego como no terminado
    GameStateValue = Menu // Establece el estado del juego como Menú
    LastUpdate = time.Now()

    // Inicia la goroutine de generación de comida.
    StartFoodGenerator()

    // Genera comida inicialmente.
    GenerateFoodAsync()
}

// Update actualiza el estado del juego en función del tiempo.
func Update() {
    currentTime := time.Now()
    elapsedTime := currentTime.Sub(LastUpdate)

    if !Paused { // Agrega esta condición para pausar la serpiente.
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
                ScoreChan <- 1 // Envía 1 al canal para incrementar el puntaje
                GenerateFood()
            } else {
                Snake = Snake[1:]
            }

            if CheckCollision(newHead) {
                GameOverValue = true // Establece el estado del juego como terminado
            }
            LastUpdate = currentTime
        }
    }
}

// GenerateFood genera una nueva ubicación para la comida de forma aleatoria.
func GenerateFood() {
    rand.Seed(time.Now().UnixNano())
    randX := rand.Intn(WinWidth/GridSize)
    randY := rand.Intn(WinHeight/GridSize)
    Food = Point{randX, randY}
}

// foodGeneratorChan se utiliza para coordinar la generación de comida de manera asíncrona.
var foodGeneratorChan = make(chan bool)

// StartFoodGenerator inicia la goroutine para generar comida.
func StartFoodGenerator() {
    go foodGenerator()
}

// foodGenerator es una goroutine que genera la ubicación de la comida cuando se le indica.
func foodGenerator() {
    rand.Seed(time.Now().UnixNano())
    <-foodGeneratorChan // Espera a que se solicite la generación de comida
    randX := rand.Intn(WinWidth/GridSize)
    randY := rand.Intn(WinHeight/GridSize)
    Food = Point{randX, randY}
}

// GenerateFoodAsync solicita la generación de comida de manera asíncrona.
func GenerateFoodAsync() {
    foodGeneratorChan <- true
}

// CheckCollision verifica si la serpiente colisiona con las paredes o con ella misma.
func CheckCollision(head Point) bool {
    if head.X < 0 || head.X >= WinWidth/GridSize || head.Y < 0 || head.Y >= WinHeight/GridSize {
        return true // Colisión con las paredes
    }
    for _, p := range Snake[:len(Snake)-1] {
        if p == head {
            return true // Colisión consigo misma
        }
    }
    return false // Sin colisión
}
