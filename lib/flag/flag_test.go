package flag_test

import (
	"testing"

	"github.com/kmhebb/serverExample/lib/flag"
)

func TestIntVar(t *testing.T) {
	defaultValue := 3

	for name, tc := range map[string]struct {
		input []string
		want  int
	}{
		"short": {
			input: []string{"-v", "4"},
			want:  4,
		},
		"long": {
			input: []string{"-value", "5"},
			want:  5,
		},
		"default": {
			input: []string{""},
			want:  defaultValue,
		},
		"long and short": {
			input: []string{"-v", "4", "-value", "5"},
			want:  5,
		},
		"short equals": {
			input: []string{"-v=4"},
			want:  4,
		},
		"long equals": {
			input: []string{"--value=5"},
			want:  5,
		},
	} {
		t.Run(name, func(t *testing.T) {
			var i int
			fs := flag.NewFlagSet("test")
			fs.IntVar(&i, "v", "value", defaultValue, "The value of the integer")
			fs.Parse(tc.input)
			if i != tc.want {
				t.Errorf("expected i=%d, but got i=%d", tc.want, i)
			}
		})
	}
}

func TestBoolVar(t *testing.T) {
	for name, tc := range map[string]struct {
		input []string
		want  bool
	}{
		"short implicit true": {
			input: []string{"-t"},
			want:  true,
		},
		"short explicit true": {
			input: []string{"-t", "true"},
			want:  true,
		},
		"long implicit true": {
			input: []string{"-test"},
			want:  true,
		},
		"long explicit true": {
			input: []string{"-test", "true"},
			want:  true,
		},
		"default": {
			input: []string{""},
			want:  false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			var b bool
			fs := flag.NewFlagSet("test")
			fs.BoolVar(&b, "t", "test", false, "The description")
			fs.Parse(tc.input)
			if b != tc.want {
				t.Errorf("expected b=%t, but got b=%t", tc.want, b)
			}
		})
	}
}

func TestStringVar(t *testing.T) {
	defaultValue := "def"

	for name, tc := range map[string]struct {
		input []string
		want  string
	}{
		"short": {
			input: []string{"-v", "short"},
			want:  "short",
		},
		"long": {
			input: []string{"-value", "long"},
			want:  "long",
		},
		"default": {
			input: []string{""},
			want:  defaultValue,
		},
		"long and short": {
			input: []string{"-v", "short", "-value", "long"},
			want:  "long",
		},
		"short equals": {
			input: []string{"-v=short"},
			want:  "short",
		},
		"long equals": {
			input: []string{"--value=long"},
			want:  "long",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var s string
			fs := flag.NewFlagSet("test")
			fs.StringVar(&s, "v", "value", defaultValue, "The value of the integer")
			fs.Parse(tc.input)
			if s != tc.want {
				t.Errorf("expected s=%q, but got s=%q", tc.want, s)
			}
		})
	}
}
