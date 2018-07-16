package diff_test

import (
	"fmt"

	"github.com/mibk/diff"
)

func Example_simpleDiffOutput() {
	a := []string{"Alice", "Bob", "Cyril", "Alice", "Bob", "Bob", "Alice", "Daniel"}
	b := []string{"Cyril", "Bob", "Alice", "Bob", "Alice", "Cyril", "Daniel"}
	eds := diff.StringSlices(a, b)

	var i int
	for _, ed := range eds {
		for ; i < ed.Index; i++ {
			fmt.Printf(" %s\n", a[i])
		}
		if ed.Op == diff.Delete {
			fmt.Printf("-%s\n", a[i])
			i++
		} else {
			fmt.Printf("+%s\n", b[ed.Arg])
		}
	}
	for ; i < len(a); i++ {
		fmt.Printf(" %s\n", a[i])
	}

	// output:
	// -Alice
	// -Bob
	//  Cyril
	// +Bob
	//  Alice
	//  Bob
	// -Bob
	//  Alice
	// +Cyril
	//  Daniel
}
