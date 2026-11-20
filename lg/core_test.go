package lg

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

type testCapture struct {
	r *os.File
	w *os.File
}

func newTestLogger(t *testing.T, ser Serializer, opts ...LoggerOption) (*Logger, *testCapture) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	base := []LoggerOption{WithPipe(NewPipe(
		WithSerializer(ser),
		WithWriter(NewConsoleWriter(WithCustomStdout(w))),
	)), WithLogLevel(LogLevelDebug)}
	logger := NewLogger(append(base, opts...)...)
	return logger, &testCapture{r: r, w: w}
}

func (c *testCapture) read(t *testing.T, logger *Logger) string {
	t.Helper()
	logger.Close()
	_ = c.w.Close()
	b, err := io.ReadAll(c.r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	_ = c.r.Close()
	return string(b)
}

func decodeJSONLine(t *testing.T, out string) map[string]any {
	t.Helper()
	line := strings.TrimSpace(out)
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("json unmarshal: %v\n--- output ---\n%s", err, out)
	}
	return m
}

func jsonLevel(m map[string]any) string {
	s, _ := m["level"].(string)
	return strings.ToLower(s)
}

func jsonMsg(m map[string]any) string {
	s, _ := m["msg"].(string)
	return s
}

func jsonCtx(m map[string]any) map[string]any {
	c, _ := m["ctx"].(map[string]any)
	return c
}

func TestCoreBasicLogging(t *testing.T) {
	l, c := newTestLogger(t, NewConsoleSerializer(WithDisabledColors()))
	l.Info("hello", Sync{}, C{"a": "b"})
	out := c.read(t, l)
	for _, want := range []string{"Info", "hello", `"a":"b"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q\n%s", want, out)
		}
	}
}

func TestCoreFormattedLogging(t *testing.T) {
	l, c := newTestLogger(t, NewJSONSerializer())
	l.Infof("hello %s #%d", "world", 7, Sync{}, C{"k": "v"})
	out := c.read(t, l)
	m := decodeJSONLine(t, out)
	if got := jsonMsg(m); got != "hello world #7" {
		t.Fatalf("msg = %q, want %q", got, "hello world #7")
	}
	if got := jsonLevel(m); got != "info" {
		t.Fatalf("level = %q, want info", got)
	}
	if got, _ := jsonCtx(m)["k"].(string); got != "v" {
		t.Fatalf("ctx[k] = %v, want v", jsonCtx(m)["k"])
	}
}

func TestCoreLogAndLogfCustomLevel(t *testing.T) {
	l, c := newTestLogger(t, NewJSONSerializer())
	l.Log(LogLevelWarn, "plain", Sync{})
	l.Logf(LogLevelError, "fmt-%s", "ok", Sync{})
	out := c.read(t, l)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d:\n%s", len(lines), out)
	}
	m1 := decodeJSONLine(t, lines[0])
	m2 := decodeJSONLine(t, lines[1])
	if jsonLevel(m1) != "warn" || jsonMsg(m1) != "plain" {
		t.Fatalf("line1 = level:%q msg:%q", jsonLevel(m1), jsonMsg(m1))
	}
	if jsonLevel(m2) != "error" || jsonMsg(m2) != "fmt-ok" {
		t.Fatalf("line2 = level:%q msg:%q", jsonLevel(m2), jsonMsg(m2))
	}
}

func TestCoreLevelFiltering(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	l := NewLogger(
		WithPipe(NewPipe(WithSerializer(NewJSONSerializer()), WithWriter(NewConsoleWriter(WithCustomStdout(w))))),
		WithLogLevel(LogLevelWarn),
	)
	l.Debug("d", Sync{})
	l.Info("i", Sync{})
	l.Warn("w", Sync{})
	l.Error("e", Sync{})
	l.Close()
	_ = w.Close()
	b, _ := io.ReadAll(r)
	_ = r.Close()
	out := string(b)
	if strings.Contains(out, `"msg":"d"`) || strings.Contains(out, `"msg":"i"`) {
		t.Fatalf("filtered messages should not appear:\n%s", out)
	}
	if !strings.Contains(out, `"msg":"w"`) || !strings.Contains(out, `"msg":"e"`) {
		t.Fatalf("expected warn/error messages:\n%s", out)
	}
}

func TestCoreContextAndErrors(t *testing.T) {
	l, c := newTestLogger(t, NewJSONSerializer())
	l.Error("failed", Sync{}, errors.New("boom"), C{"x": 1})
	out := c.read(t, l)
	m := decodeJSONLine(t, out)
	ctx := jsonCtx(m)
	if got, _ := ctx["error"].(string); got != "boom" {
		t.Fatalf("ctx[error] = %v, want boom", ctx["error"])
	}
	if _, ok := ctx["x"]; !ok {
		t.Fatalf("ctx missing x: %v", ctx)
	}
}

func TestCoreDuplicateContextKeys(t *testing.T) {
	l, c := newTestLogger(t, NewJSONSerializer())
	l.Info("msg", Sync{}, C{"k": "v1"}, C{"k": "v2"})
	out := c.read(t, l)
	m := decodeJSONLine(t, out)
	ctx := jsonCtx(m)
	if got, _ := ctx["k"].(string); got != "v1" {
		t.Fatalf("ctx[k] = %v, want v1", ctx["k"])
	}
	if got, _ := ctx["k2"].(string); got != "v2" {
		t.Fatalf("ctx[k2] = %v, want v2", ctx["k2"])
	}
}

func TestCoreConsoleSerializer(t *testing.T) {
	l, c := newTestLogger(t, NewConsoleSerializer(WithDisabledColors()))
	l.Warn("console-msg", Sync{}, C{"k": "v"})
	out := c.read(t, l)
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("expected no ANSI colors:\n%q", out)
	}
	for _, want := range []string{"Warn", "console-msg", `"k":"v"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q\n%s", want, out)
		}
	}
}

