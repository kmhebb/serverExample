package flag

import (
	"flag"
)

func NewFlagSet(name string) FlagSet {
	return FlagSet{
		fs: flag.NewFlagSet(name, flag.ContinueOnError),
	}
}

type FlagSet struct {
	fs    *flag.FlagSet
	flags []reconciler
}

func (fs *FlagSet) BoolVar(dest *bool, short, long string, defaultValue bool, description string) {
	boolf := &boolFlag{
		dest:         dest,
		short:        short,
		long:         long,
		defaultValue: defaultValue,
		description:  description,
	}
	if boolf.short != "" {
		fs.fs.BoolVar(&boolf.shortVal, boolf.short, false, boolf.description)
	}
	if boolf.long != "" {
		fs.fs.BoolVar(&boolf.longVal, boolf.long, false, boolf.description)
	}
	fs.flags = append(fs.flags, boolf)
}

func (fs *FlagSet) IntVar(dest *int, short, long string, defaultValue int, description string) {
	intf := &intFlag{
		dest:         dest,
		short:        short,
		long:         long,
		defaultValue: defaultValue,
		description:  description,
	}
	fs.fs.IntVar(&intf.shortVal, intf.short, 0, intf.description)
	fs.fs.IntVar(&intf.longVal, intf.long, 0, intf.description)
	fs.flags = append(fs.flags, intf)
}

func (fs *FlagSet) StringVar(dest *string, short, long string, defaultValue string, description string) {
	stringf := &stringFlag{
		dest:         dest,
		short:        short,
		long:         long,
		defaultValue: defaultValue,
		description:  description,
	}
	fs.fs.StringVar(&stringf.shortVal, stringf.short, "", stringf.description)
	fs.fs.StringVar(&stringf.longVal, stringf.long, "", stringf.description)
	fs.flags = append(fs.flags, stringf)
}

func (fs *FlagSet) Parse(args []string) {
	fs.fs.Parse(args)
	for _, f := range fs.flags {
		f.Reconcile()
	}
}

type reconciler interface {
	Reconcile()
}
