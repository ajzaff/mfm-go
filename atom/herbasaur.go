package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// HerbasaurInfo contains herbavore atom info.
var HerbasaurInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   'H',
	Fg:     termbox.ColorYellow,
	USeq:   "🦕",
	String: "Herbasaur",
}

func init() {
	HerbasaurInfo.Func = Herbasaur
}

type herbasaurMode int

const (
	herbasaurModeStart herbasaurMode = 1 << iota
	herbasaurModeDiffuse
	herbasaurModeHasMate
	herbasaurModeHasEaten
)

// Herbasaur implements resource atom behavior.
func Herbasaur(w *mfm.Window) {}
