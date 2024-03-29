package structures

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkStack(b *testing.B) {
	pop := func(slc []int) (int, []int) {
		return slc[len(slc)-1], slc[0 : len(slc)-1]
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
		b.Run(fmt.Sprintf("stack %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				s := NewStack[int]()
				for _, i := range tc.Slice {
					s.Push(i)
				}
				for i := 0; i < s.Len(); i++ {
					s.Pop()
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
					_, slc = pop(slc)
				}
			}
		})
		b.Run(fmt.Sprintf("stack push and pop %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				s := NewStack[int]()
				for j, i := range tc.Slice {
					s.Push(i)
					if j%2 == 0 {
						s.Pop()
					}
				}
				j := 0
				for s.Len() > 0 {
					el, _ := s.Pop()
					if j%2 == 0 {
						s.Push(el)
					}
					j += 1
				}
				for i := 0; i < s.Len(); i++ {
					s.Pop()
				}
			}
		})
		b.Run(fmt.Sprintf("slice push and pop %v", desc), func(b *testing.B) {
			for r := 0; r < b.N; r++ {
				var slc []int
				for j, i := range tc.Slice {
					slc = append(slc, i)
					if j%2 == 0 {
						_, slc = pop(slc)
					}
				}
				j := 0
				for len(slc) > 0 {
					var el int
					el, slc = pop(slc)
					if j%2 == 0 {
						slc = append(slc, el)
					}
					j += 1
				}
			}
		})
	}
}
