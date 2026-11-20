# lg

**lg** is a structured, asynchronous logger for Go with multi-output support, configurable log levels, and a clean pipe-based architecture.

## How It Works

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Logger                                    │
│                                                                     │
│  logger.Info("msg", lg.C{"key": "val"})                             │
│       │                                                             │
│       ▼                                                             │
│  buildMessage()  ──►  SafeChannel[Message]  (queue, default 8192)   │
│                              │                                      │
│                        worker goroutine                             │
│                              │                                      │
│              ┌───────────────┼───────────────┐                      │
│              ▼               ▼               ▼                      │
│           Pipe 1          Pipe 2          Pipe N                    │
│        ┌────────┐      ┌────────┐                                   │
│        │  Ser.  │      │  Ser.  │   Serializer → string             │
│        │  Wri.  │      │  Wri.  │   Writer    → stdout / file       │
│        └────────┘      └────────┘                                   │
└─────────────────────────────────────────────────────────────────────┘
```

Log calls are non-blocking by default — messages are pushed to an internal `SafeChannel` and processed by a single worker goroutine. Each `Pipe` combines one `Serializer` and one `Writer`, and the logger can fan out to any number of pipes simultaneously. Use `lg.Sync{}` as an argument to make a specific call block until it is fully processed.

## Installation

```bash
go get -u github.com/brezzgg/go-packages/lg@latest
```

## Quick Start

```go
package main

import "github.com/brezzgg/go-packages/lg"

func main() {
    defer lg.Close()

    lg.Info("server started", lg.C{"port": 8080})
    lg.Warn("high memory usage", lg.C{"used_mb": 512})
    lg.Error("db query failed", lg.C{"error": err, "query": "SELECT ..."})
}
```

Output (console, colored):

```
2026/04/28 21:00:00+3  Info   main.main:12  server started {"port": 8080}
2026/04/28 21:00:00+3  Warn   main.main:13  high memory usage {"used_mb": 512}
2026/04/28 21:00:00+3  Error  main.main:14  db query failed {"error": "...", "query": "..."}
```

## Logger

### Creating a Logger

```go
logger := lg.NewLogger()
```

By default, a new logger has one pipe: `ConsoleSerializer` + `ConsoleWriter` (stdout). All options are applied via functional options:

```go
logger := lg.NewLogger(
    lg.WithPipe(
        lg.NewPipe(
            lg.WithSerializer(lg.NewJSONSerializer()),
            lg.WithWriter(lg.NewFileWriter("app.log")),
        ),
    ),
    lg.WithQueueSize(16384),
    lg.WithLogLevel(lg.LogLevelWarn),
    lg.WithSplitFilePrefix(true),
)
```

### Logger Options

| Option | Description |
|--------|-------------|
| `WithPipe(pipe)` | Add a single pipe |
| `WithPipes(pipes...)` | Add multiple pipes |
| `WithQueueSize(n)` | Set internal message queue size (default: 8192) |
| `WithLogLevel(level)` | Ignore messages with priority below this level |
| `WithLogLevelPriority(n)` | Same as above but by raw priority number |
| `WithConstantLevelOptions(opts...)` | Apply `LogLevelOption` to every message |
| `WithSplitFilePrefix(bool)` | Trim module prefix from caller file path |
| `WithCustomTypeConverter(c)` | Override how arguments are converted to body/context |

### Log Methods

```go
logger.Debug(args...)
logger.Debugf(format, args...)
logger.Info(args...)
logger.Infof(format, args...)
logger.Warn(args...)
logger.Warnf(format, args...)
logger.Error(args...)
logger.Errorf(format, args...)
logger.Fatal(args...)    // runs EndTasks, closes logger, os.Exit(1)
logger.Fatalf(format, args...)
logger.Panic(args...)    // runs EndTasks, closes logger, prints stack, os.Exit(1)
logger.Panicf(format, args...)
logger.Log(level, args...)
logger.Logf(level, format, args...)
logger.Exit(code)        // runs EndTasks, closes logger, os.Exit(code)
logger.Invoked()         // logs "Invoked" at Debug level — useful for tracing
logger.Close()           // flushes queue and closes all pipes (idempotent)
```

### Formatted Methods (`*f`)

Formatted methods use `fmt.Sprintf`-compatible format strings:

- The first argument is the format string.
- Arguments required by format verbs are consumed from the remaining arguments in order.
- Any arguments after the formatted values are treated as regular log arguments (for example, `lg.C{...}` for context or `lg.Sync{}`).

```go
logger.Infof("user %s logged in", userID, lg.C{"ip": ip})
logger.Errorf("query failed: %v", err, lg.C{"query": q})
logger.Logf(lg.LogLevelWarn, "retry %d/%d", attempt, maxAttempts)
```

### Global Logger

`lg.GlobalLogger` is a package-level `*Logger` created with default settings. Logger methods are mirrored as package-level functions, including formatted variants:

```go
lg.Debug(args...)
lg.Debugf(format, args...)
lg.Info(args...)
lg.Infof(format, args...)
lg.Warn(args...)
lg.Warnf(format, args...)
lg.Error(args...)
lg.Errorf(format, args...)
lg.Fatal(args...)
lg.Fatalf(format, args...)
lg.Panic(args...)
lg.Panicf(format, args...)
lg.Log(level, args...)
lg.Logf(level, format, args...)
lg.Close()
lg.Invoked()
lg.Exit(code)
```

```go
lg.Infof("server started on port %d", 8080, lg.C{"env": "prod"})
lg.Errorf("db query failed: %v", err, lg.C{"query": q})
lg.Logf(lg.LogLevelWarn, "cache miss for key %s", key)
lg.Close()
```

### Utility Functions

```go
lg.F("hello %s", name)        // alias for fmt.Sprintf
lg.T("  trimmed  ")           // alias for strings.TrimSpace
lg.Ef("wrap: %w", err)        // alias for fmt.Errorf
```

## Context and Arguments

The first argument to any log call becomes the **message body**. Additional arguments are treated as **context**.

### Context map (`lg.C`)

```go
lg.Info("user logged in", lg.C{"user_id": 42, "ip": "127.0.0.1"})
```

`lg.C` is an alias for `map[string]any`:

```go
type C map[string]any
```

### Direct key-value context

You can pass context as alternating key-value pairs directly in log arguments (without wrapping them in `lg.C`). String keys are paired with the following value and merged into the final context. This can be combined with `lg.C` maps.

```go
lg.Info("request handled", "request_id", "id8235", lg.C{"database": "postgres"})
// context: {"request_id": "id8235", "database": "postgres"}

