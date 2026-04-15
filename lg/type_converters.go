package lg

import (
	"fmt"
)

type defaultTypeConverter struct{}

func NewDefaultTypeConverter() *defaultTypeConverter {
	return &defaultTypeConverter{}
}

func (c defaultTypeConverter) ConvAndPushBody(item any, push BodyConverterFunc) {
	switch value := item.(type) {
	case string:
		push(value)
	case fmt.Stringer:
		push(value.String())
	case error:
		push(value.Error())
	default:
		push("<unknown>")
	}
}

func (c defaultTypeConverter) ConvAndPushContext(item any, push ContextConverterFunc) {
	switch value := item.(type) {
	case C:
		c.convAndPushContextC(value, push)
	case error:
		push("error", value.Error())
	case Sync:
		return
	default:
		push("arg", value)
	}
}

func (c defaultTypeConverter) convAndPushContextC(context C, push ContextConverterFunc) {
	for key, val := range context {
		switch v := val.(type) {
		case error:
			push(key, v.Error())
		default:
			push(key, v)
		}
	}
}
