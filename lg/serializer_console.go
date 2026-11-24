package lg

import (
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type (
	consoleSerializer struct {
		disableColors bool
		indentContext bool
		cachedOffset  int
		cachedTime    time.Time
	}

	ConsoleSerializerOption func(*consoleSerializer)
)

func NewConsoleSerializer(opts ...ConsoleSerializerOption) *consoleSerializer {
	s := &consoleSerializer{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithDisabledColors() ConsoleSerializerOption {
	return func(s *consoleSerializer) {
		s.disableColors = true
	}
}

func WithIndentContext() ConsoleSerializerOption {
	return func(s *consoleSerializer) {
		s.indentContext = true
	}
}

func (c *consoleSerializer) Serialize(m Message) string {
	context := ""
	if m.Context != nil {
		var (
			b   []byte
			err error
		)
		if c.indentContext {
			b, err = json.MarshalIndent(m.Context, "", "  ")
		} else {
			b, err = json.Marshal(m.Context)
		}
		str := string(b)
		if err == nil && len(str) > 0 && str != "{}" {
			context = str
		}
	}

	var caller, level string
	levelSpaces := ""

	if form := m.Level.Formatter(); len(form) < 5 {
		for len(form)+len(levelSpaces) < 5 {
			levelSpaces += " "
		}
	}
	caller = c.getCaller(m)
	level = m.Level.Formatter()
	if !c.disableColors {
		caller = ColorizeString(caller, ClrFgCyan)
		level = ColorizeString(m.Level.Formatter(), m.Level.Color())
	}
	level += levelSpaces

	return c.getTime(m) + level + "  " + caller + m.Text + " " + context
}

func (c *consoleSerializer) getTime(m Message) string {
	const layout = time.RFC3339

	res, ok := m.Level.HandleOptions(func(opt LogLevelOption) (string, bool) {
		switch opt {
		case LevelOptionDisableTime:
			return "", true
		case LevelOptionTimeDisableOffset:
			return m.Time.UTC().Format(layout) + "  ", true
		}
		return "", false
	})

	if ok {
		return res
	}

	offset := c.getTimeOffset(time.Now())
	offsetStr := strconv.Itoa(offset / 3600)
	if !strings.HasPrefix(offsetStr, "-") {
		offsetStr = "+" + offsetStr
	}
	return m.Time.UTC().Format(layout) + offsetStr + "  "
}

func (c *consoleSerializer) getCaller(m Message) string {
	res, ok := m.Level.HandleOptions(func(opt LogLevelOption) (string, bool) {
		switch opt {
		case LevelOptionDisableCaller:
			return "", true
		case LevelOptionCallerOnlyFunc:
			return m.Caller.Method + "()  ", true
		case LevelOptionCallerOnlyFile:
			return m.Caller.File + "  ", true
		case LevelOptionCallerDisableFunc:
			return m.Caller.File + ":" + strconv.Itoa(m.Caller.Line) + "  ", true
		case LevelOptionCallerDisableFile:
			return m.Caller.Method + "():" + strconv.Itoa(m.Caller.Line) + "  ", true
		case LevelOptionCallerDisableLine:
			return m.Caller.File + "." + m.Caller.Method + "()  ", true
		}
		return "", false
	})

	if !ok {
		res = m.Caller.Formatter()
	}

	if m.Job != "" {
		res = m.Job + "*"
	}

	return res + "  "
}

func (c *consoleSerializer) getTimeOffset(t time.Time) int {
	if c.cachedTime.YearDay() != t.YearDay() {
		_, c.cachedOffset = t.Zone()
		c.cachedTime = t
	}
	return c.cachedOffset
}
