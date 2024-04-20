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
	stepChan chan bool,
	MaxX *int,
	MaxY *int,
) {
	paused := false
	liveCells := initialCells
	printBoard(liveCells, MaxX, MaxY)
	var anyWithinBoundaries, changed bool

	for {
		select {
		case <-exitChan:
			return
		case <-controlChan:
			paused = !paused
			if paused {
				fmt.Println("Game paused. [Right Arrow] Move forward a step.")
			}
		case <-stepChan:
			if paused {
				liveCells, anyWithinBoundaries, changed = updateWorld(
					liveCells, MaxX, MaxY, gameOverChan, paused)
				if !changed || !anyWithinBoundaries {
					return
				}
			}
		default:
			if !paused {
				liveCells, anyWithinBoundaries, changed = updateWorld(
					liveCells, MaxX, MaxY, gameOverChan, paused)
				if !changed || !anyWithinBoundaries {
					return
				}
			}
		}
	}
}

func updateWorld(
	liveCells map[model.Cell]bool,
	MaxX *int, MaxY *int,
	gameOverChan chan bool,
	paused bool,
) (map[model.Cell]bool, bool, bool) {
	printBoard(liveCells, MaxX, MaxY)
	liveCells, anyWithinBoundaries, changed := UpdateCells(liveCells, MaxX, MaxY)
	fmt.Println("[Space] Pause/Resume the game. [Ctrl+C] Exit the game.")
	if paused {
		fmt.Println("Game paused. [Right Arrow] Move forward a step.")
	}

	// Pause the time a bit for visibility
	time.Sleep(150 * time.Millisecond)

	// If there are no more changes or no live cells within the boundaries, stop the game.
	if !changed || !anyWithinBoundaries {
		printBoard(liveCells, MaxX, MaxY)
		if !changed {
			fmt.Println("No more changes, stopping the game.")
		} else {
			fmt.Println("No more live cells within the boundaries, stopping the game.")
		}
		gameOverChan <- true
		return nil, false, false // Return nil map and false to indicate game over.
	}

	return liveCells, anyWithinBoundaries, changed
}
