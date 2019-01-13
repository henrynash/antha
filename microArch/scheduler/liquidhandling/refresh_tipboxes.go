package liquidhandling

import "github.com/antha-lang/antha/laboratory/effects/id"

// no longer need to supply tipboxes after the fact
func (lh *Liquidhandler) Refresh_tipboxes_tipwastes(idGen *id.IDGenerator, rq *LHRequest) {

	// dead simple

	lh.FinalProperties.RemoveTipBoxes()

	for pos := range lh.Properties.PosLookup {
		tb, ok := lh.Properties.Tipboxes[pos]

		if ok {
			newTb := tb.Dup(idGen)
			lh.FinalProperties.AddTipBoxTo(pos, newTb)
			lh.plateIDMap[tb.ID] = newTb.ID
			tb.Refresh(idGen)
			continue
		}

		tw, ok := lh.Properties.Tipwastes[pos]

		if ok {
			// swap the wastes
			tw2 := lh.FinalProperties.Tipwastes[pos]
			tw2.Contents = tw.Contents
			tw.Empty()
			lh.plateIDMap[tw.ID] = tw2.ID
		}
	}
}
