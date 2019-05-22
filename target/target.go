// Package target provides the construction of a target machine from a
// collection of devices
package target

import (
	"fmt"

	"github.com/antha-lang/antha/instructions"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/workflow"
)

const (
	// DriverSelectorV1Name is the basic selector name for device plugins
	// (drivers)
	DriverSelectorV1Name = "antha.driver.v1.TypeReply.type"
)

// Well known device plugins (drivers) selectors
var (
	DriverSelectorV1Human = instructions.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.human.v1.Human",
	}
	DriverSelectorV1ShakerIncubator = instructions.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.shakerincubator.v1.ShakerIncubator",
	}
	DriverSelectorV1Mixer = instructions.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.mixer.v1.Mixer",
	}
	DriverSelectorV1Prompter = instructions.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.prompter.v1.Prompter",
	}
	DriverSelectorV1DataSource = instructions.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.datasource.v1.DataSource",
	}
	DriverSelectorV1WriteOnlyPlateReader = instructions.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.platereader.v1.PlateReader",
	}
	DriverSelectorV1QPCRDevice = instructions.NameValue{
		Name:  DriverSelectorV1Name,
		Value: "antha.quantstudio.v1.QuantStudioService",
	}
)

// Target for execution (collection of devices).
type Target struct {
	Devices map[workflow.DeviceInstanceID]effects.Device
}

// New creates a new target
func New() *Target {
	return &Target{
		Devices: make(map[workflow.DeviceInstanceID]effects.Device),
	}
}

func (a *Target) canCompile(d effects.Device, reqs ...instructions.Request) bool {
	for _, req := range reqs {
		if !d.CanCompile(req) {
			return false
		}
	}
	return true
}

// CanCompile returns the devices that can compile the given set of requests
func (a *Target) CanCompile(reqs ...instructions.Request) (r []effects.Device) {
	for _, d := range a.Devices {
		if a.canCompile(d, reqs...) {
			r = append(r, d)
		}
	}
	return
}

// AddDevice adds a device to the target configuration
func (a *Target) AddDevice(d effects.Device) error {
	id := d.Id()
	if _, found := a.Devices[id]; found {
		return fmt.Errorf("Device with id %v already added", id)
	}
	a.Devices[id] = d
	return nil
}

func (a *Target) Close() {
	for _, dev := range a.Devices {
		dev.Close()
	}
}

func (a *Target) Connect(wf *workflow.Workflow) error {
	for _, dev := range a.Devices {
		if err := dev.Connect(wf); err != nil {
			return err
		}
	}
	return nil
}
