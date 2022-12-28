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
		b.Run(fmt.Sprintf("queue push and pop %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				q := NewQueue[int]()
				for j, i := range tc.Slice {
					q.Enqueue(i)
					if j%2 == 0 {
						q.Dequeue()
					}
				}
				j := 0
				for q.Len() > 0 {
					el, _ := q.Dequeue()
					if j%2 == 0 {
						q.Prepend(el)
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
						slc = append([]int{el}, slc...)
					}
					j += 1
				}
			}
		})
	}
}