func TestCoreJSONSerializerIndent(t *testing.T) {
	l, c := newTestLogger(t, NewJSONSerializer(WithIndent()))
	l.Info("indent", Sync{})
	out := c.read(t, l)
	if !strings.Contains(out, "\n  \"msg\"") && !strings.Contains(out, "\n  \"level\"") {
		t.Fatalf("expected indented JSON:\n%s", out)
	}
	m := decodeJSONLine(t, out)
	if jsonMsg(m) != "indent" || jsonLevel(m) != "info" {
		t.Fatalf("unexpected decoded msg: %+v", m)
	}
}

func TestCoreMultiPipe(t *testing.T) {
	r1, w1, _ := os.Pipe()
	r2, w2, _ := os.Pipe()
	l := NewLogger(
		WithPipes(
			NewPipe(WithSerializer(NewJSONSerializer()), WithWriter(NewConsoleWriter(WithCustomStdout(w1)))),
			NewPipe(WithSerializer(NewConsoleSerializer(WithDisabledColors())), WithWriter(NewConsoleWriter(WithCustomStdout(w2)))),
		),
		WithLogLevel(LogLevelDebug),
	)
	l.Info("fanout", Sync{}, C{"k": "v"})
	l.Close()
	_ = w1.Close()
	_ = w2.Close()
	b1, _ := io.ReadAll(r1)
	b2, _ := io.ReadAll(r2)
	_ = r1.Close()
	_ = r2.Close()
	m := decodeJSONLine(t, string(b1))
	if jsonMsg(m) != "fanout" || jsonLevel(m) != "info" {
		t.Fatalf("json pipe unexpected: %s", string(b1))
	}
	out2 := string(b2)
	for _, want := range []string{"Info", "fanout", `"k":"v"`} {
		if !strings.Contains(out2, want) {
			t.Fatalf("console pipe missing %q\n%s", want, out2)
		}
	}
}

func TestCoreFileWriter(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "lg-file-*")
	if err != nil {
		t.Fatal(err)
	}
	name := tmp.Name()
	_ = tmp.Close()

	l := NewLogger(WithPipe(NewPipe(
		WithSerializer(NewJSONSerializer()),
		WithWriter(NewFileWriter(name)),
	)), WithLogLevel(LogLevelDebug))
	l.Info("to-file", Sync{}, C{"k": "v"})
	l.Close()

	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	content := strings.TrimSpace(string(b))
	line := content
	if idx := strings.Index(content, "\n"); idx >= 0 {
		line = strings.TrimSpace(content[:idx])
	}
	m := decodeJSONLine(t, line)
	if jsonMsg(m) != "to-file" || jsonLevel(m) != "info" {
		t.Fatalf("file content unexpected: %s", string(b))
	}
	if got, _ := jsonCtx(m)["k"].(string); got != "v" {
		t.Fatalf("ctx[k] = %v, want v", jsonCtx(m)["k"])
	}
}

