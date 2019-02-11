package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// DRegInfo implements a dynamic space regulator atom.
// Roles include producing Res, random destruction,
// and occasional reproduction.
var DRegInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   'D',
	Fg:     termbox.ColorMagenta,
	USeq:   "D",
	String: "DReg",
}

func init() {
	DRegInfo.Func = DReg
}

var (
	newRes  = mfm.Atom{Type: &ResInfo, Value: 0}
	newDReg = mfm.Atom{Type: &DRegInfo, Value: 0}
)

// DReg defines the update function for DReg atoms.
func DReg(w *mfm.Window) {
	a, s := w.Atom(), w.Site()
	if a.Type == nil {
		if w.Roll(511) { // create Res
			w.Set(newRes)
		} else if w.Roll(2047) { // create Dreg
			w.Set(newDReg)
		} else {
			w.Swap(0)
		}
		w.Stop()
	} else if s != 0 && ((a.Type.ID == DRegInfo.ID && w.Roll(7)) || w.Roll(255)) {
		w.Remove() // random destruction
		w.Stop()
	}
}
