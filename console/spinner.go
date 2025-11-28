package console

import (
	"fmt"
	"time"
)

// StartSpinner starts a spinner animation until a signal is received on the stop channel.
// Parameters:
//
//	quiet (bool)           - Suppress Spinner in quiet mode (true / false)
//	stop (<-chan struct{}) – placeholder for input channel
//
// Returns:
//
//	none
func StartSpinner(quiet bool, stop <-chan struct{}) {
	if quiet {
		return
	}

	frames := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	i := 0

	// disable cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // safe enable cursor at return

	fmt.Print("\r\033[K") // delete to EOL

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
//	quiet (bool)         - Suppress Spinner in quiet mode (true / false)
//	stop (chan struct{}) – placeholder for input channel to signal stop
//
// Returns:
//
//	none
func StopSpinner(quiet bool, stop chan struct{}) {
	if quiet {
		return
	}

	select {
	case <-stop:
		return //channel already closed, do nothing
	default:
		close(stop)
		fmt.Print("\r\033[K") // delete to EOL
	}
}
