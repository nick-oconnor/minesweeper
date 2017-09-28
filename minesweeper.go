package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/fatih/color"
)

const (
	unknown  = iota
	revealed = iota
	flagged  = iota
)

type fieldSpace struct {
	row   int
	col   int
	value int
	state int
}

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

func allNeighbors(target *fieldSpace, field [][]*fieldSpace, spaceFunc func(space *fieldSpace)) {
	for _, row := range field[max(0, target.row-1) : min(len(field)-1, target.row+1)+1] {
		for _, space := range row[max(0, target.col-1) : min(len(row)-1, target.col+1)+1] {
			spaceFunc(space)
		}
	}
}

func randSpace(field [][]*fieldSpace, spaceFunc func(space *fieldSpace) bool) (space *fieldSpace) {
	for {
		row := rand.Intn(len(field))
		col := rand.Intn(len(field[0]))
		space = field[row][col]
		if spaceFunc(space) {
			break
		}
	}
	return
}

func newBoard(width, height, mines int) (field [][]*fieldSpace) {
	for i := 0; i < height; i++ {
		row := []*fieldSpace{}
		for j := 0; j < width; j++ {
			row = append(row, &fieldSpace{i, j, 0, unknown})
		}
		field = append(field, row)
	}
	for i := 0; i < mines; i++ {
		space := randSpace(field, func(space *fieldSpace) bool {
			return space.value != -1
		})
		space.value = -1
		allNeighbors(space, field, func(space *fieldSpace) {
			if space.value != -1 {
				space.value++
			}
		})
	}
	return
}

func printBoard(field [][]*fieldSpace) {
	black := color.New(color.BgBlack, color.FgWhite).PrintfFunc()
	white := color.New(color.BgWhite).PrintfFunc()
	yellow := color.New(color.BgYellow).PrintFunc()
	red := color.New(color.BgRed).PrintFunc()
	fmt.Println()
	for _, row := range field {
		for _, space := range row {
			fmt.Print("|")
			switch space.state {
			case revealed:
				switch space.value {
				case -1:
					red(" ")
				case 0:
					black(" ")
				default:
					black(strconv.Itoa(space.value))
				}
			case flagged:
				yellow(" ")
			case unknown:
				white(" ")
			}
		}
		fmt.Println("|")
	}
	fmt.Println()
}

func allKnown(field [][]*fieldSpace) bool {
	for _, row := range field {
		for _, space := range row {
			if space.state == unknown {
				return false
			}
		}
	}
	return true
}

func revealNeighbors(space *fieldSpace, field [][]*fieldSpace) {
	allNeighbors(space, field, func(space *fieldSpace) {
		if space.state == unknown {
			space.state = revealed
			if space.value == 0 {
				revealNeighbors(space, field)
			}
		}
	})
}

func nextMove(field [][]*fieldSpace) *fieldSpace {
	for _, row := range field {
		for _, space := range row {
			if space.state == revealed && space.value > 0 {
				unknowns := []*fieldSpace{}
				mines := []*fieldSpace{}
				allNeighbors(space, field, func(s *fieldSpace) {
					switch s.state {
					case unknown:
						unknowns = append(unknowns, s)
					case flagged:
						mines = append(mines, s)
					}
				})
				if len(unknowns) == 0 {
					continue
				}
				s := unknowns[0]
				switch space.value - len(mines) {
				case 0:
					s.state = revealed
					return s
				case len(unknowns):
					s.state = flagged
					return s
				}
			}
		}
	}
	space := randSpace(field, func(space *fieldSpace) bool {
		return space.state == unknown
	})
	space.state = revealed
	return space
}

func main() {
	width := flag.Int("width", 20, "width of the field")
	height := flag.Int("height", 20, "height of the field")
	mines := flag.Int("mines", 50, "number of mines")
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		field := newBoard(*width, *height, *mines)
		for done := false; !done; {
			space := nextMove(field)
			printBoard(field)
			if space.state == revealed {
				switch space.value {
				case -1:
					fmt.Println("Game lost.")
					done = true
				case 0:
					revealNeighbors(space, field)
				}
			}
			if !done && allKnown(field) {
				fmt.Println("Game won!")
				done = true
			}
			time.Sleep(time.Second / 4)
		}
	}
}
