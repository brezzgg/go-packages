package diag

func From(src DiagnosticSource) *Diagnostic {
	return &Diagnostic{
		dname: src.DiagName(),
		dtype: src.DiagType(),
	}
}

func (d *Diagnostic) Child(dname, dtype string) *Diagnostic {
	n := New(dname, dtype)
	n.parent = d
	d.children = append(d.children, n)
	return n
}

func (d *Diagnostic) ChildSource(src DiagnosticSource) *Diagnostic {
	return d.Child(src.DiagName(), src.DiagType())
}

func (d *Diagnostic) ProcessSource(src DiagnosticSource) {
	src.Diagnose(d.ChildSource(src))
}
