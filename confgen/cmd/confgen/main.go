package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/brezzgg/go-packages/confgen/internal/generator"
	"github.com/brezzgg/go-packages/confgen/internal/parser"
	"github.com/brezzgg/go-packages/confgen/internal/writer"
)

func main() {
	recursive := flag.Bool("r", false, "walk dir recursive")
	flag.Parse()
	paths := flag.Args()

	if len(paths) == 0 {
		fmt.Println("usage: confgen <schema.hcl> [schema2.hcl ...]")
	}

	schemas, err := parser.Parse(paths, *recursive)
	if err != nil {
		fmt.Printf("parse file error: %s\n", err)
	}

	for file, schema := range schemas {
		go spin(file)

		gen, err := generator.Generate(schema)
		if err != nil {
			fmt.Printf("\r\033[K")
			fmt.Printf("✗ %s: generate error: %s\n", file, err)
			continue
		}

		isStdout, err := writer.Write(schema, gen)
		if err != nil {
			fmt.Printf("\r\033[K")
			fmt.Printf("✗ %s: write error: %s\n", file, err)
			continue
		}

		if !isStdout {
			fmt.Printf("\r\033[K")
			fmt.Printf("✓ %s\n", file)
		} else {
			fmt.Printf("\r\033[K")
			fmt.Printf("✓ %s:\n", file)
			fmt.Println(gen)
		}
	}
}

func spin(msg string) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		fmt.Printf("\r%s %s", frames[i%len(frames)], msg)
		i++
		time.Sleep(100 * time.Millisecond)
	}
}
