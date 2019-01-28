
package data
// Code generated by gen.py. DO NOT EDIT.

import (
	"reflect"
	"github.com/pkg/errors"
)

/*
 * specializations for more efficient Extend operations
 */

// Float64 adds a float64 col using float64 inputs.  Null on any null inputs.
func (e *ExtendOn) Float64(f func(v ...float64) float64) *Table {
	// TODO move from lazy to eager type validation
	return NewTable(append(e.extension.series, &Series{
		col: e.extension.newCol,
		typ: reflect.TypeOf(float64(0)),
		read: func(cache *seriesIterCache) iterator {
			// Every series must be cast or converted
			colReader := make([]iterFloat64, len(e.inputs))
			var err error
			for i, ser := range e.inputs {
				iter := cache.Ensure(ser)
				colReader[i], err = ser.iterateFloat64(iter) // note colReader[i] is not itself in the cache!
				if err != nil {
					// TODO non-panic option?
					// TODO test coverage
					panic(errors.Wrapf(err, "when projecting new column %v", e.extension.newCol))
				}
			}
			// end when table exhausted
			e.extension.extensionSource(cache)
			return &extendFloat64{f: f, source: colReader}
		}},
	))
}

var _ iterFloat64 = (*extendFloat64)(nil)

type extendFloat64 struct {
	f      func(v ...float64) float64
	source []iterFloat64
}

func (x *extendFloat64) Next() bool {
	return true
}
func (x *extendFloat64) Value() interface{} {
	v, ok := x.Float64()
	if !ok {
		return nil
	}
	return v
}
func (x *extendFloat64) Float64() (float64, bool) {
	args := make([]float64, len(x.source))
	var ok bool
	for i, s := range x.source {
		args[i], ok = s.Float64()
		if !ok {
			return float64(0), false
		}
	}
	v := x.f(args...)
	return v, true
}

// Int64 adds a int64 col using int64 inputs.  Null on any null inputs.
func (e *ExtendOn) Int64(f func(v ...int64) int64) *Table {
	// TODO move from lazy to eager type validation
	return NewTable(append(e.extension.series, &Series{
		col: e.extension.newCol,
		typ: reflect.TypeOf(int64(0)),
		read: func(cache *seriesIterCache) iterator {
			// Every series must be cast or converted
			colReader := make([]iterInt64, len(e.inputs))
			var err error
			for i, ser := range e.inputs {
				iter := cache.Ensure(ser)
				colReader[i], err = ser.iterateInt64(iter) // note colReader[i] is not itself in the cache!
				if err != nil {
					// TODO non-panic option?
					// TODO test coverage
					panic(errors.Wrapf(err, "when projecting new column %v", e.extension.newCol))
				}
			}
			// end when table exhausted
			e.extension.extensionSource(cache)
			return &extendInt64{f: f, source: colReader}
		}},
	))
}

var _ iterInt64 = (*extendInt64)(nil)

type extendInt64 struct {
	f      func(v ...int64) int64
	source []iterInt64
}

func (x *extendInt64) Next() bool {
	return true
}
func (x *extendInt64) Value() interface{} {
	v, ok := x.Int64()
	if !ok {
		return nil
	}
	return v
}
func (x *extendInt64) Int64() (int64, bool) {
	args := make([]int64, len(x.source))
	var ok bool
	for i, s := range x.source {
		args[i], ok = s.Int64()
		if !ok {
			return int64(0), false
		}
	}
	v := x.f(args...)
	return v, true
}

// String adds a string col using string inputs.  Null on any null inputs.
func (e *ExtendOn) String(f func(v ...string) string) *Table {
	// TODO move from lazy to eager type validation
	return NewTable(append(e.extension.series, &Series{
		col: e.extension.newCol,
		typ: reflect.TypeOf(""),
		read: func(cache *seriesIterCache) iterator {
			// Every series must be cast or converted
			colReader := make([]iterString, len(e.inputs))
			var err error
			for i, ser := range e.inputs {
				iter := cache.Ensure(ser)
				colReader[i], err = ser.iterateString(iter) // note colReader[i] is not itself in the cache!
				if err != nil {
					// TODO non-panic option?
					// TODO test coverage
					panic(errors.Wrapf(err, "when projecting new column %v", e.extension.newCol))
				}
			}
			// end when table exhausted
			e.extension.extensionSource(cache)
			return &extendString{f: f, source: colReader}
		}},
	))
}

var _ iterString = (*extendString)(nil)

type extendString struct {
	f      func(v ...string) string
	source []iterString
}

func (x *extendString) Next() bool {
	return true
}
func (x *extendString) Value() interface{} {
	v, ok := x.String()
	if !ok {
		return nil
	}
	return v
}
func (x *extendString) String() (string, bool) {
	args := make([]string, len(x.source))
	var ok bool
	for i, s := range x.source {
		args[i], ok = s.String()
		if !ok {
			return "", false
		}
	}
	v := x.f(args...)
	return v, true
}

// Bool adds a bool col using bool inputs.  Null on any null inputs.
func (e *ExtendOn) Bool(f func(v ...bool) bool) *Table {
	// TODO move from lazy to eager type validation
	return NewTable(append(e.extension.series, &Series{
		col: e.extension.newCol,
		typ: reflect.TypeOf(false),
		read: func(cache *seriesIterCache) iterator {
			// Every series must be cast or converted
			colReader := make([]iterBool, len(e.inputs))
			var err error
			for i, ser := range e.inputs {
				iter := cache.Ensure(ser)
				colReader[i], err = ser.iterateBool(iter) // note colReader[i] is not itself in the cache!
				if err != nil {
					// TODO non-panic option?
					// TODO test coverage
					panic(errors.Wrapf(err, "when projecting new column %v", e.extension.newCol))
				}
			}
			// end when table exhausted
			e.extension.extensionSource(cache)
			return &extendBool{f: f, source: colReader}
		}},
	))
}

var _ iterBool = (*extendBool)(nil)

type extendBool struct {
	f      func(v ...bool) bool
	source []iterBool
}

func (x *extendBool) Next() bool {
	return true
}
func (x *extendBool) Value() interface{} {
	v, ok := x.Bool()
	if !ok {
		return nil
	}
	return v
}
func (x *extendBool) Bool() (bool, bool) {
	args := make([]bool, len(x.source))
	var ok bool
	for i, s := range x.source {
		args[i], ok = s.Bool()
		if !ok {
			return false, false
		}
	}
	v := x.f(args...)
	return v, true
}

