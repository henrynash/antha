// anthalib//liquidhandling/liquidhandling_test.go: Part of the Antha language
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

package liquidhandling

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func TestStockConcs(*testing.T) {
	rand := wutil.GetRandom()
	names := []string{"tea", "milk", "sugar"}

	minrequired := make(map[string]float64, len(names))
	maxrequired := make(map[string]float64, len(names))
	Smax := make(map[string]float64, len(names))
	T := make(map[string]wunit.Volume, len(names))
	vmin := 10.0

	for _, name := range names {
		r := rand.Float64() + 1.0
		r2 := rand.Float64() + 1.0
		r3 := rand.Float64() + 1.0

		minrequired[name] = r * r2 * 20.0
		maxrequired[name] = r * r2 * 30.0
		Smax[name] = r * r2 * r3 * 70.0
		T[name] = wunit.NewVolume(100.0, "ul")
	}

	choose_stock_concentrations(minrequired, maxrequired, Smax, vmin, T)
	/*for k, v := range cncs {
		logger.Debug(fmt.Sprintln(k, " ", minrequired[k], " ", maxrequired[k], " ", T[k], " ", v))
	}*/
}

func configure_request_simple(ctx context.Context, rq *LHRequest) {
	water := GetComponentForTest(ctx, "water", wunit.NewVolume(100.0, "ul"))
	mmx := GetComponentForTest(ctx, "mastermix_sapI", wunit.NewVolume(100.0, "ul"))
	part := GetComponentForTest(ctx, "dna", wunit.NewVolume(50.0, "ul"))

	for k := 0; k < 9; k++ {
		ins := wtype.NewLHMixInstruction()
		ws := mixer.Sample(water, wunit.NewVolume(8.0, "ul"))
		mmxs := mixer.Sample(mmx, wunit.NewVolume(8.0, "ul"))
		ps := mixer.Sample(part, wunit.NewVolume(1.0, "ul"))

		ins.AddComponent(ws)
		ins.AddComponent(mmxs)
		ins.AddComponent(ps)
		ins.AddProduct(GetComponentForTest(ctx, "water", wunit.NewVolume(17.0, "ul")))
		rq.Add_instruction(ins)
	}

}

func configure_request_bigger(ctx context.Context, rq *LHRequest) {
	water := GetComponentForTest(ctx, "water", wunit.NewVolume(2000.0, "ul"))
	mmx := GetComponentForTest(ctx, "mastermix_sapI", wunit.NewVolume(2000.0, "ul"))
	part := GetComponentForTest(ctx, "dna", wunit.NewVolume(1000.0, "ul"))

	for k := 0; k < 99; k++ {
		ins := wtype.NewLHMixInstruction()
		ws := mixer.Sample(water, wunit.NewVolume(8.0, "ul"))
		mmxs := mixer.Sample(mmx, wunit.NewVolume(8.0, "ul"))
		ps := mixer.Sample(part, wunit.NewVolume(1.0, "ul"))

		ins.AddComponent(ws)
		ins.AddComponent(mmxs)
		ins.AddComponent(ps)
		ins.AddProduct(GetComponentForTest(ctx, "water", wunit.NewVolume(17.0, "ul")))
		rq.Add_instruction(ins)
	}

}

func configureMultiChannelTestRequest(ctx context.Context, rq *LHRequest) {
	water := GetComponentForTest(ctx, "multiwater", wunit.NewVolume(2000.0, "ul"))
	/*	mmx := GetComponentForTest(ctx, "mastermix_sapI", wunit.NewVolume(2000.0, "ul"))
		part := GetComponentForTest(ctx, "dna", wunit.NewVolume(1000.0, "ul"))
	*/
	for k := 0; k < 9; k++ {
		ins := wtype.NewLHMixInstruction()
		ws := mixer.Sample(water, wunit.NewVolume(50.0, "ul"))
		/*		mmxs := mixer.Sample(mmx, wunit.NewVolume(40.0, "ul"))
				ps := mixer.Sample(part, wunit.NewVolume(100.0, "ul"))
		*/
		ins.AddComponent(ws)
		/*		ins.AddComponent(mmxs)
				ins.AddComponent(ps)
		*/
		ins.AddProduct(GetComponentForTest(ctx, "water", wunit.NewVolume(50, "ul")))
		rq.Add_instruction(ins)
	}

}

