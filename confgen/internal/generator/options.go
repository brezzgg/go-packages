package generator

import (
	"reflect"

	"github.com/brezzgg/go-packages/confgen/internal/schema"
	"github.com/brezzgg/go-packages/confgen/templates/functions"
	"github.com/brezzgg/go-packages/confgen/templates/template"
)

type (
	Options struct {
		template string
		funcs    map[string]any
		data     map[string]any
	}
	Option func(o *Options)
)

func NewOptions(schema *schema.Schema) *Options {
	return &Options{
		template: template.Templates[template.Default],
		funcs:    functions.Default,
		data:     structToMap(*schema),
	}
}

func WithTemplateFunctions(funcs map[string]any) Option {
	return func(o *Options) {
		for name, fn := range funcs {
			if _, ok := o.funcs[name]; !ok {
				o.funcs[name] = fn
			}
		}
	}
}

func WithTemplateData(data map[string]any) Option {
	return func(o *Options) {
		for k, v := range data {
			if _, ok := o.data[k]; !ok {
				o.funcs[k] = v
			}
		}
	}
}

func WithCustomTemplate(tmpl string) Option {
	return func(o *Options) {
		o.template = tmpl
	}
}

func WithTemplate(tmplName string) Option {
	return func(o *Options) {
		o.template = template.Templates[tmplName]
	}
}

func structToMap(v any) map[string]any {
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	m := make(map[string]any, rt.NumField())
	for i := range rt.NumField() {
		m[rt.Field(i).Name] = rv.Field(i).Interface()
	}
	return m
}
