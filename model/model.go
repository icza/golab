// Package model contains types modeling the game.
package model

import (
	"image"
	"sync"
)

// Model is the model of the game.
type Model struct {
	// Mutex to protect the model from concurrent access.
	sync.RWMutex

	// Game counter. Must be increased by one when a new game is initialized.
	// Can be used to invalidate caches when its value changes.
	Counter int

	// Blocks of the lab. First indexed by row, then by column.
	Lab [][]Block

	// Our well-beloved hero Gopher
	Gopher *MovingObj

	// The ancient enemies of Gopher: the bloodthirsty Bulldogs.
	Bulldogs []*MovingObj
}

// Block is a square unit of the Labyrinth
type Block int

const (
	// BlockEmpty is the empty, free-to-walk block
	BlockEmpty = iota
	// BlockWall designates an unpassable wall.
	BlockWall
)

// MovingObj describes moving objects in the labyrinth.
type MovingObj struct {
	// Position in pixel coordinates.
	PosX, PosY float64

	// Direction this object is facing to
	Dir Dir

	// Target position this object is moving to
	TargetPos image.Point
}

// Dir represents directions
type Dir int

const (
	// DirDown .
	DirDown = iota
	// DirUp .
	DirUp
	// DirRight .
	DirRight
	// DirLeft .
	DirLeft
)
