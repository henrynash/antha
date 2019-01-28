package laboratory

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/antha-lang/antha/codegen"
	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/utils"
)

type Element interface {
	Name() string
	TypeName() string

	Setup(*Laboratory)
	Steps(*Laboratory)
	Analysis(*Laboratory)
	Validation(*Laboratory)
}

type LaboratoryBuilder struct {
	outDir   string
	workflow *composer.Workflow

	elemLock   sync.Mutex
	elemsUnrun int64
	elements   map[Element]*ElementBase

	errLock sync.Mutex
	errors  utils.ErrorSlice
	Errored chan struct{}

	Completed chan struct{}

	*lineMapManager
	Logger      *Logger
	FileManager *FileManager

	*effects.LaboratoryEffects
}

func NewLaboratoryBuilder(jobId string, fh io.Reader) *LaboratoryBuilder {
	labBuild := &LaboratoryBuilder{
		elements:  make(map[Element]*ElementBase),
		Errored:   make(chan struct{}),
		Completed: make(chan struct{}),

		lineMapManager: NewLineMapManager(),
		Logger:         NewLogger(),
	}

	if wf, err := composer.WorkflowFromReaders(fh); err != nil {
		labBuild.Fatal(err)
	} else {
		labBuild.workflow = wf
	}

	labBuild.LaboratoryEffects = effects.NewLaboratoryEffects(jobId, labBuild.workflow)

	flag.StringVar(&labBuild.outDir, "outdir", "", "Path to directory in which to write output files")
	flag.Parse()

	if labBuild.outDir == "" {
		if d, err := ioutil.TempDir("", fmt.Sprintf("antha-%s", jobId)); err != nil {
			labBuild.Fatal(err)
		} else {
			labBuild.outDir = d
			labBuild.Logger.Log("outdir", d)
		}
	}
	for _, leaf := range []string{"elements", "data"} {
		if err := os.MkdirAll(filepath.Join(labBuild.outDir, leaf), 0700); err != nil {
			labBuild.Fatal(err)
		}
	}
	labBuild.FileManager = NewFileManager(filepath.Join(labBuild.outDir, "data"))

	return labBuild
}

// Only use this before you call run.
func (labBuild *LaboratoryBuilder) InstallElement(e Element) {
	eb := NewElementBase(e)
	labBuild.elemLock.Lock()
	defer labBuild.elemLock.Unlock()
	labBuild.elements[e] = eb
	labBuild.elemsUnrun++
}

func (labBuild *LaboratoryBuilder) AddConnection(src, dst Element, fun func()) error {
	if ebSrc, found := labBuild.elements[src]; !found {
		return fmt.Errorf("Unknown src element: %v", src)
	} else if ebDst, found := labBuild.elements[dst]; !found {
		return fmt.Errorf("Unknown dst element: %v", dst)
	} else {
		ebDst.AddBlockedInput()
		ebSrc.AddOnExit(func() {
			fun()
			ebDst.InputReady()
		})
		return nil
	}
}

// Run all the installed elements.
func (labBuild *LaboratoryBuilder) RunElements() error {
	labBuild.elemLock.Lock()
	if labBuild.elemsUnrun == 0 {
		labBuild.elemLock.Unlock()
		close(labBuild.Completed)
		return nil

	} else {
		for _, eb := range labBuild.elements {
			eb.AddOnExit(labBuild.elementCompleted)
		}
		for e, eb := range labBuild.elements {
			go eb.Run(labBuild.makeLab(e))
		}
		labBuild.elemLock.Unlock()
		<-labBuild.Completed

		select {
		case <-labBuild.Errored:
			return labBuild.errors
		default:
			return nil
		}
	}
}

func (labBuild *LaboratoryBuilder) Compile(target *target.Target) ([]effects.Node, []effects.Inst, error) {
	if nodes, err := labBuild.Maker.MakeNodes(labBuild.Trace.Instructions()); err != nil {
		return nil, nil, err
	} else if instrs, err := codegen.Compile(labBuild.LaboratoryEffects, target, nodes); err != nil {
		return nil, nil, err
	} else {
		return nodes, instrs, nil
	}
}

func (labBuild *LaboratoryBuilder) elementCompleted() {
	labBuild.elemLock.Lock()
	defer labBuild.elemLock.Unlock()
	labBuild.elemsUnrun--
	if labBuild.elemsUnrun == 0 {
		close(labBuild.Completed)
	}
}

