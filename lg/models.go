package lg

/*
C is an alias that is equivalent to map[string]any.
It is used to add context to logs in a convenient way.
*/
type C map[string]any

/*
Sync is a context type used to synchronize log handling across goroutines.
*/
type Sync struct{}

type TypeConverter interface {
	ConvAndPushBody(item any, push BodyConverterFunc)
	ConvAndPushContext(item any, push ContextConverterFunc)
}

type (
	BodyConverterFunc    func(string)
	ContextConverterFunc func(string, any)
)

type Serializer interface {
	Serialize(m Message) string
}

type Writer interface {
	Write(Message string) error
	Flush()
	Close()
}

type ColorizedOut interface {
	Colorize(color string) string
	Colorized() string
	Standard() string
}
