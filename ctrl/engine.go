// Package ctrl contains the game logic.
package ctrl

import (
	"log"
	"time"

	"github.com/icza/golab/model"
)

// Engine calculates and controls the game.
type Engine struct {
	Model *model.Model

	// command channel to control the engine from other goroutines.
	cmdChan chan interface{}
}

// NewEngine returns a new Engine.
func NewEngine() *Engine {
	e := &Engine{
		Model:   &model.Model{},
		cmdChan: make(chan interface{}, 10),
	}

	e.NewGame(NewGameConfig{
		Difficulty: Difficulties[DifficultyDefaultIdx],
		LabSize:    LabSizes[LabSizeDefaultIdx],
		Speed:      Speeds[SpeedDefaultIdx],
	})

	return e
}

// NewGame enqueues a new game command with the given config.
func (e *Engine) NewGame(cfg NewGameConfig) {
	e.cmdChan <- &cfg
}

// Loop starts calculating the game.
// This function returns only if the user closes the app.
func (e *Engine) Loop() {
	for {
		e.Model.Lock()

		e.processCmds()

		e.Model.Unlock()
		// TODO

		time.Sleep(time.Millisecond)
	}
}

func (e *Engine) processCmds() {
	for {
		select {

		case cmd := <-e.cmdChan:
			switch cmd := cmd.(type) {
			case *NewGameConfig:
				e.initNewGame(cmd)
			default:
				log.Printf("Unhandled cmd type: %T", cmd)
			}

		default:
			return // No more commands queued
		}
	}
}

// initNewGame initializes a new game.
func (e *Engine) initNewGame(cfg *NewGameConfig) {
	m := e.Model

	m.Counter++

	m.Lab = make([][]model.Block, cfg.LabSize.rows)
	for row := range m.Lab {
		m.Lab[row] = make([]model.Block, cfg.LabSize.cols)
	}

	// TODO
}
