package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects"
)

type MessageInstruction struct {
	BaseRobotInstruction
	*InstructionType
	Message     string
	PassThrough map[string]*wtype.Liquid
}

func NewMessageInstruction(lhi *wtype.LHInstruction) *MessageInstruction {
	msi := &MessageInstruction{
		InstructionType: MSG,
	}
	msi.BaseRobotInstruction = NewBaseRobotInstruction(msi)

	pt := make(map[string]*wtype.Liquid)

	if lhi != nil {
		for i := 0; i < len(lhi.Inputs); i++ {
			pt[lhi.Inputs[i].ID] = lhi.Outputs[i]
		}
		msi.Message = lhi.Message
		msi.WaitTime = lhi.WaitTime
		msi.PassThrough = pt
	}

	return msi
}

func (ins *MessageInstruction) Visit(visitor RobotInstructionVisitor) {
	visitor.Message(ins)
}

func (msi *MessageInstruction) Generate(labEffects *effects.LaboratoryEffects, policy *wtype.LHPolicyRuleSet, prms *LHProperties) ([]RobotInstruction, error) {
	// use side effect to keep IDs straight

	prms.UpdateComponentIDs(msi.PassThrough)
	return nil, nil
}

func (msi *MessageInstruction) GetParameter(name InstructionParameter) interface{} {
	switch name {
	case MESSAGE:
		return msi.Message
	default:
		return msi.BaseRobotInstruction.GetParameter(name)
	}
}

func (msi *MessageInstruction) OutputTo(driver LiquidhandlingDriver) error {
	//level int, title, text string, showcancel bool

	if msi.Message != wtype.MAGICBARRIERPROMPTSTRING {
		return driver.Message(0, "", msi.Message, false).GetError()
	}
	return nil
}
