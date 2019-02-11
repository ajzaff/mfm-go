package mfm

import "ajz_xyz/experimental/computation/mfm-go/rand"

// Window represents an atom's event window and
// Provides the public sim interface for atoms.
type Window struct {
	rand rand.Rand
	loc  C2D
	mut  mutation
	me   Atom
	s    Site
	a    Atom
}

// Mutation defines a mutation that commits a change to the event window.
type mutation struct {
	mode mutMode
	dst  Site
	atom Atom
}

type mutMode int

const (
	mutNone mutMode = 0
	mutSet          = 1 << iota
	mutMove
	mutSwap
	mutStop
)

func (w *Window) Roll(u uint64) bool { return u&w.rand.Uint64() == 0 }
func (w *Window) Site() Site         { return w.s }
func (w *Window) Atom() Atom         { return w.a }
func (w *Window) Self() Atom         { return w.me }
func (w *Window) Reset()             { w.mut = mutation{} }
func (w *Window) Stop()              { w.mut.mode |= mutStop }
func (w *Window) Remove()            { w.Set(empty) }
func (w *Window) Set(a Atom)         { w.mut.mode = mutSet; w.mut.atom = a }
func (w *Window) Move(dst Site)      { w.mut.mode = mutMove; w.mut.dst = dst }
func (w *Window) Swap(dst Site)      { w.mut.mode = mutSwap; w.mut.dst = dst }
