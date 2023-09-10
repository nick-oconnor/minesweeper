package solver

import (
	"fmt"

	"gitlab.ocnr.org/apps/minesweeper/field"
)

type MoveType int

const (
	Constrained        MoveType = iota
	Enumeration        MoveType = iota
	EnumerationGuess   MoveType = iota
	UnconstrainedGuess MoveType = iota
)

// String returns a text version of the MoveType enum.
func (s MoveType) String() string {
	switch s {
	case Constrained:
		return "constraints"
	case Enumeration:
		return "enumeration"
	case EnumerationGuess:
		return "enumerated guess"
	case UnconstrainedGuess:
		return "unconstrained guess"
	}
	panic(fmt.Sprintf("invalid state %d", s))
}

type MoveInfo struct {
	operation field.State
	moveType  MoveType
}
