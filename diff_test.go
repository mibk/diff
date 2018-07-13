package diff

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want string
	}{
		{
			name: "delete on start",
			a:    []int{1, 2, 3, 1, 2, 2, 1},
			b:    []int{3, 2, 1, 2, 1, 3},
			want: `0d 1d 3i(1) 5d 7i(5)`,
		},
		{
			name: "common on start",
			a:    []int{1, 2, 4, 3, 2},
			b:    []int{1, 2, 3, 1},
			want: `2d 4d 5i(3)`,
		},
		{
			name: "equal",
			a:    []int{3, 1, 4, 1, 5, 9, 2},
			b:    []int{3, 1, 4, 1, 5, 9, 2},
			want: ``,
		},
		{
			name: "all changes",
			a:    []int{3, 1, 4, 1, 5},
			b:    []int{8, 6, 7, 0, -2},
			want: `0d 1d 2d 3d 4d 5i(0) 5i(1) 5i(2) 5i(3) 5i(4)`,
		},
		{
			name: "chain of repeats",
			a:    []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			b:    []int{1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1},
			want: `5i(5) 10d`,
		},
		{
			name: "short a",
			a:    []int{11},
			b:    []int{0, 11, 1, 2, 3, 4, 5, 6, 7},
			want: `0i(0) 1i(2) 1i(3) 1i(4) 1i(5) 1i(6) 1i(7) 1i(8)`,
		},
		{
			name: "start and end",
			a:    []int{0, 1, 2, 3, 4, 5, 6},
			b:    []int{1, 2, 3, 4, 5, 6, 11},
			want: `0d 7i(6)`,
		},
	}

	for _, tt := range tests {
		got := IntSlices(tt.a, tt.b)
		want, err := parseScript(tt.want)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s:\n got %s\nwant %s", tt.name, got, want)
		}
	}
}

func (e Edit) String() string {
	switch e.Op {
	case Delete:
		return fmt.Sprintf("%dd", e.Index)
	case Insert:
		return fmt.Sprintf("%di(%d)", e.Index, e.Arg)
	default:
		panic("unknown op")
	}
}

var editRx = regexp.MustCompile(`^(\d+)([di])(?:\((\d+)\))?$`)

func parseScript(s string) ([]Edit, error) {
	var eds []Edit
	for _, cmd := range strings.Fields(s) {
		m := editRx.FindStringSubmatch(cmd)
		if m == nil {
			return nil, fmt.Errorf("invalid format: %s", cmd)
		}
		i, _ := strconv.Atoi(m[1])
		op := Delete
		if m[2] == "i" {
			op = Insert
		}
		arg, _ := strconv.Atoi(m[3])
		eds = append(eds, Edit{i, op, arg})
	}
	return eds, nil
}
