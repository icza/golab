package engine

import "fmt"

// LabSize difines choosable labyrinth sizes..
type LabSize struct {
	Name    string
	Rows    int // Must be odd
	Cols    int // Must be odd
	Default bool
}

func (l *LabSize) String() string {
	return fmt.Sprintf("%s (%dx%d)", l.Name, l.Rows, l.Cols)
}

// LabSizes is a slice of all, ordered lab sizes.
var LabSizes = []*LabSize{
	&LabSize{Name: "XS", Rows: 9, Cols: 9},
	&LabSize{Name: "S", Rows: 15, Cols: 15},
	&LabSize{Name: "M", Rows: 33, Cols: 33, Default: true},
	&LabSize{Name: "L", Rows: 51, Cols: 51},
	&LabSize{Name: "XL", Rows: 99, Cols: 99},
}

// DefaultLabSizeIdx is the index of the default lab size  in LabSizes.
var DefaultLabSizeIdx int

func init() {
	for i, l := range LabSizes {
		if l.Default {
			DefaultLabSizeIdx = i
			break
		}
	}
}
