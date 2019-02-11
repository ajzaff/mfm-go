package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// GermInfo contains info related to Germ atom.
var GermInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   '7',
	Fg:     termbox.ColorGreen,
	USeq:   "🌱",
	String: "Germ",
}

func init() {
	GermInfo.Func = Germ
}

var newGerm = mfm.Atom{Type: &GermInfo, Value: 0}

var germCount = Var{10, 10}

// Germ implements germ atom behavior.
func Germ(w *mfm.Window) {
	me := w.Self()
	a := w.Atom()
	var seenGerm bool
	switch {
	case a.Type == nil:
		if w.Roll(65535) {
			w.Set(newGerm)
			w.Stop()
		}
	case w.Site() == 0 && germCount.ReadUint64(me) >= 8:
		w.Set(mfm.Atom{Type: &FernInfo, Value: 0})
	case a.Type.ID == GermInfo.ID:
		seenGerm = true
		// w.SetValue(me.Value)
	}
	if seenGerm {
		germCount.CountUint64(&me, 1, 8)
	} else {
		germCount.CountUint64(&me, 0, 8)
	}
}
