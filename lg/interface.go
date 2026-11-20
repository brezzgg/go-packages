// Package lg provide an api to work with Logger
package lg

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

var (
	GlobalLogger = NewLogger()
)

// --------------- Logger struct functions ---------------

/*
Debug logs a message at the debug level.
*/
func (l *Logger) Debug(args ...any) {
	l.Handle(args, LogLevelDebug)
}

/*
Debugf logs a formatted message at the debug level.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	logger.Debugf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func (l *Logger) Debugf(f string, args ...any) {
	l.Handlef(f, args, LogLevelDebug)
}

/*
Info logs a message at the info level.
*/
func (l *Logger) Info(args ...any) {
	l.Handle(args, LogLevelInfo)
}

/*
Infof logs a formatted message at the info level.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	logger.Infof("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func (l *Logger) Infof(f string, args ...any) {
	l.Handlef(f, args, LogLevelInfo)
}

/*
Warn logs a message at the warn level.
*/
func (l *Logger) Warn(args ...any) {
	l.Handle(args, LogLevelWarn)
}

/*
Warnf logs a formatted message at the warn level.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	logger.Warnf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func (l *Logger) Warnf(f string, args ...any) {
	l.Handlef(f, args, LogLevelWarn)
}

/*
Error logs a message at the error level.
*/
func (l *Logger) Error(args ...any) {
	l.Handle(args, LogLevelError)
}

/*
Errorf logs a formatted message at the error level.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	logger.Errorf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func (l *Logger) Errorf(f string, args ...any) {
	l.Handlef(f, args, LogLevelError)
}

/*
Fatal logs a message at the fatal level and terminates the program with exit code 1.
*/
func (l *Logger) Fatal(args ...any) {
	// execute end tasks
	endTasks.Execute()

	l.Handle(args, LogLevelFatal)

	l.Close()

	// exit with code 1
	os.Exit(1)
}

/*
Fatalf logs a formatted message at the fatal level and terminates the program with exit code 1.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	logger.Fatalf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func (l *Logger) Fatalf(f string, args ...any) {
	// execute end tasks
	endTasks.Execute()

	l.Handlef(f, args, LogLevelFatal)

	l.Close()

	// exit with code 1
	os.Exit(1)
}

/*
Panic logs a message at the panic level, prints a stack trace, and terminates the program with exit code 1.
*/
func (l *Logger) Panic(args ...any) {
	// execute end tasks
	endTasks.Execute()

	l.Handle(args, LogLevelPanic)

	l.Close()

	// print stack trace and exit with code 1
	debug.PrintStack()
	os.Exit(1)
}

/*
Panicf logs a formatted message at the panic level, prints a stack trace, and terminates the program with exit code 1.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	logger.Panicf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func (l *Logger) Panicf(f string, args ...any) {
	// execute end tasks
	endTasks.Execute()

	l.Handlef(f, args, LogLevelPanic)

	l.Close()

	// print stack trace and exit with code 1
	debug.PrintStack()
	os.Exit(1)
}

/*
Exit terminates the program with the given exit code.
*/
func (l *Logger) Exit(code int) {
	// execute end tasks
	endTasks.Execute()

	l.Close()

	// exit
	os.Exit(code)
}

/*
Log logs a message at the specified log level.
*/
func (l *Logger) Log(level LogLevel, args ...any) {
	l.Handle(args, level)
}

/*
Logf logs a formatted message at the specified log level.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	logger.Logf(lg.LogLevelInfo, "hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func (l *Logger) Logf(level LogLevel, f string, args ...any) {
	l.Handlef(f, args, level)
}

/*
Invoked logs a debug message indicating that a function has been called.
*/
func (l *Logger) Invoked() {
	l.Handle([]any{"Invoked"}, LogLevelDebug)
}

/*
End returns the global end tasks instance.
*/
func (l *Logger) End() *EndTasks {
	return endTasks
}

// --------------- Global logger functions ---------------

/*
Debug logs a message at the debug level using the global logger.
*/
func Debug(args ...any) {
	GlobalLogger.Debug(args...)
}

/*
Debugf logs a formatted message at the debug level using the global logger.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	lg.Debugf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func Debugf(f string, args ...any) {
	GlobalLogger.Debugf(f, args...)
}

/*
Info logs a message at the info level using the global logger.
*/
func Info(args ...any) {
	GlobalLogger.Info(args...)
}

/*
Infof logs a formatted message at the info level using the global logger.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	lg.Infof("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func Infof(f string, args ...any) {
	GlobalLogger.Infof(f, args...)
}

/*
Warn logs a message at the warn level using the global logger.
*/
func Warn(args ...any) {
	GlobalLogger.Warn(args...)
}

/*
Warnf logs a formatted message at the warn level using the global logger.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	lg.Warnf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func Warnf(f string, args ...any) {
	GlobalLogger.Warnf(f, args...)
}

/*
Error logs a message at the error level using the global logger.
*/
func Error(args ...any) {
	GlobalLogger.Error(args...)
}

/*
Errorf logs a formatted message at the error level using the global logger.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	lg.Errorf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func Errorf(f string, args ...any) {
	GlobalLogger.Errorf(f, args...)
}

/*
Fatal logs a message at the fatal level using the global logger and terminates the program with exit code 1.
*/
func Fatal(args ...any) {
	GlobalLogger.Fatal(args...)
}

/*
Fatalf logs a formatted message at the fatal level using the global logger and terminates the program with exit code 1.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	lg.Fatalf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func Fatalf(f string, args ...any) {
	GlobalLogger.Fatalf(f, args...)
}

/*
Panic logs a message at the panic level using the global logger, prints a stack trace, and terminates the program with exit code 1.
*/
func Panic(args ...any) {
	GlobalLogger.Panic(args...)
}

/*
Panicf logs a formatted message at the panic level using the global logger, prints a stack trace, and terminates the program with exit code 1.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	lg.Panicf("hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func Panicf(f string, args ...any) {
	GlobalLogger.Panicf(f, args...)
}

/*
Log logs a message at the specified log level using the global logger.
*/
func Log(level LogLevel, args ...any) {
	GlobalLogger.Log(level, args...)
}

/*
Logf logs a formatted message at the specified log level using the global logger.

The first argument, f, is a format string compatible with [fmt.Sprintf].
The arguments required by the format placeholders must come immediately
after f. Any remaining arguments are treated as log context.

Example:

	lg.Logf(lg.LogLevelInfo, "hello %s", "world", lg.Sync{}, lg.C{"some_arg": "some_value"})
*/
func Logf(level LogLevel, f string, args ...any) {
	GlobalLogger.Logf(level, f, args...)
}

/*
Close closes the global logger.
*/
func Close() {
	GlobalLogger.Close()
}

/*
Invoked logs a debug message indicating that a function has been called using the global logger.
*/
func Invoked() {
	GlobalLogger.Invoked()
}

/*
Exit terminates the program with the given exit code using the global logger.
*/
func Exit(code int) {
	GlobalLogger.Exit(code)
}

/*
F function is alias for fmt.Sprintf function.
*/
func F(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

/*
T function is alias for strings.TrimSpaces function.
*/
func T(str string) string {
	return strings.TrimSpace(str)
}

/*
Ef function is alias for fmt.Errorf function.
*/
func Ef(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
