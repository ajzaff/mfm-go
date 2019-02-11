package atom

import "ajz_xyz/experimental/computation/mfm-go"

// Var defines a variable encoded in the bits of an atom.
// It supports reading and writing a variety of datatypes.
type Var struct{ Pos, Size uint64 }

// Mask generates a mask at position i with size n.
// It panics with an unsigned underflow error when n > i.
func (v Var) Mask() uint64 { return (1<<v.Pos - 1) ^ (1<<(v.Pos-v.Size) - 1) }

// ReadUint64 reads an unsigned integer at the given bit
// Position i with the given bit length b. The integer
// is generally encoded using WriteUint which always
// Encodes using Big-endian and no overflow.
// ReadUint returns the value that was read.
func (v Var) ReadUint64(a mfm.Atom) uint64 { return a.Value & v.Mask() >> (v.Pos - v.Size) }

// WriteUint64 writes the value x to the given
// Bit position i with the given bit length b.
// WriteUint returns the value that was writen.
// The value can be read back using ReadUint64.
func (v Var) WriteUint64(a *mfm.Atom, x uint64) uint64 {
	x &= 1<<v.Size - 1
	a.Value = a.Value&^v.Mask() | x<<(v.Pos-v.Size)
	return x
}

// CountUint64 records x occurances of the event and advances the sliding window
// With the width w. This approximates the count of the event in the window.
// This returns the new value of the counter.
func (v Var) CountUint64(a *mfm.Atom, x, w uint64) uint64 {
	// Take (w-1) samples from the current count.
	// Take    1  sample from the new count.
	val := v.ReadUint64(*a)   // current count
	val = ((w-1)*val + x) / w // new count
	v.WriteUint64(a, val)
	return val
}
