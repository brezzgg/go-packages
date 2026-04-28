package main

import (
	"encoding/json"
	"fmt"

	"github.com/brezzgg/go-packages/confgen"
	"github.com/brezzgg/go-packages/examples/confgen/config"
	"gopkg.in/yaml.v3"
)

func main() {
	decoder := confgen.NewDecoder("config.yaml", yaml.Unmarshal)
	var cfg config.App
	if err := decoder.Decode(&cfg); err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
