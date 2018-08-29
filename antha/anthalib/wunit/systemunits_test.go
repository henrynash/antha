package wunit

import (
	"math"
	"testing"
)

func TestMetersCubed(t *testing.T) {
	reg := makeGlobalUnitRegistry()

	//test that we support m^3 and correct prefix exponent
	if volume, err := reg.NewMeasurement(1.0, "mm^3"); err != nil {
		t.Error(err)
	} else if ul, err := reg.GetUnit("ul"); err != nil {
		t.Error(err)
	} else if volumeInUl, err := volume.ConvertTo(ul); err != nil {
		t.Error(err)
	} else if g, e := volumeInUl.RawValue(), 1.0; math.Abs(g-e) > 1.0e-6 {
		t.Errorf("converting %v to ul: expected %g, got %g", volume, e, g)
	}
}

func TestMililitresPerMinute(t *testing.T) {
	reg := makeGlobalUnitRegistry()

	if fr, err := reg.NewMeasurement(60.0, "ml/min"); err != nil {
		t.Error(err)
	} else if ulPerSec, err := reg.GetUnit("ul/s"); err != nil {
		t.Error(err)
	} else if frInUlPerSec, err := fr.ConvertTo(ulPerSec); err != nil {
		t.Error(err)
	} else if g, e := frInUlPerSec.RawValue(), 1000.0; math.Abs(g-e) > 1.0e-9 {
		t.Errorf("converting %v to %v: got %g expected %g", fr, ulPerSec, g, e)
	}
}

func TestBadUnitConversion(t *testing.T) {

	type TestCase struct {
		Value      Measurement
		TargetUnit string
	}

	tests := []TestCase{
		{
			Value:      NewMeasurement(1.0, "g/l"),
			TargetUnit: "X",
		},
		{
			Value:      NewMeasurement(1.0, "um"),
			TargetUnit: "l",
		},
	}

	for _, test := range tests {
		t.Run(test.Value.ToString(), func(t *testing.T) {
			if _, err := test.Value.ConvertToByString(test.TargetUnit); err == nil {
				t.Error("failed to return error")
			}
		})
	}
}
