package app

import (
	model "LaVieEnGo/internal/model"
	"fmt"
	"time"
)

func Game(
	initialCells map[model.Cell]bool,
	controlChan chan bool,
	exitChan chan bool,
	gameOverChan chan bool,
	MaxX *int,
	MaxY *int,
) {
	paused := false

	liveCells := initialCells
	PrintBoard(liveCells, MaxX, MaxY)
	var changed, anyWithinBoundaries bool

	for {
		select {
		case <-exitChan:
			return
		case <-controlChan:
			paused = !paused
			if paused {
				fmt.Println("Game paused.")
			}
		default:
			if !paused {
				PrintBoard(liveCells, MaxX, MaxY)
				liveCells, anyWithinBoundaries, changed = UpdateCells(liveCells, MaxX, MaxY)
				fmt.Println("[P] Pause the game. [Ctrl+C] Exit the game.")

				// Pause the time a bit for visibility
				time.Sleep(250 * time.Millisecond)

				// If there are no more changes or no live cells within the boundaries, stop the game.
				if !changed || !anyWithinBoundaries {
					PrintBoard(liveCells, MaxX, MaxY)
					if !changed {
						fmt.Println("No more changes, stopping the game.")
					} else {
						fmt.Println("No more live cells within the boundaries, stopping the game.")
					}
					gameOverChan <- true
					return
				}
			}
		}
	}
}
