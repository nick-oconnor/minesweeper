package solver

import (
	"testing"

	"gitlab.ocnr.org/apps/minesweeper/field"
	"gitlab.ocnr.org/apps/minesweeper/matrix"
)

// TestReadme readme example.
func TestReadme(t *testing.T) {
	f := newTestField(t, [][]string{
		{"R", " ", " "},
		{" ", " ", "M"},
		{" ", "M", " "},
	})
	solve(t, f)
	assert(t, f, [][]string{
		{"R", "R", "R"},
		{"R", "R", "F"},
		{"R", "F", "R"},
	})
}

// TestSingleBasic tests basic deduction.
func TestSingleBasic(t *testing.T) {
	f := newTestField(t, [][]string{
		{"M", "M", "M"},
		{" ", " ", " "},
		{"R", " ", " "},
	})
	solve(t, f)
	assert(t, f, [][]string{
		{"F", "F", "F"},
		{"R", "R", "R"},
		{"R", "R", "R"},
	})
}

// TestMultiOneOne tests 1-1 pattern recognition.
func TestMultiOneOne(t *testing.T) {
	f := newTestField(t, [][]string{
		{" ", " ", " "},
		{" ", "M", " "},
		{"R", "R", "R"},
	})
	solve(t, f)
	assert(t, f, [][]string{
		{"R", "R", "R"},
		{"R", "F", "R"},
		{"R", "R", "R"},
	})
}

// TestMultiOneTwo tests 1-2 pattern recognition.
func TestMultiOneTwo(t *testing.T) {
	f := newTestField(t, [][]string{
		{" ", " ", " "},
		{"M", " ", "M"},
		{"R", "R", "R"},
	})
	solve(t, f)
	assert(t, f, [][]string{
		{"R", "R", "R"},
		{"F", "R", "F"},
		{"R", "R", "R"},
	})
}

// TestMultiOneTwoCorner tests 1-2 pattern recognition in a corner.
func TestMultiOneTwoCorner(t *testing.T) {
	f := newTestField(t, [][]string{
		{" ", "M", "M", " ", " ", "M"},
		{" ", " ", " ", " ", " ", "M"},
		{"R", " ", " ", " ", " ", " "},
	})
	solve(t, f)
	assert(t, f, [][]string{
		{"R", "F", "F", "R", "R", "F"},
		{"R", "R", "R", "R", "R", "F"},
		{"R", "R", "R", "R", "R", "R"},
	})
}

// TestTwoTwoCorner tests 2-2 pattern recognition in a corner.
func TestTwoTwoCorner(t *testing.T) {
	f := newTestField(t, [][]string{
		{"M", "M", " ", "M", " "},
		{" ", " ", " ", " ", "M"},
		{" ", " ", " ", " ", " "},
		{" ", " ", " ", " ", "M"},
		{"R", " ", " ", " ", "M"},
	})
	solve(t, f)
	assert(t, f, [][]string{
		{"F", "F", "R", "F", "R"},
		{"R", "R", "R", "R", "F"},
		{"R", "R", "R", "R", "R"},
		{"R", "R", "R", "R", "F"},
		{"R", "R", "R", "R", "F"},
	})
}

// TestDecoupled tests if layouts with decoupled matrices are correctly handled.
func TestDecoupled(t *testing.T) {
	f := newTestField(t, [][]string{
		{"R", "R", " ", " ", " "},
		{" ", "M", " ", " ", " "},
		{" ", " ", " ", " ", " "},
		{" ", " ", " ", "M", " "},
		{" ", " ", " ", "R", "R"},
	})
	solve(t, f)
	assert(t, f, [][]string{
		{"R", "R", "R", "R", "R"},
		{"R", "F", "R", "R", "R"},
		{"R", "R", "R", "R", "R"},
		{"R", "R", "R", "F", "R"},
		{"R", "R", "R", "R", "R"},
	})
}

// TestProbabilities1 tests the probabilities calculated for layout 1.
func TestProbabilities1(t *testing.T) {
	f := newTestField(t, [][]string{
		{" ", " ", "M", " ", "M"},
		{"F", "F", " ", " ", " "},
		{" ", " ", "R", " ", " "},
		{" ", " ", "M", " ", " "},
		{"R", " ", " ", " ", "R"},
	})
	solver := NewSolver(f, true)
	f.Print()
	want := make(map[int]float64)
	want[2] = 2.0 / 4.0
	want[3] = 2.0 / 4.0
	want[4] = 2.0 / 4.0
	want[7] = 2.0 / 4.0
	want[17] = 2.0 / 4.0
	want[22] = 2.0 / 4.0
	assertProbabilities(t, f, want, solver.probabilityPerSpace(matrix.NewMatrix(f)))
}

// TestProbabilities2 tests the probabilities calculated for layout 1.
func TestProbabilities2(t *testing.T) {
	f := newTestField(t, [][]string{
		{" ", "M", "R", " ", "M"},
		{" ", " ", " ", "R", "M"},
		{" ", " ", " ", "F", "F"},
		{" ", " ", " ", " ", " "},
		{"R", " ", " ", " ", " "},
	})
	solver := NewSolver(f, true)
	f.Print()
	want := make(map[int]float64)
	want[0] = 1.0 / 3.0
	want[1] = 2.0 / 3.0
	want[3] = 1.0 / 3.0
	want[4] = 1.0 / 3.0
	want[9] = 1.0 / 3.0
	assertProbabilities(t, f, want, solver.probabilityPerSpace(matrix.NewMatrix(f)))
}

