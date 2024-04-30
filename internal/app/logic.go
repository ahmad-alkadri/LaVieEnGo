package app

import (
	model "LaVieEnGo/internal/model"
	"fmt"
	"strings"
)

func printBoard(liveCells map[model.Cell]bool, MaxX, MaxY *int) {
	// Clear all characters on screen
	fmt.Print("\033[H\033[2J\033[3J")
	var buffer strings.Builder

	for y := 1; y <= *MaxY; y++ {
		for x := 1; x <= *MaxX; x++ {
			if liveCells[model.Cell{X: x, Y: y}] {
				buffer.WriteRune('#')
			} else {
				buffer.WriteRune(' ')
			}
		}
		buffer.WriteRune('\n')
	}

	fmt.Print(buffer.String())
}

func UpdateCells(liveCells map[model.Cell]bool, MaxX, MaxY *int) (
	map[model.Cell]bool, bool, bool,
) {
	NeighborOffsets := []model.Cell{
		{X: 0, Y: 1},   // North
		{X: 1, Y: 1},   // Northeast
		{X: 1, Y: 0},   // East
		{X: 1, Y: -1},  // Southeast
		{X: 0, Y: -1},  // South
		{X: -1, Y: -1}, // Southwest
		{X: -1, Y: 0},  // West
		{X: -1, Y: 1},  // Northwest
	}

	nextGen := make(map[model.Cell]bool)
	candidateCells := make(map[model.Cell]int)

	for cell := range liveCells {
		neighborsCount := 0
		for _, offset := range NeighborOffsets {
			neighbor := model.Cell{X: cell.X + offset.X, Y: cell.Y + offset.Y}
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
	changed := !AreMapsEqual(liveCells, nextGen)

	return nextGen, anyWithinBoundaries, changed
}

// Helper function to check if two model.Cell maps are equal.
func AreMapsEqual(a, b map[model.Cell]bool) bool {
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
