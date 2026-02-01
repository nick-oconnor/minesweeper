package field

import (
	"fmt"
	"math/rand"
)

type Field struct {
	width, height, flaggedCount, mineCount int
	spaces, revealedSpaces                 []*Space
	unknownSpaces                          map[*Space]bool
	firstMove                              bool
}

// NewField creates a new field with the given parameters.
func NewField(width int, height int, mineCount int) *Field {
	spaceCount := width * height
	f := &Field{width, height, 0, 0, make([]*Space, 0, spaceCount), make([]*Space, 0, spaceCount), make(map[*Space]bool, spaceCount), true}
	// Set up the field.
	for spaceIndex := 0; spaceIndex < f.width*f.height; spaceIndex++ {
		space := newSpace(spaceIndex)
		f.spaces = append(f.spaces, space)
		f.unknownSpaces[space] = true
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
					space.unknownNeighbors[neighbor] = true
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

// MineCount returns the number of mines in the field.
func (f *Field) MineCount() int {
	return f.mineCount
}

// FlaggedCount returns the number of spaces flagged in the field.
func (f *Field) FlaggedCount() int {
	return f.flaggedCount
}

// Spaces returns all spaces in the field.
func (f *Field) Spaces() []*Space {
	return f.spaces
}

// UnknownSpaces returns the unknown spaces in the field.
func (f *Field) UnknownSpaces() map[*Space]bool {
	return f.unknownSpaces
}

// Reveal reveals the given space in the field.
func (f *Field) Reveal(space *Space) error {
	if f.firstMove {
		f.firstMove = false
		if space.hasMine {
			space.removeMine()
			f.mineCount--
			f.addRandomMine(space)
		}
	}
	if err := space.reveal(); err != nil {
		return err
	}
	delete(f.unknownSpaces, space)
	f.revealedSpaces = append(f.revealedSpaces, space)
	f.recursiveReveal(space)
	return nil
}

// Flag flags the given space in the field.
func (f *Field) Flag(space *Space) {
	space.flag()
	delete(f.unknownSpaces, space)
	f.flaggedCount++
}

// AddMine adds a mine to the given space.
func (f *Field) AddMine(space *Space) {
	space.addMine()
	f.mineCount++
}

// RevealedEdgeSpaces retrieves the revealed edge spaces in the field.
func (f *Field) RevealedEdgeSpaces() []*Space {
	var spaces []*Space
	for _, revealedSpace := range f.revealedSpaces {
		if len(revealedSpace.unknownNeighbors) > 0 && revealedSpace.mineNeighborCount > 0 {
			spaces = append(spaces, revealedSpace)
		}
	}
	return spaces
}

// UnknownEdgeSpaces retrieves the unknown edge spaces for the field.
func (f *Field) UnknownEdgeSpaces() []*Space {
	var spaces []*Space
	for unknownSpace := range f.unknownSpaces {
		if unknownSpace.revealedNeighborCount > 0 {
			spaces = append(spaces, unknownSpace)
		}
	}
	return spaces
}

// Print prints the field to stdout.
func (f *Field) Print() {
	for spaceIndex, space := range f.spaces {
		fmt.Print("|")
		spaceContent := "   "
		switch space.state {
		case Revealed:
			if space.mineNeighborCount > 0 {
				spaceContent = fmt.Sprintf(" %d ", space.mineNeighborCount)
			}
			if space.hasMine {
				spaceContent = " * "
			}
		case Flagged:
			spaceContent = " * "
		case Unknown:
			spaceContent = " - "
		}
		fmt.Print(spaceContent)
		if spaceIndex%f.width == f.width-1 {
			fmt.Println("|")
		}
	}
	fmt.Println()
}

// recursiveReveal recursively reveals neighbors which are not touching mines.
func (f *Field) recursiveReveal(space *Space) {
	if space.mineNeighborCount == 0 {
		unknownNeighbors := make(map[*Space]bool)
		for unknownNeighbor := range space.unknownNeighbors {
			unknownNeighbors[unknownNeighbor] = true
		}
		for unknownNeighbor := range unknownNeighbors {
			if err := unknownNeighbor.reveal(); err != nil {
				panic(err)
			}
			delete(f.unknownSpaces, unknownNeighbor)
			f.revealedSpaces = append(f.revealedSpaces, unknownNeighbor)
		}
		for unknownNeighbor := range unknownNeighbors {
			f.recursiveReveal(unknownNeighbor)
		}
	}
}

// addRandomMine adds a mine to a random space which is not the given space.
func (f *Field) addRandomMine(exclude *Space) {
	availableSpaces := make([]*Space, 0, len(f.spaces))
	for _, space := range f.spaces {
		if space != exclude && !space.hasMine {
			availableSpaces = append(availableSpaces, space)
		}
	}
	if len(availableSpaces) == 0 {
		panic("no available spaces for mine placement")
	}
	f.AddMine(availableSpaces[rand.Intn(len(availableSpaces))])
}
