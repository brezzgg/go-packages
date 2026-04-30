package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brezzgg/go-packages/confgen/internal/schema"
)

const outputErr = `output block is required in hcl file.
example:
output {
  package = "config"
  output  = "internal/config"
  formats = ["yaml"]
}`

func Write(s *schema.Schema, content string) (string, error) {
	if s.Output == nil {
		return "", fmt.Errorf(outputErr)
	}

	if s.Output.Output == "stdout" {
		return s.Output.Output, nil
	}

	if err := os.MkdirAll(s.Output.Output, 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	fullPath := filepath.Join(s.Output.Output, s.Output.Package+".gen.go")
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return fullPath, nil
}
