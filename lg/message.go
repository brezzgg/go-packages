package lg

import (
	"time"
)

type Message struct {
	Time    time.Time `json:"time"`
	Caller  Caller    `json:"caller"`
	Level   LogLevel  `json:"level"`
	Text    string    `json:"msg"`
	Context C         `json:"ctx,omitempty"`

	done chan struct{}
}
