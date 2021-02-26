package solvers

import (
	"github.com/nick-oconnor/minesweeper/sim"
)

type Algorithm interface {
	MakeNextMove() (*sim.Square, string, error)
}
