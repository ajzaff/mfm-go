package atom

import (
	"ajz_xyz/experimental/computation/mfm-go"

	"github.com/nsf/termbox-go"
)

var PaperInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   '[',
	Fg:     termbox.ColorWhite,
	USeq:   "📄",
	String: "Paper",
}

func init() {
	PaperInfo.Func = Paper
}

func Paper(w *mfm.Window) {
	defer w.Stop()
	switch a := w.Atom(); a.Type {
	case nil:
		w.Swap(0)
	default:
		if a.Type.ID == RockInfo.ID {
			w.Set(mfm.Atom{Type: &PaperInfo})
		}
	}
}
