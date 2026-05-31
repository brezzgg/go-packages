package diag

import "fmt"

func (d *Diagnostic) Error(format string, a ...any) {
	d.errs = append(d.errs, Error{
		errType: ErrTypeError,
		error:   fmt.Sprintf(format, a...),
	})
}

func (d *Diagnostic) Warn(format string, a ...any) {
	d.errs = append(d.errs, Error{
		errType: ErrTypeWarn,
		error:   fmt.Sprintf(format, a...),
	})
}

type ErrType int

const (
	ErrTypeError ErrType = iota
	ErrTypeWarn
)
