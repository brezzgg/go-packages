package lg

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unicode"

	"github.com/brezzgg/go-packages/schan"
)

var endTasks = &EndTasks{}

type (
	Logger struct {
		typeConv TypeConverter
		pipes    []*Pipe

		level           uint16
		splitFilePrefix bool

		queue     *schan.SafeChannel[Message]
		wg        sync.WaitGroup
		stop      chan struct{}
		closed    atomic.Bool
		levelOpts []LogLevelOption
	}

	LoggerOption func(*Logger)
)

func NewLogger(opts ...LoggerOption) *Logger {
	l := &Logger{
		pipes:    make([]*Pipe, 0, 1),
		typeConv: NewDefaultTypeConverter(),
		queue:    schan.New[Message](8192),
		stop:     make(chan struct{}),
	}

	for _, opt := range opts {
		opt(l)
	}

	if len(l.pipes) == 0 {
		l.pipes = append(l.pipes, NewPipe())
	}

	l.wg.Add(1)
	go l.worker()
	go l.watchSignals()

	l.closed.Store(false)

	return l
}

func WithPipe(pipe *Pipe) LoggerOption {
	return func(l *Logger) {
		l.pipes = append(l.pipes, pipe)
	}
}

func WithPipes(pipes ...*Pipe) LoggerOption {
	return func(l *Logger) {
		l.pipes = append(l.pipes, pipes...)
	}
}

func WithCustomTypeConverter(c TypeConverter) LoggerOption {
	return func(l *Logger) {
		l.typeConv = c
	}
}

func WithQueueSize(size int) LoggerOption {
	return func(l *Logger) {
		l.queue = schan.New[Message](size)
	}
}

func WithConstantLevelOptions(opts ...LogLevelOption) LoggerOption {
	return func(l *Logger) {
		l.levelOpts = opts
	}
}

func WithSplitFilePrefix(split bool) LoggerOption {
	return func(l *Logger) {
		l.splitFilePrefix = split
	}
}

func WithLogLevel(level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.level = level.priority
	}
}

func WithLogLevelPriority(p uint16) LoggerOption {
	return func(l *Logger) {
		l.level = p
	}
}

func (l *Logger) Close() {
	if !l.closed.Swap(true) {
		l.queue.Close()
		close(l.stop)
		l.wg.Wait()
		for _, pipe := range l.pipes {
			pipe.Close()
		}
	}
}

func (l *Logger) worker() {
WorkerLoop:
	for {
		select {
		case <-l.stop:
			for msg := range l.queue.Chan() {
				l.handleMessages(msg)
			}
			break WorkerLoop
		case msg := <-l.queue.Chan():
			l.handleMessages(msg)
		}
	}

	l.wg.Done()
}

func (l *Logger) handleMessages(messages ...Message) {
	for _, msg := range messages {
		if msg.Text == "" {
			continue
		}
		for _, p := range l.pipes {
			p.Handle(msg)
		}
		if msg.done != nil {
			close(msg.done)
		}
	}
}

func (l *Logger) watchSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	<-sigCh
	l.Close()
}

func (l *Logger) Handle(args []any, level LogLevel) {
	if level.priority < l.level {
		return
	}

	offset := 0
	if isBadOffset(2) {
		offset = 1
	}

	msg := l.buildMessage(args, level, offset)
	if msg.done != nil {
		msg.done = make(chan struct{})
		l.handleMessages(msg)
		<-msg.done
	} else {
		l.queue.Send(msg)
	}
}

func (l *Logger) buildMessage(m []any, level LogLevel, callerOffset int) Message {
	level.opts = append(l.levelOpts, level.opts...)

	msg := Message{
		Caller: GetCallerInfo(3+callerOffset, l.splitFilePrefix),
		Level:  level,
		Time:   time.Now(),
	}

	args := make([]any, 0, len(m))
	for _, item := range m {
		if _, ok := item.(Sync); ok {
			if msg.done == nil {
				msg.done = make(chan struct{})
			}
		} else {
			args = append(args, item)
		}
	}

	if err := l.CheckArgs(m, &msg.Text, &msg.Context); err != nil {
		msg = Message{
			Caller:  msg.Caller,
			Level:   logLevelLoggerError,
			Text:    "failed to parse message context",
			Context: C{"error": err.Error()},
			Time:    msg.Time,
		}
	}

	return msg
}

func (l *Logger) CheckArgs(args []any, msgBody *string, msgCtx *C) error {
	appendBodyFunc := func(body string) {
		if len(*msgBody) > 0 {
			r := []rune((*msgBody))
			if unicode.IsUpper(r[0]) && unicode.IsLower(r[1]) {
				*msgBody = string(unicode.ToLower(r[0])) + (*msgBody)[1:]
			}
		}
		if len(*msgBody) == 0 {
			*msgBody = body
		} else {
			*msgBody = strings.TrimSuffix(
				strings.TrimSpace(body), ":",
			) + ": " + strings.TrimSpace(*msgBody)
		}
	}

	if len(args) == 0 {
		return errors.New("no arguments")
	}

	l.typeConv.ConvAndPushBody(args[0], appendBodyFunc)
	if len(args) == 1 {
		return nil
	}

	if *msgCtx == nil {
		*msgCtx = make(C)
	}
	for itemNum, item := range args {
		if itemNum == 0 {
			continue
		}
		l.typeConv.ConvAndPushContext(item, func(k string, v any) {
			(*msgCtx)[l.findKey(k, *msgCtx)] = v
		})
	}

	return nil
}

func (*Logger) findKey(k string, m C) string {
	if _, ok := m[k]; !ok {
		return k
	}
	for i := 2; i < 100; i++ {
		keyStr := k + strconv.Itoa(i)
		if _, ok := m[keyStr]; !ok {
			return keyStr
		}
	}
	panic(fmt.Sprintf("too many duplicates of key '%s'", k))
}
