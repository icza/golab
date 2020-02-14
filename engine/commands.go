package engine

// GameConfig holds config to start a new game.
type GameConfig struct {
	Difficulty *Difficulty
	LabSize    *LabSize
	Speed      *Speed
}

// Click describes a click event.
type Click struct {
	X, Y  int  // Click coordinates in the lab
	Left  bool // Tells if left button was pressed
	Right bool // Tells if right button was pressed
}

// Key describes a key event.
type Key struct {
	DirKeys map[Dir]bool // Tells if keys for the directions were pressed
}
