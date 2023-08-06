package field

import (
	"fmt"
)

// State represents the various states of a space.
type State int

const (
	Unknown  State = iota
	Flagged  State = iota
	Revealed State = iota
)

// String returns a text version of the State enum.
func (s State) String() string {
	switch s {
	case Unknown:
		return "unknown"
	case Flagged:
		return "flagged"
	case Revealed:
		return "revealed"
	}
	panic(fmt.Sprintf("invalid state %d", s))
}

// Space represents a space in the field.
type Space struct {
	id                                                             string
	hasMine                                                        bool
	state                                                          State
	mineNeighborCount, flaggedNeighborCount, revealedNeighborCount int
	neighbors, unknownNeighbors                                    []*Space
}

// newSpace creates a new space with the given ID.
func newSpace(id string) *Space {
	return &Space{id, false, Unknown, 0, 0, 0, []*Space{}, []*Space{}}
}

// Id returns the ID for the space.
func (s *Space) Id() string {
	return s.id
}

// HasMine returns whether the space contains a mine. This panics if the space
// has not been revealed.
func (s *Space) HasMine() bool {
	if s.state != Revealed {
		panic(fmt.Sprintf("retrieve mine status for an unknown space %v", s.id))
	}
	return s.hasMine
}

// State returns the state of the space.
func (s *Space) State() State {
	return s.state
}

// Neighbors returns all neighbors of the space.
func (s *Space) Neighbors() []*Space {
	return s.neighbors
}

// MineNeighborCount returns the number of mine neighbors of the space. This
// panics if the space has not been revealed.
func (s *Space) MineNeighborCount() int {
	if s.state != Revealed {
		panic(fmt.Sprintf("retrieved mine neighbors for unknown space %v", s.id))
	}
	return s.mineNeighborCount
}

// FlaggedNeighborCount returns the number of flagged neighbors of the space.
func (s *Space) FlaggedNeighborCount() int {
	return s.flaggedNeighborCount
}

// RevealedNeighborCount returns the number of revealed neighbors of the space.
func (s *Space) RevealedNeighborCount() int {
	return s.revealedNeighborCount
}

// UnknownNeighbors returns unknown neighbors of the space.
func (s *Space) UnknownNeighbors() []*Space {
	return s.unknownNeighbors
}

// HasNeighbor returns whether the given space is a neighbor.
func (s *Space) HasNeighbor(space *Space) bool {
	for _, neighbor := range s.neighbors {
		if space == neighbor {
			return true
		}
	}
	return false
}

// AddMine adds a mine to the space.
func (s *Space) AddMine() {
	if s.hasMine {
		panic(fmt.Sprintf("mine added to space %v which already contains a mine", s.id))
	}
	s.hasMine = true
	for _, neighbor := range s.neighbors {
		neighbor.mineNeighborCount++
	}
}

// removeMine removes a mine from the space.
func (s *Space) removeMine() {
	if !s.hasMine {
		panic(fmt.Sprintf("mine removed from space %v which does not contain a mine", s.id))
	}
	s.hasMine = false
	for _, neighbor := range s.neighbors {
		neighbor.mineNeighborCount--
	}
}

// reveal marks the space as revealed and updates its neighbors. panicOnError
// should be true unless the move is a guess.
func (s *Space) reveal(panicOnError bool) error {
	if s.state != Unknown {
		panic(fmt.Sprintf("space %v revealed with invalid state %v", s.id, s.state))
	}
	s.state = Revealed
	if s.hasMine {
		if panicOnError {
			panic(fmt.Sprintf("space %v revealed which contains a mine", s.id))
		}
		return fmt.Errorf("space %v revealed which contains a mine", s.id)
	}
	for _, n := range s.neighbors {
		n.unknownNeighbors = remove(n.unknownNeighbors, s)
		n.revealedNeighborCount++
	}
	return nil
}

// flag marks the space as flagged and updates its neighbors. panicOnError should
// be true unless the move is a guess.
func (s *Space) flag(panicOnError bool) error {
	if s.state != Unknown {
		panic(fmt.Sprintf("space %v flagged with invalid state %v", s.id, s.state))
	}
	s.state = Flagged
	if !s.hasMine {
		if panicOnError {
			panic(fmt.Sprintf("space %v flagged which does not contain a mine", s.id))
		}
		return fmt.Errorf("space %v flagged which does not contain a mine", s.id)
	}
	for _, n := range s.neighbors {
		n.unknownNeighbors = remove(n.unknownNeighbors, s)
		n.flaggedNeighborCount++
	}
	return nil
}
