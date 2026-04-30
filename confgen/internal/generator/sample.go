package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/brezzgg/go-packages/confgen/internal/schema"
	"gopkg.in/yaml.v3"
)

func GenerateSample(s *schema.Schema, format string) (string, error) {
	data := buildMap(s.Generate, s.Objects)

	switch format {
	case "yaml":
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(2)
		if err := enc.Encode(data); err != nil {
			return "", fmt.Errorf("marshal yaml: %w", err)
		}
		return buf.String(), nil
	case "json":
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal json: %w", err)
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func buildMap(generate map[string]*schema.Generate, objects map[string]*schema.Object) map[string]any {
	result := make(map[string]any)
	for _, gen := range generate {
		for _, f := range gen.Fields {
			result[f.Name] = resolveField(f, objects)
		}
	}
	return result
}

func resolveField(f schema.Field, objects map[string]*schema.Object) any {
	if f.Object != nil {
		obj, ok := objects[*f.Object]
		if !ok {
			return nil
		}
		nested := make(map[string]any)
		for _, nf := range obj.Fields {
			nested[nf.Name] = resolveField(nf, objects)
		}
		return nested
	}

	if f.Type != nil && *f.Type == "list" {
		return []any{}
	}

	if f.Default != nil {
		return coerce(*f.Default, f.Type)
	}

	return nil
}

func coerce(val string, typ *string) any {
	if typ == nil {
		return val
	}
	switch *typ {
	case "int":
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	case "bool":
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	case "float":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return val
}