lg.Error("query failed", "query", q, "duration_ms", 42, err)
// context: {"query": q, "duration_ms": 42, "error": "..."}

lg.Infof("user %s logged in", userID, "ip", ip, lg.C{"role": "admin"})
// message: "user <userID> logged in"
// context: {"ip": ip, "role": "admin"}
```

Passing an `error` directly as a context argument is supported:

```go
lg.Error("query failed", err)
// context: {"error": "connection refused"}
```

Duplicate keys in context are resolved automatically by appending a numeric suffix (`key`, `key2`, `key3`, ...).

### Synchronous Logging

By default, log calls are asynchronous. To block until the message is processed, pass `lg.Sync{}` as any argument:

```go
lg.Info("shutting down", lg.Sync{})
// returns only after the message is written
```

## Pipe

A `Pipe` is a pair of `Serializer` + `Writer`. The logger fans out each message to all registered pipes.

```go
pipe := lg.NewPipe(
    lg.WithSerializer(lg.NewJSONSerializer()),
    lg.WithWriter(lg.NewFileWriter("app.log")),
)

logger := lg.NewLogger(lg.WithPipe(pipe))
```

### Pipe Presets

```go
lg.AsDefaultConsole()        // ConsoleSerializer + ConsoleWriter(stdout)
lg.AsDefaultFile("app.log")  // JSONSerializer + FileWriter
```

Usage:

```go
pipe := lg.NewPipe(lg.AsDefaultFile("app.log"))
```

### Multi-Pipe Example

Write to console and file simultaneously:

```go
logger := lg.NewLogger(
    lg.WithPipes(
        lg.NewPipe(lg.AsDefaultConsole()),
        lg.NewPipe(lg.AsDefaultFile("app.log")),
    ),
)
```

## Serializers

### ConsoleSerializer

Produces a human-readable colored string:

```
2026/04/28 21:00:00+3  Info   main.main:10  message {"key": "value"}
```

```go
ser := lg.NewConsoleSerializer()
ser := lg.NewConsoleSerializer(lg.WithDisabledColors())
```

| Option | Description |
|--------|-------------|
| `WithDisabledColors()` | Disable ANSI color codes |

### JSONSerializer

Produces a compact or indented JSON string:

```json
{"time":"2026-04-28T21:00:00Z","caller":{"method":"main","file":"main","line":10},"level":"Info","msg":"message","ctx":{"key":"value"}}
```

```go
ser := lg.NewJSONSerializer()
ser := lg.NewJSONSerializer(lg.WithIndent())
```

| Option | Description |
|--------|-------------|
| `WithIndent()` | Pretty-print JSON with 2-space indent |

### Custom Serializer

Implement the `Serializer` interface:

```go
type Serializer interface {
    Serialize(m lg.Message) string
}
```

## Writers

### ConsoleWriter

Writes to `os.Stdout` by default:

```go
w := lg.NewConsoleWriter()
w := lg.NewConsoleWriter(lg.WithCustomStdout(os.Stderr))
w := lg.NewConsoleWriter(lg.WithWriterDiscard()) // discard output (useful in tests/benchmarks)
```

### FileWriter

Writes to a file with an internal buffer (default size: 1024 bytes). Flushes automatically when the buffer is full or after 10 seconds of inactivity.

```go
w := lg.NewFileWriter("app.log")
w := lg.NewFileWriter("app.log", lg.WithCustomBufferSize(4096))
```

### Custom Writer

Implement the `Writer` interface:

```go
type Writer interface {
    Write(message string) error
    Flush()
    Close()
}
```

## Log Levels

### Built-in Levels

| Level | Priority | Color |
|-------|----------|-------|
| `LogLevelDebug` | 100 | Bold Green |
| `LogLevelInfo` | 200 | Bold Blue |
| `LogLevelWarn` | 300 | Bold Yellow |
| `LogLevelError` | 400 | Bold Red |
| `LogLevelFatal` | 500 | Bold Red |
| `LogLevelPanic` | 600 | Bold Red |

### Custom Levels

```go
var LogLevelTrace = lg.NewLogLevel(lg.ClrFgCyan, "Trace").WithPriority(50)

