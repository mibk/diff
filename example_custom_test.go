package diff_test

import (
	"fmt"
	"strings"

	"github.com/mibk/diff"
)

type ignoreCase struct {
	a, b []string
}

func (ic ignoreCase) Lens() (n, m int) { return len(ic.a), len(ic.b) }
func (ic ignoreCase) Equal(i, j int) bool {
	return strings.ToLower(ic.a[i]) == strings.ToLower(ic.b[j])
}

// missingIgnoreCase returns elements that are move or deleted from a
// compared to b.
func missingIgnoreCase(a, b []string) []string {
	eds := diff.Diff(ignoreCase{a, b})

	var miss []string
	for _, ed := range eds {
		if ed.Op == diff.Insert {
			miss = append(miss, b[ed.Arg])
		}
	}
	return miss
}

func Example_customType() {
	a := []string{"black", "#31ad1d", "#8923dd", "#baddad", "yellow"}
	b := []string{"#31AD1D", "#8924dd", "#BadDad", "black", "YELLOW"}

	for _, m := range missingIgnoreCase(a, b) {
		fmt.Println("-", m)
	}

	// output:
	// - #8924dd
	// - black
}
