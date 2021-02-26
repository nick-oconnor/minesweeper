package sim

const (
	Unknown  = iota
	Revealed = iota
	Marked   = iota
)

type SquaresSet map[*Square]*struct{}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func FirstMatch(squaresSet SquaresSet, f func(s *Square) bool) *Square {
	for square := range squaresSet {
		if f(square) {
			return square
		}
	}
	return nil
}

func FirstSquare() func(s *Square) bool {
	return func(s *Square) bool {
		return true
	}
}
