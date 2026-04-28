package generator

import (
	"bytes"
	"fmt"
	T "text/template"

	"github.com/brezzgg/go-packages/confgen/internal/schema"
)

func Generate(s *schema.Schema, opts ...Option) (string, error) {
	options := NewOptions(s)
	for _, fn := range opts {
		fn(options)
	}

	if s.Output == nil || s.Output.Output == "" || s.Output.Package == "" || len(s.Output.Formats) == 0 {
		return "", fmt.Errorf("output is required in hcl file. example:\n" +
			"output {\n  package = \"config\"\n  output = \"stdout\"\n  formats = [\"yaml\", \"json\"]\n}\n",
		)
	}

	tmpl, err := T.New("config").Funcs(options.funcs).Parse(options.template)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, options.data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}
