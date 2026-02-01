package solver

import (
	"fmt"

	"gitlab.ocnr.org/apps/minesweeper/field"
	"gitlab.ocnr.org/apps/minesweeper/matrix"
)

type MoveQueue map[*field.Space]*MoveInfo

type GameResult struct {
	Won                   bool
	MoveCount, GuessCount int
}

type Solver struct {
	field     *field.Field
	visualize bool
	moveQueue MoveQueue
}

var guessTypes = map[MoveType]bool{EnumerationGuess: true, UnconstrainedGuess: true}

// NewSolver creates a new solver.
func NewSolver(f *field.Field, visualize bool) *Solver {
	return &Solver{f, visualize, make(MoveQueue)}
}

// Solve solves the field.
func (s *Solver) Solve() *GameResult {
	gameResult := &GameResult{}
	for {
		space, moveType, err := s.findAndExecuteMove()
		if space != nil {
			gameResult.MoveCount++
			if guessTypes[moveType] {
				gameResult.GuessCount++
			}
		}
		if err != nil || space == nil || len(s.field.UnknownSpaces()) == 0 {
			gameResult.Won = len(s.field.UnknownSpaces()) == 0 && err == nil
			return gameResult
		}
	}
}

// findAndExecuteMove finds the next move to make and executes it.
func (s *Solver) findAndExecuteMove() (*field.Space, MoveType, error) {
	for {
		space, moveType, err := s.nextMove()
		if err != nil || space != nil {
			return space, moveType, err
		}
		baseMatrix := matrix.NewMatrix(s.field)
		s.addConstrainedMoves(baseMatrix)
		if len(s.moveQueue) > 0 {
			continue
		}
		if s.visualize {
			fmt.Printf("no moves found by resolving constraints\n\n")
		}
		s.addEnumeratedMoves(baseMatrix)
		if len(s.moveQueue) > 0 {
			continue
		}
		if s.visualize {
			fmt.Printf("no moves found by enumeration\n\n")
		}
		s.addUnconstrainedMove(baseMatrix)
		if len(s.moveQueue) == 0 {
			panic("no moves added")
		}
	}
}

// nextMove executes the next move from the move queue.
func (s *Solver) nextMove() (*field.Space, MoveType, error) {
	for space, moveInfo := range s.moveQueue {
		delete(s.moveQueue, space)
		if space.State() == moveInfo.operation {
			continue
		}
		var err error
		switch moveInfo.operation {
		case field.Flagged:
			s.field.Flag(space)
		case field.Revealed:
			err = s.field.Reveal(space)
			if err != nil && !guessTypes[moveInfo.moveType] {
				panic(err)
			}
		default:
			panic(fmt.Sprintf("invalid move %s for space %d", moveInfo.operation, space.Index()))
		}
		if s.visualize {
			fmt.Printf("space %d %s by %s\n\n", space.Index(), moveInfo.operation, moveInfo.moveType)
			s.field.Print()
			if err != nil {
				fmt.Println(err)
			}
		}
		return space, moveInfo.moveType, err
	}
	return nil, 0, nil
}

// addConstrainedMoves adds fully constrained moves to the move queue.
func (s *Solver) addConstrainedMoves(baseMatrix matrix.Matrix) {
	if len(baseMatrix) == 0 {
		return
	}
	resolvedMatrix, err := baseMatrix.Resolve()
	if err != nil {
		panic(err)
	}
	if s.visualize {
		fmt.Printf("resolving constraints\n\n")
		resolvedMatrix.Print()
	}
	if len(resolvedMatrix) == 0 {
		return
	}
	for _, cell := range resolvedMatrix[0].Lhs() {
		switch cell.State() {
		case field.Flagged:
			s.moveQueue[cell.Space()] = &MoveInfo{field.Flagged, Constrained}
		case field.Revealed:
			s.moveQueue[cell.Space()] = &MoveInfo{field.Revealed, Constrained}
		case field.Unknown:
			// Ignore
		}
	}
}

// addEnumeratedMoves calculates all possible solutions to the currently
// constrained spaces. It flags or reveals spaces which are mines or not mines in
// every possible solution. If no spaces are mines or not mines in every
// solution, it calculates the space with the highest probability of not being a
// mine. If that probability is less than that of unconstrained spaces, it
// reveals it. This algorithm is based on
// https://www.cs.toronto.edu/~cvs/minesweeper/minesweeper.pdf.
func (s *Solver) addEnumeratedMoves(baseMatrix matrix.Matrix) {
	bestSpace, bestProbability := s.findBestMove(baseMatrix)
	if len(s.moveQueue) > 0 {
		return
	}
	if bestSpace != nil {
		s.addEnumerationGuess(bestSpace, bestProbability)
	}
}

