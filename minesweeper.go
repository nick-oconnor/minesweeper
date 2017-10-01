package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/buger/goterm"
)

const (
	unknown  = iota
	revealed = iota
	flagged  = iota
)

type fieldSpace struct {
	value            int
	state            int
	neighbors        []*fieldSpace
	unknownNeighbors []*fieldSpace
	mineNeighbors    []*fieldSpace
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

func newField(width, height, mines int) (field [][]*fieldSpace) {
	for i := 0; i < height; i++ {
		row := []*fieldSpace{}
		for j := 0; j < width; j++ {
			row = append(row, &fieldSpace{0, unknown, []*fieldSpace{}, []*fieldSpace{}, []*fieldSpace{}})
		}
		field = append(field, row)
	}
	for i, row := range field {
		for j, space := range row {
			for _, row := range field[max(0, i-1) : min(height-1, i+1)+1] {
				for _, n := range row[max(0, j-1) : min(width-1, j+1)+1] {
					if n != space {
						space.neighbors = append(space.neighbors, n)
						space.unknownNeighbors = append(space.unknownNeighbors, n)
					}
				}
			}
		}
	}
	for i := 0; i < mines; i++ {
		var space *fieldSpace
		for {
			space = field[rand.Intn(height)][rand.Intn(width)]
			if space.value != -1 {
				break
			}
		}
		space.value = -1
		for _, n := range space.neighbors {
			if n.value != -1 {
				n.value++
			}
		}
	}
	return
}

func remove(slice []*fieldSpace, space *fieldSpace) []*fieldSpace {
	for i, e := range slice {
		if space == e {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func changeState(space *fieldSpace, state int, known *int) {
	space.state = state
	*known++
	for _, n := range space.neighbors {
		if state != unknown {
			n.unknownNeighbors = remove(n.unknownNeighbors, space)
		}
		if state == flagged {
			n.mineNeighbors = append(n.mineNeighbors, space)
		}
	}
}

func nextAction(edgeSpaces []*fieldSpace, field [][]*fieldSpace, known *int) *fieldSpace {
	for i := 0; i < len(edgeSpaces); i++ {
		e := edgeSpaces[i]
		if len(e.unknownNeighbors) == 0 {
			edgeSpaces = append(edgeSpaces[:i], edgeSpaces[i+1:]...)
			i--
			continue
		}
		n := e.unknownNeighbors[len(e.unknownNeighbors)-1]
		switch e.value - len(e.mineNeighbors) {
		case 0:
			changeState(n, revealed, known)
			return n
		case len(e.unknownNeighbors):
			changeState(n, flagged, known)
			return n
		}
	}
	var n *fieldSpace
	for {
		n = field[rand.Intn(len(field))][rand.Intn(len(field[0]))]
		if n.state == unknown {
			break
		}
	}
	changeState(n, revealed, known)
	return n
}

func revealNeighbors(space *fieldSpace, edgeSpaces []*fieldSpace, known *int) []*fieldSpace {
	unknownNeighbors := space.unknownNeighbors
	space.unknownNeighbors = []*fieldSpace{}
	for _, n := range unknownNeighbors {
		changeState(n, revealed, known)
		if n.value > 0 {
			edgeSpaces = append(edgeSpaces, n)
		}
	}
	for _, n := range unknownNeighbors {
		if n.value == 0 {
			edgeSpaces = revealNeighbors(n, edgeSpaces, known)
		}
	}
	return edgeSpaces
}

func showField(field [][]*fieldSpace) {
	goterm.MoveCursor(1, 1)
	for _, row := range field {
		for _, space := range row {
			goterm.Print("|")
			switch space.state {
			case revealed:
				switch space.value {
				case -1:
					goterm.Print(goterm.Background(goterm.Color("*", goterm.BLACK), goterm.RED))
				case 0:
					goterm.Print(" ")
				default:
					goterm.Print(space.value)
				}
			case flagged:
				goterm.Print(goterm.Background(goterm.Color("*", goterm.BLACK), goterm.YELLOW))
			case unknown:
				goterm.Print(goterm.Background(" ", goterm.WHITE))
			}
		}
		goterm.Println("|")
	}
	goterm.Flush()
}

func main() {
	width := flag.Int("width", 20, "width of the field")
	height := flag.Int("height", 20, "height of the field")
	mines := flag.Int("mines", 50, "number of mines")
	duration := flag.Duration("duration", time.Second/2, "action duration")
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	field := newField(*width, *height, *mines)
	edgeSpaces := []*fieldSpace{}
	startTime := time.Now()
	known := 0
	defer func() {
		goterm.Println("Time elapsed:", time.Now().Sub(startTime))
		goterm.Flush()
	}()
	goterm.Clear()
	for known < *width**height {
		space := nextAction(edgeSpaces, field, &known)
		if space.state == revealed {
			switch space.value {
			case -1:
				showField(field)
				goterm.Println("Game lost.")
				return
			case 0:
				edgeSpaces = revealNeighbors(space, edgeSpaces, &known)
			default:
				edgeSpaces = append(edgeSpaces, space)
			}
		}
		time.Sleep(*duration)
		showField(field)
	}
	goterm.Println("Game won!")
}
