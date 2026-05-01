package main

import (
	"fmt"
	"os"

	"github.com/brezzgg/go-packages/examples/lg/options"
	"github.com/brezzgg/go-packages/lg"
)

func main() {
	ExampleSetupLogger()
	options.ExampleLevelOptions()
	ExampleLogLevels()
}

func ExampleSetupLogger() {
	lg.Log(lg.NewLogLevel(lg.ClrFgYellow, "!Example!").WithOptions(
		lg.LevelOptionDisableCaller, lg.LevelOptionDisableTime),
		"Setup logger example",
	)

	lg.GlobalLogger = lg.NewLogger(
		lg.WithPipe(
			lg.NewPipe(
				lg.WithSerializer(lg.NewConsoleSerializer()),
				lg.WithWriter(lg.NewConsoleWriter()),
			),
		),
	)

	lg.Info("some log", lg.C{"with some": "context"})
	// output in console: 2025/01/01 00:00:00+0  Info   main.ExampleSetupLogger:29  some log {"with some":"context"}

	logger := lg.NewLogger(
		lg.WithPipe(
			lg.NewPipe(
				lg.WithSerializer(lg.NewJSONSerializer()),
				lg.WithWriter(lg.NewFileWriter("log.json")),
			),
		),
	)

	logger.End().Append(func() {
		//lg.Debug("end tasks executed")
		// is bad practice, because the logger will most likely fail to output the message in time

		b, _ := os.ReadFile("log.json")
		fmt.Println(
			// trim spaces
			lg.T(
				// sprintf equivalent
				lg.F(
					"log.json: %s", string(b),
				),
			),
		)
		_ = os.Remove("log.json")
	})

	logger.Info("some log", lg.C{"with some": "context"})
	// output in log.json: {"time":"2025-01-01T00:00:00.0000000+00:00","caller":{"method":"ExampleSetupLogger","file":"main","line":20},"level":"info","msg":"some log","ctx":{"with some":"context"}}
}

func ExampleLogLevels() {
	lg.Log(lg.NewLogLevel(lg.ClrFgYellow, "!Example!").WithOptions(
		lg.LevelOptionDisableCaller, lg.LevelOptionDisableTime),
		"Log levels example",
	)

	lg.GlobalLogger.Close()
	lg.GlobalLogger = lg.NewLogger()

	var (
		str = "Hello world"
		ctx = lg.C{"hello": "world"}
	)

	// create custom log level
	logLevelProd := lg.NewLogLevel(lg.ClrFgCyan, "Prod")

	// write log with custom log level
	lg.Log(logLevelProd, str, ctx)

	lg.Debug(str, ctx)
	lg.Info(str, ctx)
	lg.Warn(str, ctx)
	lg.Error(str, ctx)
	//lg.Fatal(str, ctx)
	lg.Panic(str, ctx)
}
