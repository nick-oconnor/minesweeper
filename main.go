package main

import (
	"flag"
	"fmt"
	"github.com/buger/goterm"
	"github.com/nick-oconnor/minesweeper/sim"
	"github.com/nick-oconnor/minesweeper/solvers"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type GameResult struct {
	won       bool
	moveCount int
}

func playGame(width int, height int, numMines int, visualize bool, duration time.Duration) (bool, int) {
	field := sim.NewField(width, height)
	field.Init(numMines)
	algorithms := []solvers.Algorithm{solvers.NewPrimary(field), solvers.NewSecondary(field)}
	moveCount := 0
	if visualize {
		goterm.Clear()
	}
	for len(field.UnknownSquares()) > 0 {
		for _, algorithm := range algorithms {
			square, reason, err := algorithm.MakeNextMove()
			if visualize {
				field.Print(square, reason, err)
				time.Sleep(duration)
			}
			if err != nil {
				return false, moveCount
			}
			if square != nil {
				break
			}
		}
		moveCount++
	}
	return true, moveCount
}

func main() {
	width := flag.Int("width", 20, "width of the field")
	height := flag.Int("height", 20, "height of the field")
	numMines := flag.Int("mines", 50, "number of mines")
	numGames := flag.Int("games", 10000, "number of games")
	visualize := flag.Bool("visualize", false, "visualize gameplay")
	defaultMoveDuration, _ := time.ParseDuration("1s")
	duration := flag.Duration("duration", defaultMoveDuration, "visualize move duration")
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	if *visualize {
		playGame(*width, *height, *numMines, true, *duration)
	} else {
		parallelism := runtime.NumCPU()
		gameQueue := make(chan int)
		resultQueue := make(chan *GameResult)
		workerWg := sync.WaitGroup{}
		workerWg.Add(parallelism)
		for i := 0; i < parallelism; i++ {
			go func() {
				defer workerWg.Done()
				for {
					_, ok := <-gameQueue
					if !ok {
						return
					}
					won, moveCount := playGame(*width, *height, *numMines, false, *duration)
					resultQueue <- &GameResult{won, moveCount}
				}
			}()
		}
		resultWg := sync.WaitGroup{}
		resultWg.Add(1)
		go func() {
			defer func() {
				fmt.Println()
				resultWg.Done()
			}()
			numResults := 0
			numWon := 0
			totalMoves := 0
			for {
				result, ok := <-resultQueue
				if !ok {
					return
				}
				numResults++
				if result.won {
					numWon++
					totalMoves += result.moveCount
				}
				fmt.Printf("\rGames Played: %d, Percent Won: %.2f, Average Moves: %.2f", numResults, float64(numWon)/float64(numResults)*100, float64(totalMoves)/float64(numWon))
			}
		}()
		for i := 0; i < *numGames; i++ {
			gameQueue <- i
		}
		close(gameQueue)
		workerWg.Wait()
		close(resultQueue)
		resultWg.Wait()
	}
}
