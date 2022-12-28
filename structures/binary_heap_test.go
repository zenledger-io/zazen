package structures

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sort"
	"testing"
)

func TestBinaryHeap(t *testing.T) {
	tcs := map[string]struct {
		Size int
	}{
		"binary head of size 5k": {
			Size: 5000,
		},
	}
	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			bh := &binaryHeap[int]{
				compFn: func(t1 int, t2 int) int {
					return t1 - t2
				},
			}
			for i := 0; i < tc.Size; i++ {
				bh.Push(rand.Int())
			}
			for i := 0; i < bh.Len(); i++ {
				li := 2*i + 1
				ri := 2*i + 2
				if li < bh.Len() {
					require.LessOrEqual(t, bh.slc[li], bh.slc[i])
				}
				if ri < bh.Len() {
					require.LessOrEqual(t, bh.slc[ri], bh.slc[i])
				}
			}
			var prev int
			for i := 0; i < bh.Len(); i++ {
				el := bh.Pop()
				if i == 0 {
					prev = el
					continue
				}

				require.LessOrEqual(t, el, prev)
				prev = el
			}
		})
	}
}

func BenchmarkBinaryHeap(b *testing.B) {
	shift := func(slc []int) (int, []int) {
		// simulate current tax calc hifo
		sort.Slice(slc, func(i, j int) bool {
			return slc[i] > slc[j]
		})
		return slc[0], slc[1:]
	}

	createSlice := func(limit int) []int {
		slc := make([]int, limit)
		for i := 0; i < limit; i++ {
			slc[i] = rand.Int()
		}
		return slc
	}
	tcs := map[string]struct {
		Slice []int
	}{
		"5k items": {
			Slice: createSlice(5_000),
		},
		"10k items": {
			Slice: createSlice(10_000),
		},
	}

	for desc, tc := range tcs {
		b.Run(fmt.Sprintf("binary heap %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				bh := NewBinaryHeap(func(t1 int, t2 int) int {
					return t1 - t2
				})
				for _, i := range tc.Slice {
					bh.Push(i)
				}
				for i := 0; i < bh.Len(); i++ {
					bh.Pop()
				}
			}
		})
		b.Run(fmt.Sprintf("slice %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				var slc []int
				for _, i := range tc.Slice {
					slc = append(slc, i)
				}
				for i := 0; i < len(slc); i++ {
					_, slc = shift(slc)
				}
			}
		})
		b.Run(fmt.Sprintf("binary heap push and pop %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				bh := NewBinaryHeap(func(t1 int, t2 int) int {
					return t1 - t2
				})
				for j, i := range tc.Slice {
					bh.Push(i)
					if j%2 == 0 {
						bh.Pop()
					}
				}
				j := 0
				for bh.Len() > 0 {
					el := bh.Pop()
					if j%2 == 0 {
						bh.Push(el)
					}
					j += 1
				}
			}
		})
		b.Run(fmt.Sprintf("slice push and pop %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				var slc []int
				for j, i := range tc.Slice {
					slc = append(slc, i)
					if j%2 == 0 {
						_, slc = shift(slc)
					}
				}
				j := 0
				for len(slc) > 0 {
					var el int
					el, slc = shift(slc)
					if j%2 == 0 {
						slc = append(slc, el)
					}
					j += 1
				}
			}
		})
	}
}
