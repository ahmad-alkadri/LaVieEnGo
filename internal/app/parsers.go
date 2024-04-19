package app

import (
	model "LaVieEnGo/internal/model"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ParseCoordinates(input string, MaxX, MaxY *int) map[model.Cell]bool {
	// Prepare the mapping of live cells
	liveCells := make(map[model.Cell]bool)

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
		liveCells[model.Cell{X: x, Y: y}] = true
	}
	return liveCells
}

func ReadInitialCoordinates(MaxX, MaxY *int) map[model.Cell]bool {
	coordsFlag := flag.String(
		"c", "", "Initial live cells coordinates (e.g., -c \"1 3, 3 4\")")
	flag.Parse()

	// Command line coordinates provided
	if *coordsFlag != "" {
		return ParseCoordinates(*coordsFlag, MaxX, MaxY)
	}

	// Check if data is being piped into stdin
	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == 0 {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		return ParseCoordinates(string(input), MaxX, MaxY)
	}

	// Interactive mode: prompt for input
	fmt.Printf(
		"Enter live cell coordinates (x y), and then press Enter twice to start. Max X=%d, Max Y=%d:",
		*MaxX, *MaxY)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Read the first line of input
	return ParseCoordinates(scanner.Text(), MaxX, MaxY)
}
