package mfm

// Site is a relative offset in the event window.
//             38
//          31 22 33
//       25 15 10 17 27
//    29 13  5  2  7 19 35
// 37 21  9  1 *0  4 12 24 40
//    30 14  6  3  8 20 36
//       26 16 11 18 28
//          32 23 34
//             39
type Site int

// Me represents the event window origin at 0.
var Me Site

var xs = []int{0, -1, 0, 0, 1, -1, -1, 1, 1, -2, 0, 0, 2, -2, -2, -1, -1, 1, 1, 2, 2, -3, 0, 0, 3, -2, -2, 2, 2, -3, -3, -1, -1, 1, 1, 3, 3, -4, 0, 0, 4}
var ys = []int{0, 0, 1, -1, 0, 1, -1, 1, -1, 0, 2, -2, 0, 1, -1, 2, -2, 2, -2, 1, -1, 0, 3, -3, 0, 2, -2, 2, -2, 1, -1, 3, -3, 3, -3, 1, -1, 0, 4, -4, 0}
var dist = []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}

// Dist returns the distance to another site.
func (s Site) Dist(s1 Site) int {
	d := (xs[s1] - xs[s]) + (ys[s1] - ys[s])
	if d < 0 {
		return -d
	}
	return d
}

// C2D used to represent an absolute coordinate.
// This should not be leaked to atom implementations.
type C2D struct{ X, Y int }

// SetSite updates the value of c with the value of the site converted to a C2D.
func (c *C2D) SetSite(s Site) { c.X = xs[s]; c.Y = ys[s] }

// SetC2D updates the value of c with the value of other.
func (c *C2D) SetC2D(other C2D) { c.X = other.X; c.Y = other.Y }

// Set updates the value of c with x, y.
func (c *C2D) Set(x, y int) { c.X = x; c.Y = y }

// Add updates the value of c to c1+c2 and returns the result.
func (c *C2D) Add(c1, c2 C2D) C2D {
	c.X = c1.X + c2.X
	c.Y = c1.Y + c2.Y
	return *c
}
