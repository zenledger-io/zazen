package sortx

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestMerge(t *testing.T) {
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
			slc := Merge(tc.Slice, func(a, b int) bool {
				return a < b
			})
			require.Equal(t, len(tc.Expected), len(slc))
			for i, el := range tc.Expected {
				require.Equal(t, el, slc[i])
			}
			slc = Merge(tc.Slice, func(a, b int) bool {
				return a > b
			})
			require.Equal(t, len(tc.Expected), len(slc))
			for i := len(tc.Expected) - 1; i >= 0; i-- {
				require.Equal(t, tc.Expected[i], slc[len(tc.Expected)-1-i])
			}
		})
	}
}

func BenchmarkMerge(b *testing.B) {
	createSlice := func(limit int) []int {
		slc := make([]int, limit, limit+1)
		for i := 0; i < limit; i++ {
			slc[i] = i
		}
		return slc
	}
	tcs := map[string]struct {
		Slice    []int
		Expected []int
	}{
		"10k items": {
			Slice: createSlice(10_000),
		},
		"100k items": {
			Slice: createSlice(100_000),
		},
		"1M items": {
			Slice: createSlice(1_000_000),
		},
	}

	for desc, tc := range tcs {
		b.Run(fmt.Sprintf("merge sort %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				_ = Merge(tc.Slice, func(a, b int) bool {
					return a < b
				})
			}
		})
		b.Run(fmt.Sprintf("built in sort %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				sort.Ints(tc.Slice)
			}
		})
	}
}
