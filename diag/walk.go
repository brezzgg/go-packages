package diag

import "fmt"

func (d *Diagnostic) Walk(fn func(msg string, etype ErrType, dtype, dname string)) {
	for _, err := range d.errs {
		fn(err.error, err.errType, d.dtype, d.dname)
	}
	for _, diag := range d.children {
		diag.Walk(fn)
	}
}

func (d *Diagnostic) WalkWarnings(fn func(msg, dtype, dname string)) {
	d.Walk(func(msg string, etype ErrType, dtype, dname string) {
		if etype == ErrTypeWarn {
			fn(msg, dtype, dname)
		}
	})
}

func (d *Diagnostic) WalkErrors(fn func(msg, dtype, dname string)) {
	d.Walk(func(msg string, errType ErrType, objType, objName string) {
		if errType == ErrTypeError {
			fn(msg, objType, objName)
		}
	})
}

func (d *Diagnostic) HasError() bool {
	for _, err := range d.errs {
		if err.errType == ErrTypeError {
			return true
		}
	}
	for _, child := range d.children {
		r := child.HasError()
		if r {
			return true
		}
	}
	return false
}

func (d *Diagnostic) FirstError() error {
	for _, err := range d.errs {
		if err.errType == ErrTypeError {
			return fmt.Errorf(err.error)
		}
	}
	for _, child := range d.children {
		if err := child.FirstError(); err != nil {
			return err
		}
	}
	return nil
}
