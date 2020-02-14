// This file contains models of the game.

package engine

import (
	"fmt"
	"image"
	"math"
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

	// ExitPos: the position Gopher has to reach to win the game.
	ExitPos image.Point

	// Our well-beloved hero Gopher
	Gopher *MovingObj

	// The ancient enemies of Gopher: the bloodthirsty Bulldogs.
	Bulldogs []*MovingObj

	// Dead tells if Gopher is dead.
	Dead bool

	// Won tells if we won
	Won bool

	// For Gopher we maintain multiple target positions which specify a path on which Gopher will move along
	TargetPoss []image.Point
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

// steps steps the MovingObj.
func (m *MovingObj) step() {
	x, y := int(m.Pos.X), int(m.Pos.Y)

	// Only horizontal or vertical movement is allowed!
	if x != m.TargetPos.X {
		dx := math.Min(dt*v, math.Abs(float64(m.TargetPos.X)-m.Pos.X))
		if x > m.TargetPos.X {
			dx = -dx
			m.Dir = DirLeft
		} else {
			m.Dir = DirRight
		}
		m.Pos.X += dx
	} else if y != m.TargetPos.Y {
		dy := math.Min(dt*v, math.Abs(float64(m.TargetPos.Y)-m.Pos.Y))
		if y > m.TargetPos.Y {
			dy = -dy
			m.Dir = DirUp
		} else {
			m.Dir = DirDown
		}
		m.Pos.Y += dy
	}
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
