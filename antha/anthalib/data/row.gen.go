package data

// Code generated by gen.py. DO NOT EDIT.

// float64

// MustFloat64 extracts a value of type float64 from a Value. Panics if the value is null or the types mismatch.
func (v Value) MustFloat64() float64 {
	return v.value.(float64)
}

// Float64 extracts a value of type float64 from a Value. Returns false if the value is null. Panics if the types mismatch.
func (v Value) Float64() (float64, bool) {
	if v.IsNull() {
		return float64(0), false
	} else {
		return v.MustFloat64(), true
	}
}

// int64

// MustInt64 extracts a value of type int64 from a Value. Panics if the value is null or the types mismatch.
func (v Value) MustInt64() int64 {
	return v.value.(int64)
}

// Int64 extracts a value of type int64 from a Value. Returns false if the value is null. Panics if the types mismatch.
func (v Value) Int64() (int64, bool) {
	if v.IsNull() {
		return int64(0), false
	} else {
		return v.MustInt64(), true
	}
}

// int

// MustInt extracts a value of type int from a Value. Panics if the value is null or the types mismatch.
func (v Value) MustInt() int {
	return v.value.(int)
}

// Int extracts a value of type int from a Value. Returns false if the value is null. Panics if the types mismatch.
func (v Value) Int() (int, bool) {
	if v.IsNull() {
		return 0, false
	} else {
		return v.MustInt(), true
	}
}

// string

// MustString extracts a value of type string from a Value. Panics if the value is null or the types mismatch.
func (v Value) MustString() string {
	return v.value.(string)
}

// String extracts a value of type string from a Value. Returns false if the value is null. Panics if the types mismatch.
func (v Value) String() (string, bool) {
	if v.IsNull() {
		return "", false
	} else {
		return v.MustString(), true
	}
}

// bool

// MustBool extracts a value of type bool from a Value. Panics if the value is null or the types mismatch.
func (v Value) MustBool() bool {
	return v.value.(bool)
}

// Bool extracts a value of type bool from a Value. Returns false if the value is null. Panics if the types mismatch.
func (v Value) Bool() (bool, bool) {
	if v.IsNull() {
		return false, false
	} else {
		return v.MustBool(), true
	}
}

// TimestampMillis

// MustTimestampMillis extracts a value of type TimestampMillis from a Value. Panics if the value is null or the types mismatch.
func (v Value) MustTimestampMillis() TimestampMillis {
	return v.value.(TimestampMillis)
}

// TimestampMillis extracts a value of type TimestampMillis from a Value. Returns false if the value is null. Panics if the types mismatch.
func (v Value) TimestampMillis() (TimestampMillis, bool) {
	if v.IsNull() {
		return TimestampMillis(0), false
	} else {
		return v.MustTimestampMillis(), true
	}
}

// TimestampMicros

// MustTimestampMicros extracts a value of type TimestampMicros from a Value. Panics if the value is null or the types mismatch.
func (v Value) MustTimestampMicros() TimestampMicros {
	return v.value.(TimestampMicros)
}

// TimestampMicros extracts a value of type TimestampMicros from a Value. Returns false if the value is null. Panics if the types mismatch.
func (v Value) TimestampMicros() (TimestampMicros, bool) {
	if v.IsNull() {
		return TimestampMicros(0), false
	} else {
		return v.MustTimestampMicros(), true
	}
}
