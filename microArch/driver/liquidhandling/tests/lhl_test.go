package tests

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func getTestBlow(idGen *id.IDGenerator, ch *wtype.LHChannelParameter, multi int, tipType wtype.TipType) liquidhandling.RobotInstruction {
	cti := &liquidhandling.ChannelTransferInstruction{
		Multi:      multi,
		What:       make([]string, 0, multi),
		PltTo:      make([]string, 0, multi),
		WellTo:     make([]string, 0, multi),
		Volume:     make([]wunit.Volume, 0, multi),
		TPlateType: make([]string, 0, multi),
		TVolume:    make([]wunit.Volume, 0, multi),
		Prms:       make([]*wtype.LHChannelParameter, 0, multi),
		TipType:    make([]wtype.TipType, 0, multi),
		Component:  make([]string, 0, multi),
	}

	for i := 0; i < multi; i++ {
		cti.What = append(cti.What, "soup")
		cti.PltTo = append(cti.PltTo, "position_4")
		cti.WellTo = append(cti.WellTo, "A1")
		cti.Volume = append(cti.Volume, wunit.NewVolume(10.0, "ul"))
		cti.TPlateType = append(cti.TPlateType, "DWST12")
		cti.TVolume = append(cti.TVolume, wunit.NewVolume(20.0, "ul"))
		cti.TipType = append(cti.TipType, tipType)
		cti.Prms = append(cti.Prms, ch)
		cti.Component = append(cti.Component, "minestrone")
	}

	return liquidhandling.NewBlowInstruction(idGen, cti)
}

func getTestSuck(idGen *id.IDGenerator, ch *wtype.LHChannelParameter, multi int, tipType wtype.TipType, plateType string, wellCoords []string) liquidhandling.RobotInstruction {
	cti := &liquidhandling.ChannelTransferInstruction{
		Multi:      multi,
		What:       make([]string, 0, multi),
		PltFrom:    make([]string, 0, multi),
		WellFrom:   wellCoords,
		Volume:     make([]wunit.Volume, 0, multi),
		FPlateType: make([]string, 0, multi),
		FVolume:    make([]wunit.Volume, 0, multi),
		Prms:       make([]*wtype.LHChannelParameter, 0, multi),
		TipType:    make([]wtype.TipType, 0, multi),
		Component:  make([]string, 0, multi),
	}

	for i := 0; i < multi; i++ {
		cti.What = append(cti.What, "soup")
		cti.PltFrom = append(cti.PltFrom, "position_4")
		cti.Volume = append(cti.Volume, wunit.NewVolume(10.0, "ul"))
		cti.FPlateType = append(cti.FPlateType, "plateType")
		cti.FVolume = append(cti.FVolume, wunit.NewVolume(100.0, "ul"))
		cti.TipType = append(cti.TipType, tipType)
		cti.Prms = append(cti.Prms, ch)
		cti.Component = append(cti.Component, "minestrone")
	}

	return liquidhandling.NewSuckInstruction(idGen, cti)
}

// what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int, Components []string, policies []wtype.LHPolicy
func getTestTransfer(vol wunit.Volume) liquidhandling.RobotInstruction {
	v2 := wunit.NewVolume(5000.0, "ul")
	v3 := wunit.NewVolume(0.0, "ul")
	return liquidhandling.NewTransferInstruction(
		[]string{"water"},
		[]string{"position_4"},
		[]string{"position_8"},
		[]string{"A1"},
		[]string{"G5"},
		[]string{"DWST12"},
		[]string{"DSW96"},
		[]wunit.Volume{vol},
		[]wunit.Volume{v2},
		[]wunit.Volume{v3},
		[]int{8},
		[]int{12},
		[]int{8},
		[]int{12},
		[]string{"water"},
		[]wtype.LHPolicy{nil},
	)
}

