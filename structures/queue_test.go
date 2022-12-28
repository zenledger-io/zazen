package structures

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkQueue(b *testing.B) {
	shift := func(slc []int) (int, []int) {
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
		b.Run(fmt.Sprintf("queue %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				q := NewQueue[int]()
				for _, i := range tc.Slice {
					q.Enqueue(i)
				}
				for i := 0; i < q.Len(); i++ {
					q.Dequeue()
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
	}
}
