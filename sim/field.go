package sim

import (
	"github.com/buger/goterm"
	"math/rand"
	"strconv"
	"strings"
)

type Field struct {
	width, height                                  int
	squares                                        []*Square
	unknownSquares, revealedSquares, markedSquares SquaresSet
}

func NewField(width int, height int) *Field {
	return &Field{width, height, []*Square{}, SquaresSet{}, SquaresSet{}, SquaresSet{}}
}

func (f *Field) Width() int {
	return f.width
}

func (f *Field) Height() int {
	return f.height
}

func (f *Field) UnknownSquares() SquaresSet {
	return f.unknownSquares
}

func (f *Field) RevealedSquares() SquaresSet {
	return f.revealedSquares
}

func (f *Field) Init(numMines int) {
	for i := 0; i < f.width*f.height; i++ {
		square := NewSquare(f)
		f.squares = append(f.squares, square)
		f.unknownSquares[square] = nil
	}
	for i, square := range f.squares {
		row := i / f.width
		col := i % f.height
		for neighborRow := max(0, row-1); neighborRow <= min(f.height-1, row+1); neighborRow++ {
			for neighborCol := max(0, col-1); neighborCol <= min(f.width-1, col+1); neighborCol++ {
				neighbor := f.squares[neighborRow*f.width+neighborCol]
				if neighbor != square {
					square.neighbors = append(square.neighbors, neighbor)
					square.unknownNeighbors[neighbor] = nil
				}
			}
		}
	}
	for i := 0; i < numMines; i++ {
		var square *Square
		for {
			square = f.squares[rand.Intn(len(f.squares))]
			if !square.hasMine {
				break
			}
		}
		square.hasMine = true
	}
}

func (f *Field) Print(actionSquare *Square, reason string, err error) {
	goterm.MoveCursor(1, 1)
	for i, square := range f.squares {
		_, _ = goterm.Print("|")
		isCriticalSquare := square.state == Revealed && len(square.unknownNeighbors) > 0
		squareColor := goterm.BLACK
		squareContent := " "
		if square == actionSquare {
			squareColor = goterm.MAGENTA
		} else if isCriticalSquare {
			squareColor = goterm.BLUE
		}
		switch square.State() {
		case Revealed:
			squareNumber := square.Number()
			squareContent = strconv.Itoa(squareNumber)
			if squareNumber == 0 {
				squareContent = " "
			}
			if square.hasMine {
				squareColor = goterm.RED
				squareContent = "*"
			}
		case Marked:
			if square != actionSquare {
				squareColor = goterm.YELLOW
			}
			squareContent = "*"
		case Unknown:
			squareColor = goterm.WHITE
		}
		_, _ = goterm.Print(goterm.Background(squareContent, squareColor))
		if i%f.width == f.width-1 {
			_, _ = goterm.Println("|")
		}
	}
	_, _ = goterm.Print(strings.Repeat(" ", goterm.Width()) + "\r")
	_, _ = goterm.Println(reason)
	if err != nil {
		_, _ = goterm.Println(err)
	}
	goterm.Flush()
}

func (f *Field) revealSquare(square *Square) {
	square.state = Revealed
	delete(f.unknownSquares, square)
	f.revealedSquares[square] = nil
	for _, neighbor := range square.neighbors {
		delete(neighbor.unknownNeighbors, square)
		neighbor.revealedNeighbors[square] = nil
	}
}

func (f *Field) markSquare(square *Square) {
	square.state = Marked
	delete(f.unknownSquares, square)
	f.markedSquares[square] = nil
	for _, neighbor := range square.neighbors {
		delete(neighbor.unknownNeighbors, square)
		neighbor.markedNeighbors[square] = nil
	}
}
