// Package model contains types modeling the game.
package model

import "sync"

// Model is the model of the game.
type Model struct {
	// Mu must be locked when Model is accessed concurrently.
	Mu sync.RWMutex
}
