// Package ctrl contains the game logic.
//
// The engine's Loop() method should be launched as a goroutine,
// and it can be controlled with opaque commands safely from other
// goroutines.
package ctrl

import (
	"log"
	"math"
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

var (
	// dt is the delta time between iterations.
	// We keep this fixed to simulate slower / faster game speeds.
	dt = (50 * time.Millisecond).Seconds()

	// v is the moving speed of Gopher and the Buddlogs in pixel/sec.
	v = 2.0 * BlockSize
)

// Engine calculates and controls the game.
type Engine struct {
	Model *model.Model

	// command channel to control the engine from other goroutines.
	cmdChan chan interface{}

	// invalidate is called by the engine to request a new view frame.
	invalidate func()

	// Current game config
	cfg *GameConfig

	// directions is a reused slice of all directions
	directions []model.Dir
}

// NewEngine returns a new Engine.
// invalidate is a func which will be called by the engine to request a new view frame.
func NewEngine(invalidate func()) *Engine {
	e := &Engine{
		Model:      &model.Model{},
		cmdChan:    make(chan interface{}, 10),
		invalidate: invalidate,
		directions: make([]model.Dir, model.DirCount),
	}

	// Populate the directions slice
	for i := range e.directions {
		e.directions[i] = model.Dir(i)
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

		e.stepGopher()
		e.stepBulldogs()

		e.Model.Unlock()

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

	m.ExitPos.X, m.ExitPos.Y = (m.Cols-2)*BlockSize+BlockSize/2, (m.Rows-2)*BlockSize+BlockSize/2

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

		bd.Pos.X = float64(col*BlockSize + BlockSize/2)
		bd.Pos.Y = float64(row*BlockSize + BlockSize/2)

		bd.TargetPos.X, bd.TargetPos.Y = int(bd.Pos.X), int(bd.Pos.Y)
	}

	// TODO
}

// stepGopher handles moving the Gopher and also handles the multiple target positions of Gopher.
func (e *Engine) stepGopher() {
}

// stepBulldogs iterates over all Bulldogs, generates new random target if they reached their current, and steps them.
func (e *Engine) stepBulldogs() {
	m := e.Model

	// Gopher's position:
	gpos := m.Gopher.Pos

	dirs := e.directions

	for _, bd := range m.Bulldogs {
		x, y := int(bd.Pos.X), int(bd.Pos.Y)

		if bd.TargetPos.X == x && bd.TargetPos.Y == y {
			row, col := y/BlockSize, x/BlockSize
			// Generate new, random target.
			// For this we shuffle all the directions, and check them sequentially.
			// Firts one in which direction there is a free path wins (such path surely exists).

			// Shuffle the directions slice:
			for i := len(dirs) - 1; i > 0; i-- { // last is already random, no use switching with itself
				r := rand.Intn(i + 1)
				dirs[i], dirs[r] = dirs[r], dirs[i]
			}

			var drow, dcol int
			for _, dir := range dirs {
				switch dir {
				case model.DirLeft:
					dcol = -1
				case model.DirRight:
					dcol = 1
				case model.DirUp:
					drow = -1
				case model.DirDown:
					drow = 1
				}
				if m.Lab[row+drow][col+dcol] == model.BlockEmpty {
					// Direction is good, check if we can even step 2 bocks in this way:
					if m.Lab[row+drow*2][col+dcol*2] == model.BlockEmpty {
						drow *= 2
						dcol *= 2
					}
					break
				}
				drow, dcol = 0, 0
			}

			bd.TargetPos.X += dcol * BlockSize
			bd.TargetPos.Y += drow * BlockSize
		}

		e.stepMovingObj(bd)

		if !m.Dead {
			// Check if this Bulldog reached Gopher
			if math.Abs(gpos.X-bd.Pos.X) < BlockSize*0.75 && math.Abs(gpos.Y-bd.Pos.Y) < BlockSize*0.75 {
				m.Dead = true // OK, we just died
			}
		}
	}
}

// stepMovingObj steps the given MovingObj.
func (e *Engine) stepMovingObj(m *model.MovingObj) {
	x, y := int(m.Pos.X), int(m.Pos.Y)

	// Only horizontal or vertical movement is allowed!
	if x != m.TargetPos.X {
		dx := math.Min(dt*v, math.Abs(float64(m.TargetPos.X)-m.Pos.X))
		if x > m.TargetPos.X {
			dx = -dx
			m.Dir = model.DirLeft
		} else {
			m.Dir = model.DirRight
		}
		m.Pos.X += dx
	} else if y != m.TargetPos.Y {
		dy := math.Min(dt*v, math.Abs(float64(m.TargetPos.Y)-m.Pos.Y))
		if y > m.TargetPos.Y {
			dy = -dy
			m.Dir = model.DirUp
		} else {
			m.Dir = model.DirDown
		}
		m.Pos.Y += dy
	}
}