// log with it
logger.Log(LogLevelTrace, "trace message")

// or attach it to the global logger
LogLevelTrace.LogGlobal("trace message")
```

### Level Priority Filter

Messages with a priority strictly below the configured threshold are silently dropped:

```go
logger := lg.NewLogger(lg.WithLogLevel(lg.LogLevelWarn))
// Debug and Info calls will be ignored
```

### Level Options

`LogLevelOption` constants modify how the caller and timestamp are rendered for a specific level. Apply them per-level or globally via `WithConstantLevelOptions`.

**Caller options:**

| Option | Result |
|--------|--------|
| `LevelOptionDisableCaller` | No caller info |
| `LevelOptionCallerOnlyFunc` | `main()` |
| `LevelOptionCallerOnlyFile` | `main/app` |
| `LevelOptionCallerDisableFunc` | `main/app:10` |
| `LevelOptionCallerDisableFile` | `main():10` |
| `LevelOptionCallerDisableLine` | `main/app.main()` |

**Time options:**

| Option | Result |
|--------|--------|
| `LevelOptionDisableTime` | No timestamp |
| `LevelOptionTimeDisableOffset` | UTC only, no timezone offset |

**Example:**

```go
var LogLevelAudit = lg.NewLogLevel(lg.ClrFgCyan, "Audit").
    WithPriority(250).
    WithOptions(lg.LevelOptionDisableCaller, lg.LevelOptionTimeDisableOffset)
```

## EndTasks

`EndTasks` is a list of cleanup functions that run automatically before `Fatal`, `Panic`, and `Exit`. Register tasks with `logger.End().Append(...)`:

```go
logger.End().Append(func() {
    db.Close()
    cache.Flush()
})

lg.Fatal("unrecoverable error") // db.Close() and cache.Flush() run first
```

## Signal Handling

Each `Logger` starts a background goroutine that listens for `SIGTERM` and `SIGINT`. On signal receipt, the logger drains its queue and closes all pipes gracefully — no explicit `defer lg.Close()` is required for these cases, but it is still recommended.

## Message Structure

```go
type Message struct {
    Time    time.Time `json:"time"`
    Caller  Caller    `json:"caller"`
    Level   LogLevel  `json:"level"`
    Text    string    `json:"msg"`
    Context C         `json:"ctx,omitempty"`
}

type Caller struct {
    Method string `json:"method"`
    File   string `json:"file"`
    Line   int    `json:"line"`
}
```

## Custom Type Converter

`TypeConverter` controls how log arguments are converted to message body and context. Override it to support additional types:

```go
type TypeConverter interface {
    ConvAndPushBody(item any, push BodyConverterFunc)
    ConvAndPushContext(item any, push ContextConverterFunc)
}
```

```go
type myConverter struct{ lg.defaultTypeConverter }

func (c myConverter) ConvAndPushBody(item any, push lg.BodyConverterFunc) {
    if s, ok := item.(MyStringer); ok {
        push(s.MyString())
        return
    }
    // fall back to default
}

logger := lg.NewLogger(lg.WithCustomTypeConverter(&myConverter{}))
```

## Benchmarks

Run benchmarks:

```bash
go test -bench=. -benchmem ./...
```

Tested scenarios: single-thread async, multi-thread async — for each combination of serializer (`ConsoleSerializer`, `JSONSerializer`) and writer (`ConsoleWriter` discard, `FileWriter`), plus multi-pipe.
