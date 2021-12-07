package flag

type intFlag struct {
	defaultValue int
	description  string
	dest         *int
	long         string
	longVal      int
	short        string
	shortVal     int
}

func (f *intFlag) Reconcile() {
	if f.longVal != 0 {
		*f.dest = f.longVal
		return
	}
	if f.shortVal != 0 {
		*f.dest = f.shortVal
		return
	}
	*f.dest = f.defaultValue
}
