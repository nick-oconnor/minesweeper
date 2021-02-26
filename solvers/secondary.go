package solvers

import (
	"github.com/nick-oconnor/minesweeper/sim"
)

type Secondary struct {
	field *sim.Field
}

func NewSecondary(field *sim.Field) *Secondary {
	return &Secondary{field}
}

func (p *Secondary) MakeNextMove() (*sim.Square, string, error) {
	unknownCornerSquare := sim.FirstMatch(p.field.UnknownSquares(), func(s *sim.Square) bool {
		return len(s.Neighbors()) == 3
	})
	if unknownCornerSquare != nil {
		return unknownCornerSquare, "revealing a corner square", unknownCornerSquare.Reveal()
	}
	unknownEdgeSquare := sim.FirstMatch(p.field.UnknownSquares(), func(s *sim.Square) bool {
		return len(s.Neighbors()) == 5
	})
	if unknownEdgeSquare != nil {
		return unknownEdgeSquare, "revealing an edge square", unknownEdgeSquare.Reveal()
	}
	unknownSquare := sim.FirstMatch(p.field.UnknownSquares(), sim.FirstSquare())
	if unknownSquare != nil {
		return unknownSquare, "revealing a random square", unknownSquare.Reveal()
	}
	return nil, "no secondary moves available", nil
}
