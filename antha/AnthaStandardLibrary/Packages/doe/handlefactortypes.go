// Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package doe

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/solutions"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
)

// HandleConcFactor parses a factor name and value and returns an antha concentration.
// If the value cannot be converted to a valid concentration an error is returned.
// If the header contains a valid concentration unit a number can be specified as the value.
func HandleConcFactor(header string, value interface{}) (anthaConc wunit.Concentration, err error) {

	var floatValue float64
	var floatFound bool

	rawconcfloat, found := value.(float64)

	if found {
		floatValue = rawconcfloat
		floatFound = true
	} else {
		rawconcstring, found := value.(string)
		var floatParseErr error
		if floatValue, floatParseErr = strconv.ParseFloat(rawconcstring, 64); found && floatParseErr == nil {
			floatFound = true
		}
	}

	if floatFound {

		// handle floating point imprecision
		floatValue, err = wutil.Roundto(floatValue, 6)

		if err != nil {
			return anthaConc, err
		}
		containsconc, conc, _ := wunit.ParseConcentration(header)

		if containsconc {

			concunit := conc.Unit().PrefixedSymbol()

			anthaConc = wunit.NewConcentration(floatValue, concunit)
		} else {
			err = fmt.Errorf("No valid conc found in component %s so can't assign a concentration unit to value", header)
			return anthaConc, err
		}

	} else if rawconcstring, found := value.(string); found {

		containsconc, conc, _ := wunit.ParseConcentration(rawconcstring)

		if containsconc {
			anthaConc = conc
		} else {
			err = fmt.Errorf("No valid conc found in %s", rawconcstring)
			return anthaConc, err
		}

		// if float use conc unit from header component
	} else {
		err = fmt.Errorf("problem with type of %T expected string or float", value)
		return anthaConc, err
	}

	return
}

// HandleComponentWithConcentration returns both LHComponent and Concentration from a component name with concentration in a DOE design.
// If no valid concentration is found or an invalid component name is specifed an error is returned.
func HandleComponentWithConcentration(ctx context.Context, header string, value interface{}) (component *wtype.LHComponent, concentration wunit.Concentration, err error) {

	concentration, err = HandleConcFactor(header, value)

	if err != nil {
		return
	}

	componentName := solutions.NormaliseName(header)

	component, err = inventory.NewComponent(ctx, componentName)

	if err == nil {
		// continue
	} else if err == inventory.ErrUnknownType {
		component, err = inventory.NewComponent(ctx, inventory.WaterType)
		if err != nil {
			return
		}
		component.CName = componentName
	} else {
		return
	}

	component.SetConcentration(concentration)

	return
}

// HandleVolumeFactor parses a factor name and value and returns an antha Volume.
// If the value cannot be converted to a valid Volume an error is returned.
func HandleVolumeFactor(header string, value interface{}) (anthaVolume wunit.Volume, err error) {

	var floatValue float64
	var floatFound bool
	var volUnit string

	rawVolFloat, found := value.(float64)

	if found {
		floatValue = rawVolFloat
		floatFound = true
	} else if rawVolInt, intFound := value.(int); intFound {
		floatValue = float64(rawVolInt)
		floatFound = true
	} else {
		rawvolstring, found := value.(string)
		var floatParseErr error
		if floatValue, floatParseErr = strconv.ParseFloat(rawvolstring, 64); found && floatParseErr == nil {
			floatFound = true
		}
	}
	if floatFound {

		// handle floating point imprecision
		floatValue, err = wutil.Roundto(floatValue, 6)

		if err != nil {
			return anthaVolume, err
		}

		fields := strings.Fields(header)

		for _, field := range fields {
			if _, validUnitFound := wunit.UnitMap["Volume"][strings.Trim(field, "()")]; validUnitFound {
				volUnit = strings.Trim(field, "()")
			}
		}

		if volUnit == "" {
			volUnit = "ul"
		}

		anthaVolume = wunit.NewVolume(floatValue, volUnit)

	} else if rawVolString, found := value.(string); found {

		vol, err := wunit.ParseVolume(rawVolString)

		if err == nil {
			anthaVolume = vol
		} else {
			err = fmt.Errorf("No valid Volume found in %s: %s", rawVolString, err.Error())
			return anthaVolume, err
		}

	} else {
		err = fmt.Errorf("problem with type of %v expected string or float", value)
		return anthaVolume, err
	}

	return
}

// HandleLHComponentFactor parses a factor name and value and returns an
// LHComponent.
//
// If the value cannot be converted to a valid component an error is returned.
func HandleLHComponentFactor(ctx context.Context, header string, value interface{}) (*wtype.LHComponent, error) {
	str, found := value.(string)
	if !found {
		return nil, fmt.Errorf("value %T is not a string", value)
	}

	component, err := inventory.NewComponent(ctx, str)
	if err == nil {
		return component, nil
	}

	if err == inventory.ErrUnknownType {
		component, err = inventory.NewComponent(ctx, inventory.WaterType)
		component.CName = str
		return component, err
	}

	concentration, concErr := HandleConcFactor(header, value)

	if concErr == nil {
		component.SetConcentration(concentration)
	}

	return nil, err
}

