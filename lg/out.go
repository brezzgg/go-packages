package lg

import (
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

type (
	LogLevel struct {
		Level    string
		color    string
		opts     []LogLevelOption
		priority uint16
	}
	LogLevelOption string
)

func NewLogLevel(asciiClr string, level string, opts ...LogLevelOption) LogLevel {
	return LogLevel{
		Level:    level,
		color:    asciiClr,
		opts:     opts,
		priority: LevelDefaultPriority,
	}
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Level)
}

func (l *LogLevel) UnmarshalJSON(b []byte) error {
	var levelStr string
	if err := json.Unmarshal(b, &levelStr); err != nil {
		return err
	}
	l.Level = levelStr
	return nil
}

func (l LogLevel) WithOptions(opts ...LogLevelOption) LogLevel {
	l.opts = append(l.opts, opts...)
	return l
}

func (l LogLevel) WithPriority(priority uint16) LogLevel {
	l.priority = priority
	return l
}

/*
LogGlobal function is alias for GlobalLogger.Log
*/
func (l LogLevel) LogGlobal(args ...any) {
	if GlobalLogger == nil {
		return
	}
	GlobalLogger.Log(l, args...)
}

/*
Log function is alias for Logger.Log
*/
func (l LogLevel) Log(logger *Logger, args ...any) {
	if logger == nil {
		return
	}
	logger.Log(l, args...)
}

func (l LogLevel) HandleOptions(f func(LogLevelOption) (string, bool)) (string, bool) {
	for _, opt := range l.opts {
		str, ok := f(opt)
		return str, ok
	}
	return "", false
}

func (l LogLevel) Color() string {
	return l.color
}

func (l LogLevel) Formatter() string {
	return l.Level
}

func (l LogLevel) Equal(other LogLevel) bool {
	return l.Level == other.Level
}

type Caller struct {
	Method string `json:"method"`
	File   string `json:"file"`
	Line   int    `json:"line"`
}

func (c Caller) Formatter() string {
	var sb strings.Builder
	sb.Grow(len(c.File) + len(c.Method))
	sb.WriteString(c.File)
	sb.WriteString(".")
	sb.WriteString(c.Method)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(c.Line))
	return sb.String()
}

func (c Caller) Equal(other Caller) bool {
	return c.Line == other.Line &&
		c.Method == other.Method &&
		c.File == other.File
}

func ColorizeString(str string, ascii string) string {
	var sb strings.Builder
	sb.Grow(len(str) + len(ascii) + len(ClrReset))
	sb.WriteString(ascii)
	sb.WriteString(str)
	sb.WriteString(ClrReset)
	return sb.String()
}
