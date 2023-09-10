package matrix

import (
	"gitlab.ocnr.org/apps/minesweeper/field"
)

type Cell struct {
	value float64
	space *field.Space
	state field.State
}

func (c *Cell) State() field.State {
	return c.state
}

func (c *Cell) Space() *field.Space {
	return c.space
}

// copy returns a copy the cell.
func (c *Cell) copy() *Cell {
	return &Cell{c.value, c.space, c.state}
}
