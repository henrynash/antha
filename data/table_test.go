package data

import (
	"fmt"
	"reflect"
	"testing"
	// TODO "github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	testEquals(t, nativeSeries)
	testEquals(t, arrowSeries)
}

func testEquals(t *testing.T, typ seriesType) {
	tab := NewTable([]*Series{
		newSeries(typ, "measure", []int64{1, 1000}),
		newSeries(typ, "label", []string{"abcdef", "abcd"}),
	})
	assertEqual(t, tab, tab, "not self equal")

	tab2 := NewTable([]*Series{
		newSeries(typ, "measure", []int64{1, 1000}),
	})
	assertEqual(t, tab2, tab.Project("measure"), "not equal by value")

	if tab2.Equal(tab.Project("label")) {
		t.Error("equal with mismatched schema")
	}

	if tab2.Equal(tab2.Filter(Eq("measure", 1000))) {
		t.Error("equal with mismatched data")
	}
}

func assertEqual(t *testing.T, expected, actual *Table, msg string) {
	if !actual.Equal(expected) {
		t.Error(msg)
		t.Log("actual", actual.ToRows())
	}
}

func TestSlice(t *testing.T) {
	testSlice(t, nativeSeries)
	testSlice(t, arrowSeries)
}

func testSlice(t *testing.T, typ seriesType) {
	a := NewTable([]*Series{
		newSeries(typ, "a", []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
	})
	assertEqual(t, a, a.Slice(0, 100), "slice all")

	slice00 := a.Slice(1, 1)
	assertEqual(t, NewTable([]*Series{
		newSeries(typ, "a", []int64{}),
	}), slice00, "slice00")

	slice04 := a.Head(4)
	assertEqual(t, NewTable([]*Series{
		newSeries(typ, "a", []int64{1, 2, 3, 4}),
	}), slice04, "slice04")

	slice910 := a.Slice(9, 10)
	assertEqual(t, NewTable([]*Series{
		newSeries(typ, "a", []int64{10}),
	}), slice910, "slice910")
}

func TestExtend(t *testing.T) {
	testExtend(t, nativeSeries)
	testExtend(t, arrowSeries)
}

func testExtend(t *testing.T, typ seriesType) {
	a := NewTable([]*Series{
		newSeries(typ, "a", []int64{1, 2, 3}),
	})
	extended := a.Extend("e").By(func(r Row) interface{} {
		a, _ := r.Observation("a")
		return float64(a.MustInt64()) / 2.0
	},
		reflect.TypeOf(float64(0)))
	assertEqual(t, NewTable([]*Series{
		newSeries(typ, "e", []float64{0.5, 1.0, 1.5}),
	}), extended.Project("e"), "extend")

	floats := NewTable([]*Series{
		newSeries(typ, "floats", []float64{1, 2, 3}),
	})
	extendedStatic := floats.
		Extend("e_static").
		On("floats").
		Float64(func(v ...float64) float64 {
			return v[0] * 2.0
		})
	// extended2 := extendedStatic.Extend("another").On("floats", "e_static").Float64(func(v ...float64) float64 { return v[0] + v[1] })

	// fmt.Println(extended2.ToRows())
	assertEqual(t, NewTable([]*Series{
		Must().NewSliceSeries("e_static", []float64{2, 4, 6}),
	}), extendedStatic.Project("e_static"), "extend static")
}

func TestFilterEq(t *testing.T) {
	testFilterEq(t, nativeSeries)
	testFilterEq(t, arrowSeries)
}

func testFilterEq(t *testing.T, typ seriesType) {
	a := NewTable([]*Series{
		newSeries(typ, "a", []int64{1, 2, 3}),
	})
	filtered := a.Filter(Eq("a", 2))
	assertEqual(t, filtered, a.Slice(1, 2), "filter")
}

func TestSize(t *testing.T) {
	testSize(t, nativeSeries)
	testSize(t, arrowSeries)
}

func testSize(t *testing.T, typ seriesType) {
	empty := NewTable([]*Series{})
	if empty.Size() != 0 {
		t.Errorf("should be empty. %d", empty.Size())
	}
	a := NewTable([]*Series{
		newSeries(typ, "a", []int64{1, 2, 3}),
	})
	if a.Size() != 3 {
		t.Errorf("size? %d", a.Size())
	}
	// a filter is of unbounded size
	filtered := a.Filter(Eq("a", 1))
	if filtered.Size() != -1 {
		t.Errorf("filtered.Size()? %d", filtered.Size())
	}
	// a slice is of bounded size as long as its dependencies are
	slice1 := filtered.Head(1)
	if slice1.Size() != -1 {
		t.Errorf(" slice1.Size()? %d", slice1.Size())
	}
	if a.Head(0).Size() != 0 {
		t.Errorf("a.Head(0).Size()? %d", a.Head(0).Size())
	}
	slice2 := a.Slice(1, 4)
	if slice2.Size() != 2 {
		t.Errorf("slice2.Size()? %d", slice2.Size())
	}
}

type seriesType int

const (
	nativeSeries seriesType = iota
	arrowSeries
)

func newSeries(typ seriesType, col ColumnName, values interface{}) *Series {
	switch typ {
	case nativeSeries:
		return Must().NewSliceSeries(col, values)
	case arrowSeries:
		return Must().NewArrowSeriesFromSlice(col, values, nil)
	default:
		panic(fmt.Errorf("Unknown series type: %d", typ))
	}
}