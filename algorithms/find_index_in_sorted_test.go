package algorithms

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindIndexInSorted(t *testing.T) {
	createSortedSlc := func(size int) []int {
		slc := make([]int, size)
		for i := 0; i < size; i++ {
			slc[i] = i + 1
		}
		return slc
	}

	tcs := map[string]struct {
		Slc      []int
		Val      int
		Expected int
	}{
		"1M consecutive values": {
			Slc:      createSortedSlc(1_000_000),
			Val:      50_234,
			Expected: 50_233,
		},
		"1M consecutive values - out of bounds at end": {
			Slc:      createSortedSlc(1_000_000),
			Val:      1_000_001,
			Expected: 1_000_000,
		},
		"1M consecutive values - out of bounds at start": {
			Slc:      createSortedSlc(1_000_000),
			Val:      0,
			Expected: 0,
		},
		"6 values missing one": {
			Slc:      []int{1, 2, 3, 4, 6, 7},
			Val:      5,
			Expected: 4,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			idx := FindIndexInSorted(tc.Slc, tc.Val)
			require.Equal(t, tc.Expected, idx)
		})
	}
}
