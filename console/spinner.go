package console

import (
	"fmt"
	"time"

	"picochat/config"
)

func StartSpinner(stop <-chan struct{}) {
	cfg := config.Get()
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

func StopSpinner(stop chan struct{}) {
	cfg := config.Get()
	if cfg.Quiet {
		return
	}

	close(stop)
	ClearLine()
}

func ClearLine() {
	fmt.Print("\r\033[K") // delete to EOL
}
