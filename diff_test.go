package diff

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestIntSlices(t *testing.T) {
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

func TestFloat64Slices(t *testing.T) {
	tests := []struct {
		name string
		a, b []float64
		want string
	}{
		{
			name: "equal",
			a:    []float64{74.3, math.Inf(1), -784.0, -12.13, math.Inf(-1), -959.7485},
			b:    []float64{74.3, math.Inf(1), -784.0, -12.13, math.Inf(-1), -959.7485},
			want: ``,
		},
		{
			name: "totally different",
			a:    []float64{74.3, -784.0, math.NaN(), math.Inf(-1), -959.7485},
			b:    []float64{74.8, math.Inf(1), math.NaN(), -856.0, -959.7415, -9.19},
			want: `0d 1d 2d 3d 4d 5i(0) 5i(1) 5i(2) 5i(3) 5i(4) 5i(5)`,
		},
		{
			name: "somewhat similar",
			a:    []float64{74.3, math.Inf(1), -932.15, math.NaN(), math.Inf(-1), 3.16},
			b:    []float64{74.3, 23, -784.0, math.NaN(), math.Inf(-1), -959.7485},
			want: `1d 2d 3d 4i(1) 4i(2) 4i(3) 5d 6i(5)`,
		},
	}

	for _, tt := range tests {
		got := Float64Slices(tt.a, tt.b)
		want, err := parseScript(tt.want)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s:\n got %s\nwant %s", tt.name, got, want)
		}
	}
}

func TestStringSlices(t *testing.T) {
	tests := []struct {
		name string
		a, b []string
		want string
	}{
		{
			name: "equal",
			a:    []string{"", "R-BDYZWYuwtf4LSc4fepTtddAWN-_b", "AA", "BB", "\r\n\a\t"},
			b:    []string{"", "R-BDYZWYuwtf4LSc4fepTtddAWN-_b", "AA", "BB", "\r\n\a\t"},
			want: ``,
		},
		{
			name: "totally different",
			a:    []string{"", "a", "l2wbOryHIR6dqNh", "zXDL7Z6lEAuzEbO", "lKUDks1r6BQuiHF"},
			b:    []string{"MuV3BeGPyV0mZaC", "No7hw2kmcMNa_CB", "zZ-Ofu9Cgngz_9P", "03ffA1DOPGQge2O", "tTMNvwE59_Cs6Lb", "xy"},
			want: `0d 1d 2d 3d 4d 5i(0) 5i(1) 5i(2) 5i(3) 5i(4) 5i(5)`,
		},
		{
			name: "somewhat similar",
			a:    []string{"all", "work", "and", "KXZYNbiN5kvFajd", "play"},
			b:    []string{" ", "work", "and", "play", "makes", "FoROCdv", "a", "--"},
			want: `0d 1i(0) 3d 5i(4) 5i(5) 5i(6) 5i(7)`,
		},
	}

	for _, tt := range tests {
		got := StringSlices(tt.a, tt.b)
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
		return fmt.Sprintf("%di(%d)", e.Index, e.Bindex)
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

func BenchmarkSimilarIntSlices(b *testing.B) {
	x := []int{9, 17, 6, 60, 3, 0, 17, 4, 19, 20, 1, 21, 4, 6}
	y := []int{0, 7, 72, 60, 3, 1, 34, 9, 68, 7, 8, 17, 4, 19, 20, 76, 7, 21, 4, 6}
	for i := 0; i < b.N; i++ {
		IntSlices(x, y)
	}
}

func BenchmarkEqualIntSlices1K(b *testing.B) {
	const n = 1 << 10
	x := make([]int, n)
	for i := range x {
		x[i] = i + 10
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntSlices(x, x)
	}
}

func BenchmarkEqualIntSlices64K(b *testing.B) {
	const n = 1 << 16
	x := make([]int, n)
	for i := range x {
		x[i] = i + 16
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntSlices(x, x)
	}
}

func BenchmarkTotallyDifferentIntSlices1K(b *testing.B) {
	const n = 1 << 10
	x, y := make([]int, n), make([]int, n)
	for i := range x {
		x[i] = i - 7
		y[i] = i + 13
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntSlices(x, y)
	}
}

func BenchmarkTotallyDifferentIntSlices64K(b *testing.B) {
	const n = 1 << 16
	x, y := make([]int, n), make([]int, n)
	for i := range x {
		x[i] = i - 7
		y[i] = i + 13
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IntSlices(x, y)
	}
}
