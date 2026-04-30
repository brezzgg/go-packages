package schema

import "strings"

type FieldType struct {
	Base string
	List bool
}

func ParseFieldType(raw string) FieldType {
	if rest, ok := strings.CutPrefix(raw, "list "); ok {
		return FieldType{List: true, Base: rest}
	}
	return FieldType{List: false, Base: raw}
}
