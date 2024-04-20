package app

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func Controller(
	keyErrorChan chan error,
	controlChan chan bool,
	exitChan chan bool,
	stepChan chan bool,
) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			keyErrorChan <- err
			return
		}

		if key == keyboard.KeySpace {
			controlChan <- true // Toggle pause/resume
		}

		if key == keyboard.KeyArrowRight {
			stepChan <- true // Move forward one step
		}

		if char == 'p' || char == 'P' {
			controlChan <- true // Toggle pause/resume
		}

		if key == keyboard.KeyCtrlC {
			fmt.Println("\nExiting...")
			exitChan <- true
			return
		}
	}
}
