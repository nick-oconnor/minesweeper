package solver

import (
	"fmt"

	. "gitlab.ocnr.org/apps/minesweeper/field"
)

// Matrix represents a 2D matrix.
type Matrix [][]*Cell

// Cell represents a cell in the Matrix.
type Cell struct {
	value float64
	space *Space
}

// newMatrix creates a new matrix of revealed edge spaces (rows) and unknown edge
// spaces (cols). Only unknown spaces with non-zero probabilities are added.
func newMatrix(revealedEdgeSpaces, unknownEdgeSpaces []*Space) Matrix {
	var m Matrix
	for _, revealedEdgeSpace := range revealedEdgeSpaces {
		var row []*Cell
		for _, unknownEdgeSpace := range unknownEdgeSpaces {
			possibleMine := 0.0
			if revealedEdgeSpace.HasNeighbor(unknownEdgeSpace) {
				possibleMine = 1
			}
			row = append(row, &Cell{possibleMine, unknownEdgeSpace})
		}
		rhs := float64(revealedEdgeSpace.MineNeighborCount() - revealedEdgeSpace.FlaggedNeighborCount())
		if rhs == 0 {
			panic(fmt.Sprintf("space %v already solved", revealedEdgeSpace.Id()))
		}
		row = append(row, &Cell{rhs, revealedEdgeSpace})
		m = append(m, row)
	}
	return m
}

// reduce transforms the matrix to reduced-row echelon form. This algorithm is
// pilfered from https://rosettacode.org/wiki/Reduced_row_echelon_form.
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

// addMoves adds moves determined by the matrix to the provided move queue,
// inspired by
// https://massaioli.wordpress.com/2013/01/12/solving-minesweeper-with-matricies/.
func (m Matrix) addMoves(moveQueue MoveQueue) {
	for _, r := range m {
		lhs := make(map[float64][]*Space)
		rhs := r[len(r)-1].value
		for _, c := range r[:len(r)-1] {
			if c.value != 0 {
				lhs[c.value] = append(lhs[c.value], c.space)
			}
		}
		lhsNegativeSum := 0.0
		lhsPositiveSum := 0.0
		for value, spaces := range lhs {
			if value < 0 {
				lhsNegativeSum += value * float64(len(spaces))
			} else {
				lhsPositiveSum += value * float64(len(spaces))
			}
		}
		if rhs == lhsNegativeSum {
			for value, spaces := range lhs {
				if value < 0 {
					moveQueue.add(spaces, Flagged)
				} else {
					moveQueue.add(spaces, Revealed)
				}
			}
			continue
		}
		if rhs == lhsPositiveSum {
			for value, spaces := range lhs {
				if value < 0 {
					moveQueue.add(spaces, Revealed)
				} else {
					moveQueue.add(spaces, Flagged)
				}
			}
		}
	}
}

// add adds the given spaces with the given operation to the move queue.
func (m MoveQueue) add(spaces []*Space, operation State) {
	for _, space := range spaces {
		m[space] = &MoveInfo{operation, false}
	}
}
