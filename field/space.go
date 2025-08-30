package field

import (
	"fmt"
)

type Space struct {
	index                                                          int
	state                                                          State
	hasMine                                                        bool
	mineNeighborCount, flaggedNeighborCount, revealedNeighborCount int
	neighbors                                                      []*Space
	unknownNeighbors                                               map[*Space]bool
}

// Index returns the index for the space.
func (s *Space) Index() int {
	return s.index
}

// State returns the state of the space.
func (s *Space) State() State {
	return s.state
}

// Neighbors returns all neighbors of the space.
func (s *Space) Neighbors() []*Space {
	return s.neighbors
}

// UnknownNeighbors returns all unknown neighbors of the space.
func (s *Space) UnknownNeighbors() map[*Space]bool {
	return s.unknownNeighbors
}

// MineNeighborCount returns the space's number of mine neighbors.
func (s *Space) MineNeighborCount() int {
	if s.state != Revealed {
		print(fmt.Sprintf("retrieved mine neighbors for unknown space %d", s.index))
	}
	return s.mineNeighborCount
}

// FlaggedNeighborCount returns the space's number of flagged neighbors.
func (s *Space) FlaggedNeighborCount() int {
	return s.flaggedNeighborCount
}

// newSpace creates a new space with the given index.
func newSpace(index int) *Space {
	return &Space{index, Unknown, false, 0, 0, 0, make([]*Space, 0, 8), make(map[*Space]bool, 8)}
}

// addMine adds a mine to the space.
func (s *Space) addMine() {
	if s.hasMine {
		panic(fmt.Sprintf("mine added to space %d which already contains a mine", s.index))
	}
	s.hasMine = true
	for _, neighbor := range s.neighbors {
		neighbor.mineNeighborCount++
	}
}

// removeMine removes a mine from the space.
func (s *Space) removeMine() {
	if !s.hasMine {
		panic(fmt.Sprintf("mine removed from space %d which does not contain a mine", s.index))
	}
	s.hasMine = false
	for _, neighbor := range s.neighbors {
		neighbor.mineNeighborCount--
	}
}

// reveal marks the space as revealed and updates its neighbors.
func (s *Space) reveal() error {
	if s.state != Unknown {
		panic(fmt.Sprintf("space %d revealed with invalid state %s", s.index, s.state))
	}
	s.state = Revealed
	if s.hasMine {
		return fmt.Errorf("space %d revealed which contains a mine", s.index)
	}
	for _, neighbor := range s.neighbors {
		delete(neighbor.unknownNeighbors, s)
		neighbor.revealedNeighborCount++
	}
	return nil
}

// flag marks the space as flagged and updates its neighbors.
func (s *Space) flag() {
	if s.state != Unknown {
		panic(fmt.Sprintf("space %d flagged with invalid state %s", s.index, s.state))
	}
	s.state = Flagged
	if !s.hasMine {
		panic(fmt.Sprintf("space %d flagged which does not contain a mine", s.index))
	}
	for _, neighbor := range s.neighbors {
		delete(neighbor.unknownNeighbors, s)
		neighbor.flaggedNeighborCount++
	}
}
