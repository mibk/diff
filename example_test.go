package diff_test

import (
	"fmt"

	"github.com/mibk/diff"
)

func Example_simpleDiffOutput() {
	a := []int{'A', 'B', 'C', 'A', 'B', 'B', 'A', 'D'}
	b := []int{'C', 'B', 'A', 'B', 'A', 'C', 'D'}
	eds := diff.IntSlices(a, b)

	var i int
	for _, ed := range eds {
		for ; i < ed.Index; i++ {
			fmt.Printf(" %c\n", a[i])
		}
		if ed.Op == diff.Delete {
			fmt.Printf("-%c\n", a[i])
			i++
		} else {
			fmt.Printf("+%c\n", b[ed.Arg])
		}
	}
	for ; i < len(a); i++ {
		fmt.Printf(" %c\n", a[i])
	}

	// output:
	// -A
	// -B
	//  C
	// +B
	//  A
	//  B
	// -B
	//  A
	// +C
	//  D
}
