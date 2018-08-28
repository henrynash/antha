package wunit

import "math"

type baseUnit struct {
	Name     string
	Symbol   string
	SISymbol string //the canonincal form for the unit which can include a prefix, defaults to Symbol
	Prefixes []string
	Exponent int
}

type baseUnits map[string][]baseUnit

func (self baseUnits) AddTo(reg *UnitRegistry) error {
	for mType, defs := range self {
		for _, unit := range defs {
			SISymbol := unit.SISymbol
			if SISymbol == "" {
				SISymbol = unit.Symbol
			}
			if err := reg.DeclareUnit(mType, unit.Name, unit.Symbol, SISymbol, unit.Prefixes, unit.Exponent); err != nil {
				return err
			}
		}
	}
	return nil
}

var systemUnits = baseUnits{
	"Concentration": {
		{
			Name:     "grams per litre",
			Symbol:   "g/l",
			SISymbol: "kg/l",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
		{
			Name:     "moles per litre",
			Symbol:   "Mol/l",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
		{
			Name:     "units per litre",
			Symbol:   "U/l",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
		{
			Name:   "relative concentration",
			Symbol: "X",
		},
		{
			Name:   "volume ratio",
			Symbol: "v/v",
		},
	},
	"Volume": {
		{
			Name:     "litre",
			Symbol:   "l",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Mass": {
		{
			Name:     "gram",
			Symbol:   "g",
			SISymbol: "kg",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Density": {
		{
			Name:     "grams per meter cubed",
			Symbol:   "g/m^3",
			SISymbol: "kg/m^3",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Length": {
		{
			Name:     "metre",
			Symbol:   "m",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Area": {
		{
			Name:     "metre squared",
			Symbol:   "m^2",
			Prefixes: allPrefixes,
			Exponent: 2,
		},
	},
	"Temperature": {
		{
			Name:   "celsius",
			Symbol: "℃",
		},
	},
	"Time": {
		{
			Name:     "seconds",
			Symbol:   "s",
			Prefixes: []string{"y", "z", "a", "f", "p", "n", "u", "m"},
			Exponent: 1,
		},
	},
	"Moles": {
		{
			Name:     "moles",
			Symbol:   "Mol",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Angle": {
		{
			Name:   "radians",
			Symbol: "rad",
		},
	},
	"AngularVelocity": {
		{
			Name:   "radians per second",
			Symbol: "rad/s",
		},
	},
	"Energy": {
		{
			Name:     "joules",
			Symbol:   "J",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Force": {
		{
			Name:     "newtons",
			Symbol:   "N",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Pressure": {
		{
			Name:     "pascals",
			Symbol:   "Pa",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"SpecificHeatCapacity": {
		{
			Name:     "joules per kilogram per degrees celsius",
			Symbol:   "J/kg*C",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Velocity": {
		{
			Name:     "meters per second",
			Symbol:   "m/s",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"FlowRate": {
		{
			Name:     "litres per second",
			Symbol:   "l/s",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
	"Rate": {
		{
			Name:   "per second",
			Symbol: "/s",
		},
	},
	"Voltage": {
		{
			Name:     "volts",
			Symbol:   "V",
			Prefixes: allPrefixes,
			Exponent: 1,
		},
	},
}

type derivedUnit struct {
	Name             string
	Symbol           string
	Prefixes         []string
	Exponent         int
	TargetSymbol     string
	NewUnitInTargets float64
}

type derivedUnits map[string][]derivedUnit

func (self derivedUnits) AddTo(reg *UnitRegistry) error {
	for mType, defs := range self {
		for _, du := range defs {
			if err := reg.DeclareDerivedUnit(mType, du.Name, du.Symbol, du.Prefixes, du.Exponent, du.TargetSymbol, du.NewUnitInTargets); err != nil {
				return err
			}
		}
	}
	return nil
}

var systemDerivedUnits = derivedUnits{
	"Concentration": {
		{
			Name:             "grams per mililitre",
			Symbol:           "g/ml",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "kg/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "grams per microlitre",
			Symbol:           "g/ul",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "Mg/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "grams per nanolitre",
			Symbol:           "g/nl",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "Gg/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "moles per mililitre",
			Symbol:           "Mol/ml",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "kMol/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "moles per microlitre",
			Symbol:           "Mol/ul",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "MMol/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "moles per nanolitre",
			Symbol:           "Mol/nl",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "GMol/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "units per mililitre",
			Symbol:           "U/ml",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "kU/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "units per microlitre",
			Symbol:           "U/ul",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "MU/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "units per nanolitre",
			Symbol:           "U/nl",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "GU/l",
			NewUnitInTargets: 1.0,
		},
		{
			Name:             "percentage weight of solution",
			Symbol:           "% w/v",
			TargetSymbol:     "g/l",
			NewUnitInTargets: 10.0,
		},
	},
	"Volume": {
		{
			Name:             "meters cubed",
			Symbol:           "m^3",
			Prefixes:         allPrefixes,
			Exponent:         3,
			TargetSymbol:     "l",
			NewUnitInTargets: 1000.0,
		},
	},
	"Length": {
		{
			Name:             "inches",
			Symbol:           "in",
			TargetSymbol:     "mm",
			NewUnitInTargets: 25.4,
		},
	},
	"Time": {
		{
			Name:             "minutes",
			Symbol:           "min",
			TargetSymbol:     "s",
			NewUnitInTargets: 60.0,
		},
		{
			Name:             "hours",
			Symbol:           "h",
			TargetSymbol:     "s",
			NewUnitInTargets: 3600.0,
		},
		{
			Name:             "days",
			Symbol:           "days",
			TargetSymbol:     "s",
			NewUnitInTargets: 86400.0,
		},
	},
	"Angle": {
		{
			Name:             "degrees",
			Symbol:           "°",
			TargetSymbol:     "rad",
			NewUnitInTargets: (2.0 * math.Pi) / 360.0,
		},
	},
	"AngularVelocity": {
		{
			Name:             "radians per minute",
			Symbol:           "rad/min",
			TargetSymbol:     "rad/s",
			NewUnitInTargets: 1.0 / 60.0,
		},
		{
			Name:             "revolutions per minute",
			Symbol:           "rpm",
			TargetSymbol:     "rad/s",
			NewUnitInTargets: 2.0 * math.Pi / 60.0,
		},
	},
	"Pressure": {
		{
			Name:             "bar",
			Symbol:           "bar",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "kPa",
			NewUnitInTargets: 100.0,
		},
	},
	"FlowRate": {
		{
			Name:             "litres per minute",
			Symbol:           "l/min",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "l/s",
			NewUnitInTargets: 1 / 60.0,
		},
		{
			Name:             "litres per hour",
			Symbol:           "l/h",
			Prefixes:         allPrefixes,
			Exponent:         1,
			TargetSymbol:     "l/s",
			NewUnitInTargets: 1 / 3600.0,
		},
	},
	"Rate": {
		{
			Name:             "per minute",
			Symbol:           "/min",
			TargetSymbol:     "/s",
			NewUnitInTargets: 1 / 60.0,
		},
		{
			Name:             "per hour",
			Symbol:           "/h",
			TargetSymbol:     "/s",
			NewUnitInTargets: 1 / 3600.0,
		},
	},
}

type unitAlias struct {
	BaseSymbol string
	BaseTarget string
	Prefixes   []string
}

type unitAliass map[string][]unitAlias

func (self unitAliass) AddTo(reg *UnitRegistry) error {
	for mType, defs := range self {
		for _, a := range defs {
			if err := reg.DeclareAlias(mType, a.BaseSymbol, a.BaseTarget, a.Prefixes); err != nil {
				return err
			}
		}
	}
	return nil
}

var systemAliases = unitAliass{
	"Concentration": {
		{
			BaseSymbol: "g/L",
			BaseTarget: "g/l",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "M",
			BaseTarget: "Mol/l",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "Molar",
			BaseTarget: "Mol/l",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "M/l",
			BaseTarget: "Mol/l",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "M/ml",
			BaseTarget: "Mol/ml",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "M/ul",
			BaseTarget: "Mol/ul",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "Mol/L",
			BaseTarget: "Mol/l",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "Mol/mL",
			BaseTarget: "Mol/ml",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "Mol/uL",
			BaseTarget: "Mol/ul",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "M/L",
			BaseTarget: "Mol/l",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "M/mL",
			BaseTarget: "Mol/ml",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "M/uL",
			BaseTarget: "Mol/ul",
			Prefixes:   allPrefixes,
		},
		{
			BaseSymbol: "x",
			BaseTarget: "X",
		},
		{
			BaseSymbol: "w/v",
			BaseTarget: "% w/v",
		},
	},
	"Volume": {
		{
			BaseSymbol: "L",
			BaseTarget: "l",
			Prefixes:   allPrefixes,
		},
	},
	"Length": {
		{
			BaseSymbol: "\"",
			BaseTarget: "in",
		},
	},
	"Temperature": {
		{
			BaseSymbol: "C",
			BaseTarget: "℃",
		},
		{
			BaseSymbol: "˚C",
			BaseTarget: "℃",
		},
		{
			BaseSymbol: "°C",
			BaseTarget: "℃",
		},
	},
	"Time": {
		{
			BaseSymbol: "minutes",
			BaseTarget: "min",
		},
	},
	"Angle": {
		{
			BaseSymbol: "radians",
			BaseTarget: "rad",
		},
		{
			BaseSymbol: "deg",
			BaseTarget: "°",
		},
		{
			BaseSymbol: "degrees",
			BaseTarget: "°",
		},
		{
			BaseSymbol: "˚",
			BaseTarget: "°",
		},
	},
}