// findBestMove finds the space with the highest probability of not containing a mine.
func (s *Solver) findBestMove(baseMatrix matrix.Matrix) (*field.Space, float64) {
	var bestSpace *field.Space
	bestProbability := 0.0
	for space, probability := range s.probabilityPerSpace(baseMatrix) {
		switch probability {
		case 0:
			s.moveQueue[space] = &MoveInfo{field.Flagged, Enumeration}
		case 1:
			s.moveQueue[space] = &MoveInfo{field.Revealed, Enumeration}
		}
		if probability > bestProbability {
			bestSpace = space
			bestProbability = probability
		}
	}
	return bestSpace, bestProbability
}

// addEnumerationGuess adds an enumeration guess if the probability is better than unconstrained spaces.
func (s *Solver) addEnumerationGuess(space *field.Space, probability float64) {
	minesRemaining := s.field.MineCount() - s.field.FlaggedCount()
	fieldProbability := 1 - float64(minesRemaining)/float64(len(s.field.UnknownSpaces()))
	if s.visualize {
		fmt.Printf("unconstrained mine-free probability\n\n %.2f\n\n", fieldProbability)
	}
	if probability > fieldProbability {
		s.moveQueue[space] = &MoveInfo{field.Revealed, EnumerationGuess}
	}
}

// addUnconstrainedMove reveals an unconstrained corner, edge, then center space
// in preferential order.
func (s *Solver) addUnconstrainedMove(baseMatrix matrix.Matrix) {
	var cornerSpace, edgeSpace, centerSpace *field.Space
	constrainedSpaces := baseMatrix.ConstrainedSpaces()
	unconstrainedSpacesExist := len(s.field.UnknownSpaces()) > len(constrainedSpaces)
	for unknownSpace := range s.field.UnknownSpaces() {
		if unconstrainedSpacesExist && constrainedSpaces[unknownSpace] {
			continue
		}
		if len(unknownSpace.Neighbors()) == 3 {
			cornerSpace = unknownSpace
		}
		if len(unknownSpace.Neighbors()) == 5 {
			edgeSpace = unknownSpace
		}
		centerSpace = unknownSpace
	}
	if cornerSpace != nil {
		s.moveQueue[cornerSpace] = &MoveInfo{field.Revealed, UnconstrainedGuess}
		return
	}
	if edgeSpace != nil {
		s.moveQueue[edgeSpace] = &MoveInfo{field.Revealed, UnconstrainedGuess}
		return
	}
	if centerSpace != nil {
		s.moveQueue[centerSpace] = &MoveInfo{field.Revealed, UnconstrainedGuess}
	}
}

// probabilityPerSpace returns the probability of each edge space not containing
// a mine.
func (s *Solver) probabilityPerSpace(baseMatrix matrix.Matrix) map[*field.Space]float64 {
	var solutionsByPart [][]*matrix.Solution
	for partIndex, part := range baseMatrix.SplitCoupled() {
		if s.visualize {
			fmt.Printf("possible solutions for constraint group %d\n\n", partIndex)
		}
		solutionsByPart = append(solutionsByPart, part.Solve(s.visualize))
	}
	minFlaggedSum := 0
	minFlaggedByPart := make(map[int]int)
	for partIndex, solutions := range solutionsByPart {
		minFlagged := 0
		for _, solution := range solutions {
			if len(solution.Flagged) < minFlagged {
				minFlagged = len(solution.Flagged)
			}
		}
		minFlaggedSum += minFlagged
		minFlaggedByPart[partIndex] = minFlagged
	}
	probabilities := make(map[*field.Space]float64)
	minesRemaining := s.field.MineCount() - s.field.FlaggedCount()
	for partIndex, solutionPart := range solutionsByPart {
		minFlaggedSumOtherParts := minFlaggedSum - minFlaggedByPart[partIndex]
		solutionCount := 0
		revealedSpaceCounts := make(map[*field.Space]int)
		for _, solution := range solutionPart {
			if minFlaggedSumOtherParts+len(solution.Flagged) > minesRemaining {
				continue
			}
			for _, space := range solution.Revealed {
				revealedSpaceCounts[space]++
			}
			solutionCount++
		}
		if s.visualize {
			fmt.Printf("constrained mine-free probabilities for constraint group %d\n\n", partIndex)
		}
		for space, revealedCount := range revealedSpaceCounts {
			probabilities[space] = float64(revealedCount) / float64(solutionCount)
			if s.visualize {
				fmt.Printf("%4d %.2f\n", space.Index(), probabilities[space])
			}
		}
		if s.visualize {
			fmt.Println()
		}
	}
	return probabilities
}
