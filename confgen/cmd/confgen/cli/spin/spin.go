package spin

import (
	"fmt"
	"time"
)

var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func Spin(msg string, result <-chan string) <-chan struct{} {
	doneCh := make(chan struct{})
	go spin(msg, result, doneCh)
	return doneCh
}

func spin(msg string, result <-chan string, doneCh chan struct{}) {
	i := 0
	for {
		select {
		case r := <-result:
			fmt.Printf("\r\033[K")
			fmt.Println(r)
			doneCh <- struct{}{}
			return
		default:
			fmt.Printf("\r%s %s", frames[i%len(frames)], msg)
			i++
			time.Sleep(100 * time.Millisecond)
		}
	}
}
