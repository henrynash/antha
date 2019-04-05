package laboratory

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/laboratory/compare"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/utils"
	"github.com/antha-lang/antha/workflow"
)

// Compare compares output generated with any supplied test data in the workflow
func (labBuild *LaboratoryBuilder) Compare() error {

	if labBuild.workflow.Testing == nil {
		labBuild.Logger.Log("msg", "No comparison test data supplied.")
		return nil
	}

	mixIdx := 0
	errs := make(utils.ErrorSlice, 0, len(labBuild.instrs))

	for i, instr := range labBuild.instrs {
		if t, ok := instr.(*target.Mix); ok {
			labBuild.Logger.Log("msg", fmt.Sprintf("[%d] Checking mix instruction.", i))

			if expected, err := expectedMix(labBuild.workflow, mixIdx); err != nil {
				errs = append(errs, err)
			} else {
				if err := labBuild.compareTimings(t, expected); err != nil {
					errs = append(errs, err)
				}
				errs = append(errs, labBuild.compareOutputs(t, expected)...)
			}
			mixIdx++
		}
	}

	if len(errs) > 0 {
		if ef, err := labBuild.errorsToFile(errs); err != nil {
			return errors.New("errors in test comparisons, and failed to create error file")
		} else {
			return fmt.Errorf("errors in comparison tests, details in %s", ef)
		}
	}

	labBuild.Logger.Log("msg", "Comparison test data passed.")
	return nil
}

func (labBuild *LaboratoryBuilder) errorsToFile(errs utils.ErrorSlice) (string, error) {
	filename := filepath.Join(labBuild.outDir, "comparisons.txt")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400)
	if err != nil {
		return "", err
	}

	defer f.Close()
	errorsToWriter(f, errs)
	return filename, nil
}

func errorsToWriter(w io.Writer, errs utils.ErrorSlice) {
	for i, e := range errs {
		fmt.Fprintf(w, "[%d] %v\n", i, e)
	}
}

func expectedMix(w *workflow.Workflow, idx int) (*workflow.MixTaskCheck, error) {
	if idx >= len(w.Testing.MixTaskChecks) {
		return nil, fmt.Errorf("mix comparison %d not found, only %d mixes are expected", idx, len(w.Testing.MixTaskChecks))
	}

	return &w.Testing.MixTaskChecks[idx], nil
}

func (labBuild *LaboratoryBuilder) compareTimings(m *target.Mix, expectedMix *workflow.MixTaskCheck) error {
	const timeAccuracyPercent = 10

	if err := compareToPercent(expectedMix.TimeEstimate.Seconds(), m.GetTimeEstimate(), timeAccuracyPercent); err != nil {
		return fmt.Errorf("timing check failed, %v", err)
	}

	labBuild.Logger.Log("msg", fmt.Sprintf("Passed timing check. Expected %.3gs, found %.3gs.", expectedMix.TimeEstimate.Seconds(), m.GetTimeEstimate()))
	return nil
}

func compareToPercent(expected float64, actual float64, percent float64) error {
	const onePercent = 0.01
	if math.Abs(expected-actual) > math.Abs(onePercent*percent*expected) {
		return fmt.Errorf("expected %.2g but found %.2g (checked to %.2g%%)", expected, actual, percent)
	}

	return nil
}

func (labBuild *LaboratoryBuilder) compareOutputs(m *target.Mix, expectedMix *workflow.MixTaskCheck) utils.ErrorSlice {
	if expectedMix.Outputs == nil || len(expectedMix.Outputs) == 0 {
		labBuild.Logger.Log("msg", "No output comparison data supplied for mix task.")
		return nil
	}

	return compare.Plates(labBuild.effects.IDGenerator, expectedMix.Outputs, m.FinalProperties.Plates)
}
