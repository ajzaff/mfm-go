package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// TreeInfo contains tree atom info.
var TreeInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   '9',
	Fg:     termbox.ColorGreen,
	USeq:   "🌳",
	String: "Tree",
}

func init() {
	TreeInfo.Func = Tree
}

// Tree implements tree atom behavior.
func Tree(w *mfm.Window) {}
