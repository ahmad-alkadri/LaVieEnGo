package app

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func KeyboardWatch(keyErrorChan chan error, controlChan chan bool, exitChan chan bool) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			keyErrorChan <- err
			return
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
