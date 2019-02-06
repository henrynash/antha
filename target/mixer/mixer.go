package mixer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	planner "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/target"
)

var (
	_ effects.Device = &Mixer{}
)

// A Mixer is a device plugin for mixer devices
type Mixer struct {
	driver     driver.LiquidhandlingDriver
	properties *driver.LHProperties // Prototype to create fresh properties
	cfg        *GlobalMixerConfig
}

func (a *Mixer) String() string {
	return "Mixer"
}

// CanCompile implements a Device
func (a *Mixer) CanCompile(req effects.Request) bool {
	// TODO: Add specific volume constraints
	can := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1Mixer,
			target.DriverSelectorV1Prompter,
		},
	}
	if a.properties.CanPrompt() {
		can.Selector = append(can.Selector, target.DriverSelectorV1Prompter)
	}
	return can.Contains(req)
}

// FileType returns the file type for generated files
func (a *Mixer) FileType() (ftype string) {
	if m := a.properties.Mnfr; len(m) != 0 {
		ftype = fmt.Sprintf("application/%s", strings.ToLower(m))
	}
	return
}

type lhreq struct {
	*planner.LHRequest     // A request
	*driver.LHProperties   // ... its state
	*planner.Liquidhandler // ... and its associated planner
}

func (a *Mixer) makeLhreq(labEffects *effects.LaboratoryEffects) (*lhreq, error) {
	// MIS -- this might be a hole. We probably need to invoke the sample tracker here
	addPlate := func(req *planner.LHRequest, ip *wtype.Plate) error {
		if _, seen := req.InputPlates[ip.ID]; seen {
			return fmt.Errorf("plate %q already added", ip.ID)
		}
		req.AddUserPlate(labEffects.IDGenerator, ip)
		return nil
	}

	req := planner.NewLHRequest(labEffects.IDGenerator)

	if set := a.cfg.CustomPolicyRuleSet; set != nil {
		req.AddUserPolicies(set)
	}

	if err := req.PolicyManager.SetOption("USE_DRIVER_TIP_TRACKING", a.cfg.UseDriverTipTracking); err != nil {
		return nil, err
	}

	prop := a.properties.Dup(labEffects.IDGenerator)
	plan := planner.Init(prop)

	/* TODO
	if p := a.opt.MaxPlates; p != nil {
		req.InputSetupWeights["MAX_N_PLATES"] = *p
	}

	if p := a.opt.MaxWells; p != nil {
		req.InputSetupWeights["MAX_N_WELLS"] = *p
	}

	if p := a.opt.ResidualVolumeWeight; p != nil {
		req.InputSetupWeights["RESIDUAL_VOLUME_WEIGHT"] = *p
	}

	// TODO -- error check here to prevent nil values

	if p := a.opt.InputPlateTypes; len(p) != 0 {
		for _, v := range p {
			p, err := labEffects.Inventory.PlateTypes.NewPlate(v)
			if err != nil {
				return nil, err
			}

			req.InputPlatetypes = append(req.InputPlatetypes, p)
		}
	}

	if p := a.opt.OutputPlateTypes; len(p) != 0 {
		for _, v := range p {
			p, err := labEffects.Inventory.PlateTypes.NewPlate(v)
			if err != nil {
				return nil, err
			}
			req.OutputPlatetypes = append(req.OutputPlatetypes, p)
		}
	}

	if p := a.opt.TipTypes; len(p) != 0 {
		for _, v := range p {
			t, err := labEffects.Inventory.TipBoxes.NewTipbox(v)
			if err != nil {
				return nil, err
			}
			req.Tips = append(req.Tips, t)
		}
	}
	*/

	if ips := a.opt.InputPlates; len(ips) != 0 {
		for _, ip := range ips {
			if err := addPlate(req, ip); err != nil {
				return nil, err
			}
		}
	}

	// add plates requested via protocol
	for _, p := range labEffects.SampleTracker.GetInputPlates() {
		if err := addPlate(req, p); err != nil {
			return nil, err
		}

	}

	// print instructions?
	req.Options.PrintInstructions = a.cfg.PrintInstructions

	//physical simulation override
	req.Options.IgnorePhysicalSimulation = a.cfg.IgnorePhysicalSimulation

	return &lhreq{
		LHRequest:     req,
		LHProperties:  prop,
		Liquidhandler: plan,
	}, nil
}

