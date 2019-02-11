package mfm

import "errors"

// Window represents an atom's event window and
// Provides the public sim interface for atoms.
type Window struct {
	s    *Sim
	loc  C2D
	abuf [41]Atom
	sbuf [41]Site
}

// NeighborDeltas are the neighbors around the origin Site.
var NeighborDeltas = []Site{1, 2, 3, 4, 5, 6, 7, 8}

// InitialWindow is the slice of initial window Site results.
var InitialWindow = []Site{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}

// ErrStopIter instructs At to stop iterating and return.
var ErrStopIter = errors.New("stop iter")

// WriteWindow writes the atoms and sites to event window.
// Note it is not possible to write outside the event window
// using this function because it would panic.
func (w *Window) Write(sites []Site, atoms []Atom) {
	var c C2D
	for i, s := range sites {
		c.SetSite(s)
		c.Add(c, w.loc)
		w.s.set(c, atoms[i])
	}
}

// Diffuse moves the origin atom to a random adjacent valid empty site.
func (w *Window) Diffuse() (site Site, ok bool) {
	if !w.Empty(Me) {
		w.At(NeighborDeltas, func(s Site, a Atom) error {
			if w.Valid(s) && a == nil {
				site = s
				ok = true
				return ErrStopIter
			}
			return nil
		})
	}
	return
}

// At iterates over the sites in random order calling f on each.
// Use StopIter to stop iterating and return immediately.
func (w *Window) At(s []Site, f func(s Site, a Atom) error) error {
	var c C2D
	if len(s) == 1 {
		c.SetSite(s[0])
		c.Add(c, w.loc)
		return f(s[0], w.s.Sites[c])
	}
	copy(w.sbuf[:], s)
	for i := len(s) - 1; i > 0; i-- {
		j := uint32((uint64(w.s.r.Uint32()) * uint64(i+1)) >> 32)
		w.sbuf[i], w.sbuf[j] = w.sbuf[j], w.sbuf[i]
		c.SetSite(w.sbuf[i])
		c.Add(c, w.loc)
		if err := f(w.sbuf[i], w.s.Sites[c]); err != nil {
			if err == ErrStopIter {
				return nil
			}
			return err
		}
	}
	return nil
}

// Valid returns whether the site is within the sim bounds.
func (w Window) Valid(s Site) bool {
	var c C2D
	c.SetSite(s)
	c.Add(c, w.loc)
	return w.s.Bounds == nil ||
		c.X >= 0 && c.Y >= 0 && c.X < w.s.Bounds.X && c.Y < w.s.Bounds.Y
}

// Empty returns whether the event window is empty at the relative site.
func (w *Window) Empty(s Site) bool {
	var c C2D
	c.SetSite(s)
	c.Add(c, w.loc)
	_, ok := w.s.Sites[c]
	return !ok
}

// Roll returns true if u&rand.Int64 != 0 and false otherwise.
func (w *Window) Roll(u uint64) bool {
	return u&w.s.r.Uint64() != 0
}
