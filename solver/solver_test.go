package solver

import (
	"fmt"
	"testing"

	. "gitlab.ocnr.org/apps/minesweeper/field"
)

// TestSingleBasic tests basic deduction.
func TestSingleBasic(t *testing.T) {
	solver := newTestSolver(t, [][]string{
		{"M", "M", "M"},
		{" ", " ", " "},
		{"R", " ", " "},
	})
	solver.solve()
	solver.assert([][]string{
		{"F", "F", "F"},
		{"R", "R", "R"},
		{"R", "R", "R"},
	})
}

// TestMultiOneOne tests 1-1 pattern recognition.
func TestMultiOneOne(t *testing.T) {
	solver := newTestSolver(t, [][]string{
		{" ", " ", " ", " ", " "},
		{"M", " ", " ", "M", " "},
		{" ", " ", " ", " ", " "},
		{"R", " ", " ", " ", " "},
	})
	solver.solve()
	solver.assert([][]string{
		{"R", "R", "R", "R", "R"},
		{" ", " ", "R", " ", " "},
		{"R", "R", "R", "R", "R"},
		{"R", "R", "R", "R", "R"},
	})
}

// TestMultiOneTwo tests 1-2 pattern recognition.
func TestMultiOneTwo(t *testing.T) {
	solver := newTestSolver(t, [][]string{
		{" ", " ", " "},
		{"M", " ", "M"},
		{" ", " ", " "},
		{"R", " ", " "},
	})
	solver.solve()
	solver.assert([][]string{
		{"R", "R", "R"},
		{"F", "R", "F"},
		{"R", "R", "R"},
		{"R", "R", "R"},
	})
}

// TestMultiOneTwoCorner tests 1-2 pattern recognition in a corner.
func TestMultiOneTwoCorner(t *testing.T) {
	solver := newTestSolver(t, [][]string{
		{" ", "M", "M", " ", " ", "M"},
		{" ", " ", " ", " ", " ", "M"},
		{"R", " ", " ", " ", " ", " "},
	})
	solver.solve()
	solver.assert([][]string{
		{"R", "F", "F", "R", "R", "F"},
		{"R", "R", "R", "R", "R", "F"},
		{"R", "R", "R", "R", "R", "R"},
	})
}

// TestTwoTwoCorner tests 2-2 pattern recognition in a corner.
func TestTwoTwoCorner(t *testing.T) {
	solver := newTestSolver(t, [][]string{
		{"M", "M", " ", "M", " "},
		{" ", " ", " ", " ", "M"},
		{" ", " ", " ", " ", " "},
		{" ", " ", " ", " ", "M"},
		{"R", " ", " ", " ", "M"},
	})
	solver.solve()
	solver.assert([][]string{
		{"F", "F", "R", "F", "R"},
		{"R", "R", "R", "R", "F"},
		{"R", "R", "R", "R", "R"},
		{"R", "R", "R", "R", "F"},
		{"R", "R", "R", "R", "F"},
	})
}

// TestMultiFractional tests if matrices containing fractions are correctly handled.
func TestMultiFractional(t *testing.T) {
	solver := newTestSolver(t, [][]string{
		{"R", " ", " ", " ", " "},
		{" ", " ", " ", " ", " "},
		{"F", "F", "F", "R", "F"},
		{"M", " ", " ", "M", " "},
		{"R", " ", "R", " ", "R"},
	})
	solver.solve()
	solver.assert([][]string{
		{"R", "R", "R", "R", "R"},
		{"R", "R", "R", "R", "R"},
		{"F", "F", "F", "R", "F"},
		{" ", " ", " ", " ", " "},
		{"R", " ", "R", " ", "R"},
	})
}

func TestRandomWin(t *testing.T) {
	for {
		if _, _, err := NewSolver(NewField(5, 5, 5)).WithFieldPrinter(printField).WithMatrixPrinter(printMatrix).Solve(); err == nil {
			break
		}
	}
}

func TestRandomLoss(t *testing.T) {
	for {
		if _, _, err := NewSolver(NewField(5, 5, 5)).WithFieldPrinter(printField).WithMatrixPrinter(printMatrix).Solve(); err != nil {
			break
		}
	}
}

