// Package model contains types modeling the game.
package model

import "sync"

// Model is the model of the game.
type Model struct {
	// Mutex to protect the model from concurrent access.
	sync.RWMutex

	// Game counter. Must be increased by one when a new game is initialized.
	// Can be used to invalidate caches when its value changes.
	Counter int

	// Blocks of the lab. First indexed by row, then by column.
	Lab [][]Block
}

// Block is a square unit of the Labyrinth
type Block int

const (
	// BlockEmpty is the empty, free-to-walk block
	BlockEmpty = iota
	// BlockWall designates an unpassable wall.
	BlockWall
)
