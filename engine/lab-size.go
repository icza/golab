package engine

import "fmt"

// LabSize difines choosable labyrinth sizes..
type LabSize struct {
	Name    string
	rows    int // Must be odd
	cols    int // Must be odd
	Default bool
}

func (l *LabSize) String() string {
	return fmt.Sprintf("%s (%dx%d)", l.Name, l.rows, l.cols)
}

// LabSizes is a slice of all, ordered lab sizes.
var LabSizes = []*LabSize{
	&LabSize{Name: "XS", rows: 9, cols: 9},
	&LabSize{Name: "S", rows: 15, cols: 15},
	&LabSize{Name: "M", rows: 33, cols: 33, Default: true},
	&LabSize{Name: "L", rows: 51, cols: 51},
	&LabSize{Name: "XL", rows: 99, cols: 99},
}

// LabSizeDefaultIdx is the index of the default lab size  in LabSizes.
var LabSizeDefaultIdx int

func init() {
	for i, l := range LabSizes {
		if l.Default {
			LabSizeDefaultIdx = i
			break
		}
	}
}