// TestSolver solves and asserts on given field layouts.
type TestSolver struct {
	field *Field
	*testing.T
}

// newTestSolver creates a new test solver with the given layout.
func newTestSolver(t *testing.T, layout [][]string) *TestSolver {
	width := len(layout[0])
	height := len(layout)
	field := NewField(width, height, 0)
	for rowIndex, row := range layout {
		for colIndex, char := range row {
			if char == "M" || char == "F" {
				field.AddMine(field.Spaces()[rowIndex*width+colIndex])
			}
		}
	}
	for rowIndex, row := range layout {
		for colIndex, char := range row {
			space := field.Spaces()[rowIndex*width+colIndex]
			switch char {
			case "F":
				if err := field.Flag(space, true); err != nil {
					t.Fatal(err)
				}
			case "R":
				if err := field.Reveal(space, true); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	return &TestSolver{field, t}
}

// solve solves the given layout.
func (t *TestSolver) solve() {
	solver := NewSolver(t.field).WithGuessing(false).WithFieldPrinter(printField).WithMatrixPrinter(printMatrix)
	printField(t.field, nil, nil)
	if _, _, err := solver.Solve(); err != nil {
		t.Fatal(err)
	}
}

// assert ensures that the solved field matches the given layout.
func (t *TestSolver) assert(layout [][]string) {
	assertions := make(map[string]string)
	height := len(layout)
	width := len(layout[0])
	if width != t.field.Width() {
		t.Fatalf("solution width: want %v, found %v", t.field.Width(), width)
	}
	if height != t.field.Height() {
		t.Fatalf("solution height: want %v, found %v", t.field.Height(), height)
	}
	for rowIndex, row := range layout {
		for colIndex, char := range row {
			assertions[IndexToId(t.field.Width(), rowIndex*t.field.Width()+colIndex)] = char
		}
	}
	results := make(map[string]*Space)
	for _, space := range t.field.Spaces() {
		results[space.Id()] = space
	}
	for spaceId, char := range assertions {
		space := results[spaceId]
		var state State
		switch char {
		case " ":
			state = Unknown
		case "F":
			state = Flagged
		case "R":
			state = Revealed
		default:
			t.Fatalf("space %v: unknown assertion %v", spaceId, state)
		}
		if space.State() != state {
			t.Fatalf("space %v state: want %v, found %v", spaceId, state, space.State())
		}
		delete(assertions, spaceId)
	}
}

// printField provides a test-friendly field printer.
func printField(field *Field, changedSpace *Space, err error) {
	for spaceIndex, space := range field.Spaces() {
		fmt.Print("|")
		spaceContent := "   "
		switch space.State() {
		case Revealed:
			if space.MineNeighborCount() > 0 {
				spaceContent = fmt.Sprintf(" %d ", space.MineNeighborCount())
			}
			if space.HasMine() {
				spaceContent = " * "
			}
		case Flagged:
			spaceContent = " * "
		case Unknown:
			spaceContent = " - "
		}
		fmt.Print(spaceContent)
		if spaceIndex%field.Width() == field.Width()-1 {
			fmt.Println("|")
		}
	}
	if changedSpace != nil {
		if err == nil {
			fmt.Printf("%v %v: success\n", changedSpace.Id(), changedSpace.State())
		} else {
			fmt.Printf("%v %v: %v\n", changedSpace.Id(), changedSpace.State(), err)
		}
	}
	fmt.Println()
}

// printMatrix provides a test-friendly matrix printer.
func printMatrix(matrix Matrix) {
	for rowIndex, row := range matrix {
		if rowIndex == 0 {
			fmt.Print(" ")
			for _, cell := range row[:len(row)-1] {
				fmt.Printf("%4s ", cell.space.Id())
			}
			fmt.Println()
		}
		for cellIndex, cell := range row {
			if cellIndex%len(row) == len(row)-1 {
				fmt.Printf(" | %4s %s", fmt.Sprintf("%1.1f", cell.value), cell.space.Id())
			} else {
				fmt.Printf(" %4s", fmt.Sprintf("%1.1f", cell.value))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
