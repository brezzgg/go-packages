package schema

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type File struct {
	Output   *Output    `hcl:"output,block"`
	Objects  []Object   `hcl:"object,block"`
	Generate []Generate `hcl:"generate,block"`
	Remain   hcl.Body   `hcl:",remain"`
}

type Output struct {
	Package string   `hcl:"package"`
	Output  string   `hcl:"output"`
	Formats []string `hcl:"formats"`
	Remain  hcl.Body `hcl:",remain"`
}

type Object struct {
	Name   string   `hcl:"name,label"`
	Fields []Field  `hcl:"field,block"`
	Remain hcl.Body `hcl:",remain"`
}

type Generate struct {
	Name   string   `hcl:"name,label"`
	Fields []Field  `hcl:"field,block"`
	Remain hcl.Body `hcl:",remain"`
}

type Schema struct {
	Output   *Output
	Objects  map[string]*Object
	Generate map[string]*Generate
}

type Field struct {
	Name     string   `hcl:"name,label"`
	Default  *string  `hcl:"default"`
	Type     *string  `hcl:"type"`
	Object   *string  `hcl:"object"`
	Required *bool    `hcl:"required"`
	Desc     *string  `hcl:"desc"`
	Remain   hcl.Body `hcl:",remain"`
}

func Decode(body hcl.Body) (*Schema, error) {
	var raw File
	if diags := gohcl.DecodeBody(body, nil, &raw); diags.HasErrors() {
		return nil, fmt.Errorf("decode hcl body: %w", diags)
	}

	s := &Schema{
		Output:   raw.Output,
		Objects:  make(map[string]*Object),
		Generate: make(map[string]*Generate),
	}

	for i := range raw.Objects {
		s.Objects[raw.Objects[i].Name] = &raw.Objects[i]
	}

	for i := range raw.Generate {
		s.Generate[raw.Generate[i].Name] = &raw.Generate[i]
	}

	return s, nil
}

func (f *Field) ParsedType() *FieldType {
	if f.Type == nil {
		return nil
	}
	t := ParseFieldType(*f.Type)
	return &t
}
