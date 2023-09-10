package matrix

import (
	"fmt"
)

type Row []*Cell

// Lhs returns the left-hand-side of the matrix.
func (r Row) Lhs() Row {
	return r[:len(r)-1]
}

// Rhs returns the right-hand-side of the matrix.
func (r Row) Rhs() *Cell {
	return r[len(r)-1]
}

// ConstrainedLhs returns the non-zero left-hand-side cells of the matrix.
func (r Row) ConstrainedLhs() Row {
	var row Row
	for _, cell := range r.Lhs() {
		if cell.value != 0 {
			row = append(row, cell)
		}
	}
	return row
}

// validate validates that the left-hand-cells in the row are consistent with the
// right-hand-side.
func (r Row) validate() error {
	lhsNegativeSum := 0.0
	lhsPositiveSum := 0.0
	for _, cell := range r.ConstrainedLhs() {
		if cell.value < 0 {
			lhsNegativeSum += cell.value
			continue
		}
		lhsPositiveSum += cell.value
	}
	if r.Rhs().value > lhsPositiveSum || r.Rhs().value < lhsNegativeSum {
		return fmt.Errorf("invalid constraint for row %d", r.Rhs().space.Index())
	}
	return nil
}
