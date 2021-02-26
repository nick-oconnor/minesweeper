package solvers

import (
	"github.com/nick-oconnor/minesweeper/sim"
)

type Primary struct {
	field *sim.Field
}

func NewPrimary(field *sim.Field) *Primary {
	return &Primary{field}
}

func (p *Primary) MakeNextMove() (*sim.Square, string, error) {
	for square := range p.field.RevealedSquares() {
		unknownNeighbor := sim.FirstMatch(square.UnknownNeighbors(), sim.FirstSquare())
		if unknownNeighbor == nil {
			continue
		}
		squareNumber := square.Number()
		if squareNumber == len(square.MarkedNeighbors()) {
			return unknownNeighbor, "revealing a square that does not contain a mine", unknownNeighbor.Reveal()
		}
		if squareNumber-len(square.MarkedNeighbors()) == len(square.UnknownNeighbors()) {
			unknownNeighbor.Mark()
			return unknownNeighbor, "marking a square that contains a mine", nil
		}
	}
	return nil, "no primary moves available", nil
}