// TestProbabilities3 tests the probabilities calculated for layout 2.
func TestProbabilities3(t *testing.T) {
	f := newTestField(t, [][]string{
		{"M", "M", " ", " ", " "},
		{"M", " ", " ", " ", " "},
		{" ", " ", "R", "M", " "},
		{" ", " ", " ", " ", " "},
		{" ", " ", "M", " ", "R"},
	})
	solver := NewSolver(f, true)
	f.Print()
	want := make(map[int]float64)
	want[6] = 6.0 / 7.0
	want[7] = 6.0 / 7.0
	want[8] = 6.0 / 7.0
	want[11] = 6.0 / 7.0
	want[13] = 6.0 / 7.0
	want[14] = 1.0 / 7.0
	want[16] = 6.0 / 7.0
	want[17] = 6.0 / 7.0
	want[22] = 1.0 / 7.0
	assertProbabilities(t, f, want, solver.probabilityPerSpace(matrix.NewMatrix(f)))
}

// TestProbabilities4 tests the probabilities calculated for layout 4.
func TestProbabilities4(t *testing.T) {
	f := newTestField(t, [][]string{
		{"R", " ", " ", " ", "R"},
		{"R", "M", " ", " ", " "},
		{"R", "R", " ", " ", " "},
		{" ", "M", "R", " ", "M"},
		{"R", "M", "R", "F", "R"},
	})
	solver := NewSolver(f, true)
	f.Print()
	want := make(map[int]float64)
	want[1] = 2.0 / 3.0
	want[6] = 1.0 / 3.0
	want[15] = 1.0 / 3.0
	want[16] = 1.0 / 3.0
	want[18] = 1.0 / 3.0
	want[19] = 2.0 / 3.0
	want[21] = 1.0 / 3.0
	assertProbabilities(t, f, want, solver.probabilityPerSpace(matrix.NewMatrix(f)))
}

// TestProbabilitiesSplit tests the probabilities calculated for
func TestProbabilitiesSplit(t *testing.T) {
	f := newTestField(t, [][]string{
		{"R", " ", " ", " ", " "},
		{" ", "M", " ", " ", " "},
		{" ", " ", " ", " ", " "},
		{" ", " ", " ", "M", " "},
		{" ", " ", " ", " ", "R"},
	})
	solver := NewSolver(f, true)
	f.Print()
	want := make(map[int]float64)
	want[1] = 2.0 / 3.0
	want[5] = 2.0 / 3.0
	want[6] = 2.0 / 3.0
	want[18] = 2.0 / 3.0
	want[19] = 2.0 / 3.0
	want[23] = 2.0 / 3.0
	assertProbabilities(t, f, want, solver.probabilityPerSpace(matrix.NewMatrix(f)))
}

// newTestField creates a new field with the given layout.
func newTestField(t *testing.T, layout [][]string) *field.Field {
	width := len(layout[0])
	height := len(layout)
	f := field.NewField(width, height, 0)
	for rowIndex, row := range layout {
		for colIndex, char := range row {
			if char == "M" || char == "F" {
				f.AddMine(f.Spaces()[rowIndex*f.Width()+colIndex])
			}
		}
	}
	for rowIndex, row := range layout {
		for colIndex, char := range row {
			space := f.Spaces()[rowIndex*f.Width()+colIndex]
			switch char {
			case "F":
				f.Flag(space)
			case "R":
				if err := f.Reveal(space); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	return f
}

// solve solves the given field.
func solve(t *testing.T, f *field.Field) *Solver {
	solver := NewSolver(f, true)
	f.Print()
	if result := solver.Solve(); !result.Won {
		t.Fatal("game lost")
	}
	return solver
}

// assert ensures that the given field matches the given layout.
func assert(t *testing.T, f *field.Field, layout [][]string) {
	assertions := make(map[int]string)
	height := len(layout)
	width := len(layout[0])
	if width != f.Width() {
		t.Fatalf("solution width: want %d, found %d", f.Width(), width)
	}
	if height != f.Height() {
		t.Fatalf("solution height: want %d, found %d", f.Height(), height)
	}
	for rowIndex, row := range layout {
		for colIndex, char := range row {
			assertions[rowIndex*f.Width()+colIndex] = char
		}
	}
	for spaceId, char := range assertions {
		space := f.Spaces()[spaceId]
		var state field.State
		switch char {
		case " ":
			state = field.Unknown
		case "F":
			state = field.Flagged
		case "R":
			state = field.Revealed
		default:
			t.Fatalf("space %d: unknown assertion %s", spaceId, state)
		}
		if space.State() != state {
			t.Fatalf("space %d state: want %s, found %s", spaceId, state, space.State())
		}
		delete(assertions, spaceId)
	}
}

func assertProbabilities(t *testing.T, f *field.Field, want map[int]float64, found map[*field.Space]float64) {
	for k, v := range want {
		value := found[f.Spaces()[k]]
		if value != v {
			t.Fatalf("probability for %d: want %.2f, found %.2f", k, v, value)
		}
	}
}