// Compile implements a Device
func (a *Mixer) Compile(labEffects *effects.LaboratoryEffects, nodes []effects.Node) ([]effects.Inst, error) {
	var mixes []*wtype.LHInstruction
	for _, node := range nodes {
		if c, ok := node.(*effects.Command); !ok {
			return nil, fmt.Errorf("cannot compile %T", node)
		} else if m, ok := c.Inst.(*wtype.LHInstruction); !ok {
			return nil, fmt.Errorf("cannot compile %T", c.Inst)
		} else {
			mixes = append(mixes, m)
		}
	}

	mix, err := a.makeMix(labEffects, mixes)
	if err != nil {
		return nil, err
	}

	return []effects.Inst{mix}, nil
}

func (a *Mixer) saveFile(name string) ([]byte, error) {
	data, status := a.driver.GetOutputFile()
	if err := status.GetError(); err != nil {
		return nil, err
	} else if len(data) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	bs := []byte(data)

	if err := tw.WriteHeader(&tar.Header{
		Name:    name,
		Mode:    0644,
		Size:    int64(len(bs)),
		ModTime: time.Now(),
	}); err != nil {
		return nil, err
	} else if _, err := tw.Write(bs); err != nil {
		return nil, err
	} else if err := tw.Close(); err != nil {
		return nil, err
	} else if err := gw.Close(); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func mergePolicies(basePolicy, priorityPolicy wtype.LHPolicy) (newPolicy wtype.LHPolicy) {

	newPolicy = make(wtype.LHPolicy)

	for key, value := range priorityPolicy {
		newPolicy[key] = value
	}

	for key, value := range basePolicy {
		if _, found := priorityPolicy[key]; !found {
			newPolicy[key] = value
		}
	}
	return newPolicy
}

// any customised user policies are added to the LHRequest PolicyManager here
// Any component type names with modified policies are iterated until unique i.e. SmartMix_modified_1
func addCustomPolicies(mixes []*wtype.LHInstruction, lhreq *planner.LHRequest) error {
	systemPolicyRuleSet := lhreq.GetPolicyManager().Policies()
	systemPolicies := systemPolicyRuleSet.Policies
	var userPolicies = make(map[string]wtype.LHPolicy)
	var allPolicies = make(map[string]wtype.LHPolicy)
	var liquidClassConversionMap = make(map[string]string)

	for key, value := range systemPolicies {
		allPolicies[key] = value
	}

	userPolicyRuleSet := wtype.NewLHPolicyRuleSet()

	for _, mixInstruction := range mixes {
		for _, component := range mixInstruction.Inputs {
			if len(component.Policy) > 0 {
				if matchingSystemPolicy, found := allPolicies[string(component.Type)]; found {
					mergedPolicy := mergePolicies(matchingSystemPolicy, component.Policy)
					if !wtype.EquivalentPolicies(mergedPolicy, matchingSystemPolicy) {
						num := 1
						newPolicyName := makemodifiedTypeName(component.Type, num)
						existingCustomPolicy, found := allPolicies[newPolicyName]
						for found {
							// check if existing policy with modified name is the same
							if !wtype.EquivalentPolicies(mergedPolicy, existingCustomPolicy) {
								// if not increase number and try again
								num++
								newPolicyName = makemodifiedTypeName(component.Type, num)
								existingCustomPolicy, found = allPolicies[newPolicyName]
							} else {
								// otherwise use existing modified policy
								found = false
							}
						}
						allPolicies[newPolicyName] = mergedPolicy
						userPolicies[newPolicyName] = mergedPolicy
						component.Type = wtype.LiquidType(newPolicyName)
						liquidClassConversionMap[newPolicyName] = matchingSystemPolicy.Name()
					}
				} else {
					allPolicies[string(component.Type)] = component.Policy
					userPolicies[string(component.Type)] = component.Policy
				}
			}
		}
	}

	if len(userPolicies) > 0 {
		userPolicyRuleSet, err := wtype.AddUniversalRules(userPolicyRuleSet, userPolicies)
		if err != nil {
			return err
		}
		for newClass, original := range liquidClassConversionMap {
			err := wtype.CopyRulesFromPolicy(userPolicyRuleSet, original, newClass)
			if err != nil {
				return err
			}
		}
		lhreq.AddUserPolicies(userPolicyRuleSet)
	}

	return nil
}

const modifiedPolicySuffix = "_modified_"

func makemodifiedTypeName(componentType wtype.LiquidType, number int) string {
	return string(componentType) + modifiedPolicySuffix + strconv.Itoa(number)
}

// unModifyTypeName will trim a _modified_ suffix from a LiquidType in the CSV file.
// These are added to LiquidType names when a Liquid is modified in an element.
func unModifyTypeName(componentType string) string {
	return strings.Split(componentType, modifiedPolicySuffix)[0]
}

func (a *Mixer) makeMix(labEffects *effects.LaboratoryEffects, mixes []*wtype.LHInstruction) (*target.Mix, error) {
	hasPlate := func(plates []*wtype.Plate, typ, id string) bool {
		for _, p := range plates {
			if p.Type == typ && (id == "" || p.ID == id) {
				return true
			}
		}
		return false
	}

	getID := func(mixes []*wtype.LHInstruction) (r wtype.BlockID) {
		m := make(map[wtype.BlockID]bool)
		for _, mix := range mixes {
			m[mix.BlockID] = true
		}
		for k := range m {
			r = k
			break
		}
		return
	}

	r, err := a.makeLhreq(labEffects)
	if err != nil {
		return nil, err
	}

	// any customised user policies are added to the LHRequest PolicyManager here
	// Any component type names with modified policies are iterated until unique i.e. SmartMix_modified_1
	if err := addCustomPolicies(mixes, r.LHRequest); err != nil {
		return nil, err
	}

	for _, m := range mixes {
		if m.OutPlate != nil {
			p, ok := r.LHRequest.OutputPlates[m.OutPlate.ID]
			if ok && p != m.OutPlate {
				return nil, fmt.Errorf("Mix setup error: Plate %s already requested in different state for mix.", p.ID)
			}
			r.LHRequest.OutputPlates[m.OutPlate.ID] = m.OutPlate
		}
	}

	r.LHRequest.BlockID = getID(mixes)

	for _, mix := range mixes {
		if len(mix.Platetype) != 0 && !hasPlate(r.LHRequest.OutputPlatetypes, mix.Platetype, mix.PlateID) {
			p, err := labEffects.Inventory.PlateTypes.NewPlate(mix.Platetype)
			if err != nil {
				return nil, err
			}
			p.ID = mix.PlateID
			r.LHRequest.OutputPlatetypes = append(r.LHRequest.OutputPlatetypes, p)
		}
		r.LHRequest.Add_instruction(mix)
	}

	err = r.Liquidhandler.MakeSolutions(labEffects, r.LHRequest)
	// TODO: MIS unfortunately we need to make sure this stays up to date would
	// be better to remove this and just use the ones the liquid handler holds
	r.LHProperties = r.Liquidhandler.Properties

	if err != nil {
		return nil, err
	}

	name := a.opt.DriverOutputFileName
	if len(name) == 0 {
		// TODO: Desired filename not exposed in current driver interface, so pick
		// a name. So far, at least Gilson software cares what the filename is, so
		// use .sqlite for compatibility
		name = strings.Replace(fmt.Sprintf("%s.sqlite", time.Now().Format(time.RFC3339)), ":", "_", -1)
	}

	tarball, err := a.saveFile(name)
	if err != nil {
		return nil, err
	}

	return &target.Mix{
		Dev:             a,
		Request:         r.LHRequest,
		Properties:      r.LHProperties,
		FinalProperties: r.Liquidhandler.FinalProperties,
		Final:           r.Liquidhandler.PlateIDMap(),
		Files: target.Files{
			Tarball: tarball,
			Type:    a.FileType(),
		},
	}, nil
}

// New creates a new Mixer
func New(cfg *GlobalMixerConfig, d driver.LiquidhandlingDriver) (*Mixer, error) {
	if p, status := d.GetCapabilities(); !status.Ok() {
		return nil, status.GetError()
	} else {
		p.Driver = d
		return &Mixer{driver: d, properties: &p, cfg: cfg}, nil
	}
}
