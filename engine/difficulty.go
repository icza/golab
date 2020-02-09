package engine

// Difficulty of the game.
type Difficulty struct {
	Name string

	// "Bulldog density", it tells how many Bulldogs to generate for an area of 1,000 blocks.
	// For example if this is 10.0 and rows*cols = 21*21 = 441, 10.0*441/1000 = 4.41 => 4 Bulldogs will be generated.
	BulldogDensity float64

	Default bool
}

func (d *Difficulty) String() string {
	return d.Name
}

// Difficulties is a slice of all, ordered difficulties.
var Difficulties = []*Difficulty{
	&Difficulty{Name: "Baby", BulldogDensity: 0},
	&Difficulty{Name: "Easy", BulldogDensity: 5},
	&Difficulty{Name: "Normal", BulldogDensity: 10, Default: true},
	&Difficulty{Name: "Hard", BulldogDensity: 20},
	&Difficulty{Name: "Brutal", BulldogDensity: 40},
}

// DifficultyDefaultIdx is the index of the default difficulty in Difficulties.
var DifficultyDefaultIdx int

func init() {
	for i, d := range Difficulties {
		if d.Default {
			DifficultyDefaultIdx = i
			break
		}
	}
}