func (labBuild *LaboratoryBuilder) recordError(err error) {
	labBuild.errLock.Lock()
	defer labBuild.errLock.Unlock()
	labBuild.errors = append(labBuild.errors, err)
	select { // we keep the lock here to avoid a race to close
	case <-labBuild.Errored:
	default:
		close(labBuild.Errored)
	}
}

func (labBuild *LaboratoryBuilder) Fatal(err error) {
	labBuild.Logger.Log("fatal", err.Error())
	os.Exit(1)
}

func (labBuild *LaboratoryBuilder) Errors() error {
	select {
	case <-labBuild.Errored:
		labBuild.errLock.Lock()
		defer labBuild.errLock.Unlock()
		return labBuild.errors
	default:
		return nil
	}
}

type Laboratory struct {
	*LaboratoryBuilder
	element Element
	Logger  *Logger
}

func (labBuild *LaboratoryBuilder) makeLab(e Element) *Laboratory {
	return &Laboratory{
		LaboratoryBuilder: labBuild,
		element:           e,
		Logger:            labBuild.Logger, //.With("name", e.Name(), "type", e.TypeName()),
	}
}

// Only for use when you're in an element and want to call another element.
func (lab *Laboratory) CallSteps(e Element) error {
	eb := NewElementBase(e)

	finished := make(chan struct{})
	eb.AddOnExit(func() { close(finished) })

	go eb.Run(lab.makeLab(e), eb.element.Steps)
	<-finished

	select {
	case <-lab.Errored:
		return lab.errors
	default:
		return nil
	}
}

func (lab *Laboratory) Error(err error) {
	lab.recordError(err)
	lab.Logger.Log("error", err.Error())
}

func (lab *Laboratory) Errorf(fmtStr string, args ...interface{}) {
	lab.Error(fmt.Errorf(fmtStr, args...))
}

// ElementBase
type ElementBase struct {
	// count of inputs that are not yet ready (plus 1)
	pendingCount int64
	// this gets closed when all inputs become ready
	InputsReady chan struct{}
	// funcs to run when this element is completed
	onExit []func()
	// the actual element
	element Element
}

func NewElementBase(e Element) *ElementBase {
	return &ElementBase{
		pendingCount: 1,
		InputsReady:  make(chan struct{}),
		element:      e,
	}
}

func (eb *ElementBase) Run(lab *Laboratory, funs ...func(*Laboratory)) {
	defer eb.Completed(lab)
	eb.InputReady()

	if len(funs) == 0 {
		funs = []func(*Laboratory){
			eb.element.Setup,
			eb.element.Steps,
			eb.element.Analysis,
			eb.element.Validation,
		}
	}

	defer func() {
		if res := recover(); res != nil {
			lab.Errorf("%v", res)
			// Use println because of the embedded \n in the Stack Trace
			fmt.Println(lab.lineMapManager.ElementStackTrace())
		}
	}()

	select {
	case <-eb.InputsReady:
		for _, fun := range funs {
			select {
			case <-lab.Errored:
				return
			default:
				fun(lab)
			}
		}
	case <-lab.Errored:
		return
	}
}

func (eb *ElementBase) Completed(lab *Laboratory) {
	if err := eb.Save(lab); err != nil {
		lab.Error(err)
	}
	funs := eb.onExit
	eb.onExit = nil
	for _, fun := range funs {
		fun()
	}
}

func (eb *ElementBase) Save(lab *Laboratory) error {
	if bs, err := json.Marshal(eb.element); err != nil {
		return err
	} else {
		p := filepath.Join(lab.outDir, "elements", fmt.Sprintf("%s.json", eb.element.Name()))
		return ioutil.WriteFile(p, bs, 0400)
	}
}

func (eb *ElementBase) InputReady() {
	if atomic.AddInt64(&eb.pendingCount, -1) == 0 {
		// we've done the transition from 1 -> 0. By definition, we're
		// the only routine that can do that, so we don't need to be
		// careful.
		close(eb.InputsReady)
	}
}

func (eb *ElementBase) AddBlockedInput() {
	atomic.AddInt64(&eb.pendingCount, 1)
}

func (eb *ElementBase) AddOnExit(fun func()) {
	eb.onExit = append(eb.onExit, fun)
}
