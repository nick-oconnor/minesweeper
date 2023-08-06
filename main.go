package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	. "gitlab.ocnr.org/apps/minesweeper/field"
	. "gitlab.ocnr.org/apps/minesweeper/solver"
)

// GameResult is used for keeping game results.
type GameResult struct {
	won                   bool
	moveCount, guessCount int
}

// simulate simulates the number of requested games in parallel, one game per core.
func simulate(width int, height int, mineCount int, gameCount int) {
	gameQueue := make(chan interface{}, gameCount)
	resultQueue := make(chan *GameResult)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for range gameQueue {
				moveCount, guessCount, err := NewSolver(NewField(width, height, mineCount)).Solve()
				resultQueue <- &GameResult{err == nil, moveCount, guessCount}
			}
		}()
	}
	wonCount := 0
	moveCount := 0
	guessCount := 0
	startCycles := cpuCycles()
	for i := 0; i < gameCount; i++ {
		gameQueue <- nil
	}
	close(gameQueue)
	for gamesSimulated := 1; gamesSimulated <= gameCount; gamesSimulated++ {
		result := <-resultQueue
		if result.won {
			wonCount++
			moveCount += result.moveCount
			guessCount += result.guessCount
		}
		fmt.Printf("\rGames Simulated: %d, Won: %.1f%%, Moves/Win: %.1f, Guesses/Win: %.2f, CPU Cycles/Game: %.2e", gamesSimulated, float64(wonCount)/float64(gamesSimulated)*100, float64(moveCount)/float64(wonCount), float64(guessCount)/float64(wonCount), float64(cpuCycles()-startCycles)/float64(gamesSimulated))
	}
	fmt.Println()
}

// main is the CLI entrypoint.
func main() {
	defaultMoveDuration, _ := time.ParseDuration("0.5s")
	width := flag.Int("width", 30, "width of the field")
	height := flag.Int("height", 16, "height of the field")
	mineCount := flag.Int("mines", 99, "number of mines")
	gameCount := flag.Int("games", 1000, "number of games")
	duration := flag.Duration("duration", defaultMoveDuration, "visualize move duration")
	visualize := flag.Bool("visualize", false, "visualize gameplay")
	flag.Parse()
	if *visualize {
		_, _, _ = NewSolver(NewField(*width, *height, *mineCount)).WithFieldPrinter(NewPrinter(*duration).PrintField).Solve()
		return
	}
	simulate(*width, *height, *mineCount, *gameCount)
}
