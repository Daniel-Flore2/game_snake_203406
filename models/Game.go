
package models

import (
	"sync"
	"time"
)

const (
	gridSize   = 20
	snakeSpeed = time.Millisecond * 100
)

var (
	Snake      []Point
	Direction  Point
	GameState  GameStateType
	score      *GameScore
	lastUpdate time.Time
	Mu         sync.Mutex
	gameOver   bool
	gameScore  int // Declare a variable to hold the score
)

var WinWidth = 800
var WinHeight = 600

func InitializeGame(s *GameScore) {
	score = s
	Snake = []Point{{5, 5}}
	Direction = Point{1, 0}
	generateFood()

	
	gameScore = 0

	gameOver = false
	GameState = Menu
	lastUpdate = time.Now()
}

func updateGame(s *GameScore) {
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
			gameScore++
			score.Increase(1) 
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

func UpdateSnake() {
	for {
		Mu.Lock()
		if gameOver {
			Mu.Unlock()
			time.Sleep(time.Second)
			continue
		}

		if GameState == Playing {
			updateGame(score)
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
