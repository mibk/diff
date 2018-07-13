// Package diff implements methods for comparing objects and producing
// edit scripts.
//
// The implementation is based on the algorithm described in the paper
// "An O(ND) Difference Algorithm and Its Variations" by Eugene W.
// Myers, Algorithmica Vol. 1 No. 2, 1986, p. 251.
package diff

type Operation int

const (
	Delete Operation = iota
	Insert
)

type Edit struct {
	Index int
	Op    Operation
	Arg   int // only valid for Insert
}

type Data interface {
	Lens() (n, m int)
	Equal(i, j int) bool
}

func Diff(data Data) []Edit {
	N, M := data.Lens()
	MAX := N + M

	V := map[int]int{1: 0}
	for D := 0; D <= MAX; D++ {
		for k := -D; k <= D; k += 2 {
			var x, y int
			vert := true
			if k == -D || k != D && V[k-1] < V[k+1] {
				x = V[k+1]
			} else {
				x = V[k-1] + 1
				vert = false
			}

			y = x - k
			x0, y0 := x, y
			for x < N && y < M && data.Equal(x, y) {
				x, y = x+1, y+1
			}
			V[k] = x

			if x >= N && y >= M {
				if D > 0 {
					x1, y1 := x0, y0
					if vert {
						y1--
					} else {
						x1--
					}
					ed := Diff(&bounded{data, x1, y1})

					var e Edit
					if vert {
						e = Edit{x0, Insert, y0 - 1}
					} else {
						e = Edit{x0 - 1, Delete, 0}
					}
					return append(ed, e)
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

func IntSlices(a, b []int) []Edit { return Diff(&intSlices{a, b}) }

type intSlices struct {
	a, b []int
}

func (p *intSlices) Lens() (n, m int)    { return len(p.a), len(p.b) }
func (p *intSlices) Equal(i, j int) bool { return p.a[i] == p.b[j] }
