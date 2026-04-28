package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/brezzgg/go-packages/confgen/internal/schema"
)

func Write(s *schema.Schema, content string) (bool, error) {
	if s.Output == nil {
		return false, fmt.Errorf("output block is required in hcl file.\nexample:\noutput {\n  package = \"config\"\n  output  = \"internal/config\"\n}")
	}

	if s.Output.Output == "stdout" {
		return true, nil
	}

	if err := os.MkdirAll(s.Output.Output, 0755); err != nil {
		return false, fmt.Errorf("create output dir: %w", err)
	}

	fullPath := filepath.Join(s.Output.Output, s.Output.Package+".gen.go")
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return false, fmt.Errorf("write file: %w", err)
	}

	return false, nil
}