func configureTransferRequestForZTest(liquid string, transferVol wunit.Volume, numberOfTransfers int) (rq *LHRequest, err error) {

	// set up ctx
	ctx := testinventory.NewContext(context.Background())

	// make liquid handler
	lh := GetLiquidHandlerForTest(ctx)

	// make some tipboxes
	var tipBoxes []*wtype.LHTipbox
	tpHigh, err := inventory.NewTipbox(ctx, "Gilson200")
	if err != nil {
		return rq, err
	}
	tpLow, err := inventory.NewTipbox(ctx, "Gilson20")
	if err != nil {
		return rq, err
	}
	tipBoxes = append(tipBoxes, tpHigh, tpLow)

	//initialise request
	rq = GetLHRequestForTest()

	liq := GetComponentForTest(ctx, liquid, wunit.NewVolume(2000.0, "ul"))

	for k := 0; k < numberOfTransfers; k++ {
		ins := wtype.NewLHMixInstruction()
		ws := mixer.Sample(liq, transferVol)

		ins.AddComponent(ws)

		ins.AddProduct(GetComponentForTest(ctx, liquid, transferVol))
		rq.Add_instruction(ins)
	}

	// add plates and tip boxes
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.Tips = tipBoxes

	rq.ConfigureYourself()

	if err := lh.Plan(ctx, rq); err != nil {
		return rq, fmt.Errorf("Got an error planning with no inputs: %s", err.Error())
	}
	return rq, nil
}

func configureSingleChannelTestRequest(ctx context.Context, rq *LHRequest) {
	water := GetComponentForTest(ctx, "multiwater", wunit.NewVolume(2000.0, "ul"))
	/*	mmx := GetComponentForTest(ctx, "mastermix_sapI", wunit.NewVolume(2000.0, "ul"))
		part := GetComponentForTest(ctx, "dna", wunit.NewVolume(1000.0, "ul"))
	*/
	for k := 0; k < 1; k++ {
		ins := wtype.NewLHMixInstruction()
		ws := mixer.Sample(water, wunit.NewVolume(50.0, "ul"))
		/*		mmxs := mixer.Sample(mmx, wunit.NewVolume(40.0, "ul"))
				ps := mixer.Sample(part, wunit.NewVolume(100.0, "ul"))
		*/
		ins.AddComponent(ws)
		/*		ins.AddComponent(mmxs)
				ins.AddComponent(ps)
		*/
		ins.AddProduct(GetComponentForTest(ctx, "water", wunit.NewVolume(50, "ul")))
		rq.Add_instruction(ins)
	}

}

type zOffsetTest struct {
	liquidType              string
	inPutPlateType          string
	numberOfTransfers       int
	volume                  wunit.Volume
	expectedAspirateZOffset string
	expectedDispenseZOffset string
}

var offsetTests []zOffsetTest = []zOffsetTest{
	zOffsetTest{
		liquidType:              "multiwater",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: "1.2500",
		expectedDispenseZOffset: "1.7500",
	},
	zOffsetTest{
		liquidType:              "multiwater",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: "1.2500,1.2500",
		expectedDispenseZOffset: "1.7500,1.7500",
	},
	zOffsetTest{
		liquidType:              "multiwater",
		numberOfTransfers:       8,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: "1.2500,1.2500,1.2500,1.2500,1.2500,1.2500,1.2500,1.2500",
		expectedDispenseZOffset: "1.7500,1.7500,1.7500,1.7500,1.7500,1.7500,1.7500,1.7500",
	},
	zOffsetTest{
		liquidType:              "multiwater",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(5, "ul"),
		expectedAspirateZOffset: "0.5000",
		expectedDispenseZOffset: "1.0000",
	},
	zOffsetTest{
		liquidType:              "multiwater",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(5, "ul"),
		expectedAspirateZOffset: "0.5000,0.5000",
		expectedDispenseZOffset: "1.0000,1.0000",
	},
	zOffsetTest{
		liquidType:              "water",
		numberOfTransfers:       1,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: "1.2500",
		expectedDispenseZOffset: "1.7500",
	},
	zOffsetTest{
		liquidType:              "water",
		numberOfTransfers:       2,
		volume:                  wunit.NewVolume(50, "ul"),
		expectedAspirateZOffset: "1.2500",
		expectedDispenseZOffset: "1.7500",
	},
}

