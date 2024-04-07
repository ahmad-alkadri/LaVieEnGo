package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cell struct {
	X, Y int
}

func parseCoordinates(input string, MaxX, MaxY *int) map[Cell]bool {
	liveCells := make(map[Cell]bool)
	parts := strings.Split(input, ",")
	for _, part := range parts {
		coords := strings.Fields(strings.TrimSpace(part))
		if len(coords) != 2 {
			continue
		}
		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])
		if err1 != nil || err2 != nil || x < 1 || x > *MaxX || y < 1 || y > *MaxY {
			continue
		}
		liveCells[Cell{X: x, Y: y}] = true
	}
	return liveCells
}

func areMapsEqual(a, b map[Cell]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func readInitialCoordinates(MaxX, MaxY *int) map[Cell]bool {
	coordsFlag := flag.String("c", "", "Initial live cells coordinates (e.g., -c \"1 3, 3 4\")")
	flag.Parse()

	// Command line coordinates provided
	if *coordsFlag != "" {
		return parseCoordinates(*coordsFlag, MaxX, MaxY)
	}

	// Check if data is being piped into stdin
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == 0 {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		return parseCoordinates(string(input), MaxX, MaxY)
	}

	// Interactive mode: prompt for input
	fmt.Printf(
		"Enter live cell coordinates (x y), press Enter twice to start. Max X=%d, Max Y=%d:",
		MaxX, MaxY)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Read the first line of input
	return parseCoordinates(scanner.Text(), MaxX, MaxY)
}

func updateCells(liveCells map[Cell]bool, MaxX, MaxY *int) (map[Cell]bool, bool, bool) {
	nextGen := make(map[Cell]bool)
	candidateCells := make(map[Cell]int)

	var NeighborOffsets = []Cell{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for cell := range liveCells {
		neighborsCount := 0
		for _, offset := range NeighborOffsets {
			neighbor := Cell{X: cell.X + offset.X, Y: cell.Y + offset.Y}
			if liveCells[neighbor] {
				neighborsCount++
			} else {
				candidateCells[neighbor]++
			}
		}
		if neighborsCount == 2 || neighborsCount == 3 {
			nextGen[cell] = true
		}
	}

	for cell, count := range candidateCells {
		if count == 3 {
			nextGen[cell] = true
		}
	}

	// Check if any live cells are within the boundaries.
	// This determines if the game should continue.
	anyWithinBoundaries := false
	for cell := range nextGen {
		if cell.X >= 1 && cell.X <= *MaxX && cell.Y >= 1 && cell.Y <= *MaxY {
			anyWithinBoundaries = true
			break
		}
	}

	// Check for changes between generations
	changed := !areMapsEqual(liveCells, nextGen)

	return nextGen, anyWithinBoundaries, changed
}

func printBoard(liveCells map[Cell]bool, MaxX, MaxY *int) {
	fmt.Print("\033[H\033[2J\033[3J")
	for y := 1; y <= *MaxY; y++ {
		for x := 1; x <= *MaxX; x++ {
			if liveCells[Cell{X: x, Y: y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func main() {
	MaxX, MaxY := 60, 20

	liveCells := readInitialCoordinates(&MaxX, &MaxY)
	var changed, anyWithinBoundaries bool

	for {
		printBoard(liveCells, &MaxX, &MaxY)
		liveCells, anyWithinBoundaries, changed = updateCells(liveCells, &MaxX, &MaxY)

		// If there are no more changes or no live cells within the boundaries, stop the game.
		if !changed || !anyWithinBoundaries {
			time.Sleep(250 * time.Millisecond)
			printBoard(liveCells, &MaxX, &MaxY)
			if !changed {
				fmt.Println("No more changes, stopping the game.")
			} else {
				fmt.Println("No more live cells within the boundaries, stopping the game.")
			}
			break
		}

		// Wait a bit before the next iteration.
		time.Sleep(250 * time.Millisecond)
	}
}
