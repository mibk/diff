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
const (
	Delete Operation = iota
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
	max := n + m

	endp := map[int]int{1: 0}
	for d := 0; d <= max; d++ {
		for k := -d; k <= d; k += 2 {
			var x, y int
			vert := true
			if k == -d || k != d && endp[k-1] < endp[k+1] {
				x = endp[k+1]
			} else {
				x = endp[k-1] + 1
				vert = false
			}

			y = x - k
			x0, y0 := x, y
			for x < n && y < m && data.Equal(x, y) {
				x, y = x+1, y+1
			}
			endp[k] = x

			if x >= n && y >= m {
				if d > 0 {
					x1, y1 := x0, y0
					if vert {
						y1--
					} else {
						x1--
					}
					eds := Diff(&bounded{data, x1, y1})

					var ed Edit
					if vert {
						ed = Edit{x0, Insert, y0 - 1}
					} else {
						ed = Edit{x0 - 1, Delete, 0}
					}
					return append(eds, ed)
				}
				return nil
			}
		}
	}
	panic("unreachable")
}

type bounded struct {
	Data
	n, m int
}

func (b *bounded) Lens() (n, m int) { return b.n, b.m }

// IntSlices creates an edit script for two slices of ints.
func IntSlices(a, b []int) []Edit { return Diff(&intSlices{a, b}) }

type intSlices struct {
	a, b []int
}

func (p *intSlices) Lens() (n, m int)    { return len(p.a), len(p.b) }
func (p *intSlices) Equal(i, j int) bool { return p.a[i] == p.b[j] }
