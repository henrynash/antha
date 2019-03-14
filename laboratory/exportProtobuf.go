// +build protobuf

package laboratory

import (
	runner "github.com/Synthace/antha-runner/export"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

func export(idGen *id.IDGenerator, outDir string, instrs []effects.Inst) error {
	return runner.Export(idGen, outDir, instrs)
}