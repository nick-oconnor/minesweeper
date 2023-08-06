package main

import (
	"fmt"
	"time"

	"github.com/buger/goterm"
	. "gitlab.ocnr.org/apps/minesweeper/field"
)

// Printer provides field print functions.
type Printer struct {
	duration time.Duration
}

// NewPrinter creates a new printer.
func NewPrinter(duration time.Duration) *Printer {
	return &Printer{duration}
}

// PrintField prints the given field to stdout.
func (p *Printer) PrintField(field *Field, move *Space, err error) {
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	for spaceIndex, space := range field.Spaces() {
		_, _ = goterm.Print("|")
		spaceColor := goterm.BLACK
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
			spaceColor = goterm.YELLOW
		case Unknown:
			spaceColor = goterm.BLUE
		}
		if space == move {
			spaceColor = goterm.MAGENTA
		}
		_, _ = goterm.Print(goterm.Background(spaceContent, spaceColor))
		if spaceIndex%field.Width() == field.Width()-1 {
			_, _ = goterm.Println("|")
		}
	}
	if err == nil {
		_, _ = goterm.Printf("%v %v: success\n", move.Id(), move.State())
	} else {
		_, _ = goterm.Printf("%v %v: %v\n", move.Id(), move.State(), err)
	}
	goterm.Flush()
	time.Sleep(p.duration)
}
