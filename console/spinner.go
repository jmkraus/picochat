package console

import (
	"fmt"
	"time"
)

func StartSpinner(stop <-chan struct{}) {
	frames := []rune{'|', '/', '-', '\\'}
	i := 0
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
