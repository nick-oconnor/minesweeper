package field

import (
	"math/rand"
)

// Field represents a 2D minesweeper field.
type Field struct {
	width, height, flaggedCount, mineCount int
	spaces, revealedSpaces, unknownSpaces  []*Space
	firstMove                              bool
}

// NewField creates a new field with the given parameters.
func NewField(width int, height int, mineCount int) *Field {
	f := &Field{width, height, 0, 0, []*Space{}, []*Space{}, []*Space{}, true}
	// Set up the field.
	for spaceIndex := 0; spaceIndex < f.width*f.height; spaceIndex++ {
		space := newSpace(IndexToId(f.width, spaceIndex))
		f.spaces = append(f.spaces, space)
		f.unknownSpaces = append(f.unknownSpaces, space)
	}
	// Set up neighbor references.
	for spaceIndex, space := range f.spaces {
		rowIndex := spaceIndex / f.width
		colIndex := spaceIndex % f.width
		for neighborRowIndex := max(0, rowIndex-1); neighborRowIndex <= min(f.height-1, rowIndex+1); neighborRowIndex++ {
			for neighborColIndex := max(0, colIndex-1); neighborColIndex <= min(f.width-1, colIndex+1); neighborColIndex++ {
				neighbor := f.spaces[neighborRowIndex*f.width+neighborColIndex]
				if neighbor != space {
					space.neighbors = append(space.neighbors, neighbor)
					space.unknownNeighbors = append(space.unknownNeighbors, neighbor)
				}
			}
		}
	}
	// Add mines to the field.
	for i := 0; i < mineCount; i++ {
		f.addRandomMine(nil)
	}
	return f
}

// Width returns the width of the field.
func (f *Field) Width() int {
	return f.width
}

// Height returns the height of the field.
func (f *Field) Height() int {
	return f.height
}

// Spaces returns all spaces in the field.
func (f *Field) Spaces() []*Space {
	return f.spaces
}

// RevealedSpaces returns the revealed spaces in the field.
func (f *Field) RevealedSpaces() []*Space {
	return f.revealedSpaces
}

// UnknownSpaces returns the unknown spaces in the field.
func (f *Field) UnknownSpaces() []*Space {
	return f.unknownSpaces
}

// Reveal reveals the given space in the field. This panics if the space contains
// a mine and panicOnError is true.
func (f *Field) Reveal(space *Space, panicOnError bool) error {
	if f.firstMove {
		f.firstMove = false
		if space.hasMine {
			space.removeMine()
			f.mineCount--
			f.addRandomMine(space)
		}
	}
	if err := space.reveal(panicOnError); err != nil {
		return err
	}
	f.unknownSpaces = remove(f.unknownSpaces, space)
	f.revealedSpaces = append(f.revealedSpaces, space)
	return f.recursiveReveal(space)
}

// recursiveReveal recursively reveals neighbors which are not touching mines.
func (f *Field) recursiveReveal(space *Space) error {
	if space.MineNeighborCount() == 0 {
		unknownNeighbors := append([]*Space{}, space.unknownNeighbors...)
		for _, unknownNeighbor := range unknownNeighbors {
			if err := unknownNeighbor.reveal(false); err != nil {
				return err
			}
			f.unknownSpaces = remove(f.unknownSpaces, unknownNeighbor)
			f.revealedSpaces = append(f.revealedSpaces, unknownNeighbor)
		}
		for _, unknownNeighbor := range unknownNeighbors {
			if err := f.recursiveReveal(unknownNeighbor); err != nil {
				return err
			}
		}
	}
	return nil
}

// Flag flags the given space in the field. This panics if the space does not
// contain a mine and panicOnError is true.
func (f *Field) Flag(space *Space, panicOnError bool) error {
	if err := space.flag(panicOnError); err != nil {
		return err
	}
	f.unknownSpaces = remove(f.unknownSpaces, space)
	f.flaggedCount++
	return nil
}

// AddMine adds a mine to the given space.
func (f *Field) AddMine(space *Space) {
	space.AddMine()
	f.mineCount++
}

// addRandomMine adds a mine to a random space which is not the given space.
func (f *Field) addRandomMine(excludeSpace *Space) {
	for {
		space := f.spaces[rand.Intn(f.width*f.height)]
		if space == excludeSpace || space.hasMine {
			continue
		}
		f.AddMine(space)
		return
	}
}
