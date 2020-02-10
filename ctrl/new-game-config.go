package ctrl

// NewGameConfig holds config to start a new game.
type NewGameConfig struct {
	Difficulty *Difficulty
	LabSize    *LabSize
	Speed      *Speed
}
