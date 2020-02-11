// Package ctrl contains the game logic.
//
// The engine's Loop() method should be launched as a goroutine,
// and it can be controlled with opaque commands safely from other
// goroutines.
package ctrl

import (
	"log"
	"math/rand"
	"time"

	"github.com/icza/golab/model"
)

const (
	// BlockSize is the size of the labyrinth unit in pixels.
	BlockSize = 40
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Engine calculates and controls the game.
type Engine struct {
	Model *model.Model

	// command channel to control the engine from other goroutines.
	cmdChan chan interface{}

	// invalidate is called by the engine to request a new view frame.
	invalidate func()

	// Current game config
	cfg *GameConfig
}

// NewEngine returns a new Engine.
// invalidate is a func which will be called by the engine to request a new view frame.
func NewEngine(invalidate func()) *Engine {
	e := &Engine{
		Model:      &model.Model{},
		cmdChan:    make(chan interface{}, 10),
		invalidate: invalidate,
	}

	e.initNewGame(&GameConfig{
		Difficulty: Difficulties[DifficultyDefaultIdx],
		LabSize:    LabSizes[LabSizeDefaultIdx],
		Speed:      Speeds[SpeedDefaultIdx],
	})

	return e
}

// NewGame enqueues a new game command with the given config.
func (e *Engine) NewGame(cfg GameConfig) {
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

		e.invalidate()

		time.Sleep(e.cfg.Speed.loopDelay)
	}
}

func (e *Engine) processCmds() {
	for {
		select {

		case cmd := <-e.cmdChan:
			switch cmd := cmd.(type) {
			case *GameConfig:
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
func (e *Engine) initNewGame(cfg *GameConfig) {
	e.cfg = cfg

	m := e.Model

	m.Counter++

	// Init the labyrinth
	m.Rows, m.Cols = cfg.LabSize.rows, cfg.LabSize.cols
	m.Lab = make([][]model.Block, m.Rows)
	for row := range m.Lab {
		m.Lab[row] = make([]model.Block, m.Cols)
	}
	generateLab(m.Lab)

	// Init Gopher
	m.Gopher = new(model.MovingObj)
	m.Gopher.Pos.X = BlockSize + BlockSize/2 // Position Gopher to top left corner
	m.Gopher.Pos.Y = m.Gopher.Pos.X
	m.Gopher.Dir = model.DirRight
	m.Gopher.TargetPos.X = int(m.Gopher.Pos.X)
	m.Gopher.TargetPos.Y = int(m.Gopher.Pos.Y)

	// Init bulldogs
	numBulldogs := int(float64(m.Rows*m.Cols) * cfg.Difficulty.bulldogDensity / 1000)
	m.Bulldogs = make([]*model.MovingObj, numBulldogs)
	for i := range m.Bulldogs {
		bd := new(model.MovingObj)
		m.Bulldogs[i] = bd

		// Place bulldog at a random position
		var row, col = int(m.Gopher.Pos.Y) / BlockSize, int(m.Gopher.Pos.X) / BlockSize
		// Give some space to Gopher: do not generate Bulldogs too close:
		for gr, gc := row, col; (row-gr)*(row-gr) <= 16 && (col-gc)*(col-gc) <= 16; row, col = rPassPos(0, m.Rows), rPassPos(0, m.Cols) {
		}

		bd.Pos.X = float32(col*BlockSize + BlockSize/2)
		bd.Pos.Y = float32(row*BlockSize + BlockSize/2)

		bd.TargetPos.X, bd.TargetPos.Y = int(bd.Pos.X), int(bd.Pos.Y)
	}

	// TODO
}