func TestBlowMixing(t *testing.T) {
	tenUl := wunit.NewVolume(10.0, "ul")
	tests := []*PolicyTest{
		{
			Name: "single channel",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"POST_MIX":        5,
						"POST_MIX_VOLUME": 10.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,MIX,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 5, //the Mix
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"CYCLES": []int{5},
						"VOLUME": []wunit.Volume{tenUl},
					},
				},
			},
		},
		{
			Name: "eight channel",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"POST_MIX":        5,
						"POST_MIX_VOLUME": 10.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 8, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,MIX,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 5, //the Mix
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"CYCLES": []int{5, 5, 5, 5, 5, 5, 5, 5},
						"VOLUME": []wunit.Volume{tenUl, tenUl, tenUl, tenUl, tenUl, tenUl, tenUl, tenUl},
					},
				},
			},
		},
		{
			Name: "set post mix rate",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"POST_MIX":        5,
						"POST_MIX_VOLUME": 10.0,
						"POST_MIX_RATE":   1.5,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,SPS,MOV,MIX,SPS,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 4,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   1.5,
					},
				},
				{
					Instruction: 6, //the Mix
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"CYCLES": []int{5},
						"VOLUME": []wunit.Volume{tenUl},
					},
				},
				{
					Instruction: 7,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "set post mix out of range",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"POST_MIX":        5,
						"POST_MIX_VOLUME": 10.0,
						"POST_MIX_RATE":   150.,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			Error: "setting post mix pipetting speed: value 150.000000 out of range 0.022500 - 3.750000",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestSuckMixing(t *testing.T) {
	tenUl := wunit.NewVolume(10.0, "ul")
	tests := []*PolicyTest{
		{
			Name: "single channel",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"PRE_MIX":        5,
						"PRE_MIX_VOLUME": 10.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,MIX,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 3, //the Mix
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"CYCLES": []int{5},
						"VOLUME": []wunit.Volume{tenUl},
					},
				},
			},
		},
		{
			Name: "eight channel",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"PRE_MIX":        5,
						"PRE_MIX_VOLUME": 10.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 8, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,MIX,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 3, //the Mix
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"CYCLES": []int{5, 5, 5, 5, 5, 5, 5, 5},
						"VOLUME": []wunit.Volume{tenUl, tenUl, tenUl, tenUl, tenUl, tenUl, tenUl, tenUl},
					},
				},
			},
		},
		{
			Name: "set pre mix rate",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"PRE_MIX":        5,
						"PRE_MIX_VOLUME": 10.0,
						"PRE_MIX_RATE":   1.5,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,SPS,MOV,MIX,SPS,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 2,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   1.5,
					},
				},
				{
					Instruction: 4, //the Mix
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"CYCLES": []int{5},
						"VOLUME": []wunit.Volume{tenUl},
					},
				},
				{
					Instruction: 5,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "set pre mix out of range",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"PRE_MIX":        5,
						"PRE_MIX_VOLUME": 10.0,
						"PRE_MIX_RATE":   150.,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			Error: "setting pre mix pipetting speed: value 150.000000 out of range 0.022500 - 3.750000",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestZOffset(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "blow z-offset",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSPZOFFSET": 5.4,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 2, //the move before the dispense
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"OFFSETZ": []float64{5.4},
					},
				},
			},
		},
		{
			Name: "blow z-offset unset",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 2, //the move before the dispense
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"OFFSETZ": []float64{defaultZOffset},
					},
				},
			},
		},
		{
			Name: "suck z-offset",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASPZOFFSET": 5.4,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 2, //the move before the dispense
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"OFFSETZ": []float64{5.4},
					},
				},
			},
		},
		{
			Name: "suck z-offset unset",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 2, //the move before the dispense
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"OFFSETZ": []float64{defaultZOffset},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestEntrySpeed(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "blow entry speed",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSPENTRYSPEED": 50.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,SDS,MOV,DSP,SDS,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 1,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"DRIVE": "Z",
						"SPEED": defaultZSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"DRIVE": "Z",
						"SPEED": 50.0,
					},
				},
			},
		},
		{
			Name: "suck entry speed",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASPENTRYSPEED": 50.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,SDS,MOV,ASP,MOV,SDS]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 1,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"DRIVE": "Z",
						"SPEED": defaultZSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"DRIVE": "Z",
						"SPEED": 50.0,
					},
				},
				{
					Instruction: 7,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"DRIVE": "Z",
						"SPEED": defaultZSpeed,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestDSPPipetSpeed(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "blow pipette speed unset",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "blow pipette speed unset, different default",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DEFAULTPIPETTESPEED": 1.82,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   1.82,
					},
				},
			},
		},
		{
			Name: "blow pipette speed OK",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSPSPEED": 1.5,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,SPS,DSP,SPS,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   1.5,
					},
				},
				{
					Instruction: 5,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "blow pipette speed too large, made safe",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSPSPEED":             LVMaxRate + 1.0,
						"OVERRIDEPIPETTESPEED": true,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,SPS,DSP,SPS,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   LVMaxRate,
					},
				},
				{
					Instruction: 5,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "blow pipette speed too small, made safe",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSPSPEED":             LVMinRate * 0.5,
						"OVERRIDEPIPETTESPEED": true,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,SPS,DSP,SPS,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   LVMinRate,
					},
				},
				{
					Instruction: 5,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "blow pipette speed too large, error",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSPSPEED":             4.75,
						"OVERRIDEPIPETTESPEED": false,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			Error: "setting pipette dispense speed: value 4.750000 out of range 0.022500 - 3.750000",
		},
		{
			Name: "blow pipette speed too small, error",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSPSPEED":             0.01,
						"OVERRIDEPIPETTESPEED": false,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			Error: "setting pipette dispense speed: value 0.010000 out of range 0.022500 - 3.750000",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestASPPipetSpeed(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "suck pipette speed unset",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "suck pipette speed unset, different default",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DEFAULTPIPETTESPEED": 1.56,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   1.56,
					},
				},
			},
		},
		{
			Name: "suck pipette speed OK",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASPSPEED": 1.5,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,SPS,ASP,SPS]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   1.5,
					},
				},
				{
					Instruction: 5,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "suck pipette speed too large, made safe",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASPSPEED":             LVMaxRate + 1.0,
						"OVERRIDEPIPETTESPEED": true,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,SPS,ASP,SPS]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   LVMaxRate,
					},
				},
				{
					Instruction: 5,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "suck pipette speed too small, made safe",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASPSPEED":             LVMinRate * 0.5,
						"OVERRIDEPIPETTESPEED": true,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,SPS,ASP,SPS]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 0,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
				{
					Instruction: 3,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   LVMinRate,
					},
				},
				{
					Instruction: 5,
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"HEAD":    1,
						"CHANNEL": -1,
						"SPEED":   defaultPipetteSpeed,
					},
				},
			},
		},
		{
			Name: "suck pipette speed too large, error",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASPSPEED":             4.75,
						"OVERRIDEPIPETTESPEED": false,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			Error: "setting pipette aspirate speed: value 4.750000 out of range 0.022500 - 3.750000",
		},
		{
			Name: "suck pipette speed too small, error",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASPSPEED":             0.01,
						"OVERRIDEPIPETTESPEED": false,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			Error: "setting pipette aspirate speed: value 0.010000 out of range 0.022500 - 3.750000",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestTipReuse(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "no tip reuse allowed ",
			Rules: []*Rule{
				{
					Name: "water",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "water",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"TIP_REUSE_LIMIT": 0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestTransfer(wunit.NewVolume(300.0, "ul"))
				return nil
			},
			ExpectedInstructions: "[MOV,LOD,SPS,SDS,MOV,ASP,SPS,SDS,MOV,DSP,MOV,BLO,MOV,ULD,MOV,LOD,SPS,SDS,MOV,ASP,SPS,SDS,MOV,DSP,MOV,BLO,MOV,ULD]",
			Assertions: []*InstructionAssertion{
				{},
			},
		},
		{
			Name: "tip reuse allowed ",
			Rules: []*Rule{
				{
					Name: "water",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "water",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"TIP_REUSE_LIMIT": 100,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestTransfer(wunit.NewVolume(300.0, "ul"))
				return nil
			},
			ExpectedInstructions: "[MOV,LOD,SPS,SDS,MOV,ASP,SPS,SDS,MOV,DSP,MOV,BLO,SPS,SDS,MOV,ASP,SPS,SDS,MOV,DSP,MOV,BLO,MOV,ULD]",
			Assertions: []*InstructionAssertion{
				{},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestAspWait(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "asp wait 3s, multi 1",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASP_WAIT": 3.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,ASP,WAI]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 4, //Wait
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"TIME": 3.0,
					},
				},
			},
		},
		{
			Name: "asp wait 3s, multi 8",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"ASP_WAIT": 3.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 8, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,ASP,WAI]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 4, //Wait
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"TIME": 3.0,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestDspWait(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "dsp wait 3s, multi 1",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSP_WAIT": 3.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,WAI,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 4, //Wait
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"TIME": 3.0,
					},
				},
			},
		},
		{
			Name: "dsp wait 3s, multi 8",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"DSP_WAIT": 3.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 8, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,WAI,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 4, //Wait
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"TIME": 3.0,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestTouchoff(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "test touchoff ",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"TOUCHOFF":    true,
						"TOUCHOFFSET": 2.37,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 2, //move prior to dispense
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"REFERENCE": []int{0},
						"OFFSETZ":   []float64{0.5},
					},
				},
				{
					Instruction: 4, //touchoff move
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"REFERENCE": []int{0},        // well bottom
						"OFFSETZ":   []float64{2.37}, // as set
					},
				},
			},
		},
		{
			Name: "test large touchoff",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"TOUCHOFF":    true,
						"TOUCHOFFSET": liquidhandling.MaxTouchOffset + 5.0,
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 2, //move prior to dispense
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"REFERENCE": []int{0},
						"OFFSETZ":   []float64{0.5},
					},
				},
				{
					Instruction: 4, //touchoff move
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"REFERENCE": []int{0},                                 // well bottom
						"OFFSETZ":   []float64{liquidhandling.MaxTouchOffset}, // as set
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestExtraVolumes(t *testing.T) {
	tests := []*PolicyTest{
		{
			Name: "extra asp volume ",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"EXTRA_ASP_VOLUME": wunit.NewVolume(2.0, "ul"),
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestSuck(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20", "pcrplate_skirted_riser18", []string{"A1"})
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,ASP]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 3, // ASP
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"VOLUME": []wunit.Volume{wunit.NewVolume(12.0, "ul")},
					},
				},
			},
		},
		{
			Name: "extra dsp volume ",
			Rules: []*Rule{
				{
					Name: "soup",
					Conditions: []Condition{
						&CategoryCondition{
							Attribute: "LIQUIDCLASS",
							Value:     "soup",
						},
					},
					Policy: map[liquidhandling.InstructionParameter]interface{}{
						"EXTRA_DISP_VOLUME": wunit.NewVolume(2.0, "ul"),
					},
				},
			},
			Setup: func(self *PolicyTest, lab *laboratory.Laboratory) error {
				self.Instruction = getTestBlow(lab.IDGenerator, getLVConfig(lab.IDGenerator), 1, "Gilson20")
				return nil
			},
			ExpectedInstructions: "[SPS,SDS,MOV,DSP,MOV,BLO]",
			Assertions: []*InstructionAssertion{
				{
					Instruction: 3, // dispense
					Values: map[liquidhandling.InstructionParameter]interface{}{
						"VOLUME": []wunit.Volume{wunit.NewVolume(12.0, "ul")},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}