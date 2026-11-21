package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/brezzgg/go-packages/lg"
)

var indent bool

func main() {
	flag.BoolVar(&indent, "indent", false, "")
	flag.Parse()

	scan()
}

func scan() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	var serOpts []lg.ConsoleSerializerOption
	serOpts = append(serOpts, lg.WithDisabledColors())
	if indent {
		serOpts = append(serOpts, lg.WithIndentContext())
	}

	ser := lg.NewConsoleSerializer(serOpts...)
	wri := lg.NewConsoleWriter()

	for scanner.Scan() {
		raw := scanner.Text()
		if !strings.HasPrefix(raw, "{") || !strings.HasSuffix(raw, "}") {
			wri.Write(raw)
			continue
		}

		var msg lg.Message
		err := json.Unmarshal([]byte(raw), &msg)
		if err != nil {
			wri.Write(raw)
			continue
		}

		wri.Write(ser.Serialize(msg))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "read error:", err)
		os.Exit(1)
	}
}
