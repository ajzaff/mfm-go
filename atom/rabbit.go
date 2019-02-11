package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// RabbitInfo contains a herbavore atom info.
var RabbitInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   'R',
	Fg:     termbox.ColorYellow,
	USeq:   "🐇",
	String: "Rabbit",
}

func init() {
	RabbitInfo.Func = Rabbit
}

type rabbitMode int

const (
	rabbitModeStart rabbitMode = 1 << iota
	rabbitModeDiffuse
	rabbitModeHasMate
	rabbitModeHasEaten
)

// Rabbit implements resource atom behavior.
func Rabbit(w *mfm.Window) {}
