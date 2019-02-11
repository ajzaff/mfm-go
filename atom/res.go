package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// ResInfo defines info relating to Res atoms.
var ResInfo = mfm.AtomInfo{
	Rune:   'r',
	Fg:     termbox.ColorYellow,
	USeq:   "🔋",
	String: "Res",
}
