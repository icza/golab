// Package model contains types modeling the game.
package model

import (
	"fmt"
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

	// Size of the labyrinth in blocks.
	Rows, Cols int

	// Blocks of the lab. First indexed by row, then by column.
	Lab [][]Block

	// Our well-beloved hero Gopher
	Gopher *MovingObj

	// The ancient enemies of Gopher: the bloodthirsty Bulldogs.
	Bulldogs []*MovingObj

	// Dead tells if Gopher is dead.
	Dead bool
}

// Block is a square unit of the Labyrinth
type Block int

const (
	// BlockEmpty is the empty, free-to-walk block
	BlockEmpty = iota
	// BlockWall designates an unpassable wall.
	BlockWall

	// BlockCount is not a valid block: just to tell how many blocks there are
	BlockCount
)

// MovingObj describes moving objects in the labyrinth.
type MovingObj struct {
	// The position in the labyrinth in pixel coordinates
	Pos struct {
		X, Y float64
	}

	// Direction this object is facing to
	Dir Dir

	// Target position this object is moving to
	TargetPos image.Point
}

// Dir represents directions
type Dir int

const (
	// DirRight .
	DirRight = iota
	// DirLeft .
	DirLeft
	// DirUp .
	DirUp
	// DirDown .
	DirDown

	// DirCount is not a valid direction: just to tell how many directions there are
	DirCount
)

func (d Dir) String() string {
	switch d {
	case DirRight:
		return "right"
	case DirLeft:
		return "left"
	case DirUp:
		return "up"
	case DirDown:
		return "down"
	}
	return fmt.Sprintf("Dir(%d)", d)
}