func TestCoreEndTasks(t *testing.T) {
	et := &EndTasks{}
	order := []int{}
	et.Append(func() { order = append(order, 1) })
	et.Append(func() { order = append(order, 2) })
	et.Execute()
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("execution order = %v, want [1 2]", order)
	}

	l := NewLogger()
	if l.End() == nil || l.End() != endTasks {
		t.Fatalf("Logger.End() mismatch")
	}
	l.Close()
}

func TestCoreCloseIdempotent(t *testing.T) {
	l, c := newTestLogger(t, NewJSONSerializer())
	l.Info("once", Sync{})
	l.Close()
	l.Close() // should not panic
	out := c.readNoClose(t)
	if strings.Count(out, `"msg":"once"`) != 1 {
		t.Fatalf("expected exactly one message:\n%s", out)
	}
	_ = c.w.Close()
	_ = c.r.Close()
}

func (c *testCapture) readNoClose(t *testing.T) string {
	t.Helper()
	// logger already closed by caller; drain pipe
	_ = c.w.Close()
	b, err := io.ReadAll(c.r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(b)
}

func TestCoreUtils(t *testing.T) {
	if got := F("x=%d", 3); got != "x=3" {
		t.Fatalf("F() = %q", got)
	}
	if got := T(" \t a \n"); got != "a" {
		t.Fatalf("T() = %q", got)
	}
	err := Ef("e-%s", "x")
	if err == nil || err.Error() != "e-x" {
		t.Fatalf("Ef() = %v", err)
	}
}

func TestCoreGlobalLogging(t *testing.T) {
	old := GlobalLogger
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	GlobalLogger = NewLogger(
		WithPipe(NewPipe(WithSerializer(NewJSONSerializer()), WithWriter(NewConsoleWriter(WithCustomStdout(w))))),
		WithLogLevel(LogLevelDebug),
	)
	defer func() {
		GlobalLogger.Close()
		GlobalLogger = old
	}()

	Info("g-info", Sync{}, C{"k": "v"})
	Infof("g-%s", "infof", Sync{})
	Log(LogLevelWarn, "g-log", Sync{})
	Logf(LogLevelError, "g-%s", "logf", Sync{})
	Invoked()
	GlobalLogger.Close()
	_ = w.Close()
	b, _ := io.ReadAll(r)
	_ = r.Close()

	out := string(b)
	for _, want := range []string{`"msg":"g-info"`, `"msg":"g-infof"`, `"msg":"g-log"`, `"msg":"g-logf"`, `"msg":"Invoked"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %s\n%s", want, out)
		}
	}
	if !strings.Contains(out, `"level":"warn"`) || !strings.Contains(out, `"level":"error"`) || !strings.Contains(out, `"level":"debug"`) {
		t.Fatalf("expected levels in output:\n%s", out)
	}
}

func TestCoreDirectContext(t *testing.T) {
	l, c := newTestLogger(t, NewJSONSerializer())
	l.Info("msg", "request_id", "id8235", C{"database": "postgres"}, Sync{})
	out := c.read(t, l)
	m := decodeJSONLine(t, out)
	ctx := jsonCtx(m)
	if got, _ := ctx["request_id"].(string); got != "id8235" {
		t.Fatalf("ctx[request_id] = %v, want id8235", ctx["request_id"])
	}
	if got, _ := ctx["database"].(string); got != "postgres" {
		t.Fatalf("ctx[database] = %v, want postgres", ctx["database"])
	}

	l2, c2 := newTestLogger(t, NewJSONSerializer())
	l2.Infof("hello %s", "world", "a", float64(1), "b", "two", Sync{})
	out2 := c2.read(t, l2)
	m2 := decodeJSONLine(t, out2)
	if got := jsonMsg(m2); got != "hello world" {
		t.Fatalf("msg = %q, want %q", got, "hello world")
	}
	ctx2 := jsonCtx(m2)
	if got, ok := ctx2["a"].(float64); !ok || got != 1 {
		t.Fatalf("ctx[a] = %v (%T), want 1", ctx2["a"], ctx2["a"])
	}
	if got, _ := ctx2["b"].(string); got != "two" {
		t.Fatalf("ctx[b] = %v, want two", ctx2["b"])
	}

	l3, c3 := newTestLogger(t, NewJSONSerializer())
	l3.Error("oops", "err_ctx", errors.New("boom"), Sync{})
	out3 := c3.read(t, l3)
	m3 := decodeJSONLine(t, out3)
	if got, _ := jsonCtx(m3)["err_ctx"].(string); got != "boom" {
		t.Fatalf("ctx[err_ctx] = %v, want boom", jsonCtx(m3)["err_ctx"])
	}
}
