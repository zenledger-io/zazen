package sortx

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestInsertion(t *testing.T) {
	tcs := map[string]struct {
		Slice    []int
		Expected []int
	}{
		"even number with duplicates": {
			Slice:    []int{-1, 1, 0, 10, 300, 4, -20, -20, 10, 3},
			Expected: []int{-20, -20, -1, 0, 1, 3, 4, 10, 10, 300},
		},
		"odd number with duplicates": {
			Slice:    []int{-1, 1, 0, 10, 300, 4, -20, -20, 10, 3, 50},
			Expected: []int{-20, -20, -1, 0, 1, 3, 4, 10, 10, 50, 300},
		},
		"empty": {
			Slice:    []int{},
			Expected: []int{},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			Insertion(tc.Slice, func(a, b int) bool {
				return a < b
			})
			require.Equal(t, len(tc.Expected), len(tc.Slice))
			for i, el := range tc.Expected {
				require.Equal(t, el, tc.Slice[i])
			}
			Insertion(tc.Slice, func(a, b int) bool {
				return a > b
			})
			require.Equal(t, len(tc.Expected), len(tc.Slice))
			for i := len(tc.Expected) - 1; i >= 0; i-- {
				require.Equal(t, tc.Expected[i], tc.Slice[len(tc.Expected)-1-i])
			}
		})
	}
}

func TestInsert(t *testing.T) {
	tcs := map[string]struct {
		Slice    []int
		Append   int
		Expected []int
	}{
		"even number with duplicates": {
			Slice:    []int{-1, 1, 0, 10, 300, 4, -20, -20, 10},
			Append:   3,
			Expected: []int{-20, -20, -1, 0, 1, 3, 4, 10, 10, 300},
		},
		"odd number with duplicates": {
			Slice:    []int{-1, 1, 0, 10, 300, 4, -20, -20, 10, 50},
			Append:   3,
			Expected: []int{-20, -20, -1, 0, 1, 3, 4, 10, 10, 50, 300},
		},
		"empty": {
			Slice:    []int{},
			Append:   3,
			Expected: []int{3},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			slc := Insert(tc.Slice, tc.Append, func(a, b int) bool {
				return a < b
			})
			require.Equal(t, len(tc.Expected), len(slc))
			for i, el := range tc.Expected {
				require.Equal(t, el, slc[i])
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	createSlice := func(limit int) []int {
		slc := make([]int, limit, limit+1)
		for i := 0; i < limit; i++ {
			slc[i] = i
		}
		return slc
	}
	tcs := map[string]struct {
		Slice    []int
		Append   int
		Expected []int
	}{
		"10k items": {
			Slice:  createSlice(10_000),
			Append: 3,
		},
		"100k items": {
			Slice:  createSlice(100_000),
			Append: 3,
		},
		"1M items": {
			Slice:  createSlice(1_000_000),
			Append: 3,
		},
	}

	for desc, tc := range tcs {
		b.Run(fmt.Sprintf("insertion sort %v", desc), func(b *testing.B) {
			_ = Insert(tc.Slice, tc.Append, func(a, b int) bool {
				return a < b
			})
		})
		b.Run(fmt.Sprintf("built in sort %v", desc), func(b *testing.B) {
			sort.Ints(append(tc.Slice, tc.Append))
		})
	}
}
