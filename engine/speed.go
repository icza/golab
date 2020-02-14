package engine

import (
	"fmt"
	"time"
)

// Speed of the game.
type Speed struct {
	Name string

	loopDelay time.Duration

	Default bool
}

func (s *Speed) String() string {
	fps := (time.Second + s.loopDelay/2) / s.loopDelay // add half to make it round up from half
	return fmt.Sprintf("%s (%d FPS)", s.Name, fps)
}

// Speeds is a slice of all, ordered speeds.
var Speeds = []*Speed{
	&Speed{Name: "Slow", loopDelay: 67 * time.Millisecond},                  // ~15 FPS
	&Speed{Name: "Normal", loopDelay: 50 * time.Millisecond, Default: true}, // ~20 FPS
	&Speed{Name: "Fast", loopDelay: 37 * time.Millisecond},                  // ~27 FPS
}

// SpeedDefaultIdx is the index of the default speed in Speeds.
var SpeedDefaultIdx int

func init() {
	for i, s := range Speeds {
		if s.Default {
			SpeedDefaultIdx = i
			break
		}
	}
}
