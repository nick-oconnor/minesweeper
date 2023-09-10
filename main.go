package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"gitlab.ocnr.org/apps/minesweeper/field"
	"gitlab.ocnr.org/apps/minesweeper/solver"
)

// simulate simulates the number of requested games in parallel, one game per core.
func simulate(width int, height int, mineCount int, gameCount int, progress bool) {
	gameQueue := make(chan interface{}, gameCount)
	resultQueue := make(chan *solver.GameResult)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for range gameQueue {
				resultQueue <- solver.NewSolver(field.NewField(width, height, mineCount), false).Solve()
			}
		}()
	}
	wonCount := 0
	moveCount := 0
	guessCount := 0
	for i := 0; i < gameCount; i++ {
		gameQueue <- nil
	}
	close(gameQueue)
	for gamesSimulated := 1; gamesSimulated <= gameCount; gamesSimulated++ {
		result := <-resultQueue
		if result.Won {
			wonCount++
			moveCount += result.MoveCount
			guessCount += result.GuessCount
		}
		if progress {
			fmt.Printf("Games Simulated: %d, Win Ratio: %.1f%%, Moves/Win: %.1f, Guesses/Win: %.2f\r", gamesSimulated, float64(wonCount)/float64(gamesSimulated)*100, float64(moveCount)/float64(wonCount), float64(guessCount)/float64(wonCount))
		}
	}
	fmt.Printf("Games Simulated: %d, Win Ratio: %.1f%%, Moves/Win: %.1f, Guesses/Win: %.2f\n", gameCount, float64(wonCount)/float64(gameCount)*100, float64(moveCount)/float64(wonCount), float64(guessCount)/float64(wonCount))
}

// main is the CLI entrypoint.
func main() {
	width := flag.Int("width", 30, "width of the field")
	height := flag.Int("height", 16, "height of the field")
	mineCount := flag.Int("mines", 99, "number of mines")
	gameCount := flag.Int("games", 1000, "number of games")
	progress := flag.Bool("progress", true, "show progress")
	visualize := flag.Bool("visualize", false, "visualize gameplay")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}
	if *visualize {
		_ = solver.NewSolver(field.NewField(*width, *height, *mineCount), *visualize).Solve()
		return
	}
	simulate(*width, *height, *mineCount, *gameCount, *progress)
}
