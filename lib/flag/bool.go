package flag

type boolFlag struct {
	defaultValue bool
	description  string
	dest         *bool
	long         string
	longVal      bool
	short        string
	shortVal     bool
}

func (f *boolFlag) Reconcile() {
	*f.dest = f.longVal || f.shortVal || f.defaultValue
}
