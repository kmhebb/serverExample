// Package assert provides utilities for dealing with assertions.
package assert

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"
)

type ErrorPredicate func(err error) (string, bool)

// EmptyReporter is the interface used to determine if our value is empty.
type EmptyReporter interface {
	IsEmpty() bool
}

type ZeroReporter interface {
	IsZero() bool
}

// New creates a new Checker.
func New(t *testing.T) *Checker {
	return &Checker{t: t}
}

// Checker represents our testing structure.
type Checker struct {
	t      *testing.T
	failed bool
}

func (c *Checker) Empty(val interface{}) *Checker {
	c.t.Helper()
	c.failed = false

	switch v := val.(type) {
	case string:
		c.failed = (v != "")
	default:
		c.failed = true
		c.t.Errorf("unhandled type %T for val %v", val, val)
		return c
	}

	if c.failed {
		c.t.Errorf("expected val to be empty, but got %v", val)
	}

	return c
}

// Equals expects that the values are equal. Special handling of types
// (int64 and time.Time).
func (c *Checker) Equals(got, want interface{}) *Checker {
	c.t.Helper()
	c.failed = false

	switch w := want.(type) {
	case int64:
		g, ok := got.(int64)
		if !ok {
			c.failed = true
			c.t.Errorf("mismatched types.\nwanted:\n\t%T(%v)\ngot:\n\t%T(%v)", want, want, got, got)
			return c
		}
		if g != w {
			c.failed = true
		}
	case time.Time:
		g, ok := got.(time.Time)
		if !ok {
			c.failed = true
			c.t.Errorf("mismatched types.\nwanted:\n\t%T(%v)\ngot:\n\t%T(%v)", want, want, got, got)
			return c
		}
		if !g.Equal(w) {
			c.failed = true
		}
	default:
		if !reflect.DeepEqual(got, want) {
			c.failed = true
		}
	}

	if c.failed {
		c.t.Errorf("expected values to be equal.\nwanted:\n\t%T(%v)\ngot:\n\t%T(%v)", want, want, got, got)
	}

	return c
}

func (c *Checker) Error(err error, pred ErrorPredicate) *Checker {
	c.t.Helper()
	c.failed = false

	if err == nil {
		c.failed = true
		c.t.Error("expected error to be not nil")
		return c
	}

	if msg, ok := pred(err); !ok {
		c.failed = true
		c.t.Error(msg)
	}

	return c
}

func (c *Checker) ErrorMatches(err error, regex string) *Checker {
	c.t.Helper()
	c.failed = false
	return c.Error(err, errMatches(regex))
}

// False expects that the value is False.
func (c *Checker) False(pred bool) *Checker {
	c.t.Helper()
	c.failed = false

	if pred {
		c.failed = true
		c.t.Error("expected predicate to be false, but wasn't")
	}

	return c
}

// Fatal calls FailNow if the Checker has failed.
func (c *Checker) Fatal() {
	c.t.Helper()
	failed := c.failed
	c.failed = false

	if failed {
		c.t.FailNow()
	}
}

// NotEmpty expects that the value is not empty. Will check if the type implements
// EmptyReporter and, if so, will use EmptyReporter to determine if the value is
// empty. It will otherwise fail unless the value is a string and empty.
func (c *Checker) NotEmpty(val interface{}) *Checker {
	c.t.Helper()
	c.failed = false

	switch v := val.(type) {
	case EmptyReporter:
		c.failed = v.IsEmpty()
	case ZeroReporter:
		c.failed = v.IsZero()
	case string:
		c.failed = (v == "")
	default:
		c.failed = true
		c.t.Errorf("unhandled type %T for val %v", val, val)
		return c
	}

	if c.failed {
		c.t.Error("expected val to be not be empty, but was")
	}

	return c
}

// NotNil verifies that the error is not nil.
func (c *Checker) NotNil(err error) *Checker {
	c.t.Helper()
	c.failed = false

	if err == nil {
		c.failed = true
		c.t.Error("expected err to be not nil")
	}

	return c
}

// OK verifies that the error is nil.
func (c *Checker) OK(err error) *Checker {
	c.t.Helper()
	c.failed = false

	if err != nil {
		c.failed = true
		c.t.Error(err)
	}

	return c
}

// Then calls the encapsulated function if all previous assertions have passed.
func (c *Checker) Then(f func()) *Checker {
	c.t.Helper()

	failed := c.failed
	if !failed {
		f()
	}
	c.failed = failed

	return c
}

// True verifies the predicate to be true.
func (c *Checker) True(pred bool) *Checker {
	c.t.Helper()
	c.failed = false

	if !pred {
		c.failed = true
		c.t.Error("expected predicate to be true, but wasn't")
	}

	return c
}

func errMatches(regex string) ErrorPredicate {
	return func(err error) (string, bool) {
		matched, matchErr := regexp.MatchString(regex, err.Error())
		if matchErr != nil {
			return matchErr.Error(), false
		}
		return fmt.Sprintf("expected %q to match %q but didn't", err.Error(), regex), matched
	}
}
