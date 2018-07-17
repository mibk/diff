// Package diff implements methods for comparing objects and producing
// edit scripts.
//
// The implementation is based on the algorithm described in the paper
// "An O(ND) Difference Algorithm and Its Variations" by Eugene W.
// Myers, Algorithmica Vol. 1 No. 2, 1986, p. 251.
package diff

// Operation represents an operation in the edit script.
type Operation int

// The list of possible operations.
//
// None is never used by the package. It can be used as a sentinel
// operation for Edit.Op; the value is reserved. See the simple diff
// output example for illustration of usage.
const (
	None Operation = iota
	Delete
	Insert
)

// An Edit represents an element in the edit script, produced by the
// Diff function. An edit script represents a set of operations describing
// how to convert the first collection into the second one.
type Edit struct {
	Index int // index where the operation should occur
	Op    Operation
	Arg   int // only valid for Insert
}

// Data is the interface that is used by the Diff function to produce
// an edit script. A type that satisfies the Data interface is typically
// a wrapper around two collections. The method requires that the
// elements of the collections be enumerated by integer indexes.
type Data interface {
	// Lens returns the lengths of the underlying collections.
	Lens() (n, m int)
	// Equal reports whether the elements from the two collections
	// with indexes i and j are equal.
	Equal(i, j int) bool
}

// Diff creates an edit script for data. The edit script represents a
// set of operations describing how to convert from the first collection
// into the second one.
func Diff(data Data) []Edit {
	n, m := data.Lens()
	return diff(data, n, m, make([]int, 2*(n+m)+1))
}

func diff(data Data, n, m int, endp []int) []Edit {
	max := n + m
	if max == 0 {
		return nil
	}

	endp[max+1] = 0
	for d := 0; d <= max; d++ {
		for k := -d; k <= d; k += 2 {
			j := max + k
			var x int
			vert := true
			if k == -d || k != d && endp[j-1] < endp[j+1] {
				x = endp[j+1]
			} else {
				x = endp[j-1] + 1
				vert = false
			}

			y := x - k

			x0, y0 := x, y
			for x < n && y < m && data.Equal(x, y) {
				x, y = x+1, y+1
			}
			endp[j] = x

			if x >= n && y >= m {
				if d > 0 {
					x1, y1 := x0, y0
					if vert {
						y1--
					} else {
						x1--
					}
					eds := diff(data, x1, y1, endp[:len(endp)-2])

					var ed Edit
					if vert {
						ed = Edit{x1, Insert, y1}
					} else {
						ed = Edit{x1, Delete, 0}
					}
					return append(eds, ed)
				}
				return nil
			}
		}
	}
	panic("unreachable")
}

// IntSlices creates an edit script for two slices of ints.
func IntSlices(a, b []int) []Edit { return Diff(&intSlices{a, b}) }

type intSlices struct {
	a, b []int
}

func (p *intSlices) Lens() (n, m int)    { return len(p.a), len(p.b) }
func (p *intSlices) Equal(i, j int) bool { return p.a[i] == p.b[j] }

// Float64Slices creates an edit script for two slices of float64s.
func Float64Slices(a, b []float64) []Edit { return Diff(&float64Slices{a, b}) }

type float64Slices struct {
	a, b []float64
}

func (p *float64Slices) Lens() (n, m int)    { return len(p.a), len(p.b) }
func (p *float64Slices) Equal(i, j int) bool { return p.a[i] == p.b[j] }

// StringSlices creates an edit script for two slices of strings.
func StringSlices(a, b []string) []Edit { return Diff(&stringSlices{a, b}) }

type stringSlices struct {
	a, b []string
}

func (p *stringSlices) Lens() (n, m int)    { return len(p.a), len(p.b) }
func (p *stringSlices) Equal(i, j int) bool { return p.a[i] == p.b[j] }
