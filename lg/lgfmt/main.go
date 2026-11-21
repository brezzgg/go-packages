package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/brezzgg/go-packages/lg"
)

func main() {
	scan()
}

func scan() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	ser := lg.NewConsoleSerializer(lg.WithDisabledColors())
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
