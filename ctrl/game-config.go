package ctrl

// GameConfig holds config to start a new game.
type GameConfig struct {
	Difficulty *Difficulty
	LabSize    *LabSize
	Speed      *Speed
}