func TestMultiZOffset2(t *testing.T) {

	for _, test := range offsetTests {
		request, err := configureTransferRequestForZTest(test.liquidType, test.volume, test.numberOfTransfers)
		if err != nil {
			t.Error(err.Error())
		}

		var aspirateInstructions, dispenseInstructions []liquidhandling.Summary

		for i, instruction := range request.Instructions {
			if i > 0 {
				if liquidhandling.InstructionTypeName(instruction) == "ASP" {
					aspirateSummary, err := liquidhandling.SummariseTwoSteps(request.Instructions[i-1], instruction)
					if err != nil {
						fmt.Println(err.Error())
					}
					aspirateInstructions = append(aspirateInstructions, aspirateSummary)
				} else if liquidhandling.InstructionTypeName(instruction) == "DSP" {
					dispenseSummary, err := liquidhandling.SummariseTwoSteps(request.Instructions[i-1], instruction)
					if err != nil {
						fmt.Println(err.Error())
					}
					dispenseInstructions = append(dispenseInstructions, dispenseSummary)
				}
			}
		}
		for i, aspirationStep := range aspirateInstructions {
			if !reflect.DeepEqual(aspirationStep.OffsetZ, test.expectedAspirateZOffset) {
				t.Error("for test: ", text.PrettyPrint(aspirationStep), "\n",
					"aspiration step: ", i, "\n",
					"expected Z offset for aspirate:", test.expectedAspirateZOffset, "\n",
					"got: ", aspirationStep.OffsetZ, "\n",
				)
			}
		}

		for i, dispenseStep := range dispenseInstructions {
			if !reflect.DeepEqual(dispenseStep.OffsetZ, test.expectedDispenseZOffset) {
				t.Error(" for test: ", text.PrettyPrint(dispenseStep), "\n",
					"dispense step: ", i, "\n",
					"expected Z offset for dispense: ", test.expectedDispenseZOffset, "\n",
					"got: ", dispenseStep.OffsetZ, "\n",
				)
			}
		}

	}
}

func TestMultiZOffset(t *testing.T) {

	// set up ctx
	ctx := testinventory.NewContext(context.Background())

	// make liquid handler
	lh := GetLiquidHandlerForTest(ctx)

	// make some tipboxes
	var tipBoxes []*wtype.LHTipbox
	tpHigh, err := inventory.NewTipbox(ctx, "Gilson200")
	if err != nil {
		t.Fatal(err)
	}
	tpLow, err := inventory.NewTipbox(ctx, "Gilson20")
	if err != nil {
		t.Fatal(err)
	}
	tipBoxes = append(tipBoxes, tpHigh, tpLow)

	// set up multi

	//initialise multi request
	multiRq := GetLHRequestForTest()

	// set to Multi channel test request
	configureMultiChannelTestRequest(ctx, multiRq)
	// add plates and tip boxes
	multiRq.Input_platetypes = append(multiRq.Input_platetypes, GetPlateForTest())
	multiRq.Output_platetypes = append(multiRq.Output_platetypes, GetPlateForTest())

	multiRq.Tips = tipBoxes

	multiRq.ConfigureYourself()

	if err := lh.Plan(ctx, multiRq); err != nil {
		t.Fatalf("Got an error planning with no inputs: %s", err)
	}

	// set up single channel

	//initialise single request
	singleRq := GetLHRequestForTest()

	// set to single channel test request
	configureSingleChannelTestRequest(ctx, singleRq)
	// add plates and tip boxes
	singleRq.Input_platetypes = append(singleRq.Input_platetypes, GetPlateForTest())
	singleRq.Output_platetypes = append(singleRq.Output_platetypes, GetPlateForTest())

	singleRq.Tips = tipBoxes

	singleRq.ConfigureYourself()

	if err := lh.Plan(ctx, singleRq); err != nil {
		t.Fatalf("Got an error planning with no inputs: %s", err)
	}

	var singleAspirateInstructions, singleDispenseInstructions, multiAspirateInstructions, multiDispenseInstructions []liquidhandling.Summary

	for i, instruction := range singleRq.Instructions {
		if i > 0 {
			if liquidhandling.InstructionTypeName(instruction) == "ASP" {
				aspirateSummary, err := liquidhandling.SummariseTwoSteps(singleRq.Instructions[i-1], instruction)
				if err != nil {
					fmt.Println(err.Error())
				}
				singleAspirateInstructions = append(singleAspirateInstructions, aspirateSummary)
			} else if liquidhandling.InstructionTypeName(instruction) == "DSP" {
				dispenseSummary, err := liquidhandling.SummariseTwoSteps(singleRq.Instructions[i-1], instruction)
				if err != nil {
					fmt.Println(err.Error())
				}
				singleDispenseInstructions = append(singleDispenseInstructions, dispenseSummary)
			}
		}
	}

	for i, instruction := range multiRq.Instructions {
		if i > 0 {
			if liquidhandling.InstructionTypeName(instruction) == "ASP" {
				aspirateSummary, err := liquidhandling.SummariseTwoSteps(multiRq.Instructions[i-1], instruction)
				if err != nil {
					fmt.Println(err.Error())
				}
				multiAspirateInstructions = append(multiAspirateInstructions, aspirateSummary)
			} else if liquidhandling.InstructionTypeName(instruction) == "DSP" {
				dispenseSummary, err := liquidhandling.SummariseTwoSteps(multiRq.Instructions[i-1], instruction)
				if err != nil {
					fmt.Println(err.Error())
				}
				multiDispenseInstructions = append(multiDispenseInstructions, dispenseSummary)
			}
		}
	}
	for i, aspirationStep := range singleAspirateInstructions {
		if !reflect.DeepEqual(aspirationStep.OffsetZ, multiAspirateInstructions[i].OffsetZ) {
			t.Error(fmt.Sprintf("single Aspirate Z offset: %+v ", text.PrettyPrint(aspirationStep)), "\n",
				fmt.Sprintf("multi Aspirate Z offset: %+v ", text.PrettyPrint(multiAspirateInstructions[i])), "\n")
		}
	}

	for i, dispenseStep := range singleDispenseInstructions {
		if !reflect.DeepEqual(dispenseStep.OffsetZ, multiDispenseInstructions[i].OffsetZ) {
			t.Error("single Dispense Z offset: ", text.PrettyPrint(dispenseStep), "\n",
				fmt.Sprintf("multi Dispense Z offset: %+v ", text.PrettyPrint(multiDispenseInstructions[i])), "\n")
		}
	}

}

