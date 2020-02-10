// Package ctrl contains the game logic.
package ctrl

import (
	"time"

	"github.com/icza/golab/model"
)

// Engine calculates and controls the game.
type Engine struct {
	Model *model.Model
}

// NewEngine returns a new Engine.
func NewEngine() *Engine {
	return &Engine{
		Model: &model.Model{},
	}
}

// Loop starts calculating the game.
// This function returns only if the user closes the app.
func (e *Engine) Loop() {
	// TODO
	for {
		time.Sleep(time.Millisecond)
	}
}
