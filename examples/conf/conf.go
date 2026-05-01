package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/brezzgg/go-packages/conf"
)

func main() {
	ExampleConfigure()
}

func ExampleConfigure() {
	var (
		someString string
		someInt    conf.Converter
	)

	os.Setenv("OTHER_SOME_KEY", "other some value")
	os.Args = append(os.Args, "--some-int", "-5")

	err := conf.Parse(
		conf.Env( // parse parameters from env
			conf.DParam("SOME_KEY", "some default value"),
			// DParam: if parameter is not set, return default value
			conf.RParam("OTHER_SOME_KEY"),
			// RParam: if parameter is not set, generates an error
		),
		conf.Arg( // parse parameters from flags
			//conf.DParam("some-key", "some default value"),
			// err: parameter 'some_key' was duplicated
			conf.DParam("some-string", "some-string-value").Write(&someString),
			// write result of parse to var
			conf.RParam("some-int").WriteConv(&someInt),
			// write result`s converter to var
		),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("specific values:")
	fmt.Println("some-string:", someString)
	fmt.Println("other-some-key:", conf.Get("OTHER_SOME_KEY")) // get parsed value by key
	fmt.Println("some-int:", someInt.MustInt())                // convert value to int

	fmt.Println("\nall values:")
	b, _ := json.MarshalIndent(
		conf.Geta(), // get all parsed values
		"", "  ",
	)
	fmt.Println(string(b))
}
