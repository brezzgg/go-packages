package functions

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/brezzgg/go-packages/confgen/internal/schema"
)

var primitiveDefaults = map[string]string{
	"string":  `""`,
	"str":     `""`,
	"bool":    "false",
	"boolean": "false",
	"int":     "0",
	"int8":    "0",
	"int16":   "0",
	"int32":   "0",
	"int64":   "0",
	"uint":    "0",
	"uint8":   "0",
	"uint16":  "0",
	"uint32":  "0",
	"uint64":  "0",
	"float32": "0",
	"float64": "0",
	"byte":    "0",
}

func toType(f schema.Field) string {
	if f.Type == nil {
		return "string"
	}
	t := strings.ToLower(*f.Type)

	if base, ok := strings.CutPrefix(t, "list "); ok {
		switch base {
		case "object", "obj":
			if f.Object != nil && *f.Object != "" {
				return "[]" + toPascalCase(*f.Object)
			}
			return "[]any"
		default:
			if _, ok := primitiveDefaults[base]; ok {
				return "[]" + base
			}
			return "[]string"
		}
	}

	switch t {
	case "list", "slice", "array":
		if f.Object != nil && *f.Object != "" {
			return "[]" + toPascalCase(*f.Object)
		}
		return "[]string"
	case "object", "obj":
		if f.Object != nil {
			return toPascalCase(*f.Object)
		}
		return "any"
	default:
		if _, ok := primitiveDefaults[t]; ok {
			return t
		}
		return "string"
	}
}

func toDefault(f schema.Field) string {
	if f.Type == nil {
		if f.Default != nil && *f.Default != "" {
			return fmt.Sprintf("%q", *f.Default)
		}
		return `""`
	}
	t := strings.ToLower(*f.Type)

	if _, ok := strings.CutPrefix(t, "list "); ok {
		return "nil"
	}

	switch t {
	case "list", "slice", "array":
		return "nil"
	case "object", "obj":
		if f.Object != nil {
			return "Default" + toPascalCase(*f.Object) + "()"
		}
		return "nil"
	default:
		if f.Default != nil && *f.Default != "" {
			if t == "string" || t == "str" {
				return fmt.Sprintf("%q", *f.Default)
			}
			return *f.Default
		}
		if zero, ok := primitiveDefaults[t]; ok {
			return zero
		}
		return `""`
	}
}

func toBool(in any) bool {
	switch t := in.(type) {
	case *bool:
		return !(t == nil || *t == false)
	case bool:
		return t
	case string:
		if t == "true" || t == "yes" {
			return true
		}
		return false
	case *string:
		if t == nil {
			return false
		}
		if *t == "true" || *t == "yes" {
			return true
		}
		return false
	default:
		return false
	}
}

func toTag(name string, formats []string, required bool) string {
	var tags []string
	for _, f := range formats {
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", f, name))
	}
	if required {
		tags = append(tags, "confgen_required:\"true\"")
	}
	return "`" + strings.Join(tags, " ") + "`"
}

func deref(v any) any {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.IsNil() {
		return nil
	}
	return rv.Elem().Interface()
}
