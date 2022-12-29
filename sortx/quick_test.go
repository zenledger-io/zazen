package sortx

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sort"
	"testing"
)

func TestQuick(t *testing.T) {
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
		"one": {
			Slice:    []int{1},
			Expected: []int{1},
		},
		"empty": {
			Slice:    []int{},
			Expected: []int{},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			Quick(tc.Slice, func(a, b int) bool {
				return a < b
			})
			fmt.Printf("%v\n", tc.Slice)
			require.Equal(t, len(tc.Expected), len(tc.Slice))
			for i, el := range tc.Expected {
				require.Equal(t, el, tc.Slice[i])
			}
			Quick(tc.Slice, func(a, b int) bool {
				return a > b
			})
			require.Equal(t, len(tc.Expected), len(tc.Slice))
			for i := len(tc.Expected) - 1; i >= 0; i-- {
				require.Equal(t, tc.Expected[i], tc.Slice[len(tc.Expected)-1-i])
			}
		})
	}
}

func BenchmarkQuick(b *testing.B) {
	createSlice := func(limit int) []int {
		slc := make([]int, limit, limit+1)
		for i := 0; i < limit; i++ {
			slc[i] = rand.Int()
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
		b.Run(fmt.Sprintf("quick sort %v", desc), func(b *testing.B) {
			slcs := make([][]int, b.N)
			for i := 0; i < b.N; i++ {
				slc := make([]int, len(tc.Slice))
				copy(slc, tc.Slice)
				slcs[i] = slc
			}
			b.ResetTimer()

			for r := 0; r < b.N; r++ {
				Quick(slcs[r], func(a, b int) bool {
					return a < b
				})
			}
		})
		b.Run(fmt.Sprintf("built in sort %v", desc), func(b *testing.B) {
			slcs := make([][]int, b.N)
			for i := 0; i < b.N; i++ {
				slc := make([]int, len(tc.Slice))
				copy(slc, tc.Slice)
				slcs[i] = slc
			}
			b.ResetTimer()

			for r := 0; r < b.N; r++ {
				slc := slcs[r]
				sort.Slice(slc, func(a, b int) bool {
					return slc[a] < slc[b]
				})
			}
		})
	}
}
