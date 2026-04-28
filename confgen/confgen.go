package confgen

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Unmarshaler func(data []byte, v any) error

func Unmarshal(data []byte, v any, unmarshaler Unmarshaler) error {
	if err := unmarshaler(data, v); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	if d, ok := v.(Defaultable); ok {
		mergeDefaults(reflect.ValueOf(v), reflect.ValueOf(d.Defaults()))
	}
	return validate(reflect.ValueOf(v))
}

type Defaultable interface {
	Defaults() any
}

type Decoder struct {
	path        string
	unmarshaler Unmarshaler
}

func NewDecoder(path string, unmarshaler Unmarshaler) *Decoder {
	return &Decoder{path: path, unmarshaler: unmarshaler}
}

func (d *Decoder) Decode(v any) error {
	data, err := os.ReadFile(d.path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	return Unmarshal(data, v, d.unmarshaler)
}

func mergeDefaults(dst, def reflect.Value) {
	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			return
		}
		dst = dst.Elem()
		def = def.Elem()
	}
	if dst.Kind() != reflect.Struct {
		return
	}

	for i := range dst.NumField() {
		d := dst.Field(i)
		dv := def.Field(i)

		if d.IsZero() {
			d.Set(dv)
			continue
		}

		switch d.Kind() {
		case reflect.Ptr:
			if !d.IsNil() {
				mergeDefaults(d, dv)
			}
		case reflect.Struct:
			mergeDefaults(d, dv)
		}
	}
}

func validate(v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("confgen")
		if strings.Contains(tag, "confgen_required") && value.IsZero() {
			return fmt.Errorf("field %q is required", field.Name)
		}

		switch value.Kind() {
		case reflect.Ptr:
			if !value.IsNil() {
				if err := validate(value.Elem()); err != nil {
					return fmt.Errorf("%s.%w", field.Name, err)
				}
			}
		case reflect.Struct:
			if err := validate(value); err != nil {
				return err
			}
		}
	}

	return nil
}
