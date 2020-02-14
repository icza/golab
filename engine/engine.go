// Package engine is the game engien: it contains the game model and game logic.
//
// The engine's Loop() method should be launched as a goroutine,
// and it can be controlled with opaque commands safely from other
// goroutines.
package engine

import (
	"image"
	"log"
	"math"
	"math/rand"
	"time"
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
	Model *Model

	// command channel to control the engine from other goroutines.
	cmdChan chan interface{}

	// invalidate is called by the engine to request a new view frame.
	invalidate func()

	// Current game config
	cfg *GameConfig

	// directions is a reused slice of all directions
	directions []Dir
}

// NewEngine returns a new Engine.
// invalidate is a func which will be called by the engine to request a new view frame.
func NewEngine(invalidate func()) *Engine {
	e := &Engine{
		Model: &Model{
			TargetPoss: make([]image.Point, 0, 20), // cap defines max queueable points
		},
		cmdChan:    make(chan interface{}, 10),
		invalidate: invalidate,
		directions: make([]Dir, DirCount),
	}

	// Populate the directions slice
	for i := range e.directions {
		e.directions[i] = Dir(i)
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

// SendClick sends a click event from the user.
func (e *Engine) SendClick(c Click) {
	e.cmdChan <- &c
}

// SendKey sends a key event from the user.
func (e *Engine) SendKey(k Key) {
	e.cmdChan <- &k
}

// Loop starts calculating the game.
// This function returns only if the user closes the app.
func (e *Engine) Loop() {
	for {
		e.Model.Lock()

		e.processCmds()

		if !e.Model.Won {
			e.stepGopher()
			e.stepBulldogs()
		}

		e.Model.Unlock()

		e.invalidate()

		time.Sleep(e.cfg.Speed.loopDelay)
	}
}

// processCmds processes queued commands.
func (e *Engine) processCmds() {
	for {
		select {

		case cmd := <-e.cmdChan:
			switch cmd := cmd.(type) {
			case *GameConfig:
				e.initNewGame(cmd)
			case *Click:
				e.handleClick(cmd)
			case *Key:
				e.handleKey(cmd)
			default:
				log.Printf("Unhandled cmd type: %T", cmd)
			}

		default:
			return // No more commands queued
		}
	}
}

// handleClick handles a Click command
func (e *Engine) handleClick(c *Click) {
	m := e.Model

	if m.Dead || m.Won {
		return
	}

	if c.Right {
		m.TargetPoss = m.TargetPoss[:0]
		// Also change Gopher's current target to remain on current block:
		m.Gopher.TargetPos.X = int(m.Gopher.Pos.X)/BlockSize*BlockSize + BlockSize/2
		m.Gopher.TargetPos.Y = int(m.Gopher.Pos.Y)/BlockSize*BlockSize + BlockSize/2
		return
	}

	// If target buffer is full, do nothing:
	if len(m.TargetPoss) == cap(m.TargetPoss) {
		return
	}

	// Last target pos:
	var TargetPos image.Point
	if len(m.TargetPoss) == 0 {
		TargetPos = m.Gopher.TargetPos
	} else {
		TargetPos = m.TargetPoss[len(m.TargetPoss)-1]
	}

	// Check if new desired target is in the same row/column as the last target and if there is a free passage to there.
	pCol, pRow := TargetPos.X/BlockSize, TargetPos.Y/BlockSize
	tCol, tRow := c.X/BlockSize, c.Y/BlockSize

	// sorted simply returns its parameters in ascendant order:
	sorted := func(a, b int) (int, int) {
		if a < b {
			return a, b
		}
		return b, a
	}

	if pCol == tCol { // Same column
		for row, row2 := sorted(pRow, tRow); row <= row2; row++ {
			if m.Lab[row][tCol] == BlockWall {
				return // Wall in the route
			}
		}
	} else if pRow == tRow { // Same row
		for col, col2 := sorted(pCol, tCol); col <= col2; col++ {
			if m.Lab[tRow][col] == BlockWall {
				return // Wall in the route
			}
		}
	} else {
		return // Only the same row or column can be commanded
	}

	// Target pos is allowed and reachable.
	// Use target position rounded to the center of the target block:
	m.TargetPoss = append(m.TargetPoss, image.Pt(tCol*BlockSize+BlockSize/2, tRow*BlockSize+BlockSize/2))
}

// handleKey handles a Key command.
func (e *Engine) handleKey(k *Key) {
	m := e.Model

	if m.Dead || m.Won {
		return
	}

	Gopher := m.Gopher

	for dir := Dir(0); dir < DirCount; dir++ {
		if !k.DirKeys[dir] {
			continue
		}
		Gopher.Dir = dir
		// If Gopher's target is more than a block away, clear that target:
		dx, dy := int(Gopher.Pos.X)-Gopher.TargetPos.X, int(Gopher.Pos.Y)-Gopher.TargetPos.Y
		if dx <= -BlockSize || dx >= BlockSize || dy <= -BlockSize || dy >= BlockSize {
			m.TargetPoss = m.TargetPoss[:0]
			m.Gopher.TargetPos.X = int(m.Gopher.Pos.X)/BlockSize*BlockSize + BlockSize/2
			m.Gopher.TargetPos.Y = int(m.Gopher.Pos.Y)/BlockSize*BlockSize + BlockSize/2
		}

		col, row := Gopher.TargetPos.X/BlockSize, Gopher.TargetPos.Y/BlockSize

		var drow, dcol int
		switch dir {
		case DirLeft:
			dcol = -1
		case DirRight:
			dcol = 1
		case DirUp:
			drow = -1
		case DirDown:
			drow = 1
		}

		if m.Lab[row+drow][col+dcol] == BlockEmpty {
			m.TargetPoss = m.TargetPoss[:0]
			m.TargetPoss = append(m.TargetPoss, image.Point{
				X: (col+dcol)*BlockSize + BlockSize/2,
				Y: (row+drow)*BlockSize + BlockSize/2},
			)
			break
		}
	}
}

// initNewGame handles a GameConfig command: initializes a new game.
func (e *Engine) initNewGame(cfg *GameConfig) {
	e.cfg = cfg

	m := e.Model

	m.Counter++

	// Init the labyrinth
	m.Rows, m.Cols = cfg.LabSize.rows, cfg.LabSize.cols
	m.Lab = make([][]Block, m.Rows)
	for row := range m.Lab {
		m.Lab[row] = make([]Block, m.Cols)
	}
	generateLab(m.Lab)

	m.ExitPos.X, m.ExitPos.Y = (m.Cols-2)*BlockSize+BlockSize/2, (m.Rows-2)*BlockSize+BlockSize/2

	// Init Gopher
	m.Gopher = new(MovingObj)
	m.Gopher.Pos.X = BlockSize + BlockSize/2 // Position Gopher to top left corner
	m.Gopher.Pos.Y = m.Gopher.Pos.X
	m.Gopher.Dir = DirRight
	m.Gopher.TargetPos.X = int(m.Gopher.Pos.X)
	m.Gopher.TargetPos.Y = int(m.Gopher.Pos.Y)

	// Init bulldogs
	numBulldogs := int(float64(m.Rows*m.Cols) * cfg.Difficulty.bulldogDensity / 1000)
	m.Bulldogs = make([]*MovingObj, numBulldogs)
	for i := range m.Bulldogs {
		bd := new(MovingObj)
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

	m.Dead = false
	m.Won = false

	// Throw away queued targets
	m.TargetPoss = m.TargetPoss[:0]
}

// stepGopher handles moving the Gopher and also handles the multiple target positions of Gopher.
func (e *Engine) stepGopher() {
	m := e.Model
	Gopher := m.Gopher

	if m.Dead {
		return // Dead Gopher can't move
	}

	// Check if reached current target position:
	if int(Gopher.Pos.X) == Gopher.TargetPos.X && int(Gopher.Pos.Y) == Gopher.TargetPos.Y {
		// Check if we have more target positions in our path:
		if len(m.TargetPoss) > 0 {
			// Set the next target as the current
			Gopher.TargetPos = m.TargetPoss[0]
			// and remove it from the targets:
			m.TargetPoss = m.TargetPoss[:copy(m.TargetPoss, m.TargetPoss[1:])]
		}
	}

	// Step Gopher
	Gopher.step()

	// Check if Gopher reached the exit point
	if int(m.Gopher.Pos.X) == m.ExitPos.X && int(m.Gopher.Pos.Y) == m.ExitPos.Y {
		m.Won = true
	}
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
				case DirLeft:
					dcol = -1
				case DirRight:
					dcol = 1
				case DirUp:
					drow = -1
				case DirDown:
					drow = 1
				}
				if m.Lab[row+drow][col+dcol] == BlockEmpty {
					// Direction is good, check if we can even step 2 bocks in this way:
					if m.Lab[row+drow*2][col+dcol*2] == BlockEmpty {
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

		bd.step()

		if !m.Dead {
			// Check if this Bulldog reached Gopher (but only if not just won)
			if math.Abs(gpos.X-bd.Pos.X) < BlockSize*0.75 && math.Abs(gpos.Y-bd.Pos.Y) < BlockSize*0.75 && !m.Won {
				m.Dead = true // OK, we just died
			}
		}
	}
}