// HandleLHPlateFactor parses a factor name and value and returns an
// LHComponent.
//
// If the value cannot be converted to a valid component an error is returned.
func HandleLHPlateFactor(ctx context.Context, header string, value interface{}) (*wtype.LHPlate, error) {
	str, found := value.(string)
	if !found {
		return nil, fmt.Errorf("value %T is not a string", value)
	}

	return inventory.NewPlate(ctx, str)
}

// HandleTemperatureFactor parses a factor name and value and returns an antha Temperature.
// If the value cannot be converted to a valid Temperature an error is returned.
// A float or int value with no unit is assumed to be in C.
func HandleTemperatureFactor(header string, value interface{}) (anthaTemp wunit.Temperature, err error) {

	defaultUnit, err := lookForUnitInHeader(header, "Temperature")

	if err != nil {
		defaultUnit = "C"
	}

	switch temp := value.(type) {
	case int:
		anthaTemp = wunit.NewTemperature(float64(temp), defaultUnit)
		return
	case float64:
		anthaTemp = wunit.NewTemperature(temp, defaultUnit)
		return
	case string:
		value, unit := wunit.SplitValueAndUnit(temp)

		if unit == "" {
			unit = defaultUnit
		}

		err = wunit.ValidMeasurementUnit("Temperature", unit)

		if err != nil {
			return
		}

		anthaTemp = wunit.NewTemperature(value, unit)

		return
	default:
		return anthaTemp, fmt.Errorf("cannot convert %v of type %T to temperature!", value, temp)
	}
}

// HandleTimeFactor parses a factor name and value and returns an antha Time.
// If the value cannot be converted to a valid Time an error is returned.
// A float or int value with no unit is assumed to be in s.
func HandleTimeFactor(header string, value interface{}) (anthaTime wunit.Time, err error) {

	defaultUnit, err := lookForUnitInHeader(header, "Time")

	if err != nil {
		defaultUnit = "s"
	}

	switch time := value.(type) {
	case int:
		anthaTime = wunit.NewTime(float64(time), defaultUnit)
		return
	case float64:
		anthaTime = wunit.NewTime(time, defaultUnit)
		return
	case string:
		value, unit := wunit.SplitValueAndUnit(time)

		if unit == "" {
			unit = defaultUnit
		}

		err = wunit.ValidMeasurementUnit("Time", unit)

		if err != nil {
			return
		}

		anthaTime = wunit.NewTime(value, unit)
		return
	default:
		return anthaTime, fmt.Errorf("cannot convert %v of type %T to time!", value, time)
	}
}

// HandleRPMFactor parses a factor name and value and returns an antha Rate.
// If the value cannot be converted to a valid Rate an error is returned.
// A float or int value with no unit is assumed to be in /min.
func HandleRPMFactor(header string, value interface{}) (anthaRate wunit.Rate, err error) {

	defaultUnit, err := lookForUnitInHeader(header, "RPM")

	if err != nil {
		defaultUnit = "/min"
	}

	switch rate := value.(type) {
	case int:
		anthaRate, err = wunit.NewRate(float64(rate), defaultUnit)
		return
	case float64:
		anthaRate, err = wunit.NewRate(rate, defaultUnit)
		return
	case string:
		value, unit := wunit.SplitValueAndUnit(rate)

		if unit == "" {
			unit = defaultUnit
		}

		err = wunit.ValidMeasurementUnit("Rate", unit)

		if err != nil {
			return
		}

		anthaRate, err = wunit.NewRate(value, unit)
		return
	default:
		return anthaRate, fmt.Errorf("cannot convert %v of type %T to RPM!", value, rate)
	}
}

// lookForUnitInHeader searches for a unit flanked by ( ).
// If a measurment type is specified the unit will be checked for validity.
func lookForUnitInHeader(header, measurementType string) (unit string, err error) {

	var errs string

	fields := strings.Fields(header)

	for _, field := range fields {

		if strings.HasPrefix(field, "(") && strings.HasSuffix(field, ")") {
			trimmer := strings.Trim(field, "()")

			if measurementType != "" {
				err = wunit.ValidMeasurementUnit(measurementType, unit)

				if err == nil {
					unit = trimmed
					return
				}
			} else {
				unit = trimmed
			}
			errs = append(errs, err.Error())
		}

	}

	if len(errs) > 0 {
		return "", fmt.Errorf(strings.Join(errs, ";"))
	}

	if measurementType != "" {
		return "", fmt.Errorf("no unit found in header %s of type %s", header, measurementType)
	}
	return "", fmt.Errorf("no unit found in header %s", header)
}