func TestTipOverridePositive(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	lh := GetLiquidHandlerForTest(ctx)
	rq := GetLHRequestForTest()
	configure_request_simple(ctx, rq)
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	var tpz []*wtype.LHTipbox
	tp, err := inventory.NewTipbox(ctx, "Gilson20")
	if err != nil {
		t.Fatal(err)
	}
	tpz = append(tpz, tp)

	rq.Tips = tpz

	rq.ConfigureYourself()

	if err := lh.Plan(ctx, rq); err != nil {
		t.Fatalf("Got an error planning with no inputs: %s", err)
	}

}
func TestTipOverrideNegative(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	lh := GetLiquidHandlerForTest(ctx)
	rq := GetLHRequestForTest()
	configure_request_simple(ctx, rq)
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())
	var tpz []*wtype.LHTipbox
	tp, err := inventory.NewTipbox(ctx, "Gilson200")
	if err != nil {
		t.Fatal(err)
	}
	tpz = append(tpz, tp)

	rq.Tips = tpz

	rq.ConfigureYourself()

	err = lh.Plan(ctx, rq)

	if e, f := "No tip chosen: Volume 8 ul is too low to be accurately moved by the liquid handler (configured minimum 10 ul, tip minimum 10 ul). Low volume tips may not be available and / or the robot may need to be configured differently", err.Error(); e != f {
		t.Fatalf("expecting error %q found %q", e, f)
	}
}

func TestPlateReuse(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	lh := GetLiquidHandlerForTest(ctx)
	rq := GetLHRequestForTest()
	configure_request_simple(ctx, rq)
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()

	err := lh.Plan(ctx, rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got an error planning with no inputs: ", err))
	}

	// reset the request
	rq = GetLHRequestForTest()
	configure_request_simple(ctx, rq)

	for _, plateid := range lh.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		thing := lh.Properties.PlateLookup[plateid]

		plate, ok := thing.(*wtype.LHPlate)
		if !ok {
			continue
		}

		if strings.Contains(plate.Name(), "Output_plate") {
			// leave it out
			continue
		}

		rq.Input_plates[plateid] = plate
	}
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()

	lh = GetLiquidHandlerForTest(ctx)
	err = lh.Plan(ctx, rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got error resimulating: ", err))
	}

	// if we added nothing, input assignments should be empty

	if rq.NewComponentsAdded() {
		t.Fatal(fmt.Sprint("Resimulation failed: needed to add ", len(rq.Input_vols_wanting), " components"))
	}

	// now try a deliberate fail

	// reset the request again
	rq = GetLHRequestForTest()
	configure_request_simple(ctx, rq)

	for _, plateid := range lh.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		thing := lh.Properties.PlateLookup[plateid]

		plate, ok := thing.(*wtype.LHPlate)
		if !ok {
			continue
		}
		if strings.Contains(plate.Name(), "Output_plate") {
			// leave it out
			continue
		}
		for _, v := range plate.Wellcoords {
			if !v.Empty() {
				v.Remove(wunit.NewVolume(5.0, "ul"))
			}
		}

		rq.Input_plates[plateid] = plate
	}
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()

	lh = GetLiquidHandlerForTest(ctx)
	err = lh.Plan(ctx, rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got error resimulating: ", err))
	}

	// this time we should have added some components again
	if len(rq.Input_assignments) != 3 {
		t.Fatal(fmt.Sprintf("Error resimulating, should have added 3 components, instead added %d", len(rq.Input_assignments)))
	}
}

