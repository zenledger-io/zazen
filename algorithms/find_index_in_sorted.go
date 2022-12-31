package algorithms

import "golang.org/x/exp/constraints"

// FindIndexInSorted will find the index in a sorted slice that you would
// insert a value into. The value at that index will be greater than or equal
// to the value you are searching for. If the index is equal to the length
// of the slice, then all value in the slice are less than the value
// you are searching for.
func FindIndexInSorted[T constraints.Ordered](slc []T, val T) int {
	l, r := 0, len(slc)-1
	for l <= r {
		mid := (r + l) / 2
		if slc[mid] == val {
			return mid
		}

		if slc[mid] < val {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return r + 1
}
