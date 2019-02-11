// Package rand provides a xorshift PRNG.
//
// Based on original code written and dedicated
// to the public domain by Sebastiano Vigna (2014).
// It has the following characteristics:
//
// - Outputs 64-bit numbers.
// - Passes Dieharder and SmallCrush test batteries.
// - Does not require warm-up, no zeroland to escape.
// - Internal state is a single 64-bit integer.
// - Period is 2^64 - 1.
//
// For further analysis see:
//   <http://vigna.di.unimi.it/ftp/papers/xorshift.pdf>
package rand

// Rand implements a xorshift64star PRNG.
type Rand uint64

const mask63 = 1<<63 - 1

// Int63 generates a random 63-bit number.
func (r *Rand) Int63() int64 {
	return int64(r.Uint64() & mask63)
}

// Uint64 generates a random 64-bit number.
func (r *Rand) Uint64() uint64 {
	*r ^= Rand(uint64(*r) >> 12)
	*r ^= Rand(uint64(*r) << 25)
	*r ^= Rand(uint64(*r) >> 27)
	return uint64(*r) * 2685821657736338717
}

// Uint32 generates a random 31-bit number.
func (r *Rand) Uint32() uint32 {
	return uint32(r.Int63() >> 31)
}

// Seed implements Seed for rand.Source.
func (r *Rand) Seed(s uint64) {
	*r = Rand(s)
}
