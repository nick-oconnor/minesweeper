package matrix

import (
	"fmt"

	"gitlab.ocnr.org/apps/minesweeper/field"
)

type Matrix []Row

type Solution struct {
	Flagged, Revealed []*field.Space
}

// NewMatrix creates a new matrix of revealed edge spaces (rows) and unknown edge
// spaces (cols).
func NewMatrix(f *field.Field) Matrix {
	var matrix Matrix
	revealedEdgeSpaces := f.RevealedEdgeSpaces()
	unknownEdgeSpaces := f.UnknownEdgeSpaces()
	for _, revealedEdgeSpace := range revealedEdgeSpaces {
		var row Row
		for _, unknownEdgeSpace := range unknownEdgeSpaces {
			value := 0.0
			if revealedEdgeSpace.UnknownNeighbors()[unknownEdgeSpace] {
				value = 1
			}
			row = append(row, &Cell{value, unknownEdgeSpace, field.Unknown})
		}
		rhs := float64(revealedEdgeSpace.MineNeighborCount() - revealedEdgeSpace.FlaggedNeighborCount())
		row = append(row, &Cell{rhs, revealedEdgeSpace, field.Unknown})
		matrix = append(matrix, row)
	}
	matrix.reduce()
	return matrix.removeUnconstrained()
}

// SplitCoupled creates new matrices from rows that do not share constraints.
func (m Matrix) SplitCoupled() []Matrix {
	var matrices []Matrix
	src := m.copy()
	for len(src) > 0 {
		dst := Matrix{src[0]}
		src = src.removeRow(src[0].Rhs().space)
		src, dst = coupledRows(src, dst)
		dst = dst.removeUnconstrained()
		matrices = append(matrices, dst)
	}
	return matrices
}

// Resolve resolves states determined by constraints, inspired by
// https://massaioli.wordpress.com/2013/01/12/solving-minesweeper-with-matricies/.
func (m Matrix) Resolve() (Matrix, error) {
	matrix := m.copy()
	for _, row := range matrix {
		lhs := row.ConstrainedLhs()
		lhsNegativeSum := 0.0
		lhsPositiveSum := 0.0
		for _, cell := range lhs {
			if cell.value < 0 {
				lhsNegativeSum += cell.value
				continue
			}
			lhsPositiveSum += cell.value
		}
		rhs := row.Rhs().value
		if rhs == lhsNegativeSum {
			for _, cell := range lhs {
				if cell.value < 0 {
					if err := matrix.flag(cell.space); err != nil {
						return matrix, err
					}
				} else {
					if err := matrix.reveal(cell.space); err != nil {
						return matrix, err
					}
				}
			}
			continue
		}
		if rhs == lhsPositiveSum {
			for _, cell := range lhs {
				if cell.value < 0 {
					if err := matrix.reveal(cell.space); err != nil {
						return matrix, err
					}
				} else {
					if err := matrix.flag(cell.space); err != nil {
						return matrix, err
					}
				}
			}
			continue
		}
	}
	return matrix, nil
}

// Solve returns a list of all possible solutions for the matrix.
func (m Matrix) Solve(visualize bool) []*Solution {
	space := m.mostConstrainedSpace()
	if space == nil {
		solution := &Solution{[]*field.Space{}, []*field.Space{}}
		for _, cell := range m[0].Lhs() {
			switch cell.state {
			case field.Flagged:
				solution.Flagged = append(solution.Flagged, cell.space)
			case field.Revealed:
				solution.Revealed = append(solution.Revealed, cell.space)
			case field.Unknown:
				// Ignore
			}
		}
		if visualize {
			m.Print()
		}
		return []*Solution{solution}
	}
	var solutions []*Solution
	reveal := m.copy()
	if err := reveal.reveal(space); err == nil {
		if reveal, err = reveal.Resolve(); err == nil {
			solutions = append(solutions, reveal.Solve(visualize)...)
		}
	}
	flag := m.copy()
	if err := flag.flag(space); err == nil {
		if flag, err = flag.Resolve(); err == nil {
			solutions = append(solutions, flag.Solve(visualize)...)
		}
	}
	return solutions
}

// ConstrainedSpaces returns a map of the constrained spaces.
func (m Matrix) ConstrainedSpaces() map[*field.Space]bool {
	spaces := make(map[*field.Space]bool)
	for _, row := range m {
		for _, cell := range row.ConstrainedLhs() {
			spaces[cell.space] = true
		}
	}
	return spaces
}

