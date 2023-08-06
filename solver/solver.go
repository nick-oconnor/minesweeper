package solver

import (
	"fmt"
	. "gitlab.ocnr.org/apps/minesweeper/field"
)

// MoveInfo holds queued move info.
type MoveInfo struct {
	operation State
	isGuess   bool
}

// MoveQueue holds queued moves.
type MoveQueue map[*Space]*MoveInfo

// Solver is used for solving fields.
type Solver struct {
	field          *Field
	enableGuessing bool
	fieldPrinter   func(*Field, *Space, error)
	matrixPrinter  func(Matrix)
	queuedMoves    MoveQueue
}

// NewSolver creates a new solver.
func NewSolver(field *Field) *Solver {
	return &Solver{field, true, nil, nil, make(MoveQueue)}
}

// WithGuessing enables or disables guessing for the solver.
func (s *Solver) WithGuessing(enableGuessing bool) *Solver {
	s.enableGuessing = enableGuessing
	return s
}

// WithFieldPrinter sets the field printer to use when solving.
func (s *Solver) WithFieldPrinter(fieldPrinter func(*Field, *Space, error)) *Solver {
	s.fieldPrinter = fieldPrinter
	return s
}

// WithMatrixPrinter sets the matrix printer to use when solving.
func (s *Solver) WithMatrixPrinter(matrixPrinter func(Matrix)) *Solver {
	s.matrixPrinter = matrixPrinter
	return s
}

// Solve solves the field.
func (s *Solver) Solve() (int, int, error) {
	moveCount := 0
	guessCount := 0
	for {
		var space *Space
		var isGuess bool
		var err error
		for {
			space, isGuess, err = s.nextMove()
			if err != nil || space != nil {
				moveCount++
				if isGuess {
					guessCount++
				}
				break
			}
			s.generateSingleSpaceMoves()
			if len(s.queuedMoves) > 0 {
				continue
			}
			s.generateMultiSpaceMoves()
			if len(s.queuedMoves) > 0 {
				continue
			}
			if s.enableGuessing {
				s.generateGuess()
			}
			if len(s.queuedMoves) == 0 {
				break
			}
		}
		if s.fieldPrinter != nil {
			s.fieldPrinter(s.field, space, err)
		}
		if err != nil || space == nil || len(s.field.UnknownSpaces()) == 0 {
			return moveCount, guessCount, err
		}
	}
}

// nextMove executes the next move from the move queue.
func (s *Solver) nextMove() (*Space, bool, error) {
	for space, queuedMove := range s.queuedMoves {
		delete(s.queuedMoves, space)
		if space.State() == queuedMove.operation {
			continue
		}
		var err error
		switch queuedMove.operation {
		case Flagged:
			err = s.field.Flag(space, !queuedMove.isGuess)
		case Revealed:
			err = s.field.Reveal(space, !queuedMove.isGuess)
		default:
			panic(fmt.Sprintf("invalid move %v for space %v", queuedMove.operation, space.Id()))
		}
		return space, queuedMove.isGuess, err
	}
	return nil, false, nil
}

// generateSingleSpaceMoves adds single-space moves to the move queue.
func (s *Solver) generateSingleSpaceMoves() {
	for _, revealedEdgeSpace := range s.revealedEdgeSpaces() {
		probability := float64(revealedEdgeSpace.MineNeighborCount()-revealedEdgeSpace.FlaggedNeighborCount()) / float64(len(revealedEdgeSpace.UnknownNeighbors()))
		for _, unknownEdgeSpace := range revealedEdgeSpace.UnknownNeighbors() {
			if probability == 1 {
				s.queuedMoves[unknownEdgeSpace] = &MoveInfo{Flagged, false}
			}
			if probability == 0 {
				s.queuedMoves[unknownEdgeSpace] = &MoveInfo{Revealed, false}
			}
		}
	}
}

// generateMultiSpaceMoves adds multi-space moves to the move queue.
func (s *Solver) generateMultiSpaceMoves() {
	matrix := newMatrix(s.revealedEdgeSpaces(), s.unknownEdgeSpaces())
	if s.matrixPrinter != nil {
		s.matrixPrinter(matrix)
	}
	matrix.reduce()
	if s.matrixPrinter != nil {
		s.matrixPrinter(matrix)
	}
	matrix.addMoves(s.queuedMoves)
}

// generateGuess adds a guess to the move queue. There's room for improvement
// here. Per
// https://dash.harvard.edu/bitstream/handle/1/14398552/BECERRA-SENIORTHESIS-2015.pdf
// it should be possible to obtain a 32% win rate.
func (s *Solver) generateGuess() {
	var cornerSpace, edgeSpace, centerSpace *Space
	for _, unknownSpace := range s.field.UnknownSpaces() {
		if len(unknownSpace.Neighbors()) == 3 {
			cornerSpace = unknownSpace
		}
		if len(unknownSpace.Neighbors()) == 5 {
			edgeSpace = unknownSpace
		}
		centerSpace = unknownSpace
	}
	if cornerSpace != nil {
		s.queuedMoves[cornerSpace] = &MoveInfo{Revealed, true}
		return
	}
	if edgeSpace != nil {
		s.queuedMoves[edgeSpace] = &MoveInfo{Revealed, true}
		return
	}
	if centerSpace != nil {
		s.queuedMoves[centerSpace] = &MoveInfo{Revealed, true}
	}
}

// revealedEdgeSpaces retrieves the revealed edge spaces for the field.
func (s *Solver) revealedEdgeSpaces() []*Space {
	var spaces []*Space
	for _, revealedSpace := range s.field.RevealedSpaces() {
		if len(revealedSpace.UnknownNeighbors()) > 0 {
			spaces = append(spaces, revealedSpace)
		}
	}
	return spaces
}

// unknownEdgeSpaces retrieves the unknown edge spaces for the field.
func (s *Solver) unknownEdgeSpaces() []*Space {
	var spaces []*Space
	for _, unknownSpace := range s.field.UnknownSpaces() {
		if unknownSpace.RevealedNeighborCount() > 0 {
			spaces = append(spaces, unknownSpace)
		}
	}
	return spaces
}
