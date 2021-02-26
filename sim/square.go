package sim

import (
	"errors"
)

type Square struct {
	hasMine                                              bool
	number                                               *int
	state                                                int
	field                                                *Field
	neighbors                                            []*Square
	unknownNeighbors, revealedNeighbors, markedNeighbors SquaresSet
}

func NewSquare(field *Field) *Square {
	return &Square{false, nil, Unknown, field, []*Square{}, SquaresSet{}, SquaresSet{}, SquaresSet{}}
}

func (s *Square) Number() int {
	if s.state != Revealed {
		panic("attempted to get the number for an unknown square")
	}
	if s.number != nil {
		return *s.number
	}
	number := 0
	s.number = &number
	for _, neighbor := range s.neighbors {
		if neighbor.hasMine {
			number++
		}
	}
	return number
}

func (s *Square) State() int {
	return s.state
}

func (s *Square) Neighbors() []*Square {
	return s.neighbors
}

func (s *Square) UnknownNeighbors() SquaresSet {
	return s.unknownNeighbors
}

func (s *Square) RevealedNeighbors() SquaresSet {
	return s.revealedNeighbors
}

func (s *Square) MarkedNeighbors() SquaresSet {
	return s.markedNeighbors
}

func (s *Square) Reveal() error {
	if s.state == Revealed {
		panic("attempted to reveal a revealed square")
	}
	if s.state == Marked {
		panic("attempted to reveal a marked square")
	}
	s.field.revealSquare(s)
	if s.hasMine {
		return errors.New("revealed a square containing a mine")
	}
	if s.Number() == 0 {
		for unknownNeighbor := FirstMatch(s.unknownNeighbors, FirstSquare()); unknownNeighbor != nil; unknownNeighbor = FirstMatch(s.unknownNeighbors, FirstSquare()) {
			err := unknownNeighbor.Reveal()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Square) Mark() {
	if s.state == Revealed {
		panic("attempted to mark a revealed square")
	}
	if s.state == Marked {
		panic("attempted to mark a marked square")
	}
	s.field.markSquare(s)
}
