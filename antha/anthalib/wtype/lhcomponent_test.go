package wtype

import "testing"

func TestSampleBehaviour(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"

	c2 := NewLHComponent()
	c.CName = "cider"

	c.SetSample(true)

	if !c.IsSample() {
		t.Errorf("SetSample(true) must cause components to return true to IsSample()")
	}

	c2.SetSample(true)
	if !c2.IsSample() {
		t.Errorf("SetSample(true) must cause components to return true to IsSample()")
	}

	c.Mix(c2)

	if c.IsSample() {
		t.Errorf("Results of mixes must not be samples")
	}

	c2.SetSample(false)

	if c2.IsSample() {
		t.Errorf("SetSample(false) must cause components to return false to IsSample()")
	}

	c3 := c2.Dup()

	if c3.IsSample() {
		t.Errorf("Dup()ing a non-sample must produce a non-sample")
	}

	c3.SetSample(true)

	if !c3.IsSample() {
		t.Errorf("SetSample(true) must  cause components to return true to IsSample()... even duplicates")
	}

	if c2.IsSample() {
		t.Errorf("Duplicates must not remain linked")
	}
}
