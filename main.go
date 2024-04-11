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
	// Prepare the mapping of live cells
	liveCells := make(map[Cell]bool)

	// Expect input like `x1 y1, x2 y2, x3 y3`. Thus, split first
	// on comma and parse one by one afterwards.
	parts := strings.Split(input, ",")
	for _, part := range parts {
		coords := strings.Fields(strings.TrimSpace(part))
		if len(coords) != 2 {
			continue
		}
		// Try parsing the string coordinate to integer
		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])

		// Verify that there's no error and coordinates within boundaries
		if err1 != nil || err2 != nil || x < 1 || x > *MaxX || y < 1 || y > *MaxY {
			continue
		}

		// If all's good add the live Cell to the map
		liveCells[Cell{X: x, Y: y}] = true
	}
	return liveCells
}

func readInitialCoordinates(MaxX, MaxY *int) map[Cell]bool {
	coordsFlag := flag.String(
		"c", "", "Initial live cells coordinates (e.g., -c \"1 3, 3 4\")")
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
		"Enter live cell coordinates (x y), and then press Enter twice to start. Max X=%d, Max Y=%d:",
		MaxX, MaxY)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Read the first line of input
	return parseCoordinates(scanner.Text(), MaxX, MaxY)
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
		// If we change y, we change line
		fmt.Println()
	}
}

func updateCells(liveCells map[Cell]bool, MaxX, MaxY *int) (
	map[Cell]bool, bool, bool,
) {
	NeighborOffsets := []Cell{
		{0, 1},   // North
		{1, 1},   // Northeast
		{1, 0},   // East
		{1, -1},  // Southeast
		{0, -1},  // South
		{-1, -1}, // Southwest
		{-1, 0},  // West
		{-1, 1},  // Northwest
	}

	nextGen := make(map[Cell]bool)
	candidateCells := make(map[Cell]int)

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

	// Check if the world has become stagnant between generations.
	// This determines if the game should continue.
	changed := !areMapsEqual(liveCells, nextGen)

	return nextGen, anyWithinBoundaries, changed
}

// Helper function to check if two Cell maps are equal.
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

func main() {
	MaxX, MaxY := 60, 20

	liveCells := readInitialCoordinates(&MaxX, &MaxY)
	var changed, anyWithinBoundaries bool

	for {
		printBoard(liveCells, &MaxX, &MaxY)
		liveCells, anyWithinBoundaries, changed = updateCells(liveCells, &MaxX, &MaxY)

		// Pause the time a bit for visibility
		time.Sleep(250 * time.Millisecond)

		// If there are no more changes or no live cells within the boundaries, stop the game.
		if !changed || !anyWithinBoundaries {
			printBoard(liveCells, &MaxX, &MaxY)
			if !changed {
				fmt.Println("No more changes, stopping the game.")
			} else {
				fmt.Println("No more live cells within the boundaries, stopping the game.")
			}
			break
		}
	}
}
