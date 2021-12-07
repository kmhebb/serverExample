package flag

type stringFlag struct {
	defaultValue string
	description  string
	dest         *string
	long         string
	longVal      string
	short        string
	shortVal     string
}

func (f *stringFlag) Reconcile() {
	if f.longVal != "" {
		*f.dest = f.longVal
		return
	}
	if f.shortVal != "" {
		*f.dest = f.shortVal
		return
	}
	*f.dest = f.defaultValue
}