// Print prints the matrix to stdout.
func (m Matrix) Print() {
	if len(m) == 0 {
		fmt.Println("(empty matrix)")
	}
	for rowIndex, row := range m {
		if rowIndex == 0 {
			fmt.Print(" ")
			for _, cell := range row.Lhs() {
				fmt.Printf("%4d ", cell.space.Index())
			}
			fmt.Println()
			for _, cell := range row.Lhs() {
				switch cell.state {
				case field.Unknown:
					fmt.Print("    U")
				case field.Flagged:
					fmt.Print("    F")
				case field.Revealed:
					fmt.Print("    R")
				}
			}
			fmt.Println()
		}
		for cellIndex, cell := range row {
			if cellIndex%len(row) == len(row)-1 {
				fmt.Printf(" | %4s %d", fmt.Sprintf("%1.1f", cell.value), cell.space.Index())
			} else {
				fmt.Printf(" %4s", fmt.Sprintf("%1.1f", cell.value))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// mostConstrainedSpace returns the space with the most constraints.
func (m Matrix) mostConstrainedSpace() *field.Space {
	var mostConstrainedSpace *field.Space
	mostConstrainedCount := 0
	constraintCounts := make(map[*field.Space]int)
	for _, row := range m {
		for _, cell := range row.ConstrainedLhs() {
			constraintCounts[cell.space]++
			if constraintCounts[cell.space] > mostConstrainedCount {
				mostConstrainedSpace = cell.space
				mostConstrainedCount = constraintCounts[cell.space]
			}
		}
	}
	return mostConstrainedSpace
}

// unconstrainedSpaces returns a map of the unconstrained spaces.
func (m Matrix) unconstrainedSpaces() map[*field.Space]bool {
	spaces := make(map[*field.Space]bool)
	for _, cell := range m[0].Lhs() {
		spaces[cell.space] = true
	}
	for _, row := range m {
		for _, cell := range row.ConstrainedLhs() {
			delete(spaces, cell.space)
		}
	}
	return spaces
}

// zeroRows returns all rows with all-zero left-hand-sides.
func (m Matrix) zeroRows() map[*field.Space]bool {
	zeroRows := make(map[*field.Space]bool)
	for _, row := range m {
		if len(row.ConstrainedLhs()) == 0 {
			zeroRows[row.Rhs().space] = true
		}
	}
	return zeroRows
}

// coupledRows returns a matrix consisting of rows from source which are coupled
// to destination.
func coupledRows(src Matrix, dst Matrix) (Matrix, Matrix) {
	if len(src) == 0 {
		return src, dst
	}
	dstConstrainedSpaces := dst.ConstrainedSpaces()
	for _, srcRow := range src {
		for _, srcCell := range srcRow.ConstrainedLhs() {
			if dstConstrainedSpaces[srcCell.space] {
				dst = append(dst, srcRow)
				src = src.removeRow(srcRow.Rhs().space)
				return coupledRows(src, dst)
			}
		}
	}
	return src, dst
}

// reduce transforms the matrix to RREF. This algorithm is pilfered from
// https://rosettacode.org/wiki/Reduced_row_echelon_form.
func (m Matrix) reduce() {
	if len(m) == 0 {
		return
	}
	lead := 0
	rowCount := len(m)
	colCount := len(m[0])
	for r := range m {
		if lead >= colCount {
			return
		}
		i := r
		for m[i][lead].value == 0 {
			i++
			if i == rowCount {
				i = r
				lead++
				if lead == colCount {
					return
				}
			}
		}
		m[i], m[r] = m[r], m[i]
		div := m[r][lead].value
		for j := range m[r] {
			m[r][j].value /= div
		}
		for j := range m {
			if j != r {
				sub := m[j][lead].value
				for k := range m[r] {
					m[j][k].value -= sub * m[r][k].value
				}
			}
		}
		lead++
	}
}

// reveal reveals the given space.
func (m Matrix) reveal(space *field.Space) error {
	for _, row := range m {
		for _, cell := range row.Lhs() {
			if cell.space == space {
				cell.value = 0
				cell.state = field.Revealed
			}
		}
		if err := row.validate(); err != nil {
			return err
		}
	}
	return nil
}

// flag flags the given space.
func (m Matrix) flag(space *field.Space) error {
	for _, row := range m {
		for _, cell := range row.Lhs() {
			if cell.space == space {
				row.Rhs().value -= cell.value
				cell.value = 0
				cell.state = field.Flagged
			}
		}
		if err := row.validate(); err != nil {
			return err
		}
	}
	return nil
}

// removeUnconstrained removes rows with all zero left-hand-sides and columns
// with unconstrained spaces. This is used to clean up the output and is not
// algorithmically necessary.
func (m Matrix) removeUnconstrained() Matrix {
	if len(m) == 0 {
		return m
	}
	return m.removeRows(m.zeroRows()).removeCols(m.unconstrainedSpaces())
}

// copy returns a deep copy of the matrix.
func (m Matrix) copy() Matrix {
	var matrix Matrix
	for _, row := range m {
		var newRow Row
		for _, cell := range row {
			newRow = append(newRow, cell.copy())
		}
		matrix = append(matrix, newRow)
	}
	return matrix
}

// removeRow removes the given row from the matrix.
func (m Matrix) removeRow(row *field.Space) Matrix {
	return m.removeRows(map[*field.Space]bool{row: true})
}

// removeRows removes the given rows from the matrix.
func (m Matrix) removeRows(rows map[*field.Space]bool) Matrix {
	var matrix Matrix
	for _, row := range m {
		if !rows[row.Rhs().space] {
			matrix = append(matrix, row)
		}
	}
	return matrix
}

// removeCol removes the given spaces from the matrix.
func (m Matrix) removeCols(cols map[*field.Space]bool) Matrix {
	var matrix Matrix
	for _, row := range m {
		var newRow Row
		for _, cell := range row {
			if !cols[cell.space] {
				newRow = append(newRow, cell)
			}
		}
		matrix = append(matrix, newRow)
	}
	return matrix
}
