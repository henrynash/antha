// Code generated by gen.py. DO NOT EDIT.
package data

// float64

func (o Observation) MustFloat64() float64 {
	return o.value.(float64)
}

func (o Observation) Float64() (float64, bool) {
	if o.IsNull() {
		return float64(0), false
	} else {
		return o.MustFloat64(), true
	}
}

// int64

func (o Observation) MustInt64() int64 {
	return o.value.(int64)
}

func (o Observation) Int64() (int64, bool) {
	if o.IsNull() {
		return int64(0), false
	} else {
		return o.MustInt64(), true
	}
}

// string

func (o Observation) MustString() string {
	return o.value.(string)
}

func (o Observation) String() (string, bool) {
	if o.IsNull() {
		return "", false
	} else {
		return o.MustString(), true
	}
}

// bool

func (o Observation) MustBool() bool {
	return o.value.(bool)
}

func (o Observation) Bool() (bool, bool) {
	if o.IsNull() {
		return false, false
	} else {
		return o.MustBool(), true
	}
}
