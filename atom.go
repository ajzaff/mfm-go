package mfm

import "github.com/nsf/termbox-go"

// Atom defines the interface for an atom in the simulation
// Having a 64 bit type and 64 bit Value.
type Atom struct {
	Type  *AtomInfo
	Value uint64
}

var empty Atom

// AtomInfo defines a structure containing atom metadata.
type AtomInfo struct {
	ID     uint64
	Func   func(*Window)
	Radius int
	Rune   rune
	Fg     termbox.Attribute
	Bg     termbox.Attribute
	USeq   string
	String string
}
