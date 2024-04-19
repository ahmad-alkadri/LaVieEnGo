package main

import (
	"LaVieEnGo/internal/app"
	"fmt"
	"os"
	"os/signal"

	"github.com/eiannone/keyboard"
)

func main() {
	MaxX, MaxY := 60, 20
	initialCells := app.ReadInitialCoordinates(&MaxX, &MaxY)

	if err := keyboard.Open(); err != nil {
		fmt.Println("Failed to open keyboard:", err)
		return
	}
	defer func() {
		_ = keyboard.Close()
	}()

	controlChan := make(chan bool)
	exitChan := make(chan bool)
	gameOverChan := make(chan bool)
	keyErrorChan := make(chan error)

	// Goroutine for the game
	go app.Game(initialCells, controlChan, exitChan, gameOverChan, &MaxX, &MaxY)

	// Goroutine to read keyboard inputs
	go app.KeyboardWatch(keyErrorChan, controlChan, exitChan)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt) // Listen for Ctrl+C signal

	for {
		select {
		case <-c: // Handle Ctrl+C
			fmt.Println("\nExiting...")
			close(exitChan)
			return
		case <-gameOverChan:
			return
		case <-exitChan:
			return
		case err := <-keyErrorChan:
			fmt.Println("Error reading key:", err)
			os.Exit(1)
		}
	}
}
