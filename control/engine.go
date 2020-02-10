// Package control contains the game logic.
package control

import "github.com/icza/golab/model"

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
