package models

import "sync"

type GameScore struct {
    value int
    mu    sync.Mutex
}

func NewGameScore() GameScore {
	return GameScore{} 
}

func (s *GameScore) Get() int {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.value
}

func (s *GameScore) Increase(amount int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.value += amount
}


func (s *GameScore) Reset() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.value = 0
}
