package console

import (
	"fmt"
	"time"

	"picochat/config"
)

// StartSpinner starts a spinner animation until a signal is received on the stop channel.
// Parameters:
//
//	stop <-chan struct{} – placeholder for input channel
//
// Returns:
//
//	none
func StartSpinner(stop <-chan struct{}) {
	cfg, err := config.Get()
	if err != nil {
		//TODO: needs better err handling
		return
	}

	if cfg.Quiet {
		return
	}

	frames := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	i := 0

	// disable cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // safe enable cursor at return

	ClearLine()

	for {
		select {
		case <-stop:
			return // end routine
		default:
			fmt.Printf("\r%c", frames[i%len(frames)])
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// StopSpinner stops the spinner animation and clears the line.
// Parameters:
//
//	stop chan struct{} – placeholder for input channel to signal stop
//
// Returns:
//
//	none
func StopSpinner(stop chan struct{}) {
	cfg, err := config.Get()
	if err != nil {
		//TODO: needs better err handling
		return
	}

	if cfg.Quiet {
		return
	}

	close(stop)
	ClearLine()
}

// ClearLine clears the current terminal line.
// Parameters:
//
//	none
//
// Returns:
//
//	none
func ClearLine() {
	fmt.Print("\r\033[K") // delete to EOL
}