func TestBeforeVsAfter(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	lh := GetLiquidHandlerForTest(ctx)
	rq := GetLHRequestForTest()
	configure_request_simple(ctx, rq)
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()

	err := lh.Plan(ctx, rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got an error planning with no inputs: ", err))
	}

	for pos, _ := range lh.Properties.PosLookup {

		id1, ok1 := lh.Properties.PosLookup[pos]
		id2, ok2 := lh.FinalProperties.PosLookup[pos]

		if ok1 && !ok2 || ok2 && !ok1 {
			t.Fatal(fmt.Sprintf("Position %s inconsistent: Before %t after %t", pos, ok1, ok2))
		}

		p1 := lh.Properties.PlateLookup[id1]
		p2 := lh.FinalProperties.PlateLookup[id2]

		// check types

		t1 := reflect.TypeOf(p1)
		t2 := reflect.TypeOf(p2)

		if t1 != t2 {
			t.Fatal(fmt.Sprintf("Types of thing at position %s not same: %v %v", pos, t1, t2))
		}

		// ok nice we have some sort of consistency

		switch p1.(type) {
		case *wtype.LHPlate:
			pp1 := p1.(*wtype.LHPlate)
			pp2 := p2.(*wtype.LHPlate)
			if pp1.Type != pp2.Type {
				t.Fatal(fmt.Sprintf("Plates at %s not same type: %s %s", pos, pp1.Type, pp2.Type))
			}
			it := wtype.NewOneTimeColumnWiseIterator(pp1)

			for {
				if !it.Valid() {
					break
				}
				wc := it.Curr()
				w1 := pp1.Wellcoords[wc.FormatA1()]
				w2 := pp2.Wellcoords[wc.FormatA1()]

				if w1.Empty() && w2.Empty() {
					it.Next()
					continue
				}
				/*
					fmt.Println(pp1.PlateName, " ", pp1.Type)
					fmt.Println(pp2.PlateName, " ", pp2.Type)
					fmt.Println(wc.FormatA1())
					fmt.Println(w1.ID, " ", w1.WContents.ID, " ", w1.WContents.CName, " ", w1.WContents.Vol)
					fmt.Println(w2.ID, " ", w2.WContents.ID, " ", w2.WContents.CName, " ", w2.WContents.Vol)
				*/

				if w1.WContents.ID == w2.WContents.ID {
					t.Fatal(fmt.Sprintf("IDs before and after must differ"))
				}
				it.Next()
			}
		case *wtype.LHTipbox:
			tb1 := p1.(*wtype.LHTipbox)
			tb2 := p2.(*wtype.LHTipbox)

			if tb1.Type != tb2.Type {
				t.Fatal(fmt.Sprintf("Tipbox at changed type: %s %s", tb1.Type, tb2.Type))
			}
		case *wtype.LHTipwaste:
			tw1 := p1.(*wtype.LHTipwaste)
			tw2 := p2.(*wtype.LHTipwaste)

			if tw1.Type != tw2.Type {
				t.Fatal(fmt.Sprintf("Tipwaste changed type: %s %s", tw1.Type, tw2.Type))
			}
		}

	}

}

func TestEP3(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	lh := GetLiquidHandlerForTest(ctx)
	lh.ExecutionPlanner = ExecutionPlanner3
	rq := GetLHRequestForTest()
	configure_request_simple(ctx, rq)
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()
	err := lh.Plan(ctx, rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got planning error: ", err))
	}

}
