package diag

type Diagnoser interface {
	Diagnose(diag *Diagnostic)
}

type DiagnosticSource interface {
	Diagnoser
	DiagName() string
	DiagType() string
}

type Error struct {
	error   string
	errType ErrType
}

type Diagnostic struct {
	errs     []Error
	dname    string
	dtype    string
	parent   *Diagnostic
	children []*Diagnostic
}

func New(dname, dtype string) *Diagnostic {
	return &Diagnostic{
		dname: dname,
		dtype: dtype,
	}
}
