
package models

type GameStateType int

const (
	Menu GameStateType = iota
	Playing
	GameOver
	SnakeHitObstacle
)
